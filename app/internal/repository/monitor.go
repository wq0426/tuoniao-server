package repository

import (
	"context"

	"app/internal/model"
)

type MonitorRepository interface {
	GetMonitorList(ctx context.Context) ([]*model.MonitorResponse, error)
}

func NewMonitorRepository(
	repository *Repository,
) MonitorRepository {
	return &monitorRepository{
		Repository: repository,
	}
}

type monitorRepository struct {
	*Repository
}

func (r *monitorRepository) GetMonitorList(ctx context.Context) ([]*model.MonitorResponse, error) {
	var monitors []*model.Monitor
	if err := r.DB(ctx).Find(&monitors).Error; err != nil {
		return nil, err
	}

	var response []*model.MonitorResponse
	for _, monitor := range monitors {
		response = append(response, &model.MonitorResponse{
			ID:    monitor.ID,
			Img:   monitor.Img,
			Url:   monitor.Url,
			Title: monitor.Title,
		})
	}

	return response, nil
}
