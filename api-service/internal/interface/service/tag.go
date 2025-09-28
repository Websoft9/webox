package service

import (
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"context"
)

// TagService interface for tag business logic operations
type TagService interface {
	// Basic tag management
	CreateTag(ctx context.Context, req *request.TagCreateRequest, userID uint) (*response.TagResponse, error)
	GetTag(ctx context.Context, id uint) (*response.TagResponse, error)
	UpdateTag(ctx context.Context, id uint, req *request.TagUpdateRequest, userID uint) (*response.TagResponse, error)
	DeleteTag(ctx context.Context, id uint, userID uint) error
	ListTags(ctx context.Context, req *request.TagListRequest) ([]*response.TagResponse, error)
	SearchTags(ctx context.Context, req *request.TagNameSearchRequest) ([]*response.TagResponse, error)

	// Resource tagging operations
	AssignTags(ctx context.Context, req *request.TagAssignRequest, userID uint) (*response.TagAssignResponse, error)
	ReplaceTags(ctx context.Context, req *request.TagAssignRequest, userID uint) (*response.TagAssignResponse, error)
	UnassignTags(ctx context.Context, req *request.TagUnassignRequest, userID uint) (*response.TagUnassignResponse, error)
	GetResourceTags(ctx context.Context, req *request.TaggingListRequest) ([]*response.TagSimpleResponse, error)

	// Search operations
	SearchResourcesByTags(ctx context.Context, req *request.TagSearchRequest) (*response.TagSearchResponse, error)
}
