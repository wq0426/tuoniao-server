package service

import (
	"app/internal/model"
	"app/internal/repository"

	"github.com/gin-gonic/gin"
)

type UserAssetService interface {
	GetUserAsset(ctx *gin.Context) (*model.UserAssetResponse, error)
	RechargeBalance(ctx *gin.Context, req model.RechargeRequest) error
	WithdrawBalance(ctx *gin.Context, req model.WithdrawRequest) error
	GetBalanceRecords(ctx *gin.Context, req model.BalanceRecordQueryRequest) (*model.BalanceRecordResponse, error)
	GetWithdrawRecords(ctx *gin.Context, req model.BalanceRecordQueryRequest) (*model.BalanceRecordResponse, error)
	GetExchangeRecords(ctx *gin.Context, req model.BalanceRecordQueryRequest) (*model.PointExchangeRecordResponse, error)
}

func NewUserAssetService(
	service *Service,
	userAssetRepository repository.UserAssetRepository,
	userCouponRepository repository.UserCouponRepository,
	userAssetRecordRepository repository.UserAssetRecordRepository,
	userRepository repository.AccountRepository,
) UserAssetService {
	return &userAssetService{
		Service:                   service,
		userAssetRepository:       userAssetRepository,
		userCouponRepository:      userCouponRepository,
		userAssetRecordRepository: userAssetRecordRepository,
		userRepository:            userRepository,
	}
}

type userAssetService struct {
	*Service
	userAssetRepository       repository.UserAssetRepository
	userCouponRepository      repository.UserCouponRepository
	userAssetRecordRepository repository.UserAssetRecordRepository
	userRepository            repository.AccountRepository
}

func (s *userAssetService) GetUserAsset(ctx *gin.Context) (*model.UserAssetResponse, error) {
	userID := GetUserIdFromCtx(ctx)
	userAsset, err := s.userAssetRepository.GetUserAsset(ctx, userID)
	if err != nil {
		return nil, err
	}
	// 获取用户优惠券数量
	couponCount, err := s.userCouponRepository.GetUserCouponCount(ctx, userID)
	if err != nil {
		return nil, err
	}
	// 获取用户信息
	user, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &model.UserAssetResponse{
		UserID:      userAsset.UserID,
		Points:      userAsset.Points,
		CouponCount: couponCount,
		Balance:     userAsset.Balance,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
	}, nil
}

func (s *userAssetService) RechargeBalance(ctx *gin.Context, req model.RechargeRequest) error {
	userID := GetUserIdFromCtx(ctx)
	return s.userAssetRepository.RechargeBalance(ctx, userID, req.Amount)
}

func (s *userAssetService) WithdrawBalance(ctx *gin.Context, req model.WithdrawRequest) error {
	userID := GetUserIdFromCtx(ctx)
	return s.userAssetRepository.WithdrawBalance(ctx, userID, req.Amount)
}

func (s *userAssetService) GetBalanceRecords(ctx *gin.Context, req model.BalanceRecordQueryRequest) (*model.BalanceRecordResponse, error) {
	userID := GetUserIdFromCtx(ctx)
	return s.userAssetRecordRepository.GetBalanceRecords(ctx, userID, req.Page, req.PageSize)
}

func (s *userAssetService) GetWithdrawRecords(ctx *gin.Context, req model.BalanceRecordQueryRequest) (*model.BalanceRecordResponse, error) {
	userID := GetUserIdFromCtx(ctx)
	return s.userAssetRecordRepository.GetWithdrawRecords(ctx, userID, req.Page, req.PageSize)
}

func (s *userAssetService) GetExchangeRecords(ctx *gin.Context, req model.BalanceRecordQueryRequest) (*model.PointExchangeRecordResponse, error) {
	userID := GetUserIdFromCtx(ctx)
	return s.userAssetRecordRepository.GetExchangeRecords(ctx, userID, req.Page, req.PageSize)
}
