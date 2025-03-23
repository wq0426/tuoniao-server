package service

import (
	"app/internal/model"
	"app/internal/repository"
	"errors"

	"github.com/gin-gonic/gin"
)

type WithdrawOrderService interface {
	CreateWithdraw(ctx *gin.Context, req model.CreateWithdrawRequest) error
	GetWithdrawList(ctx *gin.Context, req model.WithdrawQueryRequest) (*model.WithdrawListResponse, error)
	GetWithdrawDetail(ctx *gin.Context, withdrawID uint64) (*model.WithdrawDetailResponse, error)
}

func NewWithdrawOrderService(
	service *Service,
	withdrawOrderRepository repository.WithdrawOrderRepository,
) WithdrawOrderService {
	return &withdrawOrderService{
		Service:                 service,
		withdrawOrderRepository: withdrawOrderRepository,
	}
}

type withdrawOrderService struct {
	*Service
	withdrawOrderRepository repository.WithdrawOrderRepository
}

// CreateWithdraw 创建提现单
func (s *withdrawOrderService) CreateWithdraw(ctx *gin.Context, req model.CreateWithdrawRequest) error {
	userID := GetUserIdFromCtx(ctx)

	// 验证请求数据
	if req.Amount <= 0 {
		return errors.New("提现金额必须大于0")
	}

	// 调用仓储层创建提现单
	return s.withdrawOrderRepository.CreateWithdraw(ctx, userID, req)
}

// GetWithdrawList 获取提现单列表
func (s *withdrawOrderService) GetWithdrawList(ctx *gin.Context, req model.WithdrawQueryRequest) (*model.WithdrawListResponse, error) {
	userID := GetUserIdFromCtx(ctx)

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 调用仓储层获取提现单列表
	return s.withdrawOrderRepository.GetWithdrawList(ctx, userID, req.Status, req.Page, req.PageSize)
}

// GetWithdrawDetail 获取提现单详情
func (s *withdrawOrderService) GetWithdrawDetail(ctx *gin.Context, withdrawID uint64) (*model.WithdrawDetailResponse, error) {
	userID := GetUserIdFromCtx(ctx)

	// 调用仓储层获取提现单详情
	return s.withdrawOrderRepository.GetWithdrawDetail(ctx, userID, withdrawID)
}
