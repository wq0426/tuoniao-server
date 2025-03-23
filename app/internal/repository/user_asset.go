package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/goccy/go-json"
	"gorm.io/gorm"

	"app/internal/common"
	"app/internal/model"
	"app/pkg/config"
)

type UserAssetRepository interface {
	Create(ctx context.Context, userResource *model.UserAsset) error
	GetUserAsset(ctx context.Context, userId string) (*model.UserAsset, error)
	UpdateUserAsset(ctx context.Context, tx *gorm.DB, userId string, bizType, actionType, assetType int8, rewardNum float32, relationId int, relationTitle string) error
	RechargeBalance(ctx context.Context, userID string, amount float64) error
	WithdrawBalance(ctx context.Context, userID string, amount float64) error
}

func NewUserAssetRepository(
	repository *Repository,
	userAssetRecordRepository UserAssetRecordRepository,
) UserAssetRepository {
	return &userAssetRepository{
		Repository:                repository,
		userAssetRecordRepository: userAssetRecordRepository,
	}
}

type userAssetRepository struct {
	*Repository
	userAssetRecordRepository UserAssetRecordRepository
}

func (r *userAssetRepository) Create(ctx context.Context, userResource *model.UserAsset) error {
	if err := r.DB(ctx).Create(userResource).Error; err != nil {
		r.logger.Error("Create error: " + err.Error())
		return err
	}
	return nil
}

func (r *userAssetRepository) GetUserAsset(ctx context.Context, userId string) (*model.UserAsset, error) {
	var userAsset model.UserAsset
	// 通过userId查询用户资产信息
	if err := r.DB(ctx).Debug().Where("user_id = ?", userId).First(&userAsset).Error; err != nil {
		r.logger.Error("GetUserAsset error: " + err.Error())
		return nil, err
	}
	return &userAsset, nil
}

func (r *userAssetRepository) UpdateUserAsset(
	ctx context.Context, tx *gorm.DB, userId string, bizType, actionType,
	assetType int8, rewardNum float32, relationId int, relationTitle string,
) error {
	var userAsset model.UserAsset
	if err := tx.Where("user_id = ?", userId).First(&userAsset).Error; err != nil {
		r.logger.Debug("First error info: " + err.Error())
		return err
	}
	// 资产是否够
	if assetType == common.ASSET_TYPE_POINT {
		if rewardNum < 0 && userAsset.Points < userAsset.Points+int(rewardNum) {
			return fmt.Errorf("积分不足")
		}
	} else if assetType == common.ASSET_TYPE_BALANCE {
		if rewardNum < 0 && userAsset.Balance < userAsset.Balance+float64(rewardNum) {
			return fmt.Errorf("余额不足")
		}
	}
	leftNum := float64(0)
	switch assetType {
	case common.ASSET_TYPE_POINT:
		userAsset.Points += int(rewardNum)
		leftNum = float64(userAsset.Points)
	case common.ASSET_TYPE_BALANCE:
		userAsset.Balance += float64(rewardNum)
		leftNum = float64(userAsset.Balance)
	}
	userAsset.UpdatedAt = time.Now()
	// 使用事务
	if err := tx.Model(&model.UserAsset{}).Where("user_id = ?", userId).Updates(
		map[string]interface{}{
			"points":     userAsset.Points,
			"balance":    userAsset.Balance,
			"updated_at": userAsset.UpdatedAt,
		},
	).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 更新资产
	if err := tx.Model(&model.UserAsset{}).Where(
		"user_id = ?", userId,
	).Save(&userAsset).Error; err != nil {
		tx.Rollback()
		r.logger.Debug("Model Updates error info: " + err.Error())
		return err
	}
	// 添加资产更新记录
	if err := r.userAssetRecordRepository.CreateUserAssetRecord(
		ctx,
		&model.UserAssetRecord{
			UserId:        userId,
			BusinessType:  bizType,
			ActionType:    actionType,
			AssetType:     assetType,
			ActionNum:     float32(rewardNum),
			LeftNum:       float32(leftNum),
			RelationId:    relationId,
			RelationTitle: relationTitle,
		},
	); err != nil {
		tx.Rollback()
		r.logger.Debug("Create error info: ", err.Error())
		return err
	}

	userAssetBytes, err := json.Marshal(userAsset)
	if err != nil {
		r.logger.Debug("json.Marshal error info: " + err.Error())
		return err
	}
	if err = config.Rdb.Set(ctx, common.PREFFIX_USER_ASSET+userId, string(userAssetBytes), 24*time.Hour).Err(); err != nil {
		r.logger.Debug("Set error info: " + err.Error())
		return err
	}
	return nil
}

func (r *userAssetRepository) RechargeBalance(ctx context.Context, userID string, amount float64) error {
	// 获取当前用户资产
	userAsset, err := r.GetUserAsset(ctx, userID)
	if err != nil {
		return err
	}

	// 保存到数据库
	tx := r.DB(ctx).Begin()
	if err := tx.Model(&model.UserAsset{}).Where("user_id = ?", userID).Updates(
		map[string]interface{}{
			"balance":    userAsset.Balance,
			"updated_at": userAsset.UpdatedAt,
		},
	).Error; err != nil {
		r.logger.Debug("充值更新余额失败: " + err.Error())
		tx.Rollback()
		return err
	}

	// 添加资产更新记录
	if err := r.UpdateUserAsset(
		ctx,
		tx,
		userID,
		common.BUSINESS_TYPE_RECHARGE,
		common.ACTION_TYPE_RECHARGE,
		common.ASSET_TYPE_BALANCE,
		float32(amount),
		0,
		"",
	); err != nil {
		r.logger.Debug("创建充值记录失败: " + err.Error())
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		r.logger.Debug("提交事务失败: " + err.Error())
		return err
	}

	// 更新缓存
	userAssetBytes, err := json.Marshal(userAsset)
	if err != nil {
		r.logger.Debug("序列化用户资产失败: " + err.Error())
		return err
	}
	if err = config.Rdb.Set(
		ctx, common.PREFFIX_USER_ASSET+userID, string(userAssetBytes), 24*time.Hour,
	).Err(); err != nil {
		r.logger.Debug("更新缓存失败: " + err.Error())
		return err
	}

	return nil
}

func (r *userAssetRepository) WithdrawBalance(ctx context.Context, userID string, amount float64) error {
	// 获取当前用户资产
	userAsset, err := r.GetUserAsset(ctx, userID)
	if err != nil {
		return err
	}

	// 检查余额是否充足
	if userAsset.Balance < amount {
		return fmt.Errorf("余额不足")
	}

	// 更新余额
	// userAsset.Balance -= float64(amount)
	userAsset.UpdatedAt = time.Now()

	// 使用事务
	// if err := r.Transaction(ctx, func(ctx context.Context) error {
	// 	// 保存到数据库
	// 	if err := r.DB(ctx).Model(&model.UserAsset{}).Where("user_id = ?", userID).Updates(
	// 		map[string]interface{}{
	// 			"balance":    userAsset.Balance,
	// 			"updated_at": userAsset.UpdatedAt,
	// 		}).Error; err != nil {
	// 		r.logger.Debug("提现更新余额失败: " + err.Error())
	// 		return err
	// 	}

	// // 添加资产更新记录
	// if err := r.UpdateUserAsset(
	// 	ctx,
	// 	userID,
	// 	common.BUSINESS_TYPE_WITHDRAW,
	// 	common.ACTION_TYPE_WITHDRAW,
	// 	common.ASSET_TYPE_BALANCE,
	// 	float32(amount),
	// 	0,
	// 	"",
	// ); err != nil {
	// 	r.logger.Debug("创建提现记录失败: " + err.Error())
	// }
	// return nil

	// }); err != nil {
	// 	r.logger.Debug("提现失败: " + err.Error())
	// 	return err
	// }
	// // 更新缓存
	// userAssetBytes, err := json.Marshal(userAsset)
	// if err != nil {
	// 	r.logger.Debug("序列化用户资产失败: " + err.Error())
	// 	return err
	// }
	// if err = config.Rdb.Set(
	// 	ctx, common.PREFFIX_USER_ASSET+userID, string(userAssetBytes), 24*time.Hour,
	// ).Err(); err != nil {
	// 	r.logger.Debug("更新缓存失败: " + err.Error())
	// 	return err
	// }

	return nil
}
