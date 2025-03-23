package repository

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"

	"app/internal/cache"
	"app/internal/common"
	"app/internal/model"
	"app/pkg/config"
	"app/pkg/oss"
)

type SettingsRepository interface {
	GetSettings(ctx context.Context, name string) (string, error)
	UploadAvatar(ctx *gin.Context, userId string, file *os.File) (string, error)
	UploadAvatarBase64(ctx *gin.Context, userId string, base64 string) (string, error)
	Logout(ctx context.Context, userId string) error
	SignOut(ctx context.Context, userId string) error
}

func NewSettingsRepository(
	repository *Repository,
) SettingsRepository {
	return &settingsRepository{
		Repository: repository,
	}
}

type settingsRepository struct {
	*Repository
	*cache.Cache
}

func (r *settingsRepository) GetCache(ctx context.Context) *cache.Cache {
	if r.Cache != nil {
		return r.Cache
	}
	r.Cache = cache.NewCache(ctx, "settings:")
	return r.Cache
}

func (r *settingsRepository) GetSettings(ctx context.Context, name string) (string, error) {
	// 从缓存中获取settings
	settings := map[string]string{}
	settingsStr, err := r.GetCache(ctx).GetString(common.PREFFIX_GLOBAL_SETTINGS)
	if err == nil && len(settingsStr) > 0 {
		if err = json.Unmarshal([]byte(settingsStr), &settings); err != nil {
			r.logger.Debug("json.Unmarshal error info: " + err.Error())
			return "", err
		}
		if v, ok := settings[name]; ok {
			r.logger.Debug("cache GetSettings", "name:", name, ", value:", v)
			return v, nil
		}
	}
	// 根据name_en查询settings表
	var settingsModels []*model.Settings
	if err = r.DB(ctx).Find(&settingsModels).Error; err != nil {
		r.logger.Debug("First error info: " + err.Error())
		return "", err
	}
	if len(settingsModels) == 0 {
		r.logger.Debug("settings not found")
		return "", errors.New("settings not found")
	}
	for _, item := range settingsModels {
		settings[item.Key] = item.Value
	}
	settingsModelBytes, err := json.Marshal(settings)
	if err != nil {
		r.logger.Debug("json.Marshal error info: " + err.Error())
		return "", err
	}
	if err = r.GetCache(ctx).Set(common.PREFFIX_GLOBAL_SETTINGS, string(settingsModelBytes)); err != nil {
		r.logger.Debug("Set error info: " + err.Error())
		return "", err
	}
	if v, ok := settings[name]; ok {
		r.logger.Debug("GetSettings", "name:", name, ", value:", v)
		return v, nil
	}
	return "", nil
}

// 通过文件base64上传头像
func (r *settingsRepository) UploadAvatarBase64(ctx *gin.Context, userId string, base64 string) (string, error) {
	// 上传文件到OSS
	ossPath, err := oss.NewOssBase64(base64).UploadBase64(ctx)
	if err != nil {
		r.logger.Debug("UploadFile failed", "err", err)
		return "", err
	}
	return ossPath, nil
}

// 通过文件附件上传
func (r *settingsRepository) UploadAvatar(ctx *gin.Context, userId string, file *os.File) (string, error) {
	// 创建临时文件
	tempFile, err := os.CreateTemp("", "upload-"+userId+"-*")
	if err != nil {
		r.logger.Debug("os.CreateTemp failed", "err", err)
		return "", err
	}
	defer tempFile.Close()

	// 将上传的文件内容复制到临时文件
	if _, err = tempFile.ReadFrom(file); err != nil {
		r.logger.Debug("tempFile.ReadFrom failed", "err", err)
		return "", err
	}
	// 获取临时文件路径
	tempFilePath := tempFile.Name()
	// 上传文件到OSS
	ossPath, err := oss.NewOss(tempFilePath, filepath.Ext(file.Name())).UploadFile(ctx)
	if err != nil {
		r.logger.Debug("UploadFile failed", "err", err)
		return "", err
	}
	// 删除临时文件
	defer os.Remove(tempFilePath)
	// 返回OSS路径
	return ossPath, nil
}

// 退出账号
func (r *settingsRepository) Logout(ctx context.Context, userId string) error {
	// 清空REDIS中的login_expire缓存标志
	if err := config.Rdb.Del(ctx, common.GetCryptoKey(userId, common.CODE_TYPE_LOGIN_EXPIRE)).Err(); err != nil {
		r.logger.Debug("config.Rdb.Del failed", "err", err)
		return err
	}
	return nil
}

// 注销账号
func (r *settingsRepository) SignOut(ctx context.Context, userId string) error {
	err := r.Transaction(
		ctx, func(ctx context.Context) error {
			// 将user表中的status字段设置为2
			if err := r.DB(ctx).Model(&model.Account{}).Where(
				"account = ? and status = ?", userId, common.STATUS_NORMAL,
			).
				Update("status", common.STATUS_CANCEL).Error; err != nil {
				r.logger.Debug("Logout failed", "err", err)
				return err
			}
			// 清空REDIS中的login_expire缓存标志
			if err := config.Rdb.Del(
				ctx, common.GetCryptoKey(userId, common.CODE_TYPE_LOGIN_EXPIRE),
			).Err(); err != nil {
				r.logger.Debug("config.Rdb.Del failed", "err", err)
				return err
			}
			return nil
		},
	)
	if err != nil {
		r.logger.Debug("Transaction failed", "err", err)
		return err
	}
	return nil
}
