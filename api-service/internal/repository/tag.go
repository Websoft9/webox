package repository

import (
	"api-service/internal/dto/request"
	"api-service/internal/interface/repository"
	"api-service/internal/model"
	"api-service/pkg/errors"
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
func (r *tagRepository) GetTagByID(ctx context.Context, id uint) (*model.Tag, error) {
	var tag model.Tag
	err := r.db.WithContext(ctx).First(&tag, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewAppError(errors.CodeRecordNotFound)
		}
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}
	return &tag, nil
}

// GetTagByName retrieves a tag by name
func (r *tagRepository) GetTagByName(ctx context.Context, name string) (*model.Tag, error) {
	var tag model.Tag
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&tag).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewAppError(errors.CodeRecordNotFound)
		}
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}
	return &tag, nil
}

// UpdateTag updates an existing tag
func (r *tagRepository) UpdateTag(ctx context.Context, tag *model.Tag) error {
	return r.db.WithContext(ctx).Save(tag).Error
}

// DeleteTag deletes a tag by ID
func (r *tagRepository) DeleteTag(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Tag{}, id).Error
}

// ListTags retrieves tags with optional search and exclusion
func (r *tagRepository) ListTags(ctx context.Context, search string, excludeIDs []uint) ([]*model.Tag, error) {
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
func (r *tagRepository) ListTagsWithUsageCount(ctx context.Context, search string, excludeIDs []uint) ([]*model.Tag, error) {
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
func (r *tagRepository) ExistsTagByNameExcludeID(ctx context.Context, name string, excludeID uint) (bool, error) {
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
func (r *tagRepository) GetTaggingsByResourceCode(ctx context.Context, resourceCode string) ([]*model.Tagging, error) {
	var taggings []*model.Tagging
	err := r.db.WithContext(ctx).
		Preload("Tag").
		Where("resource_code = ?", resourceCode).
		Find(&taggings).Error
	return taggings, err
}

// GetTaggingsByTagID retrieves all taggings for a tag
func (r *tagRepository) GetTaggingsByTagID(ctx context.Context, tagID uint) ([]*model.Tagging, error) {
	var taggings []*model.Tagging
	err := r.db.WithContext(ctx).Where("tag_id = ?", tagID).Find(&taggings).Error
	return taggings, err
}

// DeleteTagging deletes a specific tag-resource association
func (r *tagRepository) DeleteTagging(ctx context.Context, tagID uint, resourceCode string) error {
	return r.db.WithContext(ctx).
		Where("tag_id = ? AND resource_code = ?", tagID, resourceCode).
		Delete(&model.Tagging{}).Error
}

// DeleteTaggingsByResourceID deletes all taggings for a resource
func (r *tagRepository) DeleteTaggingsByResourceCode(ctx context.Context, resourceCode string) error {
	return r.db.WithContext(ctx).
		Where("resource_code = ?", resourceCode).
		Delete(&model.Tagging{}).Error
}

// DeleteTaggingsByTagIDs deletes specific tag associations for a resource
func (r *tagRepository) DeleteTaggingsByTagIDs(ctx context.Context, resourceCode string, tagIDs []uint) error {
	return r.db.WithContext(ctx).
		Where("resource_code = ? AND tag_id IN ?", resourceCode, tagIDs).
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
func (r *tagRepository) ExistsTagging(ctx context.Context, tagID uint, resourceCode string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Tagging{}).
		Where("tag_id = ? AND resource_code = ?", tagID, resourceCode).Count(&count).Error
	return count > 0, err
}

func (r *tagRepository) SearchResourcesByTags(ctx context.Context, req *request.TagSearchRequest) ([]*model.Tagging, int64, error) {
	var taggings []*model.Tagging
	var total int64

	// Collect all tag IDs from service layer
	allTagIDs := make([]uint, 0, len(req.TagIDs)+len(req.TagNames))
	allTagIDs = append(allTagIDs, req.TagIDs...)

	// Convert tag names to IDs
	for _, tagName := range req.TagNames {
		var tag model.Tag
		if err := r.db.WithContext(ctx).Where("name = ?", tagName).First(&tag).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, 0, err
			}
			// Skip non-existent tags
			continue
		}
		allTagIDs = append(allTagIDs, tag.ID)
	}

	if len(allTagIDs) == 0 {
		return []*model.Tagging{}, 0, nil
	}

	// Build base query
	query := r.db.WithContext(ctx).Model(&model.Tagging{}).
		Preload("Tag").
		Where("tag_id IN (?)", allTagIDs)

	// Apply operation logic (AND/OR)
	if req.Operation == "AND" && len(allTagIDs) > 1 {
		// For AND operation, find resources that have ALL specified tags
		subQuery := r.db.Model(&model.Tagging{}).
			Select("resource_code").
			Where("tag_id IN (?)", allTagIDs).
			Group("resource_code").
			Having("COUNT(DISTINCT tag_id) = ?", len(allTagIDs))

		query = query.Where("resource_code IN (?)", subQuery)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and sorting
	offset := req.GetOffset()
	limit := req.GetPageSize()

	err := query.Order(req.GetSortOrder()).
		Offset(offset).Limit(limit).
		Find(&taggings).Error

	if err != nil {
		return nil, 0, err
	}

	return taggings, total, nil
}
