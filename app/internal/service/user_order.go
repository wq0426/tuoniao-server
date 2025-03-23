package service

import (
	"app/internal/model"
	"app/internal/repository"

	"errors"

	"github.com/gin-gonic/gin"
)

type UserOrderService interface {
	CreateOrders(ctx *gin.Context, req model.CreateOrderRequest) error
	GetOrderList(ctx *gin.Context, req model.OrderQueryRequest) (*model.OrderListResponse, error)
	GetOrderProductDetails(ctx *gin.Context, req model.OrderProductsRequest) (*model.OrderProductsResponse, error)
	UpdateOrderStatus(ctx *gin.Context, req model.UpdateOrderStatusRequest) error
	GetOrderDetail(ctx *gin.Context, orderItemID uint64) (*model.OrderDetailResponse, error)
}

func NewUserOrderService(
	service *Service,
	userOrderRepository repository.UserOrderRepository,
) UserOrderService {
	return &userOrderService{
		Service:             service,
		userOrderRepository: userOrderRepository,
	}
}

type userOrderService struct {
	*Service
	userOrderRepository repository.UserOrderRepository
}

func (s *userOrderService) CreateOrders(ctx *gin.Context, req model.CreateOrderRequest) error {
	// Get user ID from context (assuming it's set by auth middleware)
	userID := GetUserIdFromCtx(ctx)

	// Pass the request to the repository
	return s.userOrderRepository.CreateOrders(ctx, userID, req)
}

func (s *userOrderService) GetOrderList(ctx *gin.Context, req model.OrderQueryRequest) (*model.OrderListResponse, error) {
	userID := GetUserIdFromCtx(ctx)
	return s.userOrderRepository.GetOrderList(ctx, userID, req)
}

func (s *userOrderService) GetOrderProductDetails(ctx *gin.Context, req model.OrderProductsRequest) (*model.OrderProductsResponse, error) {
	userID := GetUserIdFromCtx(ctx)

	// 获取明细
	orderItem, err := s.userOrderRepository.GetOrderDetail(ctx, userID, req.OrderItemID)
	if err != nil {
		return nil, errors.New("无权限查看该订单")
	}

	products, err := s.userOrderRepository.GetOrderProductDetails(ctx, orderItem.OrderID, req.OrderItemID)
	if err != nil {
		return nil, err
	}

	return &model.OrderProductsResponse{
		Products: products,
	}, nil
}

func (s *userOrderService) UpdateOrderStatus(ctx *gin.Context, req model.UpdateOrderStatusRequest) error {
	userID := GetUserIdFromCtx(ctx)

	// 调用仓储层更新订单状态
	return s.userOrderRepository.UpdateOrderStatus(ctx, userID, req.OrderItemID, req.Status)
}

func (s *userOrderService) GetOrderDetail(ctx *gin.Context, orderItemID uint64) (*model.OrderDetailResponse, error) {
	userID := GetUserIdFromCtx(ctx)

	// 调用仓储层获取订单详情
	return s.userOrderRepository.GetOrderDetail(ctx, userID, orderItemID)
}
