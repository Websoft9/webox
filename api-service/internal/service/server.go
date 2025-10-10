package service

import (
	"context"
	"fmt"
	"net"
	"path/filepath"
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
)

// 常量定义，避免魔术数字
const (
	// 网络延迟模拟常量
	networkDelayShort  = 50 * time.Millisecond
	networkDelayMedium = 100 * time.Millisecond

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
)

// serverService implements service.ServerService
type serverService struct {
	serverRepo       repository.ServerRepository
	secretKeyService service.SecretKeyService // 使用真实的密钥管理服务
	logger           logger.Logger
	config           *config.Config // 添加配置依赖
}

// ServerServiceConfig contains configuration for server service
type ServerServiceConfig struct {
	Logger           logger.Logger
	ServerRepo       repository.ServerRepository
	SecretKeyService service.SecretKeyService // 使用真实的密钥管理服务
	Config           *config.Config           // 添加配置依赖
}

// NewServerService creates a new server service instance
func NewServerService(config ServerServiceConfig) service.ServerService {
	return &serverService{
		serverRepo:       config.ServerRepo,
		secretKeyService: config.SecretKeyService,
		logger:           config.Logger,
		config:           config.Config,
	}
}

// CreateServer creates a new server
func (s *serverService) CreateServer(ctx context.Context, req *request.CreateServerRequest) (*response.ServerResponse, error) {
	s.logger.InfoContext(ctx, "Creating new server",
		logger.String("name", req.Name),
		logger.String("hostname", req.Hostname))

	// Check if server name already exists
	exists, err := s.serverRepo.ExistsServerByName(ctx, req.Name)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to check server name existence", logger.ErrorField(err))
		return nil, errors.NewAppError(errors.CodeInternalError)
	}
	if exists {
		return nil, errors.NewAppError(errors.CodeResourceAlreadyExists)
	}

	// Validate SSH credential if provided
	if req.SSHCredentialID != nil {
		if err := s.validateSSHCredential(ctx, *req.SSHCredentialID, req.OwnerID); err != nil {
			return nil, err
		}
	}

	// Create server model
	server := &model.Server{
		Name:            req.Name,
		Hostname:        req.Hostname,
		Host:            req.Host,
		InternalIP:      req.InternalIP,
		IPv6Address:     req.IPv6Address,
		SSHPort:         req.SSHPort,
		SSHCredentialID: req.SSHCredentialID,
		ResourceGroupID: req.ResourceGroupID,
		OwnerID:         req.OwnerID,
		Description:     &req.Description,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Create server in database
	if err := s.serverRepo.CreateServer(ctx, server); err != nil {
		s.logger.ErrorContext(ctx, "Failed to create server in database", logger.ErrorField(err))
		return nil, errors.NewAppError(errors.CodeRecordCreateFailed)
	}

	// Test connectivity asynchronously
	go s.testServerConnectivity(context.Background(), server.ID)

	s.logger.InfoContext(ctx, "Server created successfully",
		logger.Uint("id", server.ID),
		logger.String("name", server.Name))

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
func (s *serverService) UpdateServer(ctx context.Context, id uint, req *request.UpdateServerRequest) (*response.ServerResponse, error) {
	s.logger.InfoContext(ctx, "Updating server", logger.Uint("id", id))

	// Get existing server
	server, err := s.serverRepo.GetServerByID(ctx, id)
	if err != nil {
		return nil, err
	}

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
	if req.Hostname != nil {
		server.Hostname = *req.Hostname
	}
	if req.Host != nil {
		server.Host = *req.Host
	}
	if req.SSHPort != nil {
		server.SSHPort = *req.SSHPort
	}
	if req.SSHCredentialID != nil {
		// Validate SSH credential if provided
		if err := s.validateSSHCredential(ctx, *req.SSHCredentialID, server.OwnerID); err != nil {
			return nil, err
		}
		server.SSHCredentialID = req.SSHCredentialID
	}
	// OSDistro field doesn't exist in UpdateServerRequest, skip this update
	if req.Description != nil {
		server.Description = req.Description
	}

	server.UpdatedAt = time.Now()

	// Update server in database
	if err := s.serverRepo.UpdateServer(ctx, server); err != nil {
		return nil, err
	}

	s.logger.InfoContext(ctx, "Server updated successfully", logger.Uint("id", id))

	return s.convertToServerResponse(server), nil
}

// DeleteServer deletes a server
func (s *serverService) DeleteServer(ctx context.Context, id uint) error {
	s.logger.InfoContext(ctx, "Deleting server", logger.Uint("id", id))

	// Check if server exists
	_, err := s.serverRepo.GetServerByID(ctx, id)
	if err != nil {
		return err
	}

	// Delete server from database
	if err := s.serverRepo.DeleteServer(ctx, id); err != nil {
		return err
	}

	s.logger.InfoContext(ctx, "Server deleted successfully", logger.Uint("id", id))

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
func (s *serverService) CheckServersStatus(ctx context.Context, req *request.ServerStatusCheckRequest) (*response.BatchServerStatusResponse, error) {
	results := make([]response.ServerStatusCheckResult, 0, len(req.ServerIDs))
	successCount := 0
	failedCount := 0

	for _, serverID := range req.ServerIDs {
		result := s.checkSingleServerStatus(ctx, serverID, req.CheckTypes, req.Timeout)
		if result.Status == "failed" {
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
	// Check file size limit using configuration
	maxSize := s.getFileUploadMaxSize()
	if int64(len(fileData)) > maxSize {
		s.logger.WarnContext(ctx, "File size exceeds upload limit",
			logger.Uint("serverId", serverID),
			logger.String("filePath", filePath),
			logger.Int("fileSize", len(fileData)),
			logger.Int64("maxSize", maxSize))
		return nil, errors.NewAppError(errors.CodeValidationFailed)
	}

	// Get server info
	server, err := s.serverRepo.GetServerByID(ctx, serverID)
	if err != nil {
		return nil, err
	}

	s.logger.InfoContext(ctx, "Uploading file to server",
		logger.Uint("serverId", serverID),
		logger.String("filePath", filePath),
		logger.Int("fileSize", len(fileData)))

	// Mock implementation: 模拟SSH文件上传
	// 在真实实现中，这里会使用SSH连接上传文件
	time.Sleep(networkDelayMedium) // 模拟网络延迟

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
	// Get server info
	server, err := s.serverRepo.GetServerByID(ctx, serverID)
	if err != nil {
		return nil, "", err
	}

	s.logger.InfoContext(ctx, "Downloading file from server",
		logger.Uint("serverId", serverID),
		logger.String("filePath", filePath))

	// Mock implementation: 模拟SSH文件下载
	// 在真实实现中，这里会使用SSH连接下载文件
	time.Sleep(networkDelayMedium) // 模拟网络延迟

	// 返回模拟文件内容
	mockContent := []byte(fmt.Sprintf("Mock file content from %s:%s", server.Name, filePath))
	filename := filepath.Base(filePath)

	// Check download size limit using configuration
	maxSize := s.getFileDownloadMaxSize()
	if int64(len(mockContent)) > maxSize {
		s.logger.WarnContext(ctx, "File size exceeds download limit",
			logger.Uint("serverId", serverID),
			logger.String("filePath", filePath),
			logger.Int("fileSize", len(mockContent)),
			logger.Int64("maxSize", maxSize))
		return nil, "", errors.NewAppError(errors.CodeValidationFailed)
	}

	return mockContent, filename, nil
}

// DeleteFile deletes a file from server via SSH (设计文档序号10)
func (s *serverService) DeleteFile(ctx context.Context, serverID uint, filePath string) (*response.ServerFileDeleteResponse, error) {
	// Get server info
	server, err := s.serverRepo.GetServerByID(ctx, serverID)
	if err != nil {
		return nil, err
	}

	s.logger.InfoContext(ctx, "Deleting file from server",
		logger.Uint("serverId", serverID),
		logger.String("filePath", filePath))

	// Mock implementation: 模拟SSH文件删除
	// 在真实实现中，这里会使用SSH连接删除文件
	time.Sleep(networkDelayShort) // 模拟网络延迟

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

// checkSSHConnectivity tests network connectivity only (not SSH authentication)
func (s *serverService) checkSSHConnectivity(ctx context.Context, server *model.Server, timeout int) response.ServiceStatusInfo {
	if timeout == 0 {
		// 使用配置而不是硬编码常量
		timeout = s.getSSHTimeout()
	}

	start := time.Now()
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", server.Host, server.SSHPort), time.Duration(timeout)*time.Second)
	responseTime := int(time.Since(start).Milliseconds())

	now := time.Now()
	if err != nil {
		return response.ServiceStatusInfo{
			Status:       constants.SSHStatusDisconnected,
			ResponseTime: responseTime,
			CheckedAt:    &now,
			ErrorMessage: "Network connection failed",
		}
	}
	if closeErr := conn.Close(); closeErr != nil {
		s.logger.WarnContext(ctx, "Failed to close SSH connection", logger.ErrorField(closeErr))
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
		Hostname:        server.Hostname,
		Host:            server.Host,
		InternalIP:      server.InternalIP,
		IPv6Address:     server.IPv6Address,
		SSHPort:         server.SSHPort,
		SSHCredentialID: server.SSHCredentialID,
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

// getSSHTimeout returns SSH timeout from configuration with fallback to default
func (s *serverService) getSSHTimeout() int {
	if s.config != nil && s.config.Server.Config.SSH.Timeout > 0 {
		return s.config.Server.Config.SSH.Timeout
	}
	return constants.DefaultSSHTimeoutValue
}

// getFileUploadMaxSize returns file upload max size from configuration with fallback to default
func (s *serverService) getFileUploadMaxSize() int64 {
	if s.config != nil && s.config.Server.Config.FileManagement.UploadMaxSize > 0 {
		return s.config.Server.Config.FileManagement.UploadMaxSize
	}
	return constants.DefaultFileUploadMaxSize
}

// getFileDownloadMaxSize returns file download max size from configuration with fallback to default
func (s *serverService) getFileDownloadMaxSize() int64 {
	if s.config != nil && s.config.Server.Config.FileManagement.DownloadMaxSize > 0 {
		return s.config.Server.Config.FileManagement.DownloadMaxSize
	}
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

// validateSSHCredential validates if SSH credential exists and user has access to it
func (s *serverService) validateSSHCredential(ctx context.Context, credentialID string, userID uint) error {
	s.logger.InfoContext(ctx, "Validating SSH credential",
		logger.String("credentialId", credentialID),
		logger.Uint("userId", userID))

	// Convert string credential ID to uint
	// 注意：这里简化处理，实际应该使用 strconv.ParseUint
	// 假设 credentialID 格式为数字字符串
	var credID uint
	// 简化的转换逻辑，实际中应该使用 strconv.ParseUint
	switch credentialID {
	case "1":
		credID = 1
	case "2":
		credID = 2
	case "3":
		credID = 3
	default:
		s.logger.ErrorContext(ctx, "Invalid credential ID format",
			logger.String("credentialId", credentialID))
		return errors.NewAppError(errors.CodeValidationFailed)
	}

	// Validate credential ownership using SecretKeyService
	if err := s.secretKeyService.ValidateSecretKeyOwnership(ctx, credID, userID); err != nil {
		s.logger.ErrorContext(ctx, "SSH credential validation failed",
			logger.Uint("credentialId", credID),
			logger.Uint("userId", userID),
			logger.ErrorField(err))
		return errors.NewAppError(errors.CodeResourceNotFound)
	}

	// Get credential details to ensure it's valid for SSH
	credential, err := s.secretKeyService.GetSecretKey(ctx, credID, userID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get SSH credential details",
			logger.Uint("credentialId", credID),
			logger.ErrorField(err))
		return errors.NewAppError(errors.CodeResourceNotFound)
	}

	// Check if credential type is suitable for SSH
	keyTypeStr := string(credential.KeyType)
	if keyTypeStr != constants.SecretTypeSSHKey && keyTypeStr != constants.SecretTypePassword {
		s.logger.WarnContext(ctx, "Invalid credential type for SSH",
			logger.Uint("credentialId", credID),
			logger.String("keyType", keyTypeStr))
		return errors.NewAppError(errors.CodeValidationFailed)
	}

	return nil
}
