package service

import (
	"context"

	"app/internal/model"
	"app/internal/repository"
)

type NewsService interface {
	GetNewsList(ctx context.Context, keyword string) ([]*model.NewsType, error)
}

func NewNewsService(
	service *Service,
	newsRepository repository.NewsRepository,
) NewsService {
	return &NewsServices{
		Service:        service,
		newsRepository: newsRepository,
	}
}

type NewsServices struct {
	*Service
	newsRepository repository.NewsRepository
}

func (s *NewsServices) GetNewsList(ctx context.Context, keyword string) ([]*model.NewsType, error) {
	return s.newsRepository.GetNewsList(ctx, keyword)
}
