package service

import (
	"api-service/internal/model"
	"api-service/internal/repository"
)

type ApplicationService interface {
	CreateApplication(app *model.Application) error
	GetApplication(id uint) (*model.Application, error)
	UpdateApplication(app *model.Application) error
	DeleteApplication(id uint) error
	ListApplications(page, pageSize int) ([]*model.Application, int64, error)
	GetApplicationsByServer(serverID uint) ([]*model.Application, error)
	DeployApplication(appID uint) error
	StopApplication(appID uint) error
	RestartApplication(appID uint) error
}

type applicationService struct {
	appRepo repository.ApplicationRepository
}

func NewApplicationService(appRepo repository.ApplicationRepository) ApplicationService {
	return &applicationService{
		appRepo: appRepo,
	}
}

func (s *applicationService) CreateApplication(app *model.Application) error {
	return s.appRepo.Create(app)
}

func (s *applicationService) GetApplication(id uint) (*model.Application, error) {
	return s.appRepo.GetByID(id)
}

func (s *applicationService) UpdateApplication(app *model.Application) error {
	return s.appRepo.Update(app)
}

func (s *applicationService) DeleteApplication(id uint) error {
	return s.appRepo.Delete(id)
}

func (s *applicationService) ListApplications(page, pageSize int) ([]*model.Application, int64, error) {
	offset := (page - 1) * pageSize
	return s.appRepo.List(offset, pageSize)
}

func (s *applicationService) GetApplicationsByServer(serverID uint) ([]*model.Application, error) {
	return s.appRepo.GetByServerID(serverID)
}

func (s *applicationService) DeployApplication(appID uint) error {
	app, err := s.appRepo.GetByID(appID)
	if err != nil {
		return err
	}

	// TODO: 实现应用部署逻辑
	// 这里应该调用Agent的gRPC接口来部署应用
	app.Status = "running"
	return s.appRepo.Update(app)
}

func (s *applicationService) StopApplication(appID uint) error {
	app, err := s.appRepo.GetByID(appID)
	if err != nil {
		return err
	}

	// TODO: 实现应用停止逻辑
	app.Status = "stopped"
	return s.appRepo.Update(app)
}

func (s *applicationService) RestartApplication(appID uint) error {
	app, err := s.appRepo.GetByID(appID)
	if err != nil {
		return err
	}

	// TODO: 实现应用重启逻辑
	app.Status = "running"
	return s.appRepo.Update(app)
}
