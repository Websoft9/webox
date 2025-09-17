package service

import (
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"context"
)

type SystemConfigService interface {
	// Retrieve all system configurations
	ListSystemConfigs(ctx context.Context, req *request.ListSystemConfigsRequest) (*response.ListSystemConfigsResponse, error)

	// Batch update system configurations
	BatchUpdateSystemConfigs(ctx context.Context, req *request.BatchUpdateSystemConfigsRequest) error

	// Test SMTP server configuration
	TestSMTP(ctx context.Context, req *request.TestSMTPRequest) error
}
