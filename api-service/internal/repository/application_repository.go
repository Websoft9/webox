package repository

import (
	"api-service/internal/model"

	"gorm.io/gorm"
)

type ApplicationRepository interface {
	Create(app *model.Application) error
	GetByID(id uint) (*model.Application, error)
	Update(app *model.Application) error
	Delete(id uint) error
	List(offset, limit int) ([]*model.Application, int64, error)
	GetByServerID(serverID uint) ([]*model.Application, error)
}

type applicationRepository struct {
	db *gorm.DB
}

func NewApplicationRepository(db *gorm.DB) ApplicationRepository {
	return &applicationRepository{db: db}
}

func (r *applicationRepository) Create(app *model.Application) error {
	return r.db.Create(app).Error
}

func (r *applicationRepository) GetByID(id uint) (*model.Application, error) {
	var app model.Application
	err := r.db.Preload("Server").First(&app, id).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *applicationRepository) Update(app *model.Application) error {
	return r.db.Save(app).Error
}

func (r *applicationRepository) Delete(id uint) error {
	return r.db.Delete(&model.Application{}, id).Error
}

func (r *applicationRepository) List(offset, limit int) ([]*model.Application, int64, error) {
	var apps []*model.Application
	var total int64

	if err := r.db.Model(&model.Application{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Preload("Server").Offset(offset).Limit(limit).Find(&apps).Error
	return apps, total, err
}

func (r *applicationRepository) GetByServerID(serverID uint) ([]*model.Application, error) {
	var apps []*model.Application
	err := r.db.Where("server_id = ?", serverID).Find(&apps).Error
	return apps, err
}
