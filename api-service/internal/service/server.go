package service

import (
	"context"
	"fmt"
	"net"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"api-service/internal/config"
	"api-service/internal/constants"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"api-service/internal/interface/repository"
	"api-service/internal/interface/service"
	"api-service/internal/model"
	"api-service/pkg/errors"
	"api-service/pkg/logger"
	sshpkg "api-service/pkg/ssh"
)

// 常量定义，避免魔术数字
const (
	// 状态常量
	statusFailed = "failed"

	// Secret key type constants
	secretKeyTypeAccount = "ACCOUNT"

	// 服务器操作时间常量
	serverRestartTime  = 2 * time.Second
	serverShutdownTime = 3 * time.Second
	systemUpdateTime   = 10 * time.Second
	maintenanceTime    = 500 * time.Millisecond

	// 网络连接常量
	tcpConnectionTimeout = 10 * time.Second

	// 错误代码
	recordNotFoundErrorCode = 4000

	// Mock数据常量
	mockServerUptime    = 86400 // 24小时，单位秒
	mockContainersCount = 5
	mockImagesCount     = 12

	// Server code 前缀
	serverCodePrefix = "srv"
	serverCodeLength = 16 // srv + 13位时间戳 = 16字符
)

// serverService implements service.ServerService
type serverService struct {
	serverRepo       repository.ServerRepository
	secretKeyService service.SecretKeyService          // 使用真实的密钥管理服务
	systemConfigRepo repository.SystemConfigRepository // 添加系统配置仓库
	logger           logger.Logger
	config           *config.Config // 添加配置依赖
}

// ServerServiceConfig contains configuration for server service
type ServerServiceConfig struct {
	Logger           logger.Logger
	ServerRepo       repository.ServerRepository
	SecretKeyService service.SecretKeyService          // 使用真实的密钥管理服务
	SystemConfigRepo repository.SystemConfigRepository // 添加系统配置仓库
	Config           *config.Config                    // 添加配置依赖
}

// NewServerService creates a new server service instance
func NewServerService(config ServerServiceConfig) service.ServerService {
	return &serverService{
		serverRepo:       config.ServerRepo,
		secretKeyService: config.SecretKeyService,
		systemConfigRepo: config.SystemConfigRepo,
		logger:           config.Logger,
		config:           config.Config,
	}
}

// CreateServer creates a new server with the current user as owner
func (s *serverService) CreateServer(ctx context.Context, req *request.CreateServerRequest, currentUserID uint) (*response.ServerResponse, error) {
	s.logger.InfoContext(ctx, "Creating new server",
		logger.String("name", req.Name),
		logger.Uint("owner_id", currentUserID))

	// Check if server name already exists
	exists, err := s.serverRepo.ExistsServerByName(ctx, req.Name)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to check server name existence", logger.ErrorField(err))
		return nil, errors.NewAppError(errors.CodeInternalError)
	}
	if exists {
		return nil, errors.NewAppError(errors.CodeResourceAlreadyExists)
	}

	// Generate unique server code
	serverCode := s.generateServerCode()

	// Create server model with owner_id from JWT token
	server := &model.Server{
		Name:            req.Name,
		Code:            serverCode,
		Host:            req.Host,
		SSHPort:         req.SSHPort,
		ResourceGroupID: req.ResourceGroupID,
		OwnerID:         currentUserID, // Auto-populated from current user
		Description:     &req.Description,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		// Hostname, InternalIP, IPv6Address are dynamically collected by Agent
	}

	// Create server in database
	if err := s.serverRepo.CreateServer(ctx, server); err != nil {
		s.logger.ErrorContext(ctx, "Failed to create server in database", logger.ErrorField(err))
		return nil, errors.NewAppError(errors.CodeRecordCreateFailed)
	}

	// Create SSH credential secret if provided
	var credentialID string
	if req.SSHUsername != "" && (req.SSHPassword != "" || req.SSHKey != "") {
		credID, err := s.createSSHCredentialSecret(ctx, server, req, currentUserID)
		if err != nil {
			// Log warning but don't fail server creation
			s.logger.WarnContext(ctx, "Failed to create SSH credential secret",
				logger.String("serverName", server.Name),
				logger.ErrorField(err))
		} else {
			credentialID = credID
		}
	}

	// Create secret reference if credential was created
	if credentialID != "" {
		if err := s.createSecretReference(ctx, serverCode, credentialID); err != nil {
			// Log warning but don't fail server creation
			s.logger.WarnContext(ctx, "Failed to create secret reference for server",
				logger.String("serverCode", serverCode),
				logger.String("credentialId", credentialID),
				logger.ErrorField(err))
		}
	}

	// Test connectivity asynchronously
	go s.testServerConnectivity(context.Background(), server.ID)

	s.logger.InfoContext(ctx, "Server created successfully",
		logger.Uint("id", server.ID),
		logger.String("name", server.Name),
		logger.String("code", serverCode))

	// Convert to response
	return s.convertToServerResponse(server), nil
}

// GetServer retrieves a server by ID (设计文档6.3.3: MySQL基本信息 + Redis状态数据)
func (s *serverService) GetServer(ctx context.Context, id uint) (*response.ServerResponse, error) {
	// 1. 从MySQL获取基本配置信息
	server, err := s.serverRepo.GetServerByID(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get server",
			logger.Uint("id", id),
			logger.ErrorField(err))
		return nil, err
	}

	// 2. 转换基本信息
	serverResp := s.convertToServerResponse(server)

	// 3. 从Redis获取状态数据 (模拟实现)
	s.enrichServerWithStatusData(ctx, serverResp, id)

	return serverResp, nil
}

// UpdateServer updates an existing server
//
//nolint:gocyclo,gocognit // Complex business logic requires multiple conditional checks
func (s *serverService) UpdateServer(ctx context.Context, id uint, req *request.UpdateServerRequest) (*response.ServerResponse, error) {
	s.logger.InfoContext(ctx, "Updating server", logger.Uint("id", id))

	// Get existing server
	server, err := s.serverRepo.GetServerByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Track if SSH credential changed
	var credentialChanged bool
	var newCredentialID string

	// Check if new name conflicts with existing servers (excluding current server)
	if req.Name != nil && *req.Name != server.Name {
		exists, err := s.serverRepo.ExistsServerByName(ctx, *req.Name, id)
		if err != nil {
			return nil, errors.NewAppError(errors.CodeInternalError)
		}
		if exists {
			return nil, errors.NewAppError(errors.CodeResourceAlreadyExists)
		}
		server.Name = *req.Name
	}

	// Update other fields
	// Hostname, InternalIP, IPv6Address are dynamically collected by Agent, not updated via API
	if req.Host != nil {
		server.Host = *req.Host
	}
	if req.SSHPort != nil {
		server.SSHPort = *req.SSHPort
	}

	// Check if SSH credentials are being updated
	if req.SSHUsername != nil || req.SSHPassword != nil || req.SSHKey != nil {
		// Get old credential references first
		oldRefs, _ := s.secretKeyService.GetSecretReferencesByResourceCode(ctx, server.Code)
		if len(oldRefs) > 0 {
			credentialChanged = true
		}

		// Create new credential if username is provided
		if req.SSHUsername != nil && *req.SSHUsername != "" {
			// Create UpdateServerRequest wrapper for credential creation
			createReq := &request.CreateServerRequest{
				SSHUsername: *req.SSHUsername,
			}
			if req.SSHPassword != nil {
				createReq.SSHPassword = *req.SSHPassword
			}
			if req.SSHKey != nil {
				createReq.SSHKey = *req.SSHKey
			}

			credID, err := s.createSSHCredentialSecret(ctx, server, createReq, server.OwnerID)
			if err != nil {
				s.logger.WarnContext(ctx, "Failed to create new SSH credential",
					logger.String("serverName", server.Name),
					logger.ErrorField(err))
			} else {
				newCredentialID = credID
				credentialChanged = true
			}
		}
	}

	if req.Description != nil {
		server.Description = req.Description
	}

	server.UpdatedAt = time.Now()

	// Update server in database
	if err := s.serverRepo.UpdateServer(ctx, server); err != nil {
		return nil, err
	}

	// Update secret reference if SSH credential changed
	if credentialChanged {
		// Delete old reference
		s.deleteSecretReferencesByServerCode(ctx, server.Code)

		// Create new reference if new credential was created
		if newCredentialID != "" {
			if err := s.createSecretReference(ctx, server.Code, newCredentialID); err != nil {
				s.logger.WarnContext(ctx, "Failed to create new secret reference",
					logger.String("serverCode", server.Code),
					logger.String("credentialId", newCredentialID),
					logger.ErrorField(err))
			}
		}
	}

	s.logger.InfoContext(ctx, "Server updated successfully",
		logger.Uint("id", id),
		logger.Bool("credentialChanged", credentialChanged))

	return s.convertToServerResponse(server), nil
}

// DeleteServer deletes a server
func (s *serverService) DeleteServer(ctx context.Context, id uint) error {
	s.logger.InfoContext(ctx, "Deleting server", logger.Uint("id", id))

	// Get server first to retrieve server code
	server, err := s.serverRepo.GetServerByID(ctx, id)
	if err != nil {
		return err
	}

	// Delete secret references for this server
	s.deleteSecretReferencesByServerCode(ctx, server.Code)

	// Delete server from database
	if err := s.serverRepo.DeleteServer(ctx, id); err != nil {
		return err
	}

	s.logger.InfoContext(ctx, "Server deleted successfully",
		logger.Uint("id", id),
		logger.String("code", server.Code))

	return nil
}

// ListServers lists servers with pagination and filtering
func (s *serverService) ListServers(ctx context.Context, req *request.ListServersRequest) (*response.ServerListResponse, error) {
	servers, total, err := s.serverRepo.ListServers(ctx, req)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to list servers", logger.ErrorField(err))
		return nil, errors.NewAppError(errors.CodeRecordQueryFailed)
	}

	// Convert to response format
	serverResponses := make([]response.ServerResponse, len(servers))
	for i, server := range servers {
		serverResponses[i] = *s.convertToServerResponse(server)
	}

	return &response.ServerListResponse{
		Items:    serverResponses,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// CheckServersStatus checks the status of multiple servers (设计文档序号6)
// GetServerStatus checks the status of a single server (设计文档序号6.3.6)
func (s *serverService) GetServerStatus(ctx context.Context, id uint) (*response.ServerStatusCheckResult, error) {
	s.logger.InfoContext(ctx, "Checking server status",
		logger.Uint("serverId", id))

	// Use default check types (all) and default timeout
	checkTypes := []string{"ssh", "agent", "docker"}
	timeout := s.getSSHTimeout()

	result := s.checkSingleServerStatus(ctx, id, checkTypes, timeout)

	// If server not found or check failed, return error
	if result.Status == statusFailed {
		return nil, errors.NewAppError(errors.CodeRecordNotFound)
	}

	return &result, nil
}

// CheckServersStatus checks the status of multiple servers (设计文档序号6)
func (s *serverService) CheckServersStatus(ctx context.Context, req *request.ServerStatusCheckRequest) (*response.BatchServerStatusResponse, error) {
	results := make([]response.ServerStatusCheckResult, 0, len(req.ServerIDs))
	successCount := 0
	failedCount := 0

	for _, serverID := range req.ServerIDs {
		result := s.checkSingleServerStatus(ctx, serverID, req.CheckTypes, req.Timeout)
		if result.Status == statusFailed {
			failedCount++
		} else {
			successCount++
		}
		results = append(results, result)
	}

	return &response.BatchServerStatusResponse{
		SuccessCount: successCount,
		FailedCount:  failedCount,
		Results:      results,
	}, nil
}

// ExecuteServerActions executes actions on multiple servers (设计文档序号7)
func (s *serverService) ExecuteServerActions(ctx context.Context, req *request.ServerActionRequest) (*response.ServerActionResponse, error) {
	s.logger.InfoContext(ctx, "Executing server actions",
		logger.String("action", req.Action),
		logger.Int("serverCount", len(req.ServerIDs)))

	// 生成任务ID
	taskID := fmt.Sprintf("task_%d", time.Now().Unix())

	// 使用 goroutine 进行异步执行，实现解耦
	go s.executeServerActionsAsync(context.Background(), req, taskID)

	return &response.ServerActionResponse{
		TaskID:      taskID,
		Action:      req.Action,
		ServerCount: len(req.ServerIDs),
		Status:      "submitted",
		Message:     fmt.Sprintf("Action '%s' submitted for %d servers", req.Action, len(req.ServerIDs)),
	}, nil
}

// UploadFile uploads a file to server via SSH (设计文档序号8)
func (s *serverService) UploadFile(ctx context.Context, serverID uint, filePath string, fileData []byte) (*response.ServerFileUploadResponse, error) {
	s.logger.InfoContext(ctx, "Starting file upload",
		logger.Uint("serverId", serverID),
		logger.String("filePath", filePath),
		logger.Int("fileSize", len(fileData)))

	// 安全校验：路径、扩展名、文件大小
	if err := s.validateFileSecurity(ctx, filePath, int64(len(fileData)), true); err != nil {
		return nil, err
	}

	// Get server info
	server, err := s.serverRepo.GetServerByID(ctx, serverID)
	if err != nil {
		return nil, err
	}

	// Get user ID from context
	userID, ok := ctx.Value("user_id").(uint)
	if !ok {
		s.logger.ErrorContext(ctx, "User ID not found in context")
		return nil, errors.NewAppError(errors.CodeAccessDenied)
	}

	// Get SSH credentials from secret key service
	creds, err := s.getSSHCredentials(ctx, server, userID)
	if err != nil {
		return nil, err
	}

	// Create SSH client with credentials
	client, err := s.createSSHClientWithCreds(ctx, server, creds.Username, creds.Password, creds.PrivateKey)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	// Upload file via SFTP
	if err := client.UploadFile(fileData, filePath); err != nil {
		s.logger.ErrorContext(ctx, "Failed to upload file via SFTP",
			logger.Uint("serverId", serverID),
			logger.String("filePath", filePath),
			logger.ErrorField(err))
		return nil, errors.NewAppErrorWrapError(err, errors.CodeServerSSHConnectionFailed)
	}

	s.logger.InfoContext(ctx, "File uploaded successfully",
		logger.Uint("serverId", serverID),
		logger.String("filePath", filePath),
		logger.Int64("fileSize", int64(len(fileData))))

	return &response.ServerFileUploadResponse{
		ServerID:   serverID,
		FilePath:   filePath,
		FileSize:   int64(len(fileData)),
		Message:    fmt.Sprintf("File uploaded successfully to %s", server.Name),
		UploadedAt: time.Now(),
	}, nil
}

// DownloadFile downloads a file from server via SSH (设计文档序号9)
func (s *serverService) DownloadFile(ctx context.Context, serverID uint, filePath string) (fileData []byte, fileName string, err error) {
	s.logger.InfoContext(ctx, "Starting file download",
		logger.Uint("serverId", serverID),
		logger.String("filePath", filePath))

	// 安全校验：路径和扩展名
	if errValidate := s.validateFilePath(ctx, filePath); errValidate != nil {
		return nil, "", errValidate
	}
	if errValidate := s.validateFileExtension(ctx, filePath); errValidate != nil {
		return nil, "", errValidate
	}

	// Get server info
	server, err := s.serverRepo.GetServerByID(ctx, serverID)
	if err != nil {
		return nil, "", err
	}

	// Get user ID from context
	userID, ok := ctx.Value("user_id").(uint)
	if !ok {
		s.logger.ErrorContext(ctx, "User ID not found in context")
		return nil, "", errors.NewAppError(errors.CodeAccessDenied)
	}

	// Get SSH credentials from secret key service
	creds, err := s.getSSHCredentials(ctx, server, userID)
	if err != nil {
		return nil, "", err
	}

	// Create SSH client with credentials
	client, err := s.createSSHClientWithCreds(ctx, server, creds.Username, creds.Password, creds.PrivateKey)
	if err != nil {
		return nil, "", err
	}
	defer client.Close()

	// Check if file exists
	exists, err := client.FileExists(filePath)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to check file existence",
			logger.Uint("serverId", serverID),
			logger.String("filePath", filePath),
			logger.ErrorField(err))
		return nil, "", errors.NewAppErrorWrapError(err, errors.CodeServerSSHConnectionFailed)
	}
	if !exists {
		return nil, "", errors.NewAppError(errors.CodeResourceNotFound)
	}

	// Get file size for security check
	fileSize, err := client.GetFileSize(filePath)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get file size",
			logger.Uint("serverId", serverID),
			logger.String("filePath", filePath),
			logger.ErrorField(err))
		return nil, "", errors.NewAppErrorWrapError(err, errors.CodeServerSSHConnectionFailed)
	}

	// Validate file size
	if errValidate := s.validateFileSecurity(ctx, filePath, fileSize, false); errValidate != nil {
		return nil, "", errValidate
	}

	// Download file via SFTP
	fileData, err = client.DownloadFile(filePath)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to download file via SFTP",
			logger.Uint("serverId", serverID),
			logger.String("filePath", filePath),
			logger.ErrorField(err))
		return nil, "", errors.NewAppErrorWrapError(err, errors.CodeServerSSHConnectionFailed)
	}

	// Verify downloaded size matches expected size
	if int64(len(fileData)) != fileSize {
		s.logger.WarnContext(ctx, "Downloaded file size mismatch",
			logger.Uint("serverId", serverID),
			logger.String("filePath", filePath),
			logger.Int64("expected", fileSize),
			logger.Int("actual", len(fileData)))
	}

	filename := filepath.Base(filePath)
	s.logger.InfoContext(ctx, "File downloaded successfully",
		logger.Uint("serverId", serverID),
		logger.String("filePath", filePath),
		logger.Int64("fileSize", int64(len(fileData))))

	return fileData, filename, nil
}

// DeleteFile deletes a file from server via SSH (设计文档序号10)
func (s *serverService) DeleteFile(ctx context.Context, serverID uint, filePath string) (*response.ServerFileDeleteResponse, error) {
	s.logger.InfoContext(ctx, "Starting file deletion",
		logger.Uint("serverId", serverID),
		logger.String("filePath", filePath))

	// 安全校验：路径和扩展名
	if err := s.validateFilePath(ctx, filePath); err != nil {
		return nil, err
	}
	if err := s.validateFileExtension(ctx, filePath); err != nil {
		return nil, err
	}

	// Get server info
	server, err := s.serverRepo.GetServerByID(ctx, serverID)
	if err != nil {
		return nil, err
	}

	// Get user ID from context
	userID, ok := ctx.Value("user_id").(uint)
	if !ok {
		s.logger.ErrorContext(ctx, "User ID not found in context")
		return nil, errors.NewAppError(errors.CodeAccessDenied)
	}

	// Get SSH credentials from secret key service
	creds, err := s.getSSHCredentials(ctx, server, userID)
	if err != nil {
		return nil, err
	}

	// Create SSH client with credentials
	client, err := s.createSSHClientWithCreds(ctx, server, creds.Username, creds.Password, creds.PrivateKey)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	// Delete file via SFTP
	if err := client.DeleteFile(filePath); err != nil {
		s.logger.ErrorContext(ctx, "Failed to delete file via SFTP",
			logger.Uint("serverId", serverID),
			logger.String("filePath", filePath),
			logger.ErrorField(err))
		return nil, errors.NewAppErrorWrapError(err, errors.CodeServerSSHConnectionFailed)
	}

	s.logger.InfoContext(ctx, "File deleted successfully",
		logger.Uint("serverId", serverID),
		logger.String("filePath", filePath))

	return &response.ServerFileDeleteResponse{
		ServerID:  serverID,
		FilePath:  filePath,
		Message:   fmt.Sprintf("File deleted successfully from %s", server.Name),
		DeletedAt: time.Now(),
	}, nil
}

// Helper methods

// checkSingleServerStatus checks status for a single server (for batch status check)
func (s *serverService) checkSingleServerStatus(ctx context.Context, serverID uint, checkTypes []string, timeout int) response.ServerStatusCheckResult {
	// Get server info
	server, err := s.serverRepo.GetServerByID(ctx, serverID)
	if err != nil {
		return response.ServerStatusCheckResult{
			ServerID:     serverID,
			ServerName:   fmt.Sprintf("Server-%d", serverID),
			Status:       "failed",
			ErrorCode:    recordNotFoundErrorCode, // CodeRecordNotFound
			ErrorMessage: "Server not found",
			Checks:       map[string]response.ServiceStatusInfo{},
		}
	}

	checks := make(map[string]response.ServiceStatusInfo)

	// Default to all checks if not specified
	if len(checkTypes) == 0 {
		checkTypes = []string{"ssh", "agent", "docker"}
	}

	// Define check type constants
	const (
		checkTypeSSH    = "ssh"
		checkTypeAgent  = "agent"
		checkTypeDocker = "docker"
	)

	// Perform requested checks
	for _, checkType := range checkTypes {
		switch checkType {
		case checkTypeSSH:
			checks[checkTypeSSH] = s.checkSSHConnectivity(ctx, server, timeout)
		case checkTypeAgent:
			checks[checkTypeAgent] = s.checkAgentStatus(ctx, server)
		case checkTypeDocker:
			checks[checkTypeDocker] = s.checkDockerStatus(ctx, server)
		}
	}

	return response.ServerStatusCheckResult{
		ServerID:   serverID,
		ServerName: server.Name,
		Status:     "success",
		Checks:     checks,
	}
}

// checkSSHConnectivity tests SSH connectivity with real authentication
func (s *serverService) checkSSHConnectivity(ctx context.Context, server *model.Server, timeout int) response.ServiceStatusInfo {
	if timeout == 0 {
		timeout = s.getSSHTimeout()
	}

	now := time.Now()

	// Get SSH credential from secret references
	refs, refErr := s.secretKeyService.GetSecretReferencesByResourceCode(ctx, server.Code)
	var credentialID uint
	if refErr == nil && len(refs) > 0 {
		credentialID = refs[0].SecretID
	}

	// Check if credentials are configured
	if credentialID == 0 {
		// No credentials configured, do TCP connectivity test only
		start := time.Now()
		tcpConn, dialErr := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", server.Host, server.SSHPort), time.Duration(timeout)*time.Second)
		responseTime := int(time.Since(start).Milliseconds())

		if dialErr != nil {
			return response.ServiceStatusInfo{
				Status:       constants.SSHStatusDisconnected,
				ResponseTime: responseTime,
				CheckedAt:    &now,
				ErrorMessage: "SSH port not reachable (credentials not configured)",
			}
		}
		if closeErr := tcpConn.Close(); closeErr != nil {
			s.logger.WarnContext(ctx, "Failed to close TCP connection", logger.ErrorField(closeErr))
		}

		return response.ServiceStatusInfo{
			Status:       constants.SSHStatusConnected,
			ResponseTime: responseTime,
			CheckedAt:    &now,
			ErrorMessage: "TCP connectivity OK (SSH auth test requires credentials)",
		}
	}

	// Try to get user ID from context for SSH authentication
	userID, hasUserID := ctx.Value("user_id").(uint)

	// If we have SSH credentials configured and user ID, try SSH authentication
	if hasUserID {
		// Get SSH credentials from secret key service
		creds, credsErr := s.getSSHCredentials(ctx, server, userID)
		if credsErr != nil {
			s.logger.WarnContext(ctx, "Failed to get SSH credentials for connectivity test, falling back to TCP test",
				logger.Uint("serverId", server.ID),
				logger.ErrorField(credsErr))
		} else {
			// Perform SSH authentication test
			start := time.Now()
			client, err := s.createSSHClientWithCreds(ctx, server, creds.Username, creds.Password, creds.PrivateKey)
			responseTime := int(time.Since(start).Milliseconds())

			if err != nil {
				return response.ServiceStatusInfo{
					Status:       constants.SSHStatusDisconnected,
					ResponseTime: responseTime,
					CheckedAt:    &now,
					ErrorMessage: fmt.Sprintf("SSH authentication failed: %v", err),
				}
			}
			defer client.Close()

			// Test connection by executing a simple command
			if err := client.TestConnection(); err != nil {
				return response.ServiceStatusInfo{
					Status:       constants.SSHStatusDisconnected,
					ResponseTime: responseTime,
					CheckedAt:    &now,
					ErrorMessage: fmt.Sprintf("SSH connection test failed: %v", err),
				}
			}

			return response.ServiceStatusInfo{
				Status:       constants.SSHStatusConnected,
				ResponseTime: responseTime,
				CheckedAt:    &now,
			}
		}
	}

	// Fall back to simple TCP connectivity test
	start := time.Now()
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", server.Host, server.SSHPort), time.Duration(timeout)*time.Second)
	responseTime := int(time.Since(start).Milliseconds())

	if err != nil {
		return response.ServiceStatusInfo{
			Status:       constants.SSHStatusDisconnected,
			ResponseTime: responseTime,
			CheckedAt:    &now,
			ErrorMessage: "Network connection failed (user context not available for SSH auth test)",
		}
	}
	if closeErr := conn.Close(); closeErr != nil {
		s.logger.WarnContext(ctx, "Failed to close TCP connection", logger.ErrorField(closeErr))
	}

	return response.ServiceStatusInfo{
		Status:       constants.SSHStatusConnected,
		ResponseTime: responseTime,
		CheckedAt:    &now,
	}
}

// checkAgentStatus checks Agent status (mock implementation)
func (s *serverService) checkAgentStatus(_ context.Context, _ *model.Server) response.ServiceStatusInfo {
	now := time.Now()
	// Mock: Simulate Agent status check
	return response.ServiceStatusInfo{
		Status:        constants.AgentStatusOnline,
		CheckedAt:     &now,
		LastHeartbeat: &now,
		Version:       "v1.2.3",
		Uptime:        mockServerUptime,
	}
}

// checkDockerStatus checks Docker status (mock implementation)
func (s *serverService) checkDockerStatus(_ context.Context, _ *model.Server) response.ServiceStatusInfo {
	now := time.Now()
	// Mock: Simulate Docker status check
	return response.ServiceStatusInfo{
		Status:          constants.DockerStatusRunning,
		CheckedAt:       &now,
		Version:         "20.10.17",
		ContainersCount: mockContainersCount,
		ImagesCount:     mockImagesCount,
	}
}

func (s *serverService) testServerConnectivity(ctx context.Context, serverID uint) {
	server, err := s.serverRepo.GetServerByID(ctx, serverID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get server for connectivity test",
			logger.Uint("serverId", serverID),
			logger.ErrorField(err))
		return
	}

	status := s.testServerConnectivitySync(ctx, server)

	if err := s.serverRepo.UpdateServerStatus(ctx, serverID, status); err != nil {
		s.logger.ErrorContext(ctx, "Failed to update server status after connectivity test",
			logger.Uint("serverId", serverID),
			logger.ErrorField(err))
	}
}

func (s *serverService) testServerConnectivitySync(ctx context.Context, server *model.Server) string {
	// Test basic network connectivity first
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", server.Host, server.SSHPort), tcpConnectionTimeout)
	if err != nil {
		return constants.ServerStatusUnreachable
	}
	if closeErr := conn.Close(); closeErr != nil {
		s.logger.WarnContext(ctx, "Failed to close network connection", logger.ErrorField(closeErr))
	}

	return constants.ServerStatusOnline
}

func (s *serverService) convertToServerResponse(server *model.Server) *response.ServerResponse {
	resp := &response.ServerResponse{
		ID:              server.ID,
		Name:            server.Name,
		Code:            server.Code,
		Hostname:        server.Hostname,
		Host:            server.Host,
		InternalIP:      server.InternalIP,
		IPv6Address:     server.IPv6Address,
		SSHPort:         server.SSHPort,
		OSDistro:        server.OSDistro,
		OSVersion:       server.OSVersion,
		KernelVersion:   server.KernelVersion,
		CPUCores:        server.CPUCores,
		MemoryTotal:     server.MemoryTotal,
		DiskTotal:       server.DiskTotal,
		Architecture:    server.Architecture,
		ResourceGroupID: server.ResourceGroupID,
		OwnerID:         server.OwnerID,
		Description:     server.Description,
		CreatedAt:       server.CreatedAt,
		UpdatedAt:       server.UpdatedAt,
	}

	return resp
}

// sshCredentials represents SSH authentication credentials
type sshCredentials struct {
	Username   string `json:"username"`
	Password   string `json:"password,omitempty"`
	PrivateKey string `json:"private_key,omitempty"`
}

// getSSHCredentials retrieves SSH credentials from secretKeyService
func (s *serverService) getSSHCredentials(ctx context.Context, server *model.Server, userID uint) (*sshCredentials, error) {
	// Get SSH credential from secret references
	refs, err := s.secretKeyService.GetSecretReferencesByResourceCode(ctx, server.Code)
	if err != nil || len(refs) == 0 {
		return nil, errors.NewAppError(errors.CodeServerCredentialNotFound)
	}

	credentialID := refs[0].SecretID

	// Get secret key value from secret key service
	secretValue, err := s.secretKeyService.GetSecretKeyValue(ctx, credentialID, userID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get SSH credentials from secret service",
			logger.Uint("credentialId", credentialID),
			logger.Uint("userId", userID),
			logger.ErrorField(err))
		return nil, err
	}

	// Validate secret key type - should be ACCOUNT type for SSH credentials
	if secretValue.KeyType != secretKeyTypeAccount {
		s.logger.ErrorContext(ctx, "Invalid secret key type for SSH credentials",
			logger.Uint("credentialId", credentialID),
			logger.String("keyType", string(secretValue.KeyType)))
		return nil, errors.NewAppError(errors.CodeValidationFailed)
	}

	// Extract credentials from CustomFields
	var creds sshCredentials

	// Get username
	if username, ok := secretValue.CustomFields["username"].(string); ok {
		creds.Username = username
	} else {
		s.logger.ErrorContext(ctx, "SSH credentials missing username field",
			logger.Uint("credentialId", credentialID))
		return nil, errors.NewAppError(errors.CodeValidationFailed)
	}

	// Get password (optional)
	if password, ok := secretValue.CustomFields["password"].(string); ok {
		creds.Password = password
	}

	// Get private_key (optional)
	if privateKey, ok := secretValue.CustomFields["private_key"].(string); ok {
		creds.PrivateKey = privateKey
	}

	// Validate credentials - must have either password or private key
	if creds.Password == "" && creds.PrivateKey == "" {
		s.logger.ErrorContext(ctx, "SSH credentials missing both password and private_key",
			logger.Uint("credentialId", credentialID))
		return nil, errors.NewAppError(errors.CodeValidationFailed)
	}

	s.logger.InfoContext(ctx, "SSH credentials retrieved successfully",
		logger.Uint("serverId", server.ID),
		logger.String("username", creds.Username),
		logger.Bool("hasPassword", creds.Password != ""),
		logger.Bool("hasPrivateKey", creds.PrivateKey != ""))

	return &creds, nil
}

// createSSHClientWithCreds creates an SSH client with provided credentials
func (s *serverService) createSSHClientWithCreds(ctx context.Context, server *model.Server, username, password, privateKey string) (*sshpkg.Client, error) {
	sshConfig := &sshpkg.ClientConfig{
		Host:       server.Host,
		Port:       server.SSHPort,
		Username:   username,
		Password:   password,
		PrivateKey: privateKey,
		Timeout:    time.Duration(s.getSSHTimeout()) * time.Second,
	}

	client, err := sshpkg.NewClient(ctx, sshConfig)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to create SSH client",
			logger.String("host", server.Host),
			logger.Int("port", server.SSHPort),
			logger.ErrorField(err))
		return nil, errors.NewAppErrorWrapError(err, errors.CodeServerSSHConnectionFailed)
	}

	return client, nil
}

// enrichServerWithStatusData enriches server response with status data from Redis
func (s *serverService) enrichServerWithStatusData(ctx context.Context, serverResp *response.ServerResponse, serverID uint) {
	// Mock implementation: 在真实环境中会从Redis获取状态数据
	// Redis Key Pattern: server:status:{server_id}

	now := time.Now()

	// 模拟从Redis获取状态数据
	s.logger.InfoContext(ctx, "Enriching server with status data from Redis (mock)",
		logger.Uint("serverId", serverID))

	// Mock status data (在真实实现中会从Redis获取)
	sshStatus := constants.SSHStatusConnected
	agentStatus := constants.AgentStatusOnline
	dockerStatus := constants.DockerStatusRunning

	serverResp.SSHStatus = &sshStatus
	serverResp.AgentStatus = &agentStatus
	serverResp.DockerStatus = &dockerStatus
	serverResp.LastCheckedAt = &now
}

// getSSHTimeout returns SSH timeout from database configuration with fallback to default
func (s *serverService) getSSHTimeout() int {
	// 尝试从数据库获取配置
	if s.systemConfigRepo != nil {
		if config, err := s.systemConfigRepo.GetByKey(context.Background(), "server.ssh_timeout"); err == nil && config != nil {
			value := config.GetEffectiveValue()
			if timeout, parseErr := strconv.Atoi(value); parseErr == nil && timeout > 0 {
				return timeout
			}
		}
	}

	// 回退到 config.yml 配置
	if s.config != nil && s.config.Server.Config.SSH.Timeout > 0 {
		return s.config.Server.Config.SSH.Timeout
	}

	// 最后使用默认值
	return constants.DefaultSSHTimeoutValue
}

// getFileUploadMaxSize returns file upload max size from database configuration with fallback to default
func (s *serverService) getFileUploadMaxSize() int64 {
	// 尝试从数据库获取配置
	if s.systemConfigRepo != nil {
		if config, err := s.systemConfigRepo.GetByKey(context.Background(), "server.file_upload_max_size"); err == nil && config != nil {
			value := config.GetEffectiveValue()
			if maxSize, parseErr := strconv.ParseInt(value, 10, 64); parseErr == nil && maxSize > 0 {
				return maxSize
			}
		}
	}

	// 回退到 config.yml 配置
	if s.config != nil && s.config.Server.Config.FileManagement.UploadMaxSize > 0 {
		return s.config.Server.Config.FileManagement.UploadMaxSize
	}

	// 最后使用默认值
	return constants.DefaultFileUploadMaxSize
}

// getFileDownloadMaxSize returns file download max size from database configuration with fallback to default
func (s *serverService) getFileDownloadMaxSize() int64 {
	// 尝试从数据库获取配置
	if s.systemConfigRepo != nil {
		if config, err := s.systemConfigRepo.GetByKey(context.Background(), "server.file_download_max_size"); err == nil && config != nil {
			value := config.GetEffectiveValue()
			if maxSize, parseErr := strconv.ParseInt(value, 10, 64); parseErr == nil && maxSize > 0 {
				return maxSize
			}
		}
	}

	// 回退到 config.yml 配置
	if s.config != nil && s.config.Server.Config.FileManagement.DownloadMaxSize > 0 {
		return s.config.Server.Config.FileManagement.DownloadMaxSize
	}

	// 最后使用默认值
	return constants.DefaultFileDownloadMaxSize
}

// getBatchConcurrency returns batch operation concurrency from configuration with fallback to default
func (s *serverService) getBatchConcurrency() int {
	// 使用常量作为默认值，将来可以从配置中读取
	return constants.DefaultBatchConcurrency
}

// executeServerActionsAsync 异步执行服务器批量操作
func (s *serverService) executeServerActionsAsync(ctx context.Context, req *request.ServerActionRequest, taskID string) {
	s.logger.InfoContext(ctx, "Starting async server actions execution",
		logger.String("taskID", taskID),
		logger.String("action", req.Action),
		logger.Int("serverCount", len(req.ServerIDs)))

	successCount := 0
	failedCount := 0

	// 并发执行多服务器操作，控制并发数避免资源过载
	maxConcurrency := s.getBatchConcurrency()
	semaphore := make(chan struct{}, maxConcurrency)

	var wg sync.WaitGroup
	for _, serverID := range req.ServerIDs {
		wg.Add(1)
		go func(id uint) {
			defer wg.Done()
			semaphore <- struct{}{}        // 获取信号量
			defer func() { <-semaphore }() // 释放信号量

			if err := s.executeServerAction(ctx, id, req.Action); err != nil {
				failedCount++
				s.logger.ErrorContext(ctx, "Server action failed",
					logger.Uint("serverID", id),
					logger.String("action", req.Action),
					logger.ErrorField(err))
			} else {
				successCount++
				s.logger.InfoContext(ctx, "Server action succeeded",
					logger.Uint("serverID", id),
					logger.String("action", req.Action))
			}
		}(serverID)
	}

	// 等待所有操作完成
	wg.Wait()

	s.logger.InfoContext(ctx, "Async server actions completed",
		logger.String("taskID", taskID),
		logger.String("action", req.Action),
		logger.Int("successCount", successCount),
		logger.Int("failedCount", failedCount))
}

// executeServerAction 执行单个服务器操作
func (s *serverService) executeServerAction(ctx context.Context, serverID uint, action string) error {
	// 验证服务器是否存在
	server, err := s.serverRepo.GetServerByID(ctx, serverID)
	if err != nil {
		return fmt.Errorf("server not found: %w", err)
	}

	s.logger.InfoContext(ctx, "Executing server action",
		logger.Uint("serverID", serverID),
		logger.String("serverName", server.Name),
		logger.String("action", action))

	switch action {
	case constants.ServerActionRestart:
		return s.restartServer(ctx, server)
	case constants.ServerActionShutdown:
		return s.shutdownServer(ctx, server)
	case constants.ServerActionReboot:
		return s.rebootServer(ctx, server)
	case constants.ServerActionUpdate:
		return s.updateServerSystem(ctx, server)
	case constants.ServerActionMaintenance:
		return s.setMaintenanceMode(ctx, server, true)
	case constants.ServerActionOnline:
		return s.setMaintenanceMode(ctx, server, false)
	case constants.ServerActionHealthCheck:
		return s.performHealthCheck(ctx, server)
	default:
		return fmt.Errorf("unsupported action: %s", action)
	}
}

// restartServer 重启服务器（模拟实现）
func (s *serverService) restartServer(ctx context.Context, server *model.Server) error {
	s.logger.InfoContext(ctx, "Restarting server",
		logger.Uint("serverID", server.ID),
		logger.String("serverName", server.Name))

	// 模拟重启操作：检查SSH连接 -> 发送重启命令 -> 等待重启完成
	time.Sleep(serverRestartTime) // 模拟重启时间

	// 更新服务器状态（如果需要）
	// 这里可以调用 s.serverRepo.UpdateServerStatus(ctx, server.ID, "restarting")

	return nil
}

// shutdownServer 关闭服务器（模拟实现）
func (s *serverService) shutdownServer(ctx context.Context, server *model.Server) error {
	s.logger.InfoContext(ctx, "Shutting down server",
		logger.Uint("serverID", server.ID),
		logger.String("serverName", server.Name))

	// 模拟关闭操作
	time.Sleep(1 * time.Second) // 模拟关闭时间

	return nil
}

// rebootServer 重新启动服务器（模拟实现）
func (s *serverService) rebootServer(ctx context.Context, server *model.Server) error {
	s.logger.InfoContext(ctx, "Rebooting server",
		logger.Uint("serverID", server.ID),
		logger.String("serverName", server.Name))

	// 模拟重启操作，比restart时间更长
	time.Sleep(serverShutdownTime)

	return nil
}

// updateServerSystem 更新服务器系统（模拟实现）
func (s *serverService) updateServerSystem(ctx context.Context, server *model.Server) error {
	s.logger.InfoContext(ctx, "Updating server system",
		logger.Uint("serverID", server.ID),
		logger.String("serverName", server.Name))

	// 模拟系统更新操作，时间较长
	time.Sleep(systemUpdateTime)

	return nil
}

// setMaintenanceMode 设置服务器维护模式（模拟实现）
func (s *serverService) setMaintenanceMode(ctx context.Context, server *model.Server, maintenance bool) error {
	mode := "online"
	if maintenance {
		mode = "maintenance"
	}

	s.logger.InfoContext(ctx, "Setting server maintenance mode",
		logger.Uint("serverID", server.ID),
		logger.String("serverName", server.Name),
		logger.String("mode", mode))

	// 模拟设置维护模式
	time.Sleep(maintenanceTime)

	return nil
}

// performHealthCheck 执行服务器健康检查
func (s *serverService) performHealthCheck(ctx context.Context, server *model.Server) error {
	s.logger.InfoContext(ctx, "Performing health check",
		logger.Uint("serverID", server.ID),
		logger.String("serverName", server.Name))

	// 复用现有的SSH连通性检查
	result := s.checkSSHConnectivity(ctx, server, 0)

	if result.Status != constants.SSHStatusConnected {
		return fmt.Errorf("health check failed: %s", result.ErrorMessage)
	}

	return nil
}

// ==========================================
// File Security Functions
// ==========================================

// validateFileSecurity performs comprehensive security validation for file operations
func (s *serverService) validateFileSecurity(ctx context.Context, filePath string, fileSize int64, isUpload bool) error {
	// 1. 文件大小检查
	if isUpload {
		maxSize := s.getFileUploadMaxSize()
		if fileSize > maxSize {
			s.logger.WarnContext(ctx, "File size exceeds upload limit",
				logger.String("filePath", filePath),
				logger.Int64("fileSize", fileSize),
				logger.Int64("maxSize", maxSize))
			return errors.NewAppError(errors.CodeFileSizeTooLarge)
		}
		// 硬限制检查
		if fileSize > constants.MaxFileSizeLimit {
			return errors.NewAppError(errors.CodeFileSizeTooLarge)
		}
	} else {
		maxSize := s.getFileDownloadMaxSize()
		if fileSize > maxSize {
			s.logger.WarnContext(ctx, "File size exceeds download limit",
				logger.String("filePath", filePath),
				logger.Int64("fileSize", fileSize),
				logger.Int64("maxSize", maxSize))
			return errors.NewAppError(errors.CodeFileSizeTooLarge)
		}
	}

	// 2. 路径安全检查
	if err := s.validateFilePath(ctx, filePath); err != nil {
		return err
	}

	// 3. 扩展名安全检查
	if err := s.validateFileExtension(ctx, filePath); err != nil {
		return err
	}

	return nil
}

// validateFilePath validates file path against security rules
func (s *serverService) validateFilePath(ctx context.Context, filePath string) error {
	// 1. 全局黑名单检查（最高优先级，不可绕过）
	if s.matchGlobalForbiddenPaths(filePath) {
		s.logger.WarnContext(ctx, "Path matched global forbidden list",
			logger.String("filePath", filePath))
		return errors.NewAppError(errors.CodePathForbidden)
	}

	// 2. 获取安全模式
	mode := s.getFilePathMode()

	// 3. 根据模式校验
	switch mode {
	case constants.FilePathModeBlacklist:
		// 黑名单模式：检查自定义禁止列表
		if s.matchCustomForbiddenPaths(filePath) {
			s.logger.WarnContext(ctx, "Path matched custom forbidden list",
				logger.String("filePath", filePath),
				logger.String("mode", mode))
			return errors.NewAppError(errors.CodePathForbidden)
		}
		return nil // 不在黑名单则允许

	case constants.FilePathModeWhitelist:
		// 白名单模式：必须在允许列表中
		if s.matchCustomAllowedPaths(filePath) {
			return nil
		}
		s.logger.WarnContext(ctx, "Path not in allowed list",
			logger.String("filePath", filePath),
			logger.String("mode", mode))
		return errors.NewAppError(errors.CodePathForbidden)

	default:
		// 未知模式，默认使用黑名单模式
		s.logger.WarnContext(ctx, "Unknown file path mode, using blacklist as default",
			logger.String("mode", mode))
		return s.validateFilePath(ctx, filePath)
	}
}

// validateFileExtension validates file extension against security rules
func (s *serverService) validateFileExtension(ctx context.Context, fileName string) error {
	// 1. 全局黑名单检查（不可绕过）
	if s.matchGlobalForbiddenExtensions(fileName) {
		s.logger.WarnContext(ctx, "File extension matched global forbidden list",
			logger.String("fileName", fileName))
		return errors.NewAppError(errors.CodeFileExtensionForbidden)
	}

	// 2. 获取安全模式
	mode := s.getFileExtMode()

	// 3. 根据模式校验
	switch mode {
	case constants.FileExtModeBlacklist:
		// 黑名单模式
		if s.matchCustomForbiddenExtensions(fileName) {
			s.logger.WarnContext(ctx, "File extension matched custom forbidden list",
				logger.String("fileName", fileName),
				logger.String("mode", mode))
			return errors.NewAppError(errors.CodeFileExtensionForbidden)
		}
		return nil

	case constants.FileExtModeWhitelist:
		// 白名单模式
		if s.matchCustomAllowedExtensions(fileName) {
			return nil
		}
		s.logger.WarnContext(ctx, "File extension not in allowed list",
			logger.String("fileName", fileName),
			logger.String("mode", mode))
		return errors.NewAppError(errors.CodeFileExtensionForbidden)

	default:
		// 未知模式，默认使用黑名单模式
		s.logger.WarnContext(ctx, "Unknown file extension mode, using blacklist as default",
			logger.String("mode", mode))
		return s.validateFileExtension(ctx, fileName)
	}
}

// getFilePathMode returns file path security mode from configuration with three-tier fallback
func (s *serverService) getFilePathMode() string {
	// 第一层：数据库配置
	if s.systemConfigRepo != nil {
		if config, err := s.systemConfigRepo.GetByKey(context.Background(), "server.file_path_mode"); err == nil && config != nil {
			value := config.GetEffectiveValue()
			if value == constants.FilePathModeBlacklist || value == constants.FilePathModeWhitelist {
				return value
			}
		}
	}

	// 第二层：config.yml（将来实现）
	// if s.config != nil && s.config.Server.FileManagement.PathMode != "" {
	//     return s.config.Server.FileManagement.PathMode
	// }

	// 第三层：常量默认值
	return constants.FilePathModeBlacklist
}

// getFileExtMode returns file extension security mode from configuration with three-tier fallback
func (s *serverService) getFileExtMode() string {
	// 第一层：数据库配置
	if s.systemConfigRepo != nil {
		if config, err := s.systemConfigRepo.GetByKey(context.Background(), "server.file_ext_mode"); err == nil && config != nil {
			value := config.GetEffectiveValue()
			if value == constants.FileExtModeBlacklist || value == constants.FileExtModeWhitelist {
				return value
			}
		}
	}

	// 第二层：config.yml（将来实现）
	// if s.config != nil && s.config.Server.FileManagement.ExtMode != "" {
	//     return s.config.Server.FileManagement.ExtMode
	// }

	// 第三层：常量默认值
	return constants.FileExtModeBlacklist
}

// getCustomForbiddenPaths returns custom forbidden paths from configuration
func (s *serverService) getCustomForbiddenPaths() []string {
	// 第一层：数据库配置（JSON）
	if s.systemConfigRepo != nil {
		if config, err := s.systemConfigRepo.GetByKey(context.Background(), "server.custom_forbidden_paths"); err == nil && config != nil {
			var paths []string
			value := config.GetEffectiveValue()
			// 尝试解析 JSON
			//nolint:staticcheck // SA9003: TODO: implement JSON parsing when needed
			if value != "" && value != "[]" {
				// 这里先返回空，实际使用时需要 json.Unmarshal
				// Example: json.Unmarshal([]byte(value), &paths)
			}
			return paths
		}
	}

	// 第二层：config.yml（将来实现）

	// 第三层：常量默认值
	return constants.DefaultCustomForbiddenPaths
}

// getCustomAllowedPaths returns custom allowed paths from configuration
func (s *serverService) getCustomAllowedPaths() []string {
	// 第一层：数据库配置（JSON）
	if s.systemConfigRepo != nil {
		if config, err := s.systemConfigRepo.GetByKey(context.Background(), "server.custom_allowed_paths"); err == nil && config != nil {
			var paths []string
			// JSON 解析逻辑
			return paths
		}
	}

	// 第二层：config.yml（将来实现）

	// 第三层：常量默认值
	return constants.DefaultAllowedPaths
}

// getCustomForbiddenExtensions returns custom forbidden extensions from configuration
func (s *serverService) getCustomForbiddenExtensions() []string {
	// 第一层：数据库配置
	if s.systemConfigRepo != nil {
		if config, err := s.systemConfigRepo.GetByKey(context.Background(), "server.custom_forbidden_extensions"); err == nil && config != nil {
			var exts []string
			// JSON 解析逻辑
			return exts
		}
	}

	// 第二层：config.yml（将来实现）

	// 第三层：常量默认值
	return constants.DefaultCustomForbiddenExtensions
}

// getCustomAllowedExtensions returns custom allowed extensions from configuration
func (s *serverService) getCustomAllowedExtensions() []string {
	// 第一层：数据库配置
	if s.systemConfigRepo != nil {
		if config, err := s.systemConfigRepo.GetByKey(context.Background(), "server.custom_allowed_extensions"); err == nil && config != nil {
			var exts []string
			// JSON 解析逻辑
			return exts
		}
	}

	// 第二层：config.yml（将来实现）

	// 第三层：常量默认值
	return constants.DefaultAllowedExtensions
}

// ==========================================
// Path Matching Functions
// ==========================================

// matchGlobalForbiddenPaths checks if path matches global forbidden patterns
func (s *serverService) matchGlobalForbiddenPaths(path string) bool {
	for _, pattern := range constants.GlobalForbiddenPaths {
		if s.matchPath(pattern, path) {
			return true
		}
	}
	return false
}

// matchCustomForbiddenPaths checks if path matches custom forbidden patterns
func (s *serverService) matchCustomForbiddenPaths(path string) bool {
	customPaths := s.getCustomForbiddenPaths()
	for _, pattern := range customPaths {
		if s.matchPath(pattern, path) {
			return true
		}
	}
	return false
}

// matchCustomAllowedPaths checks if path matches custom allowed patterns
func (s *serverService) matchCustomAllowedPaths(path string) bool {
	customPaths := s.getCustomAllowedPaths()
	for _, pattern := range customPaths {
		if s.matchPath(pattern, path) {
			return true
		}
	}
	return false
}

// matchGlobalForbiddenExtensions checks if file extension matches global forbidden list
func (s *serverService) matchGlobalForbiddenExtensions(fileName string) bool {
	ext := s.getFileExtension(fileName)
	for _, forbiddenExt := range constants.GlobalForbiddenExtensions {
		if ext == forbiddenExt {
			return true
		}
	}
	return false
}

// matchCustomForbiddenExtensions checks if file extension matches custom forbidden list
func (s *serverService) matchCustomForbiddenExtensions(fileName string) bool {
	ext := s.getFileExtension(fileName)
	customExts := s.getCustomForbiddenExtensions()
	for _, forbiddenExt := range customExts {
		if ext == forbiddenExt {
			return true
		}
	}
	return false
}

// matchCustomAllowedExtensions checks if file extension matches custom allowed list
func (s *serverService) matchCustomAllowedExtensions(fileName string) bool {
	ext := s.getFileExtension(fileName)
	customExts := s.getCustomAllowedExtensions()
	for _, allowedExt := range customExts {
		if ext == allowedExt {
			return true
		}
	}
	return false
}

// matchPath matches file path against pattern with wildcard support
// Supports * (single level) and ** (multiple levels) wildcards
func (s *serverService) matchPath(pattern, path string) bool {
	// 简单的通配符匹配实现
	// * 匹配任意字符（不包括路径分隔符）
	// 未来可以扩展支持更复杂的模式

	// 精确匹配
	if pattern == path {
		return true
	}

	// 处理 ~ 开头的路径
	if pattern != "" && pattern[0] == '~' {
		// 简化处理：匹配任何以 .ssh 结尾的路径
		if contains(path, "/.ssh/") || endsWith(path, "/.ssh") {
			return true
		}
	}

	// 处理通配符 *
	if contains(pattern, "*") {
		return s.wildcardMatch(pattern, path)
	}

	return false
}

// wildcardMatch performs wildcard pattern matching
func (s *serverService) wildcardMatch(pattern, str string) bool {
	// 将模式分割为段
	parts := splitBy(pattern, "*")

	// 如果只有一个段，直接比较
	if len(parts) == 1 {
		return parts[0] == str
	}

	// 检查开头
	if parts[0] != "" && !startsWith(str, parts[0]) {
		return false
	}

	// 检查结尾
	if parts[len(parts)-1] != "" && !endsWith(str, parts[len(parts)-1]) {
		return false
	}

	// 检查中间部分
	currentPos := len(parts[0])
	for i := 1; i < len(parts)-1; i++ {
		part := parts[i]
		if part == "" {
			continue
		}
		idx := indexFrom(str, part, currentPos)
		if idx == -1 {
			return false
		}
		currentPos = idx + len(part)
	}

	return true
}

// getFileExtension extracts file extension from file name or path
func (s *serverService) getFileExtension(fileName string) string {
	// 从最后一个 . 开始提取扩展名
	lastDot := -1
	for i := len(fileName) - 1; i >= 0; i-- {
		if fileName[i] == '.' {
			lastDot = i
			break
		}
		// 如果遇到路径分隔符，说明没有扩展名
		if fileName[i] == '/' || fileName[i] == '\\' {
			break
		}
	}

	if lastDot == -1 || lastDot == len(fileName)-1 {
		return ""
	}

	return fileName[lastDot:]
}

// ==========================================
// String Helper Functions
// ==========================================

func contains(s, substr string) bool {
	return indexOf(s, substr) != -1
}

func indexOf(s, substr string) int {
	return indexFrom(s, substr, 0)
}

func indexFrom(s, substr string, start int) int {
	if start < 0 {
		start = 0
	}
	if substr == "" {
		return start
	}
	if len(s) < len(substr) {
		return -1
	}

	for i := start; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if s[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}

func startsWith(s, prefix string) bool {
	if len(s) < len(prefix) {
		return false
	}
	for i := 0; i < len(prefix); i++ {
		if s[i] != prefix[i] {
			return false
		}
	}
	return true
}

func endsWith(s, suffix string) bool {
	if len(s) < len(suffix) {
		return false
	}
	start := len(s) - len(suffix)
	for i := 0; i < len(suffix); i++ {
		if s[start+i] != suffix[i] {
			return false
		}
	}
	return true
}

func splitBy(s, sep string) []string {
	if sep == "" {
		return []string{s}
	}

	var result []string
	start := 0
	for i := 0; i <= len(s)-len(sep); i++ {
		match := true
		for j := 0; j < len(sep); j++ {
			if s[i+j] != sep[j] {
				match = false
				break
			}
		}
		if match {
			result = append(result, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	result = append(result, s[start:])
	return result
}

// generateServerCode generates a unique server code
// Format: srv + timestamp (srv1234567890123)
func (s *serverService) generateServerCode() string {
	// Use timestamp to ensure uniqueness
	timestamp := time.Now().UnixMilli()
	return fmt.Sprintf("%s%d", serverCodePrefix, timestamp)
}

// createSecretReference creates a secret reference for SSH credential
func (s *serverService) createSecretReference(ctx context.Context, serverCode, credentialID string) error {
	// Parse credential ID from string to uint
	credID, err := strconv.ParseUint(credentialID, 10, 64)
	if err != nil {
		s.logger.ErrorContext(ctx, "Invalid SSH credential ID format",
			logger.String("credentialId", credentialID),
			logger.ErrorField(err))
		return errors.NewAppErrorWrapError(err, errors.CodeInvalidParameterFormat)
	}

	// Create secret reference record
	reference := &model.SecretReference{
		SecretID:     uint(credID),
		ResourceCode: serverCode,
	}

	if err := s.secretKeyService.CreateSecretReference(ctx, reference); err != nil {
		s.logger.ErrorContext(ctx, "Failed to create secret reference",
			logger.String("serverCode", serverCode),
			logger.Uint("secretId", uint(credID)),
			logger.ErrorField(err))
		return err
	}

	s.logger.InfoContext(ctx, "Secret reference created successfully",
		logger.String("serverCode", serverCode),
		logger.Uint("secretId", uint(credID)))

	return nil
}

// createSSHCredentialSecret creates a secret for SSH credentials and returns the credential ID
func (s *serverService) createSSHCredentialSecret(ctx context.Context, server *model.Server, req *request.CreateServerRequest, userID uint) (string, error) {
	// Build credential name
	secretName := fmt.Sprintf("SSH-%s-%s", server.Name, server.Code)

	// Build custom fields for ACCOUNT type
	customFields := map[string]interface{}{
		"username": req.SSHUsername,
	}

	if req.SSHPassword != "" {
		customFields["password"] = req.SSHPassword
	}
	if req.SSHKey != "" {
		customFields["private_key"] = req.SSHKey
	}

	// Create secret key request
	secretReq := &request.SecretKeyCreateRequest{
		Name:         secretName,
		KeyType:      model.SecretKeyTypeAccount,
		CustomFields: customFields,
	}

	// Create secret via secret key service
	secretResp, err := s.secretKeyService.CreateSecretKey(ctx, secretReq, userID)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", secretResp.ID), nil
}

// deleteSecretReferencesByServerCode deletes all secret references for a server
func (s *serverService) deleteSecretReferencesByServerCode(ctx context.Context, serverCode string) {
	if err := s.secretKeyService.DeleteSecretReferencesByResourceCode(ctx, serverCode); err != nil {
		s.logger.WarnContext(ctx, "Failed to delete secret references for server",
			logger.String("serverCode", serverCode),
			logger.ErrorField(err))
		// Don't fail server deletion if reference cleanup fails
		return
	}

	s.logger.InfoContext(ctx, "Secret references deleted successfully",
		logger.String("serverCode", serverCode))
}
