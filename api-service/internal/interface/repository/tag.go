package repository

import (
	"api-service/internal/model"
	"context"
)

// TagRepository interface for tag data access operations
type TagRepository interface {
	// Basic CRUD operations for tags
	CreateTag(ctx context.Context, tag *model.Tag) error
	GetTagByID(ctx context.Context, id uint64) (*model.Tag, error)
	GetTagByName(ctx context.Context, name string) (*model.Tag, error)
	UpdateTag(ctx context.Context, tag *model.Tag) error
	DeleteTag(ctx context.Context, id uint64) error

	// Tag listing and searching
	ListTags(ctx context.Context, search string, excludeIDs []uint64) ([]*model.Tag, error)
	ListTagsWithUsageCount(ctx context.Context, search string, excludeIDs []uint64) ([]*model.Tag, error)
	SearchTagsByName(ctx context.Context, query string) ([]*model.Tag, error)

	// Tag existence checks
	ExistsTagByName(ctx context.Context, name string) (bool, error)
	ExistsTagByNameExcludeID(ctx context.Context, name string, excludeID uint64) (bool, error)

	// Tagging operations
	CreateTagging(ctx context.Context, tagging *model.Tagging) error
	GetTaggingsByResourceID(ctx context.Context, resourceID uint64) ([]*model.Tagging, error)
	GetTaggingsByTagID(ctx context.Context, tagID uint64) ([]*model.Tagging, error)
	DeleteTagging(ctx context.Context, tagID, resourceID uint64) error
	DeleteTaggingsByResourceID(ctx context.Context, resourceID uint64) error
	DeleteTaggingsByTagIDs(ctx context.Context, resourceID uint64, tagIDs []uint64) error

	// Bulk operations
	CreateTaggingsBatch(ctx context.Context, taggings []*model.Tagging) error
	ExistsTagging(ctx context.Context, tagID, resourceID uint64) (bool, error)

	// Search operations
	SearchResourcesByTags(ctx context.Context, tagIDs []uint64, operation string, offset, limit int) ([]*model.Tagging, int64, error)
}
