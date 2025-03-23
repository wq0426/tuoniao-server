package service

import (
	"context"

	"app/internal/model"
	"app/internal/repository"
)

type BannerService interface {
	GetBannerList(ctx context.Context) ([]*model.BannerResponse, error)
}

func NewBannerService(
	service *Service,
	bannerRepository repository.BannerRepository,
) BannerService {
	return &BannerServices{
		Service:          service,
		bannerRepository: bannerRepository,
	}
}

type BannerServices struct {
	*Service
	bannerRepository repository.BannerRepository
}

func (s *BannerServices) GetBannerList(ctx context.Context) ([]*model.BannerResponse, error) {
	return s.bannerRepository.GetBannerList(ctx)
}
