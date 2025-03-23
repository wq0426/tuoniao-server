package repository

import (
	"context"

	"app/internal/common"
	"app/internal/model"
)

type UserAssetRecordRepository interface {
	GetUserAssetRecord(ctx context.Context, id int64) (*model.UserAssetRecord, error)
	CreateUserAssetRecord(ctx context.Context, userAssetRecord *model.UserAssetRecord) error
	GetBalanceRecords(ctx context.Context, userID string, page, pageSize int) (*model.BalanceRecordResponse, error)
	GetWithdrawRecords(ctx context.Context, userID string, page, pageSize int) (*model.BalanceRecordResponse, error)
	GetExchangeRecords(ctx context.Context, userID string, page, pageSize int) (*model.PointExchangeRecordResponse, error)
}

func NewUserAssetRecordRepository(
	repository *Repository,
) UserAssetRecordRepository {
	return &userAssetRecordRepository{
		Repository: repository,
	}
}

type userAssetRecordRepository struct {
	*Repository
}

func (r *userAssetRecordRepository) GetUserAssetRecord(ctx context.Context, id int64) (*model.UserAssetRecord, error) {
	var userAssetRecord model.UserAssetRecord

	return &userAssetRecord, nil
}

// 添加一条记录
func (r *userAssetRecordRepository) CreateUserAssetRecord(
	ctx context.Context, userAssetRecord *model.UserAssetRecord,
) error {
	if err := r.DB(ctx).Create(userAssetRecord).Error; err != nil {
		r.logger.Error("Create error: " + err.Error())
		return err
	}
	return nil
}

// 实现查询余额变动记录
func (r *userAssetRecordRepository) GetBalanceRecords(ctx context.Context, userID string, page, pageSize int) (*model.BalanceRecordResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	var total int64
	var records []model.UserAssetRecord

	// 查询总数
	if err := r.DB(ctx).Model(&model.UserAssetRecord{}).
		Where("user_id = ? AND asset_type = ?", userID, common.ASSET_TYPE_BALANCE).
		Count(&total).Error; err != nil {
		r.logger.Debug("查询余额记录总数失败", "error", err)
		return nil, err
	}

	// 查询记录
	if err := r.DB(ctx).Where("user_id = ? AND asset_type = ?", userID, common.ASSET_TYPE_BALANCE).
		Order("created_at DESC").
		Limit(pageSize).Offset((page - 1) * pageSize).
		Find(&records).Error; err != nil {
		r.logger.Debug("查询余额记录失败", "error", err)
		return nil, err
	}

	// 转换为DTO
	recordDTOs := make([]model.UserAssetRecordDTO, 0, len(records))
	for _, record := range records {
		title := r.getActionTypeString(record.BusinessType, record.ActionType, record.RelationTitle)
		recordDTOs = append(recordDTOs, model.UserAssetRecordDTO{
			ID:           uint64(record.Id),
			UserId:       record.UserId,
			Title:        title,
			BusinessType: record.BusinessType,
			ActionType:   record.ActionType,
			AssetType:    record.AssetType,
			ActionNum:    record.ActionNum,
			LeftNum:      record.LeftNum,
			CreatedAt:    record.CreatedAt.Format("2006-01-02 15:04"),
		})
	}

	return &model.BalanceRecordResponse{
		Total:    total,
		Records:  recordDTOs,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// 实现查询提现记录
func (r *userAssetRecordRepository) GetWithdrawRecords(ctx context.Context, userID string, page, pageSize int) (*model.BalanceRecordResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	var total int64
	var records []model.UserAssetRecord

	// 查询总数
	if err := r.DB(ctx).Model(&model.UserAssetRecord{}).
		Where("user_id = ? AND asset_type = ? AND action_type = ?",
			userID, common.ASSET_TYPE_BALANCE, common.ACTION_TYPE_WITHDRAW).
		Count(&total).Error; err != nil {
		r.logger.Debug("查询提现记录总数失败", "error", err)
		return nil, err
	}

	// 查询记录
	if err := r.DB(ctx).Where("user_id = ? AND asset_type = ? AND action_type = ?",
		userID, common.ASSET_TYPE_BALANCE, common.ACTION_TYPE_WITHDRAW).
		Order("created_at DESC").
		Limit(pageSize).Offset((page - 1) * pageSize).
		Find(&records).Error; err != nil {
		r.logger.Debug("查询提现记录失败", "error", err)
		return nil, err
	}

	// 转换为DTO
	recordDTOs := make([]model.UserAssetRecordDTO, 0, len(records))
	for _, record := range records {
		recordDTOs = append(recordDTOs, model.UserAssetRecordDTO{
			ID:           uint64(record.Id),
			UserId:       record.UserId,
			Title:        "提现",
			BusinessType: record.BusinessType,
			ActionType:   record.ActionType,
			AssetType:    record.AssetType,
			ActionNum:    record.ActionNum,
			LeftNum:      record.LeftNum,
			CreatedAt:    record.CreatedAt.Format("2006-01-02 15:04"),
		})
	}

	return &model.BalanceRecordResponse{
		Total:    total,
		Records:  recordDTOs,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// 实现查询积分兑换记录方法
func (r *userAssetRecordRepository) GetExchangeRecords(ctx context.Context, userID string, page, pageSize int) (*model.PointExchangeRecordResponse, error) {
	var records []model.PointExchangeRecord
	var total int64

	// 计算分页偏移量
	offset := (page - 1) * pageSize

	// 查询总数
	err := r.DB(ctx).Model(&model.PointExchangeRecord{}).
		Where("user_id = ?", userID).
		Count(&total).Error
	if err != nil {
		r.logger.Error("查询积分兑换记录总数失败: " + err.Error())
		return nil, err
	}

	// 查询记录列表
	err = r.DB(ctx).Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&records).Error
	if err != nil {
		r.logger.Error("查询积分兑换记录列表失败: " + err.Error())
		return nil, err
	}

	// 转换为DTO
	recordDTOs := make([]model.PointExchangeRecordDTO, 0, len(records))
	for _, record := range records {
		recordDTOs = append(recordDTOs, model.PointExchangeRecordDTO{
			ID:             record.ID,
			UserID:         record.UserID,
			ConfigID:       record.ConfigID,
			MinAmount:      record.MinAmount,
			ExchangeAmount: record.ExchangeAmount,
			Required:       record.Required,
			Points:         record.Points,
			Type:           record.Type,
			Images:         record.Images,
			Title:          record.Title,
			CreatedAt:      record.CreatedAt.Format("2006-01-02"),
			Deadline:       record.Deadline.Format("2006-01-02"),
		})
	}

	return &model.PointExchangeRecordResponse{
		Total:    total,
		Records:  recordDTOs,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// 根据业务类型和动作类型转成字符串
func (r *userAssetRecordRepository) getActionTypeString(businessType int8, actionType int8, relationTitle string) string {
	switch businessType {
	case common.BUSINESS_TYPE_ORDER:
		switch actionType {
		case common.ACTION_TYPE_BUY:
			return "购买(" + relationTitle + ")"
		case common.ACTION_TYPE_REFUND:
			return "退款(" + relationTitle + ")"
		}
	case common.BUSINESS_TYPE_RECHARGE:
		return "充值"
	case common.BUSINESS_TYPE_WITHDRAW:
		return "提现"
	case common.BUSINESS_TYPE_EXCHANGE:
		return "积分兑换"
	}
	return ""
}
