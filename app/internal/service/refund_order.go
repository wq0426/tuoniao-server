package service

import (
	"app/internal/model"
	"app/internal/repository"

	"github.com/gin-gonic/gin"
)

type RefundOrderService interface {
	CreateRefund(ctx *gin.Context, req model.RefundItemRequest) error
	GetRefundList(ctx *gin.Context, req model.RefundQueryRequest) (*model.RefundListResponse, error)
	GetRefundDetail(ctx *gin.Context, refundID uint64) (*model.RefundDetailResponse, error)
	CancelRefund(ctx *gin.Context, refundID uint64) error
	DeleteRefund(ctx *gin.Context, refundID uint64) error
}

func NewRefundOrderService(
	service *Service,
	refundOrderRepository repository.RefundOrderRepository,
) RefundOrderService {
	return &refundOrderService{
		Service:               service,
		refundOrderRepository: refundOrderRepository,
	}
}

type refundOrderService struct {
	*Service
	refundOrderRepository repository.RefundOrderRepository
}

// CreateRefund 创建退款单
func (s *refundOrderService) CreateRefund(ctx *gin.Context, req model.RefundItemRequest) error {
	userID := GetUserIdFromCtx(ctx)

	// 调用仓储层创建退款单
	return s.refundOrderRepository.CreateRefund(ctx, userID, req)
}

// GetRefundList 获取退款单列表
func (s *refundOrderService) GetRefundList(ctx *gin.Context, req model.RefundQueryRequest) (*model.RefundListResponse, error) {
	userID := GetUserIdFromCtx(ctx)

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 调用仓储层获取退款单列表
	return s.refundOrderRepository.GetRefundList(ctx, userID, req.Status, req.Page, req.PageSize)
}

// GetRefundDetail 获取退款单详情
func (s *refundOrderService) GetRefundDetail(ctx *gin.Context, refundID uint64) (*model.RefundDetailResponse, error) {
	userID := GetUserIdFromCtx(ctx)

	// 调用仓储层获取退款单详情
	return s.refundOrderRepository.GetRefundDetail(ctx, userID, refundID)
}

// CancelRefund 撤销退款申请
func (s *refundOrderService) CancelRefund(ctx *gin.Context, refundID uint64) error {
	userID := GetUserIdFromCtx(ctx)

	// 调用仓储层撤销退款
	return s.refundOrderRepository.CancelRefund(ctx, userID, refundID)
}

// DeleteRefund 删除退款记录
func (s *refundOrderService) DeleteRefund(ctx *gin.Context, refundID uint64) error {
	userID := GetUserIdFromCtx(ctx)

	// 调用仓储层删除退款记录
	return s.refundOrderRepository.DeleteRefund(ctx, userID, refundID)
}
