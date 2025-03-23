package service

import (
	"app/internal/model"
	"app/internal/repository"

	"github.com/gin-gonic/gin"
)

type PointExchangeConfigService interface {
	GetPointExchangeConfigList(ctx *gin.Context) (*model.PointExchangeConfigListResponse, error)
	ExchangePoints(ctx *gin.Context, req model.PointExchangeRequest) error
}

func NewPointExchangeConfigService(
	service *Service,
	pointExchangeConfigRepository repository.PointExchangeConfigRepository,
) PointExchangeConfigService {
	return &pointExchangeConfigService{
		Service:                       service,
		pointExchangeConfigRepository: pointExchangeConfigRepository,
	}
}

type pointExchangeConfigService struct {
	*Service
	pointExchangeConfigRepository repository.PointExchangeConfigRepository
}

func (s *pointExchangeConfigService) GetPointExchangeConfigList(ctx *gin.Context) (*model.PointExchangeConfigListResponse, error) {
	userID := GetUserIdFromCtx(ctx)
	return s.pointExchangeConfigRepository.GetPointExchangeConfigList(ctx, userID)
}

func (s *pointExchangeConfigService) ExchangePoints(ctx *gin.Context, req model.PointExchangeRequest) error {
	userID := GetUserIdFromCtx(ctx)
	return s.pointExchangeConfigRepository.ExchangePoints(ctx, userID, req.ConfigID)
}
