package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"app/internal/cache"
	"app/internal/common"
	"app/internal/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserOrderRepository interface {
	CreateOrders(ctx *gin.Context, userID string, req model.CreateOrderRequest) error
	GetProductList(ctx context.Context, prosductIds []uint64) ([]model.ProductListItemDTO, error)
	GetOrderList(ctx context.Context, userID string, req model.OrderQueryRequest) (*model.OrderListResponse, error)
	GetUserOrderTotalAmount(ctx context.Context, userID string) (float64, float64, error)
	GetOrderProductDetails(ctx context.Context, orderID uint64, orderItemID uint64) ([]model.ProductListItemDTO, error)
	VerifyOrderOwnership(ctx context.Context, userID string, orderID uint64) (bool, error)
	UpdateOrderStatus(ctx context.Context, userID string, orderItemID uint64, status uint8) error
	GetOrderDetail(ctx context.Context, userID string, orderItemID uint64) (*model.OrderDetailResponse, error)
}

func NewUserOrderRepository(
	repository *Repository,
	userCartRepository UserCartRepository,
	userAssetRepository UserAssetRepository,
) UserOrderRepository {
	return &userOrderRepository{
		Repository:          repository,
		userCartRepository:  userCartRepository,
		userAssetRepository: userAssetRepository,
	}
}

type userOrderRepository struct {
	*Repository
	*cache.Cache
	userCartRepository  UserCartRepository
	userAssetRepository UserAssetRepository
}

func (r *userOrderRepository) GetCache(ctx context.Context, key string) *cache.Cache {
	if r.Cache != nil {
		return r.Cache
	}
	if key == "" {
		key = "user_order:"
	}
	r.Cache = cache.NewCache(ctx, key)
	return r.Cache
}

func (r *userOrderRepository) GetOrderItem(ctx context.Context, orderItemID uint64) (*model.UserOrderItem, error) {
	var orderItem model.UserOrderItem
	if err := r.DB(ctx).Where("id = ?", orderItemID).First(&orderItem).Error; err != nil {
		return nil, err
	}
	return &orderItem, nil
}

func (r *userOrderRepository) CreateOrders(ctx *gin.Context, userID string, req model.CreateOrderRequest) error {
	isDefault := 0
	if req.Address.IsDefault {
		isDefault = 1
	}
	// 计算总金额
	totalFee := float64(0)
	// 产品列表
	productIds := make([]uint64, 0, len(req.Items))
	for _, item := range req.Items {
		totalFee += item.CurrentPrice * float64(item.Quantity)
		if item.CouponID > 0 {
			couponPrice := item.CouponPrice
			totalFee -= couponPrice
		}
		if item.MemberDiscount > 0 {
			totalFee -= item.MemberDiscount
		}
		// 运费
		courierFeeMin := item.CourierFeeMin
		if courierFeeMin > 0 {
			totalFee += courierFeeMin
		}
		productIds = append(productIds, item.ProductID)
	}
	// 如果使用的余额账户，验证余额是否充足
	if req.PaymentMethod == 3 {
		balance, err := r.userAssetRepository.GetUserAsset(ctx, userID)
		if err != nil {
			return err
		}
		if balance.Balance < totalFee {
			return errors.New("余额不足")
		}
	}
	// 获取产品列表
	productList, err := r.GetProductList(ctx, productIds)
	if err != nil {
		r.logger.Debug("获取产品列表失败", "error", err)
		return err
	}
	productInfoMap := make(map[uint64]model.ProductListItemDTO, len(productList))
	for _, item := range productList {
		productInfoMap[uint64(item.ProductID)] = item
	}
	// 添加到user_order表
	tx := r.DB(ctx).Begin()
	var payTime *time.Time
	if req.Status == 1 {
		now := time.Now()
		payTime = &now
	}
	order := model.UserOrder{
		UserID:        userID,
		OrderNo:       common.GenerateOrderNo(),
		PaymentMethod: req.PaymentMethod,
		AddressID:     req.Address.ID,
		Name:          req.Address.Name,
		Phone:         req.Address.Phone,
		Province:      req.Address.Province,
		City:          req.Address.City,
		District:      req.Address.District,
		Detail:        req.Address.Detail,
		IsDefault:     uint8(isDefault),
		TotalFee:      totalFee,
	}
	if err := tx.Create(&order).Error; err != nil {
		r.logger.Debug("创建订单失败", "error", err, "order", order)
		tx.Rollback()
		return err
	}

	// Create a slice to hold all the orders
	orders := make([]model.UserOrderItem, 0, len(req.Items))
	// Convert each item to an order
	couponIds := make([]uint64, 0, len(req.Items))
	for _, item := range req.Items {
		itemTotalFee := item.CurrentPrice * float64(item.Quantity)
		couponPrice := item.CouponPrice
		memberDiscount := item.MemberDiscount
		courierFeeMin := item.CourierFeeMin
		totalFee := itemTotalFee + courierFeeMin - couponPrice - memberDiscount
		order := model.UserOrderItem{
			UserId:         userID,
			Category1Id:    productInfoMap[uint64(item.ProductID)].Category1ID,
			Category2Id:    productInfoMap[uint64(item.ProductID)].Category2ID,
			OrderID:        order.ID,
			OrderNo:        order.OrderNo,
			ProductID:      item.ProductID,
			Quantity:       item.Quantity,
			ProductName:    item.ProductName,
			HeaderImg:      item.Image,
			StoreID:        common.STORE_ID,
			StoreName:      common.STORE_NAME,
			StoreLogo:      common.STORE_LOGO,
			CurrentPrice:   item.CurrentPrice,
			CourierFeeMin:  item.CourierFeeMin,
			MemberDiscount: item.MemberDiscount,
			Note:           item.Note,
			CouponID:       item.CouponID,
			CouponPrice:    item.CouponPrice,
			TotalFee:       totalFee,
			Status:         uint8(req.Status),
			PayTime:        payTime,
		}
		orders = append(orders, order)
		if item.CouponID > 0 {
			couponIds = append(couponIds, item.CouponID)
		}
		if req.PaymentMethod == 3 {
			// 更新用户资产
			r.logger.Debug("更新用户资产", "totalFee", totalFee, "orderID", order.ID, "productName", item.ProductName)
			if err := r.userAssetRepository.UpdateUserAsset(ctx, tx, userID, common.BUSINESS_TYPE_ORDER,
				common.ACTION_TYPE_BUY, common.ASSET_TYPE_BALANCE, -float32(totalFee), int(order.ID), item.ProductName); err != nil {
				r.logger.Debug("更新用户资产失败", "error", err)
				tx.Rollback()
				return err
			}
		}
	}
	// Create all orders in a single transaction
	if err := tx.Create(&orders).Error; err != nil {
		r.logger.Debug("创建订单子项失败", "error", err, "orders", orders)
		tx.Rollback()
		return err
	}

	// 查询用户购买的商品是否在购物车中
	if r.userCartRepository != nil {
		cartList, _ := r.userCartRepository.GetCartList(ctx, userID, productIds)
		if len(cartList) > 0 {
			cartIDs := make([]uint64, 0, len(cartList))
			for _, item := range cartList {
				cartIDs = append(cartIDs, uint64(item.CartID))
			}
			if err := tx.Model(&model.UserCart{}).Where("id IN (?)", cartIDs).Update("status", 1).Error; err != nil {
				r.logger.Debug("更新购物车状态失败", "error", err, "cartIDs", cartIDs)
				tx.Rollback()
				return err
			}
		}
	}

	if req.Status == 1 {
		// 更新优惠券状态
		if len(couponIds) > 0 {
			if err := tx.Model(&model.UserCoupon{}).Where("user_id = ? AND coupon_id IN (?) AND status = ?", userID, couponIds, 0).Update("status", 1).Error; err != nil {
				r.logger.Debug("更新优惠券状态失败", "error", err)
				tx.Rollback()
				return err
			}
		}

		// 是否需要更新身份,普通会员:注册后未消费，高级会员：注册后消费记录，初级农场主：消费金额超过3万 高级农场主超过6万 资深农场主超过24万 合伙人后台指定
		// 获取用户购买消费（category1_id != 3）的总金额，以及认养鸵鸟（category1_id = 3）的总金额
		totalBuyAmount, totalAdoptAmount, err := r.GetUserOrderTotalAmount(ctx, userID)
		if err != nil {
			r.logger.Debug("获取用户消费总金额失败", "error", err)
			return err
		}
		updateRole := 0
		// 获取用户身份
		var role int
		if err := tx.Model(&model.Account{}).Where("user_id = ?", userID).Select("role").Scan(&role).Error; err != nil {
			r.logger.Debug("获取用户身份失败", "error", err)
			return err
		}
		// 如果用户消费总金额小于3万，则更新身份为高级会员
		if role == common.ROLE_NORMAL && totalBuyAmount < 30000 {
			updateRole = common.ROLE_VIP
		}
		// 如果用户消费总金额超过3万，则更新身份为初级农场主
		if role == common.ROLE_VIP && totalBuyAmount >= 30000 {
			updateRole = common.ROLE_FARMER_BEGIN
		}
		// 如果用户消费总金额超过5万，则更新身份为高级农场主
		if role == common.ROLE_FARMER_BEGIN && totalBuyAmount >= 60000 {
			updateRole = common.ROLE_FARMER_VIP
		}
		// 如果用户消费总金额超过10万，则更新身份为资深农场主
		if role == common.ROLE_FARMER_VIP && totalBuyAmount >= 240000 {
			updateRole = common.ROLE_FARMER_EXPERT
		}
		// 如果用户消费总金额超过20万，则更新身份为合伙人
		if role == common.ROLE_FARMER_EXPERT && totalAdoptAmount >= 200000 {
			updateRole = common.ROLE_FARMER_PARTNER
		}
		// 更新身份
		if updateRole > 0 {
			if err := tx.Model(&model.Account{}).Where("user_id = ?", userID).Update("role", updateRole).Error; err != nil {
				r.logger.Debug("更新身份失败", "error", err)
				tx.Rollback()
				return err
			}
			// 删除用户缓存
			r.GetCache(ctx, "account:").Del(userID)
		}
	}
	if err := tx.Commit().Error; err != nil {
		r.logger.Debug("提交事务失败", "error", err)
		tx.Rollback()
		return err
	}
	return nil
}

// 获取产品列表
func (r *userOrderRepository) GetProductList(ctx context.Context, productIds []uint64) ([]model.ProductListItemDTO, error) {
	if len(productIds) == 0 {
		return []model.ProductListItemDTO{}, nil
	}

	productList := []model.Product{}
	if err := r.DB(ctx).Model(&model.Product{}).Where("id IN ?", productIds).Find(&productList).Error; err != nil {
		r.logger.Debug("获取产品列表失败", "error", err)
		return nil, err
	}
	productListDTO := make([]model.ProductListItemDTO, 0, len(productList))
	for _, item := range productList {
		productListDTO = append(productListDTO, model.ProductListItemDTO{
			ProductID:            item.ID,
			ProductName:          item.ProductName,
			ProductCurrentPrice:  item.CurrentPrice,
			ProductOriginalPrice: item.OriginPrice,
			ProductUnit:          "",
			ProductSpec:          item.Specification,
			ProductSales:         item.Sales,
			ProductSpecification: item.Specification,
			ProductIsSpecial:     item.IsSpecial,
		})
	}

	return productListDTO, nil
}

// 获取用户购买消费（category1_id != 3）的总金额，以及认养鸵鸟（category1_id = 3）的总金额
func (r *userOrderRepository) GetUserOrderTotalAmount(ctx context.Context, userID string) (float64, float64, error) {
	if userID == "" {
		return 0, 0, nil
	}
	// 获取用户购买消费（category1_id != 3）的总金额
	var totalBuyAmount float64
	var totalBuyAmountNullable sql.NullFloat64
	if err := r.DB(ctx).Model(&model.UserOrderItem{}).Where("user_id = ? AND category1_id != ?", userID, 3).Select("COALESCE(SUM(total_fee), 0) as total_fee_all").Scan(&totalBuyAmountNullable).Error; err != nil {
		return 0, 0, err
	}
	totalBuyAmount = totalBuyAmountNullable.Float64

	// 获取用户认养鸵鸟（category1_id = 3）的总金额
	var totalAdoptAmount float64
	var totalAdoptAmountNullable sql.NullFloat64
	if err := r.DB(ctx).Model(&model.UserOrderItem{}).Where("user_id = ? AND category1_id = ?", userID, 3).Select("COALESCE(SUM(total_fee), 0) as total_fee_all").Scan(&totalAdoptAmountNullable).Error; err != nil {
		return 0, 0, err
	}
	totalAdoptAmount = totalAdoptAmountNullable.Float64
	return totalBuyAmount, totalAdoptAmount, nil
}

// 实现订单列表查询方法
func (r *userOrderRepository) GetOrderList(ctx context.Context, userID string, req model.OrderQueryRequest) (*model.OrderListResponse, error) {
	var orders []model.UserOrderItem
	var total int64

	// 计算分页偏移量
	offset := (req.Page - 1) * req.PageSize

	// 构建基础查询
	query := r.DB(ctx).Model(&model.UserOrderItem{}).Where("user_id = ?", userID)

	// 如果有状态过滤条件，添加到查询中
	if req.Status != nil && *req.Status != -1 {
		if *req.Status == 3 {
			query = query.Where("status IN (?)", []uint8{3, 4, 5, 6})
		} else {
			query = query.Where("status = ?", *req.Status)
		}
	}

	// 如果有订单号过滤条件，添加到查询中
	if req.OrderNo != nil && *req.OrderNo != "" {
		query = query.Where("order_no like ?", "%"+*req.OrderNo+"%")
	}

	// 查询总数
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("查询订单总数失败: " + err.Error())
		return nil, err
	}

	// 查询订单列表，附带地址信息和订单项
	if err := query.Order("created_at DESC").
		Limit(req.PageSize).
		Offset(offset).
		Find(&orders).Error; err != nil {
		r.logger.Error("查询订单列表失败: " + err.Error())
		return nil, err
	}

	// 遍历获取订单ID
	orderMap := make(map[uint64]bool)
	orderIDs := make([]uint64, 0, len(orders))
	for _, order := range orders {
		if _, ok := orderMap[order.OrderID]; !ok {
			orderMap[order.OrderID] = true
			orderIDs = append(orderIDs, order.OrderID)
		}
	}
	// 获取订单信息
	var orderList []model.UserOrder
	if err := r.DB(ctx).Where("id IN (?)", orderIDs).Find(&orderList).Error; err != nil {
		return nil, err
	}
	// 转成map格式
	orderInfoMap := make(map[uint64]model.UserOrder)
	for _, order := range orderList {
		orderInfoMap[order.ID] = order
	}

	// 获取店铺信息
	// 获取第一个订单项的店铺信息
	var storeName string = "鸵小妥"
	var storeIcon string = "https://txtimages.oss-cn-beijing.aliyuncs.com/icons/dianpu.png" // 默认店铺图标
	if len(orders) > 0 {
		storeName = orders[0].StoreName
		storeIcon = orders[0].StoreLogo
	}

	// 转换为 DTO
	orderDTOs := make([]model.OrderListDTO, 0, len(orders))
	for _, item := range orders {
		// 构建地址字符串
		address := ""
		orderInfo, ok := orderInfoMap[item.OrderID]
		if !ok {
			continue
		}
		if orderInfo.OrderAddress != nil {
			address = fmt.Sprintf("%s %s %s %s",
				orderInfo.OrderAddress.Province,
				orderInfo.OrderAddress.City,
				orderInfo.OrderAddress.District,
				orderInfo.OrderAddress.Street)
		}
		// 是否超过25分钟
		now := time.Now()
		if item.Status == 0 && item.CreatedAt.Add(25*time.Minute).Before(now) {
			item.Status = 6
		}
		// 构建状态文本
		statusText := r.GetOrderStatusText(int(item.Status))

		// 构建支付时间
		var payTime *string
		if item.PayTime != nil && !item.PayTime.IsZero() {
			timeStr := item.PayTime.Format("2006-01-02 15:04:05")
			payTime = &timeStr
		}

		// 构建订单商品
		product := model.OrderProductDTO{
			ItemID:      item.ID,
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			Price:       item.CurrentPrice,
			Image:       item.HeaderImg,
		}

		orderDTOs = append(orderDTOs, model.OrderListDTO{
			ID:            orderInfo.ID,
			OrderNo:       orderInfo.OrderNo,
			UserID:        orderInfo.UserID,
			TotalFee:      orderInfo.TotalFee,
			Status:        int(item.Status),
			StatusText:    statusText,
			PaymentMethod: orderInfo.PaymentMethod,
			PayTime:       payTime,
			CreatedAt:     *orderInfo.CreatedAt,
			Address:       address,
			Product:       product,
			StoreName:     storeName,
			StoreIcon:     storeIcon,
		})
	}

	return &model.OrderListResponse{
		Total:    total,
		Orders:   orderDTOs,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// 获取订单状态文本
func (r *userOrderRepository) GetOrderStatusText(status int) string {
	switch status {
	case 0:
		return "待付款"
	case 1:
		return "待发货"
	case 2:
		return "待收货"
	case 3:
		return "待评价"
	case 4:
		return "已完成"
	case 5:
		return "交易关闭"
	case 6:
		return "已失效"
	case 7:
		return "退款进行中"
	default:
		return "未知状态"
	}
}

// GetOrderProductDetails 获取订单商品详情
func (r *userOrderRepository) GetOrderProductDetails(ctx context.Context, orderID uint64, orderItemID uint64) ([]model.ProductListItemDTO, error) {
	// 查询订单是否存在
	var order model.UserOrder
	if err := r.DB(ctx).Where("id = ?", orderID).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("订单不存在")
		}
		r.logger.Error("查询订单失败", zap.Error(err))
		return nil, err
	}

	// 查询订单商品项
	var orderItems []model.UserOrderItem
	if err := r.DB(ctx).Where("order_id = ?", orderID).Find(&orderItems).Error; err != nil {
		r.logger.Error("查询订单商品项失败", zap.Error(err))
		return nil, err
	}

	// 构建商品详情列表
	products := make([]model.ProductListItemDTO, 0, len(orderItems))
	for _, item := range orderItems {
		// 查询商品详细信息
		var product model.Product
		if err := r.DB(ctx).Where("id = ?", item.ProductID).First(&product).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				r.logger.Error("查询商品详情失败", zap.Error(err))
				return nil, err
			}
			// 如果商品不存在，继续使用订单项中的信息
		}

		// 构建商品详情
		productDTO := model.ProductListItemDTO{
			ItemID:               uint(item.ID),
			ProductID:            uint(item.ProductID),
			ProductName:          item.ProductName,
			ProductCurrentPrice:  item.CurrentPrice,
			ProductOriginalPrice: product.OriginPrice,
			ProductUnit:          "", // You may need to add this field to the Product model
			ProductSpec:          product.Specification,
			HeaderImg:            item.HeaderImg,
			ProductImages:        strings.Split(product.BannerImg, ","),
			ProductIsSpecial:     product.IsSpecial,
			CourierFeeMin:        product.CourierFeeMin,
			CourierFeeMax:        product.CourierFeeMax,
			MemberDiscount:       product.MemberDiscount,
			ProductContent:       product.Content,
			Category1ID:          item.Category1Id,
			Category2ID:          item.Category2Id,
			// 优惠券信息
			CouponID:    int(item.CouponID),
			CouponPrice: item.CouponPrice,
			// 订单特有信息
			ProductQuantity: item.Quantity,
		}
		if orderItemID == item.ID {
			products = []model.ProductListItemDTO{
				productDTO,
			}
			break
		}
		products = append(products, productDTO)
	}

	return products, nil
}

// VerifyOrderOwnership 验证订单归属权
func (r *userOrderRepository) VerifyOrderOwnership(ctx context.Context, userID string, orderItemID uint64) (bool, error) {
	var count int64
	if err := r.DB(ctx).Model(&model.UserOrderItem{}).Where("id = ? AND user_id = ?", orderItemID, userID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// UpdateOrderStatus 更新订单状态
func (r *userOrderRepository) UpdateOrderStatus(ctx context.Context, userID string, orderItemID uint64, status uint8) error {
	// 验证订单是否存在且属于该用户
	var orderItem model.UserOrderItem
	if err := r.DB(ctx).Where("id = ? AND user_id = ?", orderItemID, userID).First(&orderItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("订单不存在或无权限")
		}
		r.logger.Error("查询订单失败", zap.Error(err))
		return err
	}

	// 验证状态转换是否合法
	if !isValidStatusTransition(orderItem.Status, status) {
		return errors.New("非法的状态转换")
	}

	// 更新订单状态
	updateFields := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}

	// 如果已发货，则更新发货时间
	if status == 2 && orderItem.Status == 1 {
		updateFields["shipped_at"] = time.Now()
	}

	// 如果已收货，则更新收货时间
	if status == 3 && orderItem.Status == 2 {
		updateFields["completed_at"] = time.Now()
		// 给用户添加积分，要求消费金额>10元，积分=（消费金额/10）取整
		if orderItem.TotalFee > 10 {
			integral := int(orderItem.TotalFee / 10)
			if err := r.DB(ctx).Model(&model.UserAsset{}).Where("user_id = ?", userID).Update("points", gorm.Expr("points + ?", integral)).Error; err != nil {
				r.logger.Debug("更新积分失败", "error", err)
				return err
			}
		}
	}

	// 如果状态变更为已支付，则更新支付时间
	if status == 1 && orderItem.Status == 0 {
		updateFields["pay_time"] = time.Now()
	}

	if err := r.DB(ctx).Model(&model.UserOrderItem{}).Where("id = ?", orderItemID).Updates(updateFields).Error; err != nil {
		r.logger.Error("更新订单状态失败", zap.Error(err))
		return err
	}

	return nil
}

// isValidStatusTransition 验证状态转换是否合法
func isValidStatusTransition(currentStatus, newStatus uint8) bool {
	// 状态流转规则：
	// 0(待支付) -> 1(待发货)
	// 1(待发货) -> 2(待收货)
	// 2(待收货) -> 3(待评价)
	// 3(待评价) -> 4(已完成)
	// 任何状态 -> 5(已取消)，假设5表示已取消
	// 待发货 -> 6(已退款)

	if newStatus == 5 { // 可以从任何状态取消
		return true
	}

	// 常规状态流转
	validTransitions := map[uint8][]uint8{
		0: {1},    // 待支付 -> 待发货
		1: {2, 6}, // 待发货 -> 待收货, 已退款
		2: {3},    // 待收货 -> 待评价
		3: {4},    // 待评价 -> 已完成
	}

	allowedNextStatus, exists := validTransitions[currentStatus]
	if !exists {
		return false
	}

	for _, status := range allowedNextStatus {
		if status == newStatus {
			return true
		}
	}

	return false
}

// GetOrderDetail 获取订单详情
func (r *userOrderRepository) GetOrderDetail(ctx context.Context, userID string, orderItemID uint64) (*model.OrderDetailResponse, error) {
	// 根据订单ID和用户ID查询订单
	var orderItem model.UserOrderItem
	if err := r.DB(ctx).Where("id = ? AND user_id = ?", orderItemID, userID).First(&orderItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("订单不存在或无权限")
		}
		r.logger.Error("查询订单失败", zap.Error(err))
		return nil, err
	}
	// 获取商品信息
	products, err := r.GetOrderProductDetails(ctx, orderItem.OrderID, orderItemID)
	if err != nil {
		return nil, err
	}

	// 查询订单
	var order model.UserOrder
	if err := r.DB(ctx).Where("id = ?", orderItem.OrderID).First(&order).Error; err != nil {
		r.logger.Error("查询订单失败", zap.Error(err))
		return nil, err
	}

	// 查询地址信息
	var address model.UserAddress
	if err := r.DB(ctx).Where("id = ?", order.AddressID).First(&address).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Error("查询地址信息失败", zap.Error(err))
			return nil, err
		}
		// 如果地址不存在，使用一个空地址
		address = model.UserAddress{}
	}

	// 构建地址信息
	addressInfo := &model.AddressInfo{
		ReceiverName:  address.Name,
		ReceiverPhone: address.Phone,
		Province:      address.Province,
		City:          address.City,
		District:      address.District,
		DetailAddress: address.Street,
	}

	totalPrice := orderItem.CurrentPrice * float64(orderItem.Quantity)
	couponDiscount := orderItem.CouponPrice
	memberDiscount := orderItem.MemberDiscount
	totalShippingFee := orderItem.CourierFeeMin

	storeID := orderItem.StoreID
	storeName := orderItem.StoreName
	storeLogo := orderItem.StoreLogo

	// 构建订单详情响应
	var payTime string
	if orderItem.PayTime != nil {
		payTime = orderItem.PayTime.Format("2006-01-02 15:04:05")
	}
	var shippedAt string
	if orderItem.ShippedAt != nil {
		shippedAt = orderItem.ShippedAt.Format("2006-01-02 15:04:05")
	}
	var completedAt string
	if orderItem.CompletedAt != nil {
		completedAt = orderItem.CompletedAt.Format("2006-01-02 15:04:05")
	}
	response := &model.OrderDetailResponse{
		OrderID:        order.ID,
		OrderNo:        order.OrderNo,
		UserID:         order.UserID,
		OrderAmount:    order.TotalFee,
		Status:         int(orderItem.Status),
		StatusText:     r.GetOrderStatusText(int(orderItem.Status)),
		PayMethod:      order.PaymentMethod,
		PayMethodText:  getPayMethodText(order.PaymentMethod),
		PayTime:        payTime,
		CreatedAt:      orderItem.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      orderItem.UpdatedAt.Format("2006-01-02 15:04:05"),
		ShippedAt:      shippedAt,
		CompletedAt:    completedAt,
		Address:        addressInfo,
		Products:       products,
		TotalPrice:     totalPrice,
		ShippingFee:    totalShippingFee,
		CouponDiscount: couponDiscount,
		MemberDiscount: memberDiscount,
		ActualAmount:   totalPrice - totalShippingFee - couponDiscount,
		StoreID:        storeID,
		StoreName:      storeName,
		StoreLogo:      storeLogo,
	}

	return response, nil
}

// getPayMethodText 获取支付方式文本
func getPayMethodText(payMethod int) string {
	switch payMethod {
	case 0:
		return "未支付"
	case 1:
		return "余额支付"
	case 2:
		return "微信支付"
	case 3:
		return "支付宝支付"
	default:
		return "未知方式"
	}
}
