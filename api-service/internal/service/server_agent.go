package service

import (
	"context"
	"fmt"
	"time"

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
	deploymentSimulationTime = 2 * time.Second
)

// serverAgentService implements service.ServerAgentService
type serverAgentService struct {
	agentRepo  repository.ServerAgentRepository
	serverRepo repository.ServerRepository
	logger     logger.Logger
}

// ServerAgentServiceConfig contains configuration for server agent service
type ServerAgentServiceConfig struct {
	AgentRepo  repository.ServerAgentRepository
	ServerRepo repository.ServerRepository
	Logger     logger.Logger
}

// NewServerAgentService creates a new server agent service instance
func NewServerAgentService(config ServerAgentServiceConfig) service.ServerAgentService {
	return &serverAgentService{
		agentRepo:  config.AgentRepo,
		serverRepo: config.ServerRepo,
		logger:     config.Logger,
	}
}

// DeployAgent deploys an agent to a server
func (s *serverAgentService) DeployAgent(ctx context.Context, req *request.DeployAgentRequest) (*response.ServerAgentResponse, error) {
	s.logger.InfoContext(ctx, "Deploying agent to server",
		logger.Uint("serverId", req.ServerID),
		logger.String("deploymentType", req.DeploymentType))

	// Validate server exists
	server, err := s.serverRepo.GetServerByID(ctx, req.ServerID)
	if err != nil {
		return nil, err
	}

	// Check if agent already exists for this server
	existingAgent, err := s.agentRepo.GetAgentByServerID(ctx, req.ServerID)
	if err == nil && existingAgent != nil {
		return nil, errors.NewAppError(errors.CodeResourceAlreadyExists)
	}

	// Validate deployment type
	deploymentType, err := s.validateDeploymentType(req.DeploymentType)
	if err != nil {
		return nil, err
	}

	// Create agent record
	agent := &model.ServerAgent{
		ServerID:       req.ServerID,
		DeploymentType: deploymentType,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Create agent in database
	if err := s.agentRepo.CreateAgent(ctx, agent); err != nil {
		s.logger.ErrorContext(ctx, "Failed to create agent record", logger.ErrorField(err))
		return nil, errors.NewAppError(errors.CodeRecordCreateFailed)
	}

	// Start deployment task asynchronously
	go s.executeAgentDeployment(context.Background(), agent, server)

	s.logger.InfoContext(ctx, "Agent deployment initiated",
		logger.Uint("agentId", agent.ID),
		logger.Uint("serverId", req.ServerID))

	return s.convertToAgentResponse(agent), nil
}

// GetAgent retrieves an agent by ID
func (s *serverAgentService) GetAgent(ctx context.Context, id uint) (*response.ServerAgentResponse, error) {
	agent, err := s.agentRepo.GetAgentByID(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get agent",
			logger.Uint("id", id),
			logger.ErrorField(err))
		return nil, err
	}

	return s.convertToAgentResponse(agent), nil
}

// UpdateAgent updates an agent's configuration
func (s *serverAgentService) UpdateAgent(ctx context.Context, id uint, req *request.UpdateAgentRequest) (*response.ServerAgentResponse, error) {
	s.logger.InfoContext(ctx, "Updating agent", logger.Uint("id", id))

	// Get existing agent
	agent, err := s.agentRepo.GetAgentByID(ctx, id)
	if err != nil {
		return nil, err
	}

	agent.UpdatedAt = time.Now()

	// Update agent in database
	if err := s.agentRepo.UpdateAgent(ctx, agent); err != nil {
		return nil, err
	}

	s.logger.InfoContext(ctx, "Agent updated successfully", logger.Uint("id", id))

	return s.convertToAgentResponse(agent), nil
}

// UninstallAgent uninstalls an agent from a server
func (s *serverAgentService) UninstallAgent(ctx context.Context, id uint) error {
	s.logger.InfoContext(ctx, "Uninstalling agent", logger.Uint("id", id))

	// Get agent
	agent, err := s.agentRepo.GetAgentByID(ctx, id)
	if err != nil {
		return err
	}

	// ServerAgent model uses LastHeartbeatAt for status tracking

	// Start uninstallation task asynchronously
	go s.executeAgentUninstallation(context.Background(), agent)

	s.logger.InfoContext(ctx, "Agent uninstallation initiated", logger.Uint("id", id))

	return nil
}

// RestartAgent restarts an agent
func (s *serverAgentService) RestartAgent(ctx context.Context, id uint) error {
	s.logger.InfoContext(ctx, "Restarting agent", logger.Uint("id", id))

	// Get agent
	agent, err := s.agentRepo.GetAgentByID(ctx, id)
	if err != nil {
		return err
	}

	// ServerAgent model uses LastHeartbeatAt for status tracking

	// Start restart task asynchronously
	go s.executeAgentRestart(context.Background(), agent)

	s.logger.InfoContext(ctx, "Agent restart initiated", logger.Uint("id", id))

	return nil
}

// GetAgentLogs retrieves agent logs
func (s *serverAgentService) GetAgentLogs(ctx context.Context, req *request.GetAgentLogsRequest) (*response.AgentLogsResponse, error) {
	// Get agent
	agent, err := s.agentRepo.GetAgentByID(ctx, req.AgentID)
	if err != nil {
		return nil, err
	}

	// Mock log retrieval (in real implementation, this would connect to the server and fetch logs)
	logs := s.mockGetAgentLogs(ctx, agent, req.Lines)

	return &response.AgentLogsResponse{
		AgentID:   req.AgentID,
		Logs:      logs,
		Timestamp: time.Now(),
	}, nil
}

// GetAgentStatus checks the status of an agent
func (s *serverAgentService) GetAgentStatus(ctx context.Context, id uint) (*response.AgentStatusResponse, error) {
	agent, err := s.agentRepo.GetAgentByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Mock status check (in real implementation, this would ping the agent)
	status := s.mockCheckAgentStatus(ctx, agent)

	return &response.AgentStatusResponse{
		AgentID:   id,
		Status:    status,
		CheckedAt: time.Now(),
		Message:   s.getAgentStatusMessage(status),
	}, nil
}

// Helper methods

func (s *serverAgentService) validateDeploymentType(deploymentType string) (model.AgentDeploymentType, error) {
	switch deploymentType {
	case "docker":
		return model.AgentDeploymentDocker, nil
	case "systemd":
		return model.AgentDeploymentSystemd, nil
	default:
		return "", errors.NewAppError(errors.CodeInvalidParameterFormat)
	}
}

func (s *serverAgentService) executeAgentDeployment(ctx context.Context, agent *model.ServerAgent, server *model.Server) {
	s.logger.InfoContext(ctx, "Executing agent deployment",
		logger.Uint("agentId", agent.ID),
		logger.Uint("serverId", server.ID))

	// Mock deployment process
	time.Sleep(deploymentSimulationTime) // Simulate deployment time

	// Update last heartbeat to indicate agent is online
	if err := s.agentRepo.UpdateAgentLastSeen(ctx, agent.ID); err != nil {
		s.logger.ErrorContext(ctx, "Failed to update agent last seen after deployment",
			logger.Uint("agentId", agent.ID),
			logger.ErrorField(err))
	} else {
		s.logger.InfoContext(ctx, "Agent deployment completed successfully", logger.Uint("agentId", agent.ID))
	}
}

func (s *serverAgentService) executeAgentUninstallation(ctx context.Context, agent *model.ServerAgent) {
	s.logger.InfoContext(ctx, "Executing agent uninstallation", logger.Uint("agentId", agent.ID))

	// Mock uninstallation process
	time.Sleep(1 * time.Second) // Simulate uninstallation time

	// Remove agent record from database
	if err := s.agentRepo.DeleteAgent(ctx, agent.ID); err != nil {
		s.logger.ErrorContext(ctx, "Failed to delete agent record after uninstallation",
			logger.Uint("agentId", agent.ID),
			logger.ErrorField(err))
	} else {
		s.logger.InfoContext(ctx, "Agent uninstallation completed successfully", logger.Uint("agentId", agent.ID))
	}
}

func (s *serverAgentService) executeAgentRestart(ctx context.Context, agent *model.ServerAgent) {
	s.logger.InfoContext(ctx, "Executing agent restart", logger.Uint("agentId", agent.ID))

	// Mock restart process
	time.Sleep(1 * time.Second) // Simulate restart time

	// Update last heartbeat to indicate agent is online again
	if err := s.agentRepo.UpdateAgentLastSeen(ctx, agent.ID); err != nil {
		s.logger.ErrorContext(ctx, "Failed to update agent last seen after restart",
			logger.Uint("agentId", agent.ID),
			logger.ErrorField(err))
	} else {
		s.logger.InfoContext(ctx, "Agent restart completed successfully", logger.Uint("agentId", agent.ID))
	}
}

func (s *serverAgentService) mockGetAgentLogs(_ context.Context, agent *model.ServerAgent, lines int) []string {
	// Mock log lines (in real implementation, this would fetch from server)
	mockLogs := []string{
		fmt.Sprintf("[%s] Agent started successfully", time.Now().Format("2006-01-02 15:04:05")),
		fmt.Sprintf("[%s] Connected to Websoft9 platform", time.Now().Add(-1*time.Minute).Format("2006-01-02 15:04:05")),
		fmt.Sprintf("[%s] Health check passed", time.Now().Add(-2*time.Minute).Format("2006-01-02 15:04:05")),
		fmt.Sprintf("[%s] Monitoring system metrics", time.Now().Add(-3*time.Minute).Format("2006-01-02 15:04:05")),
		fmt.Sprintf("[%s] Task execution completed", time.Now().Add(-5*time.Minute).Format("2006-01-02 15:04:05")),
	}

	if lines > 0 && lines < len(mockLogs) {
		return mockLogs[:lines]
	}

	return mockLogs
}

func (s *serverAgentService) mockCheckAgentStatus(_ context.Context, agent *model.ServerAgent) string {
	// Mock status check (in real implementation, this would ping the agent)
	// Use IsOnline method from the model
	const (
		agentStatusRunning = "running"
		agentStatusOffline = "offline"
	)

	if agent.IsOnline() {
		return agentStatusRunning
	}
	return agentStatusOffline
}

func (s *serverAgentService) getAgentStatusMessage(status string) string {
	const (
		agentStatusRunning = "running"
		agentStatusStopped = "stopped"
		agentStatusOffline = "offline"
	)

	switch status {
	case agentStatusRunning:
		return "Agent is running normally"
	case agentStatusStopped:
		return "Agent is stopped"
	case agentStatusOffline:
		return "Agent is offline or not responding"
	case "deploying":
		return "Agent deployment in progress"
	case "uninstalling":
		return "Agent uninstallation in progress"
	case "restarting":
		return "Agent restart in progress"
	case "error":
		return "Agent encountered an error"
	default:
		return "Unknown agent status"
	}
}

func (s *serverAgentService) convertToAgentResponse(agent *model.ServerAgent) *response.ServerAgentResponse {
	resp := &response.ServerAgentResponse{
		ID:             agent.ID,
		ServerID:       agent.ServerID,
		DeploymentType: string(agent.DeploymentType),
		CreatedAt:      agent.CreatedAt,
		UpdatedAt:      agent.UpdatedAt,
	}

	return resp
}
