package service

import (
	"app/internal/model"
	"app/internal/repository"

	"github.com/gin-gonic/gin"
)

type FreeMarketMineService interface {
	GetUserEggsSummary(ctx *gin.Context) (*model.FreeMarketMineResponse, error)
	UpdateEggPrice(ctx *gin.Context, price float64, id int) error
}

func NewFreeMarketMineService(
	service *Service,
	freeMarketMineRepository repository.FreeMarketMineRepository,
) FreeMarketMineService {
	return &freeMarketMineService{
		Service:                  service,
		freeMarketMineRepository: freeMarketMineRepository,
	}
}

type freeMarketMineService struct {
	*Service
	freeMarketMineRepository repository.FreeMarketMineRepository
}

func (s *freeMarketMineService) GetUserEggsSummary(ctx *gin.Context) (*model.FreeMarketMineResponse, error) {
	userID := GetUserIdFromCtx(ctx)
	return s.freeMarketMineRepository.GetUserEggsSummary(ctx, userID)
}

func (s *freeMarketMineService) UpdateEggPrice(ctx *gin.Context, price float64, id int) error {
	return s.freeMarketMineRepository.UpdateEggPrice(ctx, price, id)
}
