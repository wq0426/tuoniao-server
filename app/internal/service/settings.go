package service

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"

	"app/internal/repository"
)

type SettingsService interface {
	GetSettings(ctx *gin.Context, name string) (string, error)
	UploadAvatar(ctx *gin.Context, file *multipart.FileHeader) (string, error)
	UploadAvatarBase64(ctx *gin.Context, base64 string) (string, error)
	Logout(ctx *gin.Context) error
	SignOut(ctx *gin.Context) error
}

func NewSettingsService(
	service *Service,
	settingsRepository repository.SettingsRepository,
	userRepository repository.AccountRepository,
) SettingsService {
	return &settingsService{
		Service:            service,
		settingsRepository: settingsRepository,
		userRepository:     userRepository,
	}
}

type settingsService struct {
	*Service
	settingsRepository repository.SettingsRepository
	userRepository     repository.AccountRepository
}

func (s *settingsService) GetSettings(ctx *gin.Context, name string) (string, error) {
	return s.settingsRepository.GetSettings(ctx, name)
}

// 上传头像文件
func (s *settingsService) UploadAvatar(ctx *gin.Context, file *multipart.FileHeader) (string, error) {
	userId := GetUserIdFromCtx(ctx)
	return s.settingsRepository.UploadAvatarBase64(ctx, userId, file.Filename)
}

// 上传头像文件
func (s *settingsService) UploadAvatarBase64(ctx *gin.Context, base64 string) (string, error) {
	userId := GetUserIdFromCtx(ctx)
	return s.settingsRepository.UploadAvatarBase64(ctx, userId, base64)
}

// 退出账号
func (s *settingsService) Logout(ctx *gin.Context) error {
	userId := GetUserIdFromCtx(ctx)
	return s.settingsRepository.Logout(ctx, userId)
}

// 注销账号
func (s *settingsService) SignOut(ctx *gin.Context) error {
	userId := GetUserIdFromCtx(ctx)
	return s.settingsRepository.SignOut(ctx, userId)
}
