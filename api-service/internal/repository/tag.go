package repository

import (
	"api-service/internal/interface/repository"
	"api-service/internal/model"
	"context"

	"gorm.io/gorm"
)

// tagRepository implements tag data access operations
type tagRepository struct {
	db *gorm.DB
}

// NewTagRepository creates a new tag repository instance
func NewTagRepository(db *gorm.DB) repository.TagRepository {
	return &tagRepository{db: db}
}

// CreateTag creates a new tag
func (r *tagRepository) CreateTag(ctx context.Context, tag *model.Tag) error {
	return r.db.WithContext(ctx).Create(tag).Error
}

// GetTagByID retrieves a tag by ID
func (r *tagRepository) GetTagByID(ctx context.Context, id uint64) (*model.Tag, error) {
	var tag model.Tag
	err := r.db.WithContext(ctx).First(&tag, id).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// GetTagByName retrieves a tag by name
func (r *tagRepository) GetTagByName(ctx context.Context, name string) (*model.Tag, error) {
	var tag model.Tag
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// UpdateTag updates an existing tag
func (r *tagRepository) UpdateTag(ctx context.Context, tag *model.Tag) error {
	return r.db.WithContext(ctx).Save(tag).Error
}

// DeleteTag deletes a tag by ID
func (r *tagRepository) DeleteTag(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.Tag{}, id).Error
}

// ListTags retrieves tags with optional search and exclusion
func (r *tagRepository) ListTags(ctx context.Context, search string, excludeIDs []uint64) ([]*model.Tag, error) {
	var tags []*model.Tag
	query := r.db.WithContext(ctx)

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	if len(excludeIDs) > 0 {
		query = query.Where("id NOT IN ?", excludeIDs)
	}

	err := query.Order("name ASC").Find(&tags).Error
	return tags, err
}

// ListTagsWithUsageCount retrieves tags with usage count
func (r *tagRepository) ListTagsWithUsageCount(ctx context.Context, search string, excludeIDs []uint64) ([]*model.Tag, error) {
	var tags []*model.Tag
	query := r.db.WithContext(ctx).
		Select("tags.*, COUNT(taggings.id) as usage_count").
		Joins("LEFT JOIN taggings ON tags.id = taggings.tag_id").
		Group("tags.id")

	if search != "" {
		query = query.Where("tags.name LIKE ?", "%"+search+"%")
	}

	if len(excludeIDs) > 0 {
		query = query.Where("tags.id NOT IN ?", excludeIDs)
	}

	err := query.Order("tags.name ASC").Find(&tags).Error
	return tags, err
}

// SearchTagsByName searches for tags by name pattern
func (r *tagRepository) SearchTagsByName(ctx context.Context, query string) ([]*model.Tag, error) {
	var tags []*model.Tag

	dbQuery := r.db.WithContext(ctx).Model(&model.Tag{})

	if query != "" {
		dbQuery = dbQuery.Where("name LIKE ? OR description LIKE ?", "%"+query+"%", "%"+query+"%")
	}

	err := dbQuery.Order("name ASC").Find(&tags).Error
	return tags, err
}

// ExistsTagByName checks if a tag exists by name
func (r *tagRepository) ExistsTagByName(ctx context.Context, name string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Tag{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}

// ExistsTagByNameExcludeID checks if a tag exists by name excluding a specific ID
func (r *tagRepository) ExistsTagByNameExcludeID(ctx context.Context, name string, excludeID uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Tag{}).
		Where("name = ? AND id != ?", name, excludeID).Count(&count).Error
	return count > 0, err
}

// CreateTagging creates a new tag-resource association
func (r *tagRepository) CreateTagging(ctx context.Context, tagging *model.Tagging) error {
	return r.db.WithContext(ctx).Create(tagging).Error
}

// GetTaggingsByResourceID retrieves all taggings for a resource
func (r *tagRepository) GetTaggingsByResourceID(ctx context.Context, resourceID uint64) ([]*model.Tagging, error) {
	var taggings []*model.Tagging
	err := r.db.WithContext(ctx).
		Preload("Tag").
		Where("resource_id = ?", resourceID).
		Find(&taggings).Error
	return taggings, err
}

// GetTaggingsByTagID retrieves all taggings for a tag
func (r *tagRepository) GetTaggingsByTagID(ctx context.Context, tagID uint64) ([]*model.Tagging, error) {
	var taggings []*model.Tagging
	err := r.db.WithContext(ctx).Where("tag_id = ?", tagID).Find(&taggings).Error
	return taggings, err
}

// DeleteTagging deletes a specific tag-resource association
func (r *tagRepository) DeleteTagging(ctx context.Context, tagID, resourceID uint64) error {
	return r.db.WithContext(ctx).
		Where("tag_id = ? AND resource_id = ?", tagID, resourceID).
		Delete(&model.Tagging{}).Error
}

// DeleteTaggingsByResourceID deletes all taggings for a resource
func (r *tagRepository) DeleteTaggingsByResourceID(ctx context.Context, resourceID uint64) error {
	return r.db.WithContext(ctx).
		Where("resource_id = ?", resourceID).
		Delete(&model.Tagging{}).Error
}

// DeleteTaggingsByTagIDs deletes specific tag associations for a resource
func (r *tagRepository) DeleteTaggingsByTagIDs(ctx context.Context, resourceID uint64, tagIDs []uint64) error {
	return r.db.WithContext(ctx).
		Where("resource_id = ? AND tag_id IN ?", resourceID, tagIDs).
		Delete(&model.Tagging{}).Error
}

// CreateTaggingsBatch creates multiple tag-resource associations in batch
func (r *tagRepository) CreateTaggingsBatch(ctx context.Context, taggings []*model.Tagging) error {
	if len(taggings) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&taggings).Error
}

// ExistsTagging checks if a tag-resource association exists
func (r *tagRepository) ExistsTagging(ctx context.Context, tagID, resourceID uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Tagging{}).
		Where("tag_id = ? AND resource_id = ?", tagID, resourceID).Count(&count).Error
	return count > 0, err
}

// SearchResourcesByTags searches resources by tags with AND/OR operation
func (r *tagRepository) SearchResourcesByTags(ctx context.Context, tagIDs []uint64, operation string, offset, limit int) ([]*model.Tagging, int64, error) {
	var taggings []*model.Tagging
	var total int64

	if len(tagIDs) == 0 {
		return taggings, total, nil
	}

	query := r.db.WithContext(ctx).Model(&model.Tagging{}).Preload("Tag")

	if operation == "AND" {
		// For AND operation, find resources that have all specified tags
		resourceIDs := []uint64{}
		subQuery := r.db.WithContext(ctx).Model(&model.Tagging{}).
			Select("resource_id").
			Where("tag_id IN ?", tagIDs).
			Group("resource_id").
			Having("COUNT(DISTINCT tag_id) = ?", len(tagIDs))

		err := subQuery.Pluck("resource_id", &resourceIDs).Error
		if err != nil {
			return nil, 0, err
		}

		if len(resourceIDs) == 0 {
			return taggings, 0, nil
		}

		query = query.Where("resource_id IN ?", resourceIDs)
	} else {
		// For OR operation, find resources that have any of the specified tags
		query = query.Where("tag_id IN ?", tagIDs)
	}

	// Count total
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err = query.Offset(offset).Limit(limit).Find(&taggings).Error
	return taggings, total, err
}
