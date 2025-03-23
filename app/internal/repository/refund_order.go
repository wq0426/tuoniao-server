package repository

import (
	"context"
	"errors"
	"time"

	"app/internal/common"
	"app/internal/model"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RefundOrderRepository interface {
	CreateRefund(ctx context.Context, userID string, req model.RefundItemRequest) error
	GetRefundList(ctx context.Context, userID string, status *uint8, page, pageSize int) (*model.RefundListResponse, error)
	GetRefundDetail(ctx context.Context, userID string, refundID uint64) (*model.RefundDetailResponse, error)
	CancelRefund(ctx context.Context, userID string, refundID uint64) error
	GetRefundDetailByID(ctx context.Context, refundID uint64) (*model.RefundDetailResponse, error)
	DeleteRefund(ctx context.Context, userID string, refundID uint64) error
}

type refundOrderRepository struct {
	*Repository
}

func NewRefundOrderRepository(repository *Repository) RefundOrderRepository {
	return &refundOrderRepository{
		Repository: repository,
	}
}

// CreateRefund 创建退款单
func (r *refundOrderRepository) CreateRefund(ctx context.Context, userID string, req model.RefundItemRequest) error {
	// 校验订单是否存在
	var orderItem model.UserOrderItem
	if err := r.DB(ctx).Where("id = ? AND user_id = ?", req.OrderItemID, userID).First(&orderItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("订单不存在或无权限")
		}
		return err
	}

	// 校验订单状态是否允许退款（已支付的订单才能退款）
	if orderItem.Status < 1 {
		return errors.New("订单状态不允许退款")
	}

	// 创建退款单号
	refundNo := "RF" + common.GenerateOrderNo()
	now := time.Now()

	// 准备退款单数据
	refundOrder := model.RefundOrder{
		RefundNo:     refundNo,
		UserID:       userID,
		OrderID:      orderItem.OrderID,
		OrderItemID:  orderItem.ID,
		OrderNo:      orderItem.OrderNo,
		RefundAmount: orderItem.TotalFee,
		RefundReason: req.RefundReason,
		RefundType:   req.RefundType,
		Status:       0, // 退款中
		OriginStatus: orderItem.Status,
		Images:       req.Images,
		ApplyTime:    now.Format("2006-01-02"),
		CreatedAt:    &now,
		UpdatedAt:    &now,
		ProductID:    orderItem.ProductID,
		ProductName:  orderItem.ProductName,
		HeaderImg:    orderItem.HeaderImg,
		Quantity:     orderItem.Quantity,
		Price:        orderItem.CurrentPrice,
		StoreName:    orderItem.StoreName,
		StoreIcon:    orderItem.StoreLogo,
	}

	// 启动事务
	tx := r.DB(ctx).Begin()
	if err := tx.Error; err != nil {
		return err
	}

	// 创建退款单
	if err := tx.Create(&refundOrder).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 变更订单状态
	if err := tx.Model(&orderItem).Update("status", common.ORDER_STATUS_REFUNDED).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 校验退款数量是否合法
	if orderItem.Quantity <= 0 || orderItem.Quantity > orderItem.Quantity {
		tx.Rollback()
		return errors.New("退款数量不合法")
	}
	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// GetRefundList 获取退款单列表
func (r *refundOrderRepository) GetRefundList(ctx context.Context, userID string, status *uint8, page, pageSize int) (*model.RefundListResponse, error) {
	var total int64
	var refunds []model.RefundOrder

	// 构建查询条件
	query := r.DB(ctx).Model(&model.RefundOrder{}).Where("user_id = ?", userID)
	if status != nil {
		if *status == 1 {
			query = query.Where("status = ?", 0)
		} else if *status == 2 {
			query = query.Where("status = ?", 1)
		}
	}

	// 查询总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	// offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Find(&refunds).Error; err != nil {
		return nil, err
	}

	// 构建响应
	result := &model.RefundListResponse{
		Total: total,
		List:  make([]model.RefundOrderItem, 0, len(refunds)),
		Page:  page,
		Size:  pageSize,
	}

	// 填充列表
	for _, refund := range refunds {
		item := model.RefundOrderItem{
			RefundListItem: model.RefundListItem{
				ID:           refund.ID,
				RefundNo:     refund.RefundNo,
				OrderID:      refund.OrderID,
				OrderNo:      refund.OrderNo,
				RefundAmount: refund.RefundAmount,
				RefundType:   refund.RefundType,
				Status:       refund.Status,
				StatusText:   getRefundStatusText(refund.Status),
				ApplyTime:    refund.ApplyTime,
				CreatedAt:    refund.CreatedAt,
				StoreName:    refund.StoreName,
				StoreIcon:    refund.StoreIcon,
			},
			ProductName: refund.ProductName,
			HeaderImg:   refund.HeaderImg,
			Price:       refund.Price,
			Quantity:    refund.Quantity,
		}
		result.List = append(result.List, item)
	}

	return result, nil
}

// GetRefundDetail 获取退款单详情
func (r *refundOrderRepository) GetRefundDetail(ctx context.Context, userID string, refundID uint64) (*model.RefundDetailResponse, error) {
	// 查询退款单
	var refund model.RefundOrder
	if err := r.DB(ctx).Where("id = ? AND user_id = ?", refundID, userID).First(&refund).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("退款单不存在或无权限")
		}
		return nil, err
	}

	// 查询订单详情
	var order model.UserOrder
	if err := r.DB(ctx).Where("id = ?", refund.OrderID).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("订单不存在")
		}
		return nil, err
	}
	// 将applyTime转换为时间格式
	// 将2025-03-19T23:28:08+08:00转成2025-03-19 00:00:00
	applyTime, err := time.Parse("2006-01-02T15:04:05+08:00", refund.ApplyTime)
	if err != nil {
		return nil, err
	}
	refund.ApplyTime = applyTime.Format("2006-01-02")
	r.logger.Debug("refund.ApplyTime", refund.ApplyTime)
	// 查询订单项
	var orderItem model.UserOrderItem
	if err := r.DB(ctx).Where("id = ?", refund.OrderItemID).Find(&orderItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("订单子项不存在")
		}
		return nil, err
	}

	// 构建订单详情
	orderDetail := model.OrderDetails{
		OrderID:      order.ID,
		OrderNo:      order.OrderNo,
		OrderStatus:  orderItem.Status,
		OrderCreated: order.CreatedAt,
	}

	// 构建响应
	response := &model.RefundDetailResponse{
		Refund:      refund,
		OrderDetail: orderDetail,
	}

	return response, nil
}

// getRefundStatusText 获取退款状态文本
func getRefundStatusText(status uint8) string {
	switch status {
	case 0:
		return "退款进行中"
	case 1:
		return "已退款"
	case 3:
		return "已拒绝"
	default:
		return "未知状态"
	}
}

// getOrderStatusText 获取订单状态文本
// func getOrderStatusText(status int) string {
// 	switch status {
// 	case 0:
// 		return "待付款"
// 	case 1:
// 		return "待发货"
// 	case 2:
// 		return "待收货"
// 	case 3:
// 		return "待评价"
// 	case 4:
// 		return "已完成"
// 	case 5:
// 		return "交易关闭"
// 	default:
// 		return "未知状态"
// 	}
// }

// CancelRefund 撤销退款申请
func (r *refundOrderRepository) CancelRefund(ctx context.Context, userID string, refundID uint64) error {
	// 查询退款单是否存在且属于当前用户
	var refundOrder model.RefundOrder
	if err := r.DB(ctx).Where("id = ? AND user_id = ?", refundID, userID).First(&refundOrder).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("退款单不存在或无权限")
		}
		return err
	}

	// 校验退款单状态是否允许撤销（只有退款中状态才能撤销）
	if refundOrder.Status != 0 {
		return errors.New("当前退款单状态不允许撤销")
	}

	// 启动事务
	tx := r.DB(ctx).Begin()
	if err := tx.Error; err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 恢复订单项状态为原始状态
	if err := tx.Model(&model.UserOrderItem{}).
		Where("id = ?", refundOrder.OrderItemID).
		Update("status", refundOrder.OriginStatus).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除退款单
	if err := tx.Delete(&refundOrder).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// GetRefundDetailByID 获取退款单详情
func (r *refundOrderRepository) GetRefundDetailByID(ctx context.Context, refundID uint64) (*model.RefundDetailResponse, error) {
	// 查询退款单
	var refund model.RefundOrder
	if err := r.DB(ctx).Where("id = ?", refundID).First(&refund).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("退款单不存在或无权限")
		}
		return nil, err
	}

	// 查询订单详情
	var order model.UserOrder
	if err := r.DB(ctx).Where("id = ?", refund.OrderID).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("订单不存在")
		}
		return nil, err
	}
	// 将applyTime转换为时间格式
	// 将2025-03-19T23:28:08+08:00转成2025-03-19 00:00:00
	applyTime, err := time.Parse("2006-01-02T15:04:05+08:00", refund.ApplyTime)
	if err != nil {
		return nil, err
	}
	refund.ApplyTime = applyTime.Format("2006-01-02")
	r.logger.Debug("refund.ApplyTime", refund.ApplyTime)
	// 查询订单项
	var orderItem model.UserOrderItem
	if err := r.DB(ctx).Where("id = ?", refund.OrderItemID).Find(&orderItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("订单子项不存在")
		}
		return nil, err
	}

	// 构建订单详情
	orderDetail := model.OrderDetails{
		OrderID:      order.ID,
		OrderNo:      order.OrderNo,
		OrderStatus:  orderItem.Status,
		OrderCreated: order.CreatedAt,
	}

	// 构建响应
	response := &model.RefundDetailResponse{
		Refund:      refund,
		OrderDetail: orderDetail,
	}

	return response, nil
}

// DeleteRefund 删除退款记录
func (r *refundOrderRepository) DeleteRefund(ctx context.Context, userID string, refundID uint64) error {
	// 查询退款单是否存在并属于当前用户
	var refund model.RefundOrder
	if err := r.DB(ctx).Where("id = ? AND user_id = ?", refundID, userID).First(&refund).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("退款单不存在或无权限")
		}
		r.logger.Error("查询退款单失败", zap.Error(err))
		return err
	}

	// 删除退款单（软删除）
	if err := r.DB(ctx).Delete(&refund).Error; err != nil {
		r.logger.Error("删除退款单失败", zap.Error(err))
		return err
	}

	return nil
}
