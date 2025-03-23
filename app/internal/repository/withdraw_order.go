package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"app/internal/common"
	"app/internal/model"

	"gorm.io/gorm"
)

type WithdrawOrderRepository interface {
	CreateWithdraw(ctx context.Context, userID string, req model.CreateWithdrawRequest) error
	GetWithdrawList(ctx context.Context, userID string, status *uint8, page, pageSize int) (*model.WithdrawListResponse, error)
	GetWithdrawDetail(ctx context.Context, userID string, withdrawID uint64) (*model.WithdrawDetailResponse, error)
}

type withdrawOrderRepository struct {
	*Repository
}

func NewWithdrawOrderRepository(repository *Repository) WithdrawOrderRepository {
	return &withdrawOrderRepository{
		Repository: repository,
	}
}

// CreateWithdraw 创建提现单
func (r *withdrawOrderRepository) CreateWithdraw(ctx context.Context, userID string, req model.CreateWithdrawRequest) error {
	// 校验余额是否足够
	var userAsset model.UserAsset
	if err := r.DB(ctx).Where("user_id = ?", userID).First(&userAsset).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户资产不存在")
		}
		r.logger.Error("查询用户资产失败: " + err.Error())
		return err
	}

	if userAsset.Balance < req.Amount {
		return errors.New("余额不足")
	}

	// 计算手续费和实际到账金额（这里以1%手续费为例）
	fee := req.Amount * 0.01
	actualAmount := req.Amount - fee

	// 开始事务
	tx := r.DB(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 生成提现单号
	withdrawNo := fmt.Sprintf("W%s%s", time.Now().Format("20060102"), common.GenerateOrderNo())

	// 创建提现单
	now := time.Now()
	withdrawOrder := model.WithdrawOrder{
		WithdrawNo:   withdrawNo,
		UserID:       userID,
		Amount:       req.Amount,
		Fee:          fee,
		ActualAmount: actualAmount,
		Status:       0, // 待审核
		// BankName:     req.BankName,
		// AccountName:  req.AccountName,
		// AccountNo:    req.AccountNo,
		// Remark:       req.Remark,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	if err := tx.Create(&withdrawOrder).Error; err != nil {
		tx.Rollback()
		r.logger.Error("创建提现单失败: " + err.Error())
		return err
	}

	// 扣减用户余额
	if err := tx.Model(&model.UserAsset{}).Where("user_id = ?", userID).
		Update("balance", gorm.Expr("balance - ?", req.Amount)).Error; err != nil {
		tx.Rollback()
		r.logger.Error("扣减用户余额失败: " + err.Error())
		return err
	}

	// 记录资产变动日志
	assetRecord := model.UserAssetRecord{
		UserId:        userID,
		BusinessType:  common.BUSINESS_TYPE_WITHDRAW,
		ActionType:    common.ACTION_TYPE_WITHDRAW,
		AssetType:     common.ASSET_TYPE_BALANCE,
		ActionNum:     float32(-req.Amount), // 负数表示支出
		LeftNum:       float32(userAsset.Balance - req.Amount),
		RelationId:    int(withdrawOrder.ID),
		RelationTitle: "提现",
		CreatedAt:     now,
	}

	if err := tx.Create(&assetRecord).Error; err != nil {
		tx.Rollback()
		r.logger.Error("创建资产变动记录失败: " + err.Error())
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		r.logger.Error("提交事务失败: " + err.Error())
		return err
	}

	return nil
}

// GetWithdrawList 获取提现单列表
func (r *withdrawOrderRepository) GetWithdrawList(ctx context.Context, userID string, status *uint8, page, pageSize int) (*model.WithdrawListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	var total int64
	var withdrawOrders []model.WithdrawOrder

	// 构建查询条件
	query := r.DB(ctx).Model(&model.WithdrawOrder{}).Where("user_id = ?", userID)
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// 查询总数
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("查询提现单总数失败: " + err.Error())
		return nil, err
	}

	// 查询列表
	if err := query.Order("created_at DESC").
		Limit(pageSize).Offset(offset).
		Find(&withdrawOrders).Error; err != nil {
		r.logger.Error("查询提现单列表失败: " + err.Error())
		return nil, err
	}

	// 转换为列表项
	list := make([]model.WithdrawListItem, 0, len(withdrawOrders))
	for _, order := range withdrawOrders {
		list = append(list, model.WithdrawListItem{
			ID:           order.ID,
			Title:        "提现",
			WithdrawNo:   order.WithdrawNo,
			Amount:       order.Amount,
			Fee:          order.Fee,
			ActualAmount: order.ActualAmount,
			Status:       order.Status,
			StatusText:   getWithdrawStatusText(order.Status),
			BankName:     order.BankName,
			AccountName:  order.AccountName,
			AccountNo:    order.AccountNo,
			CreatedAt:    order.CreatedAt,
		})
	}

	return &model.WithdrawListResponse{
		Total: total,
		List:  list,
		Page:  page,
		Size:  pageSize,
	}, nil
}

// GetWithdrawDetail 获取提现单详情
func (r *withdrawOrderRepository) GetWithdrawDetail(ctx context.Context, userID string, withdrawID uint64) (*model.WithdrawDetailResponse, error) {
	var withdrawOrder model.WithdrawOrder
	if err := r.DB(ctx).Where("id = ? AND user_id = ?", withdrawID, userID).First(&withdrawOrder).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("提现单不存在或无权限")
		}
		r.logger.Error("查询提现单详情失败: " + err.Error())
		return nil, err
	}

	return &model.WithdrawDetailResponse{
		ID:           withdrawOrder.ID,
		WithdrawNo:   withdrawOrder.WithdrawNo,
		UserID:       withdrawOrder.UserID,
		Amount:       withdrawOrder.Amount,
		Fee:          withdrawOrder.Fee,
		ActualAmount: withdrawOrder.ActualAmount,
		Status:       withdrawOrder.Status,
		StatusText:   getWithdrawStatusText(withdrawOrder.Status),
		RejectReason: withdrawOrder.RejectReason,
		BankName:     withdrawOrder.BankName,
		AccountName:  withdrawOrder.AccountName,
		AccountNo:    withdrawOrder.AccountNo,
		Remark:       withdrawOrder.Remark,
		AuditTime:    withdrawOrder.AuditTime,
		CompleteTime: withdrawOrder.CompleteTime,
		CreatedAt:    withdrawOrder.CreatedAt,
		UpdatedAt:    withdrawOrder.UpdatedAt,
	}, nil
}

// getWithdrawStatusText 获取提现状态文本
func getWithdrawStatusText(status uint8) string {
	switch status {
	case 0:
		return "审核中"
	case 1:
		return "处理中"
	case 2:
		return "已打款"
	case 3:
		return "已拒绝"
	default:
		return "未知状态"
	}
}
