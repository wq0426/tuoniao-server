package repository

import (
	"context"
	"encoding/json"
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"

	v1 "app/api/v1"
	"app/internal/cache"
	"app/internal/common"
	"app/internal/model"
)

type AccountRepository interface {
	GetCache(ctx context.Context) *cache.Cache
	Create(ctx context.Context, user *model.Account) error
	Update(ctx context.Context, user *model.Account) error
	GetByID(ctx context.Context, id string) (*model.Account, error)
	GetByPhone(ctx context.Context, phone string) (*model.Account, error)
	UpdateAccount(ctx context.Context, account *model.Account) (*model.Account, error)
	GetByOpenID(ctx context.Context, openID string) (*model.Account, error)
}

func NewAccountRepository(
	repository *Repository,
	settingsRepository SettingsRepository,
) AccountRepository {
	return &accountRepository{
		Repository:         repository,
		settingsRepository: settingsRepository,
	}
}

type accountRepository struct {
	*Repository
	settingsRepository SettingsRepository
	*cache.Cache
}

func (r *accountRepository) GetCache(ctx context.Context) *cache.Cache {
	if r.Cache != nil {
		return r.Cache
	}
	r.Cache = cache.NewCache(ctx, "account:")
	return r.Cache
}

func (r *accountRepository) Create(ctx context.Context, user *model.Account) error {
	if err := r.DB(ctx).Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *accountRepository) Update(ctx context.Context, user *model.Account) error {
	// 删除缓存
	err := r.GetCache(ctx).Del(user.UserId)
	if err := r.DB(ctx).Save(user).Error; err != nil {
		return err
	}
	if err != nil {
		r.logger.Debug("Set error info: " + err.Error())
		return err
	}
	return nil
}

func (r *accountRepository) GetByID(ctx context.Context, userId string) (*model.Account, error) {
	var accountInfoModel model.Account
	accountInfo, err := r.GetCache(ctx).GetString(userId)
	if err != nil || len(accountInfo) == 0 {
		// 从redis缓存中获取用户信息
		if err = r.DB(ctx).Where(
			"user_id = ? AND status = ?", userId, common.STATUS_NORMAL,
		).First(&accountInfoModel).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, v1.ErrNotFound
			}
			r.logger.Debug("DB error info: " + err.Error())
			return nil, err
		}
		accountInfoStr, _ := json.Marshal(accountInfoModel)
		err = r.GetCache(ctx).Set(userId, accountInfoStr)
		if err != nil {
			r.logger.Debug("Set error info: " + err.Error())
			return nil, err
		}
	} else {
		// 将accountInfo转成model.Account
		if err = json.Unmarshal([]byte(accountInfo), &accountInfoModel); err != nil {
			r.logger.Debug("json.Unmarshal error info: " + err.Error())
			return nil, err
		}
	}
	return &accountInfoModel, nil
}

func (r *accountRepository) GetByPhone(ctx context.Context, phone string) (*model.Account, error) {
	var user model.Account
	if err := r.DB(ctx).Where("phone = ? AND status = ?", phone, common.STATUS_NORMAL).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *accountRepository) UpdateAccount(ctx context.Context, account *model.Account) (*model.Account, error) {
	// 更新用户信息
	if err := r.DB(ctx).Save(&account).Error; err != nil {
		return nil, err
	}
	err := r.GetCache(ctx).Del(account.UserId)
	if err != nil {
		r.logger.Debug("Set error info: " + err.Error())
		return nil, err
	}
	return account, nil
}

func (r *accountRepository) GetByOpenID(ctx context.Context, openID string) (*model.Account, error) {
	var account model.Account
	err := r.DB(ctx).Where("open_id = ? AND status = ?", openID, common.STATUS_NORMAL).First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.logger.Error("查询用户失败", zap.Error(err))
		return nil, err
	}
	return &account, nil
}
