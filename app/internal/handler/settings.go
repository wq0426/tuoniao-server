package handler

import (
	"github.com/gin-gonic/gin"

	v1 "app/api/v1"
	"app/internal/service"
)

type SettingsHandler struct {
	*Handler
	settingsService service.SettingsService
}

func NewSettingsHandler(
	handler *Handler,
	settingsService service.SettingsService,
) *SettingsHandler {
	return &SettingsHandler{
		Handler:         handler,
		settingsService: settingsService,
	}
}

func (h *SettingsHandler) UploadAvatar(ctx *gin.Context) {
	// 上传头像
	file, err := ctx.FormFile("file")
	if err != nil {
		h.logger.Error("upload avatar FormFile error: " + err.Error())
		v1.HandleError(ctx, v1.ErrUploadAvatarCode, v1.MsgUploadAvatar, nil)
		return
	}
	avatarUrl, err := h.settingsService.UploadAvatar(ctx, file)
	// 接收参数到AvatarRequest结构体
	//var avatarRequest model.AvatarRequest
	//if err := ctx.ShouldBindJSON(&avatarRequest); err != nil {
	//	v1.HandleError(ctx, v1.ErrParamFormatCode, v1.MsgParamFormateError, nil)
	//	return
	//}
	//if !avatarRequest.Validate() {
	//	v1.HandleError(ctx, v1.ErrParamFormatCode, v1.MsgParamFormateError, nil)
	//	return
	//}
	//avatarUrl, err := h.settingsService.UploadAvatarBase64(ctx, avatarRequest.File)
	if err != nil {
		h.logger.Error("upload avatar UploadAvatar error: " + err.Error())
		v1.HandleError(ctx, v1.ErrUploadAvatarCode, v1.MsgUploadAvatar, nil)
		return
	}
	v1.HandleSuccess(ctx, avatarUrl)
}

func (h *SettingsHandler) Logout(ctx *gin.Context) {
	// 注销账号
	err := h.settingsService.Logout(ctx)
	if err != nil {
		v1.HandleError(ctx, v1.ErrLogoutCode, v1.MsgLogout, nil)
		return
	}
	v1.HandleSuccess(ctx, true)
}

func (h *SettingsHandler) SignOut(ctx *gin.Context) {
	// 注销账号
	err := h.settingsService.SignOut(ctx)
	if err != nil {
		v1.HandleError(ctx, v1.ErrSignOutCode, v1.MsgSignOut, nil)
		return
	}
	v1.HandleSuccess(ctx, true)
}
