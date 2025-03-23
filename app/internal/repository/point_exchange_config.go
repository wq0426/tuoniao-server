package repository

import (
	"app/internal/model"
	"context"
	"fmt"
	"time"

	"app/internal/common"
)

type PointExchangeConfigRepository interface {
	GetPointExchangeConfigList(ctx context.Context, userID string) (*model.PointExchangeConfigListResponse, error)
	ExchangePoints(ctx context.Context, userID string, configID uint64) error
}

func NewPointExchangeConfigRepository(repository *Repository) PointExchangeConfigRepository {
	return &pointExchangeConfigRepository{
		Repository: repository,
	}
}

type pointExchangeConfigRepository struct {
	*Repository
}

func (r *pointExchangeConfigRepository) GetPointExchangeConfigList(ctx context.Context, userID string) (*model.PointExchangeConfigListResponse, error) {
	var configs []model.PointExchangeConfig
	var total int64

	// 查询总数
	if err := r.DB(ctx).Model(&model.PointExchangeConfig{}).Count(&total).Error; err != nil {
		r.logger.Error("查询积分兑换配置总数失败: " + err.Error())
		return nil, err
	}

	// 查询列表
	if err := r.DB(ctx).Find(&configs).Error; err != nil {
		r.logger.Error("查询积分兑换配置列表失败: " + err.Error())
		return nil, err
	}

	// 查询用户积分
	var userAsset model.UserAsset
	if err := r.DB(ctx).Where("user_id = ?", userID).First(&userAsset).Error; err != nil {
		r.logger.Error("查询用户积分失败: " + err.Error())
		return nil, err
	}

	// 查询用户兑换记录
	var pointExchangeRecords []model.PointExchangeRecord
	if err := r.DB(ctx).Where("user_id = ?", userID).Find(&pointExchangeRecords).Error; err != nil {
		r.logger.Error("查询用户兑换记录失败: " + err.Error())
		return nil, err
	}
	// 使用map存储兑换记录
	pointExchangeRecordsMap := make(map[uint64]model.PointExchangeRecord)
	for _, record := range pointExchangeRecords {
		pointExchangeRecordsMap[uint64(record.ConfigID)] = record
	}

	pointExchangeConfigResponse := make([]model.PointExchangeConfigResponse, 0)
	for _, config := range configs {
		isExchange := false
		if pointExchangeRecordsMap[uint64(config.ID)].ID > 0 {
			isExchange = true
		}
		pointExchangeConfigResponse = append(pointExchangeConfigResponse, model.PointExchangeConfigResponse{
			PointExchangeConfig: config,
			IsExchange:          isExchange,
		})
	}

	return &model.PointExchangeConfigListResponse{
		Total: int(total),
		List:  pointExchangeConfigResponse,
	}, nil
}

func (r *pointExchangeConfigRepository) ExchangePoints(ctx context.Context, userID string, configID uint64) error {
	// 获取兑换配置
	var config model.PointExchangeConfig
	if err := r.DB(ctx).Where("id = ?", configID).First(&config).Error; err != nil {
		r.logger.Error("获取积分兑换配置失败: " + err.Error())
		return err
	}

	// 获取用户资产
	var userAsset model.UserAsset
	if err := r.DB(ctx).Where("user_id = ?", userID).First(&userAsset).Error; err != nil {
		r.logger.Error("获取用户资产失败: " + err.Error())
		return err
	}

	// 检查用户积分是否足够
	if userAsset.Points < config.Points {
		return fmt.Errorf("积分不足")
	}

	// 开始事务
	tx := r.DB(ctx).Begin()

	// 扣除用户积分
	if err := tx.Model(&model.UserAsset{}).Where("user_id = ?", userID).
		Update("points", userAsset.Points-config.Points).Error; err != nil {
		tx.Rollback()
		r.logger.Error("扣除用户积分失败: " + err.Error())
		return err
	}

	// 记录积分变动
	assetRecord := &model.UserAssetRecord{
		UserId:       userID,
		BusinessType: common.BUSINESS_TYPE_EXCHANGE,
		ActionType:   common.ACTION_TYPE_EXCHANGE,
		AssetType:    common.ASSET_TYPE_POINT,
		ActionNum:    float32(config.Points),
		LeftNum:      float32(userAsset.Points - config.Points),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := tx.Create(assetRecord).Error; err != nil {
		tx.Rollback()
		r.logger.Error("记录积分变动失败: " + err.Error())
		return err
	}

	// 记录积分兑换记录
	pointExchangeRecord := &model.PointExchangeRecord{
		UserID:         userID,
		ConfigID:       int(configID),
		Points:         config.Points,
		MinAmount:      config.MinAmount,
		ExchangeAmount: config.ExchangeAmount,
		Required:       config.Required,
		Type:           config.Type,
		Images:         config.Images,
		Title:          config.Title,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Deadline:       time.Now().AddDate(0, 1, 0), // 设置一个月后过期
	}
	if err := tx.Create(pointExchangeRecord).Error; err != nil {
		tx.Rollback()
		r.logger.Error("记录积分兑换记录失败: " + err.Error())
		return err
	}

	// 创建用户优惠券
	userCoupon := &model.UserCoupon{
		Type:              uint8(config.Type),
		UserID:            userID,
		CouponID:          config.ID,
		Status:            0, // 未使用状态
		AvailableMinPrice: config.MinAmount,
		CouponPrice:       config.ExchangeAmount,
		Deadline:          time.Now().AddDate(0, 1, 0), // 设置一个月后过期
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	if err := tx.Create(userCoupon).Error; err != nil {
		tx.Rollback()
		r.logger.Error("创建用户优惠券失败: " + err.Error())
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		r.logger.Error("提交事务失败: " + err.Error())
		return err
	}

	return nil
}
