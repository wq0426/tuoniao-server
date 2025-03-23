package repository

import (
	"app/internal/model"
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// UserCouponRepository interface for user coupon operations
type UserCouponRepository interface {
	GetUserCoupons(ctx context.Context, userID string, productID uint64) ([]model.UserCoupon, error)
	GetAllUserCoupons(ctx context.Context, userID string) (map[uint64][]model.ProductCouponDTO, error)
	ClaimCoupon(ctx context.Context, userID string, couponID uint64) error
	GetUserCouponCount(ctx context.Context, userID string) (int, error)
	GetUserCouponByID(ctx context.Context, userID string, couponID uint64) (*model.CouponDetailResponse, error)
}

// Implementation
type userCouponRepository struct {
	*Repository
}

func NewUserCouponRepository(repository *Repository) UserCouponRepository {
	return &userCouponRepository{
		Repository: repository,
	}
}

// GetUserCoupons retrieves coupons for a specific user and product
func (r *userCouponRepository) GetUserCoupons(ctx context.Context, userID string, productID uint64) ([]model.UserCoupon, error) {
	var coupons []model.UserCoupon

	// Query coupons by user ID and product ID
	result := r.DB(ctx).Where("user_id = ? AND product_id = ?", userID, productID).Find(&coupons)
	if result.Error != nil {
		return nil, result.Error
	}

	return coupons, nil
}

// GetAllUserCoupons retrieves all coupons for a specific user
func (r *userCouponRepository) GetAllUserCoupons(ctx context.Context, userID string) (map[uint64][]model.ProductCouponDTO, error) {
	var userCoupons []*model.UserCoupon
	userCouponMap := make(map[uint64]uint8)

	// 构建查询条件
	query := r.DB(ctx).Where("user_id = ?", userID)
	// 执行查询
	result := query.Model(&model.UserCoupon{}).Find(&userCoupons)
	if result.Error != nil {
		return nil, result.Error
	}

	// 获取所有的券ID
	couponIDs := []uint64{}
	for _, userCoupon := range userCoupons {
		couponIDs = append(couponIDs, userCoupon.CouponID)
		userCouponMap[userCoupon.CouponID] = userCoupon.Status
	}

	// 如果没有优惠券，直接返回空数组
	if len(couponIDs) == 0 {
		return map[uint64][]model.ProductCouponDTO{}, nil
	}

	// 获取所有的券
	coupons := []model.ProductCoupon{}
	result = r.DB(ctx).Where("id IN (?)", couponIDs).Find(&coupons)
	if result.Error != nil {
		return nil, result.Error
	}

	productCouponsDTOMap := make(map[uint64][]model.ProductCouponDTO)
	for _, coupon := range coupons {
		deadlineStr := ""
		if coupon.Deadline != nil {
			deadlineStr = coupon.Deadline.Format("2006-01-02")
		}

		couponDTO := model.ProductCouponDTO{
			CouponID:          coupon.ID,
			CouponName:        coupon.CouponName,
			CouponPrice:       coupon.CouponPrice,
			AvailableMinPrice: coupon.AvailableMinPrice,
			Deadline:          deadlineStr,
			ProductID:         uint(coupon.ProductID),
		}
		couponStatus := 0
		// 是否过期
		if coupon.Deadline != nil && coupon.Deadline.Before(time.Now()) {
			couponStatus = 2
		} else if v, ok := userCouponMap[uint64(coupon.ID)]; ok {
			couponStatus = int(v)
		}
		couponDTO.IsReceived = couponStatus
		if _, ok := productCouponsDTOMap[uint64(couponStatus)]; !ok {
			productCouponsDTOMap[uint64(couponStatus)] = []model.ProductCouponDTO{
				couponDTO,
			}
		} else {
			productCouponsDTOMap[uint64(coupon.Status)] = append(productCouponsDTOMap[uint64(coupon.Status)], couponDTO)
		}
	}

	return productCouponsDTOMap, nil
}

// ClaimCoupon claims a coupon for a user
func (r *userCouponRepository) ClaimCoupon(ctx context.Context, userID string, couponID uint64) error {
	// Check if the user has already claimed this coupon
	var count int64
	if err := r.DB(ctx).Model(&model.UserCoupon{}).
		Where("user_id = ? AND coupone_id = ?", userID, couponID).
		Count(&count).Error; err == nil && count > 0 {
		return fmt.Errorf("用户已经领取过该优惠券")
	}

	// 检查优惠券是否存在
	var coupon model.ProductCoupon
	if err := r.DB(ctx).Where("id = ?", couponID).First(&coupon).Error; err != nil {
		return fmt.Errorf("优惠券不存在")
	}

	// Create new user coupon record
	userCoupon := model.UserCoupon{
		UserID:            userID,
		CouponID:          couponID,
		ProductID:         coupon.ProductID,
		Type:              1,
		Status:            0,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		CouponName:        coupon.CouponName,
		CouponPrice:       coupon.CouponPrice,
		AvailableMinPrice: coupon.AvailableMinPrice,
		Deadline:          *coupon.Deadline,
	}

	// Save to database
	if err := r.Transaction(ctx, func(ctx context.Context) error {
		if err := r.DB(ctx).Create(&userCoupon).Error; err != nil {
			return err
		}
		// 更新优惠券状态
		if err := r.DB(ctx).Model(&model.ProductCoupon{}).Where("id = ?", couponID).Update("status", 1).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		r.logger.Debug("Transaction error info: " + err.Error())
		return err
	}

	return nil
}

func (r *userCouponRepository) GetUserCouponCount(ctx context.Context, userID string) (int, error) {
	var count int64
	if err := r.DB(ctx).Model(&model.UserCoupon{}).Where("user_id = ? AND status = ?", userID, 0).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

// GetUserCouponByID 根据ID获取优惠券详情
func (r *userCouponRepository) GetUserCouponByID(ctx context.Context, userID string, couponID uint64) (*model.CouponDetailResponse, error) {
	var coupon model.UserCoupon
	if err := r.DB(ctx).Where("id = ? AND user_id = ?", couponID, userID).First(&coupon).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("优惠券不存在或不属于当前用户")
		}
		r.logger.Debug("查询优惠券详情失败", "error", err)
		return nil, err
	}

	return &model.CouponDetailResponse{
		ID:                coupon.ID,
		CouponName:        coupon.CouponName,
		CouponPrice:       coupon.CouponPrice,
		AvailableMinPrice: coupon.AvailableMinPrice,
		Type:              coupon.Type,
		Status:            coupon.Status,
		ProductID:         coupon.ProductID,
		Deadline:          coupon.Deadline,
		CreatedAt:         coupon.CreatedAt,
	}, nil
}
