package repository

import (
	"context"

	"gorm.io/gorm"

	"api-service/internal/dto/request"
	"api-service/internal/interface/repository"
	"api-service/internal/model"
	"api-service/pkg/errors"
)

const (
	dbConnectionCodePrefix = "database" // Prefix for database connection code
	dbConnectionCodeLength = 10         // Length of random string in connection code
)

// databaseConnectionRepository implements DatabaseConnectionRepository
type databaseConnectionRepository struct {
	db *gorm.DB
}

// NewDatabaseConnectionRepository creates a new database connection repository
func NewDatabaseConnectionRepository(db *gorm.DB) repository.DatabaseConnectionRepository {
	return &databaseConnectionRepository{
		db: db,
	}
}

// Create creates a new database connection and generates its code
func (r *databaseConnectionRepository) Create(ctx context.Context, conn *model.DatabaseConnection) error {
	// Generate code before creating
	conn.Code = model.GenerateCode(dbConnectionCodePrefix, dbConnectionCodeLength)

	// Create the connection
	if err := r.db.WithContext(ctx).Create(conn).Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordCreateFailed)
	}

	return nil
}

// GetByID retrieves a database connection by ID
func (r *databaseConnectionRepository) GetByID(ctx context.Context, id uint) (*model.DatabaseConnection, error) {
	var conn model.DatabaseConnection
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&conn).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(errors.CodeRecordNotFound)
		}
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return &conn, nil
}

// GetList retrieves a paginated list of database connections with filters
func (r *databaseConnectionRepository) GetList(ctx context.Context, req *request.GetDatabaseConnectionListRequest, ownerID uint) ([]*model.DatabaseConnection, int64, error) {
	var connections []*model.DatabaseConnection
	var total int64

	query := r.db.WithContext(ctx).Model(&model.DatabaseConnection{})

	// Filter by owner
	query = query.Where("owner_id = ?", ownerID)

	// Apply filters
	if req.DBType != nil && *req.DBType != "" {
		query = query.Where("db_type = ?", *req.DBType)
	}

	// Filter by code
	if req.Code != nil && *req.Code != "" {
		query = query.Where("code = ?", *req.Code)
	}

	if req.Keyword != "" {
		query = query.Where("name LIKE ? OR host LIKE ? OR description LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	// Apply pagination and sorting
	offset := req.GetOffset()
	pageSize := req.GetPageSize()
	orderBy := req.GetSortOrder()

	if err := query.Order(orderBy).Offset(offset).Limit(pageSize).Find(&connections).Error; err != nil {
		return nil, 0, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return connections, total, nil
}

// Update updates an existing database connection
func (r *databaseConnectionRepository) Update(ctx context.Context, conn *model.DatabaseConnection) error {
	if err := r.db.WithContext(ctx).Updates(conn).Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordUpdateFailed)
	}

	return nil
}

// Delete deletes a database connection by ID
func (r *databaseConnectionRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.DatabaseConnection{}, id).Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordDeleteFailed)
	}

	return nil
}
