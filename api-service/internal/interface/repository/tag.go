package repository

import (
	"api-service/internal/dto/request"
	"api-service/internal/model"
	"context"
)

// TagRepository interface for tag data access operations
type TagRepository interface {
	// Basic CRUD operations for tags
	CreateTag(ctx context.Context, tag *model.Tag) error
	GetTagByID(ctx context.Context, id uint) (*model.Tag, error)
	GetTagByName(ctx context.Context, name string) (*model.Tag, error)
	UpdateTag(ctx context.Context, tag *model.Tag) error
	DeleteTag(ctx context.Context, id uint) error

	// Tag listing and searching
	ListTags(ctx context.Context, search string, excludeIDs []uint) ([]*model.Tag, error)
	ListTagsWithUsageCount(ctx context.Context, search string, excludeIDs []uint) ([]*model.Tag, error)
	SearchTagsByName(ctx context.Context, query string) ([]*model.Tag, error)

	// Tag existence checks
	ExistsTagByName(ctx context.Context, name string) (bool, error)
	ExistsTagByNameExcludeID(ctx context.Context, name string, excludeID uint) (bool, error)

	// Tagging operations
	CreateTagging(ctx context.Context, tagging *model.Tagging) error
	GetTaggingsByResourceCode(ctx context.Context, resourceCode string) ([]*model.Tagging, error)
	GetTaggingsByTagID(ctx context.Context, tagID uint) ([]*model.Tagging, error)
	DeleteTagging(ctx context.Context, tagID uint, resourceCode string) error
	DeleteTaggingsByResourceCode(ctx context.Context, resourceCode string) error
	DeleteTaggingsByTagIDs(ctx context.Context, resourceCode string, tagIDs []uint) error

	// Bulk operations
	CreateTaggingsBatch(ctx context.Context, taggings []*model.Tagging) error
	ExistsTagging(ctx context.Context, tagID uint, resourceCode string) (bool, error)

	// Search operations
	SearchResourcesByTags(ctx context.Context, req *request.TagSearchRequest) ([]*model.Tagging, int64, error)
}
