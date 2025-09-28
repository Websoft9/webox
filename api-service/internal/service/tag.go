package service

import (
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"api-service/internal/interface/repository"
	"api-service/internal/interface/service"
	"api-service/internal/model"
	"api-service/pkg/errors"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
	"context"
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type tagService struct {
	tagRepo      repository.TagRepository
	db           *gorm.DB
	logger       logger.Logger
	i18nInstance *i18n.I18n
}

// NewTagService creates a new tag service instance
func NewTagService(
	tagRepo repository.TagRepository,
	db *gorm.DB,
	logger logger.Logger,
	i18nInstance *i18n.I18n,
) service.TagService {
	return &tagService{
		tagRepo:      tagRepo,
		db:           db,
		logger:       logger,
		i18nInstance: i18nInstance,
	}
}

// CreateTag creates a new tag
func (s *tagService) CreateTag(ctx context.Context, req *request.TagCreateRequest, userID uint) (*response.TagResponse, error) {
	s.logger.InfoContext(ctx, "Creating tag",
		logger.String("operation", "CreateTag"),
		logger.String("tag_name", req.Name),
		logger.Uint("user_id", uint(userID)))

	// Validate tag name uniqueness
	exists, err := s.tagRepo.ExistsTagByName(ctx, req.Name)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to check tag existence", logger.ErrorField(err))
		return nil, err
	}
	if exists {
		s.logger.WarnContext(ctx, "Tag name already exists", logger.String("tag_name", req.Name))
		return nil, errors.NewAppError(errors.CodeResourceAlreadyExists)
	}

	// Create tag entity
	tag := &model.Tag{
		Name:        req.Name,
		Color:       req.Color,
		Description: req.Description,
		CreatedBy:   userID,
	}

	// Save to database
	if err := s.tagRepo.CreateTag(ctx, tag); err != nil {
		s.logger.ErrorContext(ctx, "Failed to create tag", logger.ErrorField(err))
		return nil, err
	}

	s.logger.InfoContext(ctx, "Tag created successfully", logger.Uint("tag_id", uint(tag.ID)))
	return s.convertTagToResponse(tag), nil
}

// GetTag retrieves a tag by ID
func (s *tagService) GetTag(ctx context.Context, id uint) (*response.TagResponse, error) {
	s.logger.InfoContext(ctx, "Getting tag", logger.Uint("tag_id", uint(id)))

	tag, err := s.tagRepo.GetTagByID(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get tag", logger.ErrorField(err))
		return nil, err
	}

	return s.convertTagToResponse(tag), nil
}

// UpdateTag updates an existing tag
func (s *tagService) UpdateTag(ctx context.Context, id uint, req *request.TagUpdateRequest, userID uint) (*response.TagResponse, error) {
	s.logger.InfoContext(ctx, "Updating tag",
		logger.Uint("tag_id", uint(id)),
		logger.String("tag_name", req.Name),
		logger.Uint("user_id", uint(userID)))

	// Get existing tag
	tag, err := s.tagRepo.GetTagByID(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get tag", logger.ErrorField(err))
		return nil, err
	}

	// Validate name uniqueness if changed
	if tag.Name != req.Name {
		exists, err := s.tagRepo.ExistsTagByNameExcludeID(ctx, req.Name, id)
		if err != nil {
			s.logger.ErrorContext(ctx, "Failed to check tag existence", logger.ErrorField(err))
			return nil, err
		}
		if exists {
			return nil, errors.NewAppError(errors.CodeResourceAlreadyExists)
		}
	}

	// Update tag fields
	tag.Name = req.Name
	tag.Color = req.Color
	tag.Description = req.Description

	// Save changes
	if err := s.tagRepo.UpdateTag(ctx, tag); err != nil {
		s.logger.ErrorContext(ctx, "Failed to update tag", logger.ErrorField(err))
		return nil, err
	}

	s.logger.InfoContext(ctx, "Tag updated successfully", logger.Uint("tag_id", uint(id)))
	return s.convertTagToResponse(tag), nil
}

// DeleteTag deletes a tag
func (s *tagService) DeleteTag(ctx context.Context, id, userID uint) error {
	s.logger.InfoContext(ctx, "Deleting tag",
		logger.Uint("tag_id", uint(id)),
		logger.Uint("user_id", uint(userID)))

	// Check if tag exists
	_, err := s.tagRepo.GetTagByID(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get tag", logger.ErrorField(err))
		return err
	}

	// Delete tag (cascades to taggings due to foreign key)
	if err := s.tagRepo.DeleteTag(ctx, id); err != nil {
		s.logger.ErrorContext(ctx, "Failed to delete tag", logger.ErrorField(err))
		return err
	}

	s.logger.InfoContext(ctx, "Tag deleted successfully", logger.Uint("tag_id", uint(id)))
	return nil
}

// ListTags retrieves tags with optional filtering
func (s *tagService) ListTags(ctx context.Context, req *request.TagListRequest) ([]*response.TagResponse, error) {
	s.logger.InfoContext(ctx, "Listing tags", logger.String("search", req.Search))

	// Parse exclude IDs
	var excludeIDs []uint
	if req.ExcludeIDs != "" {
		idStrings := strings.Split(req.ExcludeIDs, ",")
		for _, idStr := range idStrings {
			if id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 32); err == nil {
				excludeIDs = append(excludeIDs, uint(id))
			}
		}
	}

	// Get tags
	tags, err := s.tagRepo.ListTags(ctx, req.Search, excludeIDs)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to list tags", logger.ErrorField(err))
		return nil, err
	}

	// Convert to response
	responses := make([]*response.TagResponse, 0, len(tags))
	for _, tag := range tags {
		responses = append(responses, s.convertTagToResponse(tag))
	}

	return responses, nil
}

// AssignTags assigns tags to a resource
func (s *tagService) AssignTags(ctx context.Context, req *request.TagAssignRequest, userID uint) (*response.TagAssignResponse, error) {
	s.logger.InfoContext(ctx, "Assigning tags",
		logger.Uint("resource_id", uint(req.ResourceID)),
		logger.Uint("user_id", uint(userID)))

	var allResults []response.TagAssignResult
	var allTagIDs []uint

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Process tag IDs
		if len(req.TagIDs) > 0 {
			results1, tagIDs1, err := s.processTagIDsForAssignment(tx, req.TagIDs, req.ResourceID, userID)
			if err != nil {
				return err
			}
			allResults = append(allResults, results1...)
			allTagIDs = append(allTagIDs, tagIDs1...)
		}

		// Process tag names
		if len(req.TagNames) > 0 {
			results2, tagIDs2, err := s.processTagNamesForAssignment(tx, req.TagNames, req.ResourceID, userID)
			if err != nil {
				return err
			}
			allResults = append(allResults, results2...)
			allTagIDs = append(allTagIDs, tagIDs2...)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Get updated resource tags
	resourceTags, err := s.GetResourceTags(ctx, &request.TaggingListRequest{ResourceID: req.ResourceID})
	if err != nil {
		return nil, err
	}

	// Convert slice to match expected type
	resourceTagsConverted := make([]response.TagSimpleResponse, 0, len(resourceTags))
	for _, tag := range resourceTags {
		resourceTagsConverted = append(resourceTagsConverted, *tag)
	}

	return &response.TagAssignResponse{
		Results:      allResults,
		ResourceTags: resourceTagsConverted,
	}, nil
}

// ReplaceTags replaces all tags for a resource
func (s *tagService) ReplaceTags(ctx context.Context, req *request.TagAssignRequest, userID uint) (*response.TagAssignResponse, error) {
	s.logger.InfoContext(ctx, "Replacing tags",
		logger.Uint("resource_id", uint(req.ResourceID)),
		logger.Uint("user_id", uint(userID)))

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Remove all existing tags using transaction
		if err := tx.Where("resource_id = ?", req.ResourceID).Delete(&model.Tagging{}).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Add new tags
	return s.AssignTags(ctx, req, userID)
}

// UnassignTags removes tag associations from a resource
func (s *tagService) UnassignTags(ctx context.Context, req *request.TagUnassignRequest, userID uint) (*response.TagUnassignResponse, error) {
	s.logger.InfoContext(ctx, "Unassigning tags",
		logger.Uint("resource_id", uint(req.ResourceID)),
		logger.Uint("user_id", uint(userID)))

	// Remove specified tag associations
	if err := s.tagRepo.DeleteTaggingsByTagIDs(ctx, req.ResourceID, req.TagIDs); err != nil {
		s.logger.ErrorContext(ctx, "Failed to remove tag associations", logger.ErrorField(err))
		return nil, err
	}

	return &response.TagUnassignResponse{
		RemovedCount: len(req.TagIDs),
		Message:      fmt.Sprintf("Removed %d tag associations", len(req.TagIDs)),
	}, nil
}

// GetResourceTags retrieves all tags for a resource
func (s *tagService) GetResourceTags(ctx context.Context, req *request.TaggingListRequest) ([]*response.TagSimpleResponse, error) {
	s.logger.InfoContext(ctx, "Getting resource tags", logger.Uint("resource_id", uint(req.ResourceID)))

	taggings, err := s.tagRepo.GetTaggingsByResourceID(ctx, req.ResourceID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get resource tags", logger.ErrorField(err))
		return nil, err
	}

	responses := make([]*response.TagSimpleResponse, 0, len(taggings))
	for _, tagging := range taggings {
		if tagging.Tag != nil {
			responses = append(responses, &response.TagSimpleResponse{
				ID:    tagging.Tag.ID,
				Name:  tagging.Tag.Name,
				Color: tagging.Tag.Color,
			})
		}
	}

	return responses, nil
}

// SearchResourcesByTags searches resources by tags
func (s *tagService) SearchResourcesByTags(ctx context.Context, req *request.TagSearchRequest) (*response.TagSearchResponse, error) {
	s.logger.InfoContext(ctx, "Searching resources by tags")

	// Collect all tag IDs
	allTagIDs, err := s.collectTagIDs(ctx, req)
	if err != nil {
		return nil, err
	}

	if len(allTagIDs) == 0 {
		return &response.TagSearchResponse{
			Total:     0,
			Page:      req.Page,
			PageSize:  req.PageSize,
			Resources: []response.TaggedResource{},
		}, nil
	}

	// Set defaults
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}
	if req.Operation == "" {
		req.Operation = "AND"
	}

	// Calculate offset
	offset := (req.Page - 1) * req.PageSize

	// Search resources
	taggings, total, err := s.tagRepo.SearchResourcesByTags(ctx, allTagIDs, req.Operation, offset, req.PageSize)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to search resources by tags", logger.ErrorField(err))
		return nil, err
	}

	// Build resource map
	resourceMap := s.buildSimpleResourceMap(taggings)

	// Convert to slice
	resources := make([]response.TaggedResource, 0, len(resourceMap))
	for _, resource := range resourceMap {
		resources = append(resources, *resource)
	}

	return &response.TagSearchResponse{
		Total:     total,
		Page:      req.Page,
		PageSize:  req.PageSize,
		Resources: resources,
	}, nil
}

// buildSimpleResourceMap builds a simple resource map from taggings
func (s *tagService) buildSimpleResourceMap(taggings []*model.Tagging) map[uint]*response.TaggedResource {
	resourceMap := make(map[uint]*response.TaggedResource)

	for _, tagging := range taggings {
		if resource, exists := resourceMap[tagging.ResourceID]; exists {
			// Add tag to existing resource
			if tagging.Tag != nil {
				resource.Tags = append(resource.Tags, response.TagSimpleResponse{
					ID:    tagging.Tag.ID,
					Name:  tagging.Tag.Name,
					Color: tagging.Tag.Color,
				})
			}
		} else {
			// Create new resource entry
			resource := &response.TaggedResource{
				ResourceID:   tagging.ResourceID,
				ResourceName: fmt.Sprintf("Resource %d", tagging.ResourceID),
				Tags:         []response.TagSimpleResponse{},
				MatchedTags:  []uint{},
				CreatedAt:    tagging.CreatedAt,
			}
			if tagging.Tag != nil {
				resource.Tags = append(resource.Tags, response.TagSimpleResponse{
					ID:    tagging.Tag.ID,
					Name:  tagging.Tag.Name,
					Color: tagging.Tag.Color,
				})
			}
			resourceMap[tagging.ResourceID] = resource
		}
	}

	return resourceMap
}

// SearchTags searches for tags by name pattern
func (s *tagService) SearchTags(ctx context.Context, req *request.TagNameSearchRequest) ([]*response.TagResponse, error) {
	tags, err := s.tagRepo.SearchTagsByName(ctx, req.Q)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to search tags", logger.String("query", req.Q), logger.ErrorField(err))
		return nil, err
	}

	tagResponses := make([]*response.TagResponse, len(tags))
	for i, tag := range tags {
		tagResponses[i] = &response.TagResponse{
			ID:          tag.ID,
			Name:        tag.Name,
			Color:       tag.Color,
			Description: tag.Description,
			CreatedAt:   tag.CreatedAt,
			UpdatedAt:   tag.UpdatedAt,
		}
	}

	return tagResponses, nil
}

// convertTagToResponse converts a tag model to response DTO
func (s *tagService) convertTagToResponse(tag *model.Tag) *response.TagResponse {
	return &response.TagResponse{
		ID:          tag.ID,
		Name:        tag.Name,
		Color:       tag.Color,
		Description: tag.Description,
		CreatedBy:   tag.CreatedBy,
		CreatedAt:   tag.CreatedAt,
		UpdatedAt:   tag.UpdatedAt,
	}
}

// processTagIDsForAssignment processes existing tag IDs and creates associations
func (s *tagService) processTagIDsForAssignment(tx *gorm.DB, tagIDs []uint, resourceID, userID uint) ([]response.TagAssignResult, []uint, error) {
	results := make([]response.TagAssignResult, 0, len(tagIDs))
	processedTagIDs := make([]uint, 0, len(tagIDs))

	for _, tagID := range tagIDs {
		// Check if tag exists using transaction
		var tag model.Tag
		err := tx.First(&tag, tagID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue // Skip non-existent tags
			}
			return nil, nil, err
		}

		// Check if already associated using transaction
		var count int64
		err = tx.Model(&model.Tagging{}).
			Where("tag_id = ? AND resource_id = ?", tagID, resourceID).
			Count(&count).Error
		if err != nil {
			return nil, nil, err
		}

		if count == 0 {
			// Create association using transaction
			tagging := &model.Tagging{
				TagID:      tagID,
				ResourceID: resourceID,
				CreatedBy:  userID,
			}
			if err := tx.Create(tagging).Error; err != nil {
				return nil, nil, err
			}
		}

		results = append(results, response.TagAssignResult{
			Name:    tag.Name,
			TagID:   tagID,
			Status:  "associated",
			Message: "Tag associated",
		})
		processedTagIDs = append(processedTagIDs, tagID)
	}

	return results, processedTagIDs, nil
}

// processTagNamesForAssignment processes tag names and creates new tags if needed
func (s *tagService) processTagNamesForAssignment(tx *gorm.DB, tagNames []string, resourceID, userID uint) ([]response.TagAssignResult, []uint, error) {
	results := make([]response.TagAssignResult, 0, len(tagNames))
	processedTagIDs := make([]uint, 0, len(tagNames))

	for _, tagName := range tagNames {
		// Try to get existing tag using transaction
		var tag model.Tag
		err := tx.Where("name = ?", tagName).First(&tag).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, err
		}

		// Create tag if not exists
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tag = model.Tag{
				Name:      tagName,
				CreatedBy: userID,
			}
			if createErr := tx.Create(&tag).Error; createErr != nil {
				return nil, nil, createErr
			}
			results = append(results, response.TagAssignResult{
				Name:    tagName,
				TagID:   tag.ID,
				Status:  "created",
				Message: "Tag created and associated",
			})
		} else {
			results = append(results, response.TagAssignResult{
				Name:    tagName,
				TagID:   tag.ID,
				Status:  "associated",
				Message: "Tag associated",
			})
		}

		// Create association if not exists using transaction
		var count int64
		err = tx.Model(&model.Tagging{}).
			Where("tag_id = ? AND resource_id = ?", tag.ID, resourceID).
			Count(&count).Error
		if err != nil {
			return nil, nil, err
		}

		if count == 0 {
			tagging := &model.Tagging{
				TagID:      tag.ID,
				ResourceID: resourceID,
				CreatedBy:  userID,
			}
			if err := tx.Create(tagging).Error; err != nil {
				return nil, nil, err
			}
		}

		processedTagIDs = append(processedTagIDs, tag.ID)
	}

	return results, processedTagIDs, nil
}

// collectTagIDs collects all tag IDs from request and tag names
func (s *tagService) collectTagIDs(ctx context.Context, req *request.TagSearchRequest) ([]uint, error) {
	allTagIDs := make([]uint, 0, len(req.TagIDs)+len(req.TagNames))
	allTagIDs = append(allTagIDs, req.TagIDs...)

	for _, tagName := range req.TagNames {
		tag, err := s.tagRepo.GetTagByName(ctx, tagName)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
			// Skip non-existent tags
			continue
		}
		allTagIDs = append(allTagIDs, tag.ID)
	}

	return allTagIDs, nil
}
