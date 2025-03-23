package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"app/internal/model"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProductReviewRepository interface {
	CreateReview(ctx context.Context, userID string, req model.CreateReviewRequest) error
	GetProductReviews(ctx context.Context, productID uint64, page, pageSize int) (*model.ReviewListResponse, error)
	GetUserReviews(ctx context.Context, userID string, page, pageSize int) (*model.ReviewListResponse, error)
	GetUserReviewsByTab(ctx context.Context, userID string, tab model.ReviewTabType) (*model.UserReviewListResponse, error)
	DeleteReview(ctx context.Context, userID string, reviewID uint64) error
	GetReviewDetail(ctx context.Context, reviewID uint64) (*model.ReviewDetailResponse, error)
	IncrementReviewCounter(ctx context.Context, reviewID uint64, counterType uint8) error
}

type productReviewRepository struct {
	*Repository
}

func NewProductReviewRepository(repository *Repository) ProductReviewRepository {
	return &productReviewRepository{
		Repository: repository,
	}
}

// CreateReview 创建商品评价
func (r *productReviewRepository) CreateReview(ctx context.Context, userID string, req model.CreateReviewRequest) error {
	// 校验商品是否在该订单中
	var orderItem model.UserOrderItem
	if err := r.DB(ctx).Where("id = ? AND user_id = ?", req.OrderItemID, userID).First(&orderItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("该订单不存在")
		}
		r.logger.Error("查询订单项失败", zap.Error(err))
		return err
	}

	// 校验该商品是否已评价
	var existingReview model.ProductReview
	if err := r.DB(ctx).Where("order_id = ? AND product_id = ? AND user_id = ?",
		req.OrderItemID, req.ProductID, userID).First(&existingReview).Error; err == nil {
		return errors.New("您已经评价过该商品")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		r.logger.Error("查询评价失败", zap.Error(err))
		return err
	}

	// 处理评价图片
	imagesStr := strings.Join(req.Images, ",")

	// 查询用户是否是回头客
	var productReview model.ProductReview
	isReturn := 0
	if err := r.DB(ctx).Where("user_id = ? AND product_id = ?", userID, req.ProductID).First(&productReview).Error; err == nil && productReview.ID > 0 {
		isReturn = 1
	}

	// 创建评价记录
	now := time.Now()
	review := model.ProductReview{
		UserID:          userID,
		OrderID:         orderItem.OrderID,
		OrderNo:         orderItem.OrderNo,
		OrderItemID:     orderItem.ID,
		ProductID:       req.ProductID,
		Rating:          req.Rating,
		FreshnessRating: req.FreshnessRating,
		PackagingRating: req.PackagingRating,
		DeliveryRating:  req.DeliveryRating,
		ServiceRating:   req.ServiceRating,
		Content:         req.Content,
		Images:          imagesStr,
		IsAnonymous:     req.IsAnonymous,
		Status:          1, // 默认为已通过
		CreatedAt:       &now,
		UpdatedAt:       &now,
		IsReturn:        int8(isReturn),
		ViewNums:        0,
		EvaluateNums:    0,
		PraiseNums:      0,
	}

	// 保存到数据库
	if err := r.DB(ctx).Create(&review).Error; err != nil {
		r.logger.Error("创建评价失败", zap.Error(err))
		return err
	}

	// 更新订单项状态
	if err := r.DB(ctx).Model(&model.UserOrderItem{}).
		Where("id = ?", req.OrderItemID).
		Update("status", 4).Error; err != nil {
		r.logger.Error("更新订单项状态失败", zap.Error(err))
		return err
	}

	return nil
}

// GetProductReviews 获取商品评价列表
func (r *productReviewRepository) GetProductReviews(ctx context.Context, productID uint64, page, pageSize int) (*model.ReviewListResponse, error) {
	var total int64
	var reviews []model.ProductReview

	// offset := (page - 1) * pageSize

	// 查询总数
	if err := r.DB(ctx).Model(&model.ProductReview{}).
		Where("product_id = ?", productID).
		Count(&total).Error; err != nil {
		r.logger.Error("查询评价总数失败", zap.Error(err))
		return nil, err
	}

	// 查询评价列表
	if err := r.DB(ctx).Where("product_id = ?", productID).
		Order("created_at DESC").
		// Offset(offset).
		// Limit(pageSize).
		Find(&reviews).Error; err != nil {
		r.logger.Error("查询评价列表失败", zap.Error(err))
		return nil, err
	}

	// 转换为DTO
	reviewDTOs, err := r.convertToReviewDTOs(ctx, reviews)
	if err != nil {
		return nil, err
	}

	return &model.ReviewListResponse{
		Total:   total,
		Reviews: reviewDTOs,
	}, nil
}

// GetUserReviews 获取用户评价列表
func (r *productReviewRepository) GetUserReviews(ctx context.Context, userID string, page, pageSize int) (*model.ReviewListResponse, error) {
	var total int64
	var reviews []model.ProductReview

	offset := (page - 1) * pageSize

	// 查询总数
	if err := r.DB(ctx).Model(&model.ProductReview{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		r.logger.Error("查询用户评价总数失败", zap.Error(err))
		return nil, err
	}

	// 查询评价列表
	if err := r.DB(ctx).Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&reviews).Error; err != nil {
		r.logger.Error("查询用户评价列表失败", zap.Error(err))
		return nil, err
	}

	// 转换为DTO
	reviewDTOs, err := r.convertToReviewDTOs(ctx, reviews)
	if err != nil {
		return nil, err
	}

	return &model.ReviewListResponse{
		Total:   total,
		Reviews: reviewDTOs,
	}, nil
}

// GetUserReviewsByTab 根据标签获取用户评价列表
func (r *productReviewRepository) GetUserReviewsByTab(ctx context.Context, userID string, tab model.ReviewTabType) (*model.UserReviewListResponse, error) {
	// 根据标签获取评价列表
	var reviews []interface{}
	if tab == model.ReviewTabAll {
		// 已评价的订单项
		productReviews := []model.ProductReview{}
		if err := r.DB(ctx).Model(&model.ProductReview{}).
			Where("user_id = ?", userID).Find(&productReviews).Error; err != nil {
			return nil, err
		}
		// 获取订单ID数组
		orderItemIDs := []uint64{}
		for _, productReview := range productReviews {
			orderItemIDs = append(orderItemIDs, productReview.OrderItemID)
		}
		// 获取订单项
		orderItems := []model.UserOrderItem{}
		if err := r.DB(ctx).Where("id IN (?)", orderItemIDs).Find(&orderItems).Error; err != nil {
			return nil, err
		}
		// 获取map
		orderItemMap := make(map[uint64]model.UserOrderItem)
		for _, orderItem := range orderItems {
			orderItemMap[orderItem.ID] = orderItem
		}
		// 待评价的订单项
		var pendingOrderItems []model.UserOrderItem
		if err := r.DB(ctx).Where("user_id = ? AND status = 3", userID).Find(&pendingOrderItems).Error; err != nil {
			return nil, err
		}
		// 合并到reviews
		for _, productReview := range productReviews {
			reviews = append(reviews, model.ReviewListDTO{
				ID:            productReview.ID,
				OrderID:       productReview.OrderID,
				OrderNo:       productReview.OrderNo,
				OrderItemID:   productReview.OrderItemID,
				ProductID:     productReview.ProductID,
				ProductName:   orderItemMap[productReview.OrderItemID].ProductName,
				ProductImage:  orderItemMap[productReview.OrderItemID].HeaderImg,
				Price:         orderItemMap[productReview.OrderItemID].CurrentPrice,
				Quantity:      orderItemMap[productReview.OrderItemID].Quantity,
				OrderTime:     orderItemMap[productReview.OrderItemID].CreatedAt.Format("2006-01-02"),
				ReviewStatus:  1,
				StatusText:    getReviewStatusText(1),
				StoreName:     orderItemMap[productReview.OrderItemID].StoreName,
				StoreIcon:     orderItemMap[productReview.OrderItemID].StoreLogo,
				ReviewContent: productReview.Content,
				ReviewImages:  strings.Split(productReview.Images, ","),
				IsAnonymous:   productReview.IsAnonymous,
			})
		}
		for _, orderItem := range pendingOrderItems {
			reviews = append(reviews, model.ReviewListDTO{
				ID:           orderItem.ID,
				OrderID:      orderItem.OrderID,
				OrderNo:      orderItem.OrderNo,
				OrderItemID:  orderItem.ID,
				ProductID:    orderItem.ProductID,
				ProductName:  orderItem.ProductName,
				ProductImage: orderItem.HeaderImg,
				Price:        orderItem.CurrentPrice,
				Quantity:     orderItem.Quantity,
				OrderTime:    orderItem.CreatedAt.Format("2006-01-02"),
				ReviewStatus: 0,
				StatusText:   getReviewStatusText(0),
				StoreName:    orderItem.StoreName,
				StoreIcon:    orderItem.StoreLogo,
			})
		}
	} else if tab == model.ReviewTabPending {
		// 待评价的订单项
		var pendingOrderItems []model.UserOrderItem
		if err := r.DB(ctx).Where("user_id = ? AND status = 3", userID).Find(&pendingOrderItems).Error; err != nil {
			return nil, err
		}
		for _, orderItem := range pendingOrderItems {
			reviews = append(reviews, model.ReviewListDTO{
				ID:           orderItem.ID,
				OrderID:      orderItem.OrderID,
				OrderNo:      orderItem.OrderNo,
				OrderItemID:  orderItem.ID,
				ProductID:    orderItem.ProductID,
				ProductName:  orderItem.ProductName,
				ProductImage: orderItem.HeaderImg,
				Price:        orderItem.CurrentPrice,
				Quantity:     orderItem.Quantity,
				OrderTime:    orderItem.CreatedAt.Format("2006-01-02"),
				ReviewStatus: 0,
				StatusText:   getReviewStatusText(0),
				StoreName:    orderItem.StoreName,
				StoreIcon:    orderItem.StoreLogo,
			})
		}
	} else if tab == model.ReviewTabCompleted {
		// 已评价的订单项
		productReviews := []model.ProductReview{}
		if err := r.DB(ctx).Model(&model.ProductReview{}).
			Where("user_id = ?", userID).Find(&productReviews).Error; err != nil {
			return nil, err
		}
		// 获取订单ID数组
		orderItemIDs := []uint64{}
		for _, productReview := range productReviews {
			orderItemIDs = append(orderItemIDs, productReview.OrderItemID)
		}
		// 获取订单项
		orderItems := []model.UserOrderItem{}
		if err := r.DB(ctx).Where("id IN (?)", orderItemIDs).Find(&orderItems).Error; err != nil {
			return nil, err
		}
		// 获取map
		orderItemMap := make(map[uint64]model.UserOrderItem)
		for _, orderItem := range orderItems {
			orderItemMap[orderItem.ID] = orderItem
		}
		for _, productReview := range productReviews {
			reviews = append(reviews, model.ReviewDoneListDTO{
				ID:            productReview.ID,
				Status:        productReview.Status,
				CreatedAt:     productReview.CreatedAt.Format("2006-01-02"),
				ProductID:     productReview.ProductID,
				ProductName:   orderItemMap[productReview.OrderItemID].ProductName,
				ProductImage:  orderItemMap[productReview.OrderItemID].HeaderImg,
				Price:         orderItemMap[productReview.OrderItemID].CurrentPrice,
				Quantity:      orderItemMap[productReview.OrderItemID].Quantity,
				ReviewContent: productReview.Content,
				ReviewImages:  strings.Split(productReview.Images, ","),
				ReviewStatus:  0,
				StatusText:    getReviewStatusText(0),
				StoreName:     orderItemMap[productReview.OrderItemID].StoreName,
				StoreIcon:     orderItemMap[productReview.OrderItemID].StoreLogo,
				IsAnonymous:   productReview.IsAnonymous,
			})
		}
	}
	return &model.UserReviewListResponse{
		Tab:  tab,
		List: reviews,
	}, nil
}

// 辅助函数：获取订单状态文本
func getOrderStatusText(status int) string {
	switch status {
	case 0:
		return "待付款"
	case 1:
		return "待发货"
	case 2:
		return "待收货"
	case 3:
		return "已完成"
	case 4:
		return "已取消"
	case 5:
		return "已退款"
	default:
		return "未知状态"
	}
}

func getReviewStatusText(status int) string {
	switch status {
	case 0:
		return "待评价"
	case 1:
		return "已评价"
	default:
		return "未知状态"
	}
}

// convertToReviewDTOs 将评价模型转为DTO
func (r *productReviewRepository) convertToReviewDTOs(ctx context.Context, reviews []model.ProductReview) ([]model.ReviewListDTO, error) {
	reviewDTOs := make([]model.ReviewListDTO, 0, len(reviews))

	for _, review := range reviews {
		// 查询商品信息
		var product model.Product
		if err := r.DB(ctx).Where("id = ?", review.ProductID).First(&product).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				r.logger.Error("查询商品信息失败", zap.Error(err))
				return nil, err
			}
		}

		// 查询用户信息
		var userOrderItem model.UserOrderItem
		if err := r.DB(ctx).Where("id = ?", review.OrderItemID).First(&userOrderItem).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				r.logger.Error("查询用户信息失败", zap.Error(err))
				return nil, err
			}
		}
		// 查询用户信息
		var user model.Account
		if err := r.DB(ctx).Where("user_id = ?", review.UserID).First(&user).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				r.logger.Error("查询用户信息失败", zap.Error(err))
				return nil, err
			}
		}

		// 构建DTO
		dto := model.ReviewListDTO{
			ID:            review.ID,
			OrderID:       review.OrderID,
			OrderNo:       review.OrderNo,
			OrderItemID:   review.OrderItemID,
			StoreName:     product.StoreName,
			StoreIcon:     product.StoreIcon,
			StatusText:    getReviewStatusText(int(review.Status)),
			ProductName:   product.ProductName,
			ProductImage:  product.HeaderImg,
			Price:         product.CurrentPrice,
			Quantity:      int(userOrderItem.Quantity),
			ReviewStatus:  review.Status,
			OrderTime:     review.CreatedAt.Format("2006-01-02"),
			ReviewContent: review.Content,
			ReviewImages:  strings.Split(review.Images, ","),
			ViewNums:      review.ViewNums,
			EvaluateNums:  review.EvaluateNums,
			PraiseNums:    review.PraiseNums,
			IsAnonymous:   review.IsAnonymous,
		}
		if review.IsAnonymous {
			dto.UserName = "匿名用户"
			dto.UserAvatar = "/images/icons/default.png"
		} else {
			dto.UserName = user.Nickname
			dto.UserAvatar = user.Avatar
		}

		reviewDTOs = append(reviewDTOs, dto)
	}

	return reviewDTOs, nil
}

// DeleteReview 删除评价
func (r *productReviewRepository) DeleteReview(ctx context.Context, userID string, reviewID uint64) error {
	// 首先查询评价是否存在且属于当前用户
	var review model.ProductReview
	if err := r.DB(ctx).Where("id = ?", reviewID).First(&review).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("评价不存在")
		}
		r.logger.Error("查询评价失败", zap.Error(err))
		return err
	}

	// 验证评价是否属于当前用户
	if review.UserID != userID {
		return errors.New("无权限删除该评价")
	}

	// 执行软删除（如果使用gorm的软删除功能）
	if err := r.DB(ctx).Delete(&review).Error; err != nil {
		r.logger.Error("删除评价失败", zap.Error(err))
		return err
	}

	return nil
}

// GetReviewDetail 获取评价详情
func (r *productReviewRepository) GetReviewDetail(ctx context.Context, reviewID uint64) (*model.ReviewDetailResponse, error) {
	var review model.ProductReview
	// 查询评价信息
	if err := r.DB(ctx).Where("id = ?", reviewID).First(&review).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("评价不存在")
		}
		r.logger.Error("查询评价详情失败", zap.Error(err))
		return nil, err
	}

	// 查询商品信息
	var product model.Product
	if err := r.DB(ctx).Where("id = ?", review.ProductID).First(&product).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Error("查询商品信息失败", zap.Error(err))
			return nil, err
		}
	}

	// 查询订单项信息
	var orderItem model.UserOrderItem
	if err := r.DB(ctx).Where("id = ?", review.OrderItemID).First(&orderItem).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Error("查询订单项信息失败", zap.Error(err))
			return nil, err
		}
	}

	// 查询评论和回复列表
	evaluateList := []model.ProductEvaluate{}
	if err := r.DB(ctx).Where("review_id = ?", review.ID).Find(&evaluateList).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Error("查询评论和回复列表失败", zap.Error(err))
			return nil, err
		}
	}
	evaluateListDTO := []model.ProductEvaluateList{}
	for _, evaluate := range evaluateList {
		evaluateListDTO = append(evaluateListDTO, model.ProductEvaluateList{
			ProductEvaluate: evaluate,
			CreatedAtStr:    evaluate.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	// 将evaluateList中的created_at转换为标准时间字符串
	// 按照Map组织parent_id和子级
	evaluateMap := make(map[int8][]model.ProductEvaluateList)
	for _, evaluate := range evaluateListDTO {
		if evaluate.ParentID > 0 {
			evaluateMap[evaluate.ParentID] = append(evaluateMap[evaluate.ParentID], evaluate)
		}
	}
	var newEvaluateListDTO []model.ProductEvaluateList
	for _, evaluate := range evaluateListDTO {
		if evaluate.ParentID == 0 {
			evaluate.Children = evaluateMap[int8(evaluate.ID)]
			newEvaluateListDTO = append(newEvaluateListDTO, evaluate)
		}
	}

	// 查询用户信息
	var user model.Account
	if err := r.DB(ctx).Where("user_id = ?", review.UserID).First(&user).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Error("查询用户信息失败", zap.Error(err))
			return nil, err
		}
	}

	// 处理评价图片
	var images []string
	if review.Images != "" {
		images = strings.Split(review.Images, ",")
	}

	// 组装评价详情响应
	response := &model.ReviewDetailResponse{
		ID:              review.ID,
		UserID:          review.UserID,
		UserName:        user.Nickname,
		UserAvatar:      user.Avatar,
		OrderID:         review.OrderID,
		OrderNo:         review.OrderNo,
		OrderItemID:     review.OrderItemID,
		ProductID:       review.ProductID,
		ProductName:     product.ProductName,
		ProductImage:    product.HeaderImg,
		ProductSpec:     product.Specification,
		Price:           product.CurrentPrice,
		Quantity:        int(orderItem.Quantity),
		StoreName:       product.StoreName,
		StoreIcon:       product.StoreIcon,
		Rating:          review.Rating,
		FreshnessRating: review.FreshnessRating,
		PackagingRating: review.PackagingRating,
		DeliveryRating:  review.DeliveryRating,
		ServiceRating:   review.ServiceRating,
		Content:         review.Content,
		Images:          images,
		IsAnonymous:     review.IsAnonymous,
		Status:          review.Status,
		StatusText:      getReviewStatusText(int(review.Status)),
		CreatedAt:       review.CreatedAt.Format("2006-01-02"),
		CommentNums:     review.EvaluateNums,
		EvaluateList:    newEvaluateListDTO,
	}

	// 如果是匿名评价，则隐藏用户信息
	if review.IsAnonymous {
		response.UserName = "匿名用户"
		response.UserAvatar = ""
	} else {
		response.UserName = user.Nickname
		response.UserAvatar = user.Avatar
	}

	return response, nil
}

// IncrementReviewCounter 更新评价计数器
func (r *productReviewRepository) IncrementReviewCounter(ctx context.Context, reviewID uint64, counterType uint8) error {
	// 首先查询评价是否存在
	var review model.ProductReview
	if err := r.DB(ctx).Where("id = ?", reviewID).First(&review).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("评价不存在")
		}
		r.logger.Error("查询评价失败", zap.Error(err))
		return err
	}

	// 根据计数器类型更新对应字段
	var field string
	var logMsg string

	switch counterType {
	case 0: // 浏览人数
		field = "view_nums"
		logMsg = "浏览人数"
	case 1: // 评价人数
		field = "evaluate_nums"
		logMsg = "评价人数"
	case 2: // 点赞人数
		field = "praise_nums"
		logMsg = "点赞人数"
	case 3: // 取消点赞
		field = "praise_nums"
		logMsg = "取消点赞"
	default:
		return errors.New("无效的计数器类型")
	}

	// 更新计数器（+1）
	if counterType == 3 {
		if err := r.DB(ctx).Model(&model.ProductReview{}).
			Where("id = ?", reviewID).
			Update(field, gorm.Expr(field+" - ?", 1)).Error; err != nil {
			r.logger.Error("更新"+logMsg+"失败", zap.Error(err))
			return err
		}
	} else {
		if err := r.DB(ctx).Model(&model.ProductReview{}).
			Where("id = ?", reviewID).
			Update(field, gorm.Expr(field+" + ?", 1)).Error; err != nil {
			r.logger.Error("更新"+logMsg+"失败", zap.Error(err))
			return err
		}
	}

	return nil
}
