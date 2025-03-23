package service

import (
	"app/internal/model"
	"app/internal/repository"
	"errors"

	"github.com/gin-gonic/gin"
)

type UserEarningService interface {
	AddEarning(ctx *gin.Context, req model.AddEarningRequest) error
	GetEarningList(ctx *gin.Context, req model.QueryEarningRequest) (*model.EarningListResponse, error)
}

type userEarningService struct {
	*Service
	userEarningRepository repository.UserEarningRepository
}

func NewUserEarningService(
	service *Service,
	userEarningRepository repository.UserEarningRepository,
) UserEarningService {
	return &userEarningService{
		Service:               service,
		userEarningRepository: userEarningRepository,
	}
}

// AddEarning 添加用户收益
func (s *userEarningService) AddEarning(ctx *gin.Context, req model.AddEarningRequest) error {
	userID := GetUserIdFromCtx(ctx)

	// 验证收益类型
	if req.EarningType < 1 || req.EarningType > 3 {
		return errors.New("无效的收益类型")
	}

	// 调用仓储层添加收益
	return s.userEarningRepository.AddEarning(ctx, userID, req)
}

// GetEarningList 获取用户收益列表
func (s *userEarningService) GetEarningList(ctx *gin.Context, req model.QueryEarningRequest) (*model.EarningListResponse, error) {
	userID := GetUserIdFromCtx(ctx)

	// 调用仓储层获取收益列表
	return s.userEarningRepository.GetEarningList(ctx, userID, req)
}
