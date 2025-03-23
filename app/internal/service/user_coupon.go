package service

import (
	"app/internal/model"
	"app/internal/repository"

	"github.com/gin-gonic/gin"
)

type UserCouponService interface {
	GetUserCoupons(ctx *gin.Context, productID uint64) ([]model.UserCoupon, error)
	GetAllUserCoupons(ctx *gin.Context) (map[uint64][]model.ProductCouponDTO, error)
	ClaimCoupon(ctx *gin.Context, req model.ClaimCouponRequest) error
	GetUserCouponCount(ctx *gin.Context) (int, error)
	GetUserCouponByID(ctx *gin.Context, couponID uint64) (*model.CouponDetailResponse, error)
}

func NewUserCouponService(
	service *Service,
	userCouponRepository repository.UserCouponRepository,
) UserCouponService {
	return &userCouponService{
		Service:              service,
		userCouponRepository: userCouponRepository,
	}
}

type userCouponService struct {
	*Service
	userCouponRepository repository.UserCouponRepository
}

func (s *userCouponService) GetUserCoupons(ctx *gin.Context, productID uint64) ([]model.UserCoupon, error) {
	userID := GetUserIdFromCtx(ctx)
	return s.userCouponRepository.GetUserCoupons(ctx, userID, productID)
}

func (s *userCouponService) GetAllUserCoupons(ctx *gin.Context) (map[uint64][]model.ProductCouponDTO, error) {
	userID := GetUserIdFromCtx(ctx)
	return s.userCouponRepository.GetAllUserCoupons(ctx, userID)
}

func (s *userCouponService) ClaimCoupon(ctx *gin.Context, req model.ClaimCouponRequest) error {
	userID := GetUserIdFromCtx(ctx)
	return s.userCouponRepository.ClaimCoupon(ctx, userID, req.CouponID)
}

func (s *userCouponService) GetUserCouponCount(ctx *gin.Context) (int, error) {
	userID := GetUserIdFromCtx(ctx)
	return s.userCouponRepository.GetUserCouponCount(ctx, userID)
}

func (s *userCouponService) GetUserCouponByID(ctx *gin.Context, couponID uint64) (*model.CouponDetailResponse, error) {
	userID := GetUserIdFromCtx(ctx)
	return s.userCouponRepository.GetUserCouponByID(ctx, userID, couponID)
}
