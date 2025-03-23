package service

import (
	"context"

	"app/internal/model"
	"app/internal/repository"
)

type MonitorService interface {
	GetMonitorList(ctx context.Context) ([]*model.MonitorResponse, error)
}

func NewMonitorService(
	service *Service,
	monitorRepository repository.MonitorRepository,
) MonitorService {
	return &MonitorServices{
		Service:           service,
		monitorRepository: monitorRepository,
	}
}

type MonitorServices struct {
	*Service
	monitorRepository repository.MonitorRepository
}

func (s *MonitorServices) GetMonitorList(ctx context.Context) ([]*model.MonitorResponse, error) {
	return s.monitorRepository.GetMonitorList(ctx)
}
