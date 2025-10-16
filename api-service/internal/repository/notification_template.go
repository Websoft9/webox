package repository

import (
	"context"

	"gorm.io/gorm"

	"api-service/internal/dto/request"
	"api-service/internal/interface/repository"
	"api-service/internal/model"
	"api-service/pkg/errors"
)

// notificationTemplateRepository implements NotificationTemplateRepository
type notificationTemplateRepository struct {
	db *gorm.DB
}

// NewNotificationTemplateRepository creates a new notification template repository
func NewNotificationTemplateRepository(db *gorm.DB) repository.NotificationTemplateRepository {
	return &notificationTemplateRepository{
		db: db,
	}
}

// Create creates a new notification template
func (r *notificationTemplateRepository) Create(ctx context.Context, template *model.NotificationTemplate) error {
	if err := r.db.WithContext(ctx).Create(template).Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordCreateFailed)
	}

	return nil
}

// GetByID retrieves a notification template by ID
func (r *notificationTemplateRepository) GetByID(ctx context.Context, id uint) (*model.NotificationTemplate, error) {
	var template model.NotificationTemplate
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&template).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(errors.CodeRecordNotFound)
		}
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return &template, nil
}

// GetList retrieves a paginated list of notification templates with filters
func (r *notificationTemplateRepository) GetList(ctx context.Context, req *request.GetNotificationTemplateListRequest) ([]*model.NotificationTemplate, int64, error) {
	var templates []*model.NotificationTemplate
	var total int64

	query := r.db.WithContext(ctx).Model(&model.NotificationTemplate{})

	// Apply filters
	if req.TemplateType != "" {
		query = query.Where("template_type = ?", req.TemplateType)
	}

	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	if req.IsSystem != nil {
		query = query.Where("is_system = ?", *req.IsSystem)
	}

	if req.Keyword != "" {
		query = query.Where("name LIKE ? OR content LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	// Apply pagination and sorting
	offset := req.GetOffset()
	pageSize := req.GetPageSize()
	orderBy := req.GetSortOrder()

	if err := query.Order(orderBy).Offset(offset).Limit(pageSize).Find(&templates).Error; err != nil {
		return nil, 0, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return templates, total, nil
}

// Update updates an existing notification template
func (r *notificationTemplateRepository) Update(ctx context.Context, template *model.NotificationTemplate) error {
	if err := r.db.WithContext(ctx).Updates(template).Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordUpdateFailed)
	}

	return nil
}

// Delete deletes a notification template by ID
func (r *notificationTemplateRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.NotificationTemplate{}, id).Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordDeleteFailed)
	}

	return nil
}

// ExistsByName checks if a template with the given name exists
func (r *notificationTemplateRepository) ExistsByName(ctx context.Context, name string, excludeID ...uint) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&model.NotificationTemplate{}).Where("name = ?", name)

	if len(excludeID) > 0 && excludeID[0] > 0 {
		query = query.Where("id != ?", excludeID[0])
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return count > 0, nil
}
