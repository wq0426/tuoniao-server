package repository

import (
	"context"

	"app/internal/model"
)

type BannerRepository interface {
	GetBannerList(ctx context.Context) ([]*model.BannerResponse, error)
}

func NewBannerRepository(
	repository *Repository,
) BannerRepository {
	return &bannerRepository{
		Repository: repository,
	}
}

type bannerRepository struct {
	*Repository
}

func (r *bannerRepository) GetBannerList(ctx context.Context) ([]*model.BannerResponse, error) {
	var banners []*model.Banner
	if err := r.DB(ctx).Where("path = ?", "home").Find(&banners).Error; err != nil {
		return nil, err
	}

	var response []*model.BannerResponse
	for _, banner := range banners {
		response = append(response, &model.BannerResponse{
			ID:  banner.ID,
			Img: banner.Img,
			Url: banner.Url,
		})
	}

	return response, nil
}
