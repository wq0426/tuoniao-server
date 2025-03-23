package service

import (
	"app/internal/model"
	"app/internal/repository"
	"errors"

	"github.com/gin-gonic/gin"
)

type ProductReviewService interface {
	CreateReview(ctx *gin.Context, req model.CreateReviewRequest) error
	GetProductReviews(ctx *gin.Context, productID uint64, page, pageSize int) (*model.ReviewListResponse, error)
	GetUserReviews(ctx *gin.Context, page, pageSize int) (*model.ReviewListResponse, error)
	GetUserReviewsByTab(ctx *gin.Context, req model.UserReviewListRequest) (*model.UserReviewListResponse, error)
	DeleteReview(ctx *gin.Context, reviewID uint64) error
	GetReviewDetail(ctx *gin.Context, reviewID uint64) (*model.ReviewDetailResponse, error)
	IncrementReviewCounter(ctx *gin.Context, req model.UpdateReviewCounterRequest) error
}

type productReviewService struct {
	*Service
	productReviewRepository repository.ProductReviewRepository
}

func NewProductReviewService(
	service *Service,
	productReviewRepository repository.ProductReviewRepository,
) ProductReviewService {
	return &productReviewService{
		Service:                 service,
		productReviewRepository: productReviewRepository,
	}
}

// CreateReview 创建商品评价
func (s *productReviewService) CreateReview(ctx *gin.Context, req model.CreateReviewRequest) error {
	userID := GetUserIdFromCtx(ctx)

	// 校验评分参数
	if req.Rating < 1 || req.Rating > 5 {
		return errors.New("商品评分必须在1-5之间")
	}
	if req.FreshnessRating < 1 || req.FreshnessRating > 5 {
		return errors.New("新鲜程度评分必须在1-5之间")
	}
	if req.PackagingRating < 1 || req.PackagingRating > 5 {
		return errors.New("包装评分必须在1-5之间")
	}
	if req.DeliveryRating < 1 || req.DeliveryRating > 5 {
		return errors.New("配送评分必须在1-5之间")
	}
	if req.ServiceRating < 1 || req.ServiceRating > 5 {
		return errors.New("服务态度评分必须在1-5之间")
	}

	// 调用仓储层创建评价
	return s.productReviewRepository.CreateReview(ctx, userID, req)
}

// GetProductReviews 获取商品评价列表
func (s *productReviewService) GetProductReviews(ctx *gin.Context, productID uint64, page, pageSize int) (*model.ReviewListResponse, error) {
	// 设置默认值
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	// 调用仓储层获取商品评价
	return s.productReviewRepository.GetProductReviews(ctx, productID, page, pageSize)
}

// GetUserReviews 获取用户评价列表
func (s *productReviewService) GetUserReviews(ctx *gin.Context, page, pageSize int) (*model.ReviewListResponse, error) {
	userID := GetUserIdFromCtx(ctx)

	// 设置默认值
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	// 调用仓储层获取用户评价
	return s.productReviewRepository.GetUserReviews(ctx, userID, page, pageSize)
}

// GetUserReviewsByTab 根据标签获取用户评价列表
func (s *productReviewService) GetUserReviewsByTab(ctx *gin.Context, req model.UserReviewListRequest) (*model.UserReviewListResponse, error) {
	userID := GetUserIdFromCtx(ctx)

	// 验证标签类型
	switch req.Tab {
	case model.ReviewTabAll, model.ReviewTabPending, model.ReviewTabCompleted:
		// 有效的标签类型
	default:
		return nil, errors.New("无效的标签类型")
	}

	// 调用仓储层获取用户评价
	return s.productReviewRepository.GetUserReviewsByTab(ctx, userID, req.Tab)
}

// DeleteReview 删除商品评价
func (s *productReviewService) DeleteReview(ctx *gin.Context, reviewID uint64) error {
	userID := GetUserIdFromCtx(ctx)

	// 调用仓储层删除评价
	return s.productReviewRepository.DeleteReview(ctx, userID, reviewID)
}

// GetReviewDetail 获取商品评价详情
func (s *productReviewService) GetReviewDetail(ctx *gin.Context, reviewID uint64) (*model.ReviewDetailResponse, error) {
	// 调用仓储层获取评价详情
	return s.productReviewRepository.GetReviewDetail(ctx, reviewID)
}

// IncrementReviewCounter 更新评论计数器
func (s *productReviewService) IncrementReviewCounter(ctx *gin.Context, req model.UpdateReviewCounterRequest) error {
	// 调用仓储层更新计数器
	return s.productReviewRepository.IncrementReviewCounter(ctx, req.ReviewID, req.Type)
}
