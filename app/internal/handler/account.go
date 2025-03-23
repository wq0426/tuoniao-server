package handler

import (
	"github.com/gin-gonic/gin"

	v1 "app/api/v1"
	"app/internal/common"
	"app/internal/model"
	"app/internal/service"
)

type AccountHandler struct {
	*Handler
	userService service.AccountService
}

func NewAccountHandler(
	handler *Handler,
	accountService service.AccountService,
) *AccountHandler {
	return &AccountHandler{
		Handler:     handler,
		userService: accountService,
	}
}

func (h *AccountHandler) Register(ctx *gin.Context) {
	req := new(v1.RegisterRequest)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Debug("ShouldBindJSON error info: " + err.Error())
		v1.HandleError(ctx, v1.ErrParamFormatCode, v1.MsgParamFormateError, nil)
		return
	}
	if err := req.Validate(); err != nil {
		h.logger.Debug("Validate error info: " + err.Error())
		v1.HandleError(ctx, v1.ErrParamFormatCode, err.Error(), nil)
		return
	}
	if err := h.userService.Register(ctx, req); err != nil {
		h.logger.Debug("error info: " + err.Error())
		v1.HandleError(ctx, v1.ErrRegisterCode, v1.MsgRegister, nil)
		return
	}

	v1.HandleSuccess(ctx, true)
}

// Login godoc
// @Summary 账号登录
// @Schemes
// @Description 支持手机验证码和密码登录
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param request body v1.LoginRequest true "params"
// @Success 0 {object} v1.LoginResponse
// @Router /account/login [post]
func (h *AccountHandler) Login(ctx *gin.Context) {
	var req v1.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, v1.ErrParamFormatCode, v1.MsgParamFormateError, nil)
		return
	}
	if err := req.Validate(); err != nil {
		v1.HandleError(ctx, v1.ErrParamFormatCode, v1.MsgParamFormateError, nil)
		return
	}
	token, err := h.userService.Login(ctx, &req)
	if err != nil {
		h.logger.Debug("error info: " + err.Error())
		v1.HandleError(ctx, v1.ErrLoginCode, err.Error(), nil)
		return
	}
	v1.HandleSuccess(
		ctx, v1.LoginResponseData{
			AccessToken: token,
			Expired:     common.ONE_DAY_SECONDS,
		},
	)
}

// GetProfile godoc
// @Summary 获取用户信息
// @Schemes
// @Description
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 0 {object} v1.GetProfileResponse
// @Router /account/profile [get]
func (h *AccountHandler) GetProfile(ctx *gin.Context) {
	user, err := h.userService.GetProfile(ctx)
	if err != nil {
		v1.HandleError(ctx, v1.ErrGetProfileCode, v1.MsgGetProfile, nil)
		return
	}
	v1.HandleSuccess(ctx, user)
}

// ResetPassword godoc
// @Summary 重置用户密码
// @Schemes
// @Description 通过手机号和新密码重置用户密码
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body v1.ResetPasswordRequest true "params"
// @Success 0 {object} v1.Response{data=bool}
// @Router /account/reset [post]
func (h *AccountHandler) ResetPassword(ctx *gin.Context) {
	var req v1.ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, v1.ErrParamFormatCode, v1.MsgParamFormateError, nil)
		return
	}
	if err := req.Validate(); err != nil {
		v1.HandleError(ctx, v1.ErrParamFormatCode, v1.MsgParamFormateError, nil)
		return
	}
	if err := h.userService.ResetPassword(ctx, &req); err != nil {
		v1.HandleError(ctx, v1.ErrResetPasswordCode, v1.MsgResetPassword, nil)
		return
	}

	v1.HandleSuccess(ctx, true)
}

// UpdateProfileSettings godoc
// @Summary 修改昵称和头像
// @Schemes
// @Description 修改昵称和头像
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body model.NicknameAvatar true "修改昵称和头像"
// @Success 0 {object} v1.Response{data=bool}
// @Router /account/profile [put]
func (h *AccountHandler) UpdateProfileSettings(ctx *gin.Context) {
	// 获取请求参数
	var req model.NicknameAvatar
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, v1.ErrParamFormatCode, v1.MsgParamFormateError, nil)
		return
	}
	if !req.Validate() {
		v1.HandleError(ctx, v1.ErrParamFormatCode, v1.MsgParamFormateError, nil)
		return
	}
	// 修改昵称和头像
	err := h.userService.UpdateProfile(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, v1.ErrUpdateProfileSettingsCode, err.Error(), nil)
		return
	}
	v1.HandleSuccess(ctx, true)
}

// WeChatLoginCallback godoc
// @Summary 微信登录回调
// @Schemes
// @Description 处理微信登录回调
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param code query string true "微信授权临时票据"
// @Param state query string false "状态参数"
// @Success 0 {object} model.LoginResponseData
// @Router /account/wechat/login/callback [get]
func (h *AccountHandler) WeChatLoginCallback(ctx *gin.Context) {
	code := ctx.Query("code")
	if code == "" {
		v1.HandleError(ctx, v1.ErrParamFormatCode, "code is required", nil)
		return
	}

	// 使用code向微信服务器请求access_token和openid
	accessToken, openID, err := h.userService.GetWeChatAccessTokenAndOpenID(code)
	if err != nil {
		h.logger.Debug("error info: " + err.Error())
		v1.HandleError(ctx, v1.ErrWeChatLoginCallbackCode, err.Error(), nil)
		return
	}

	// 使用access_token和openid拉取用户信息
	bindInfo, err := h.userService.GetWeChatUserInfo(accessToken, openID)
	if err != nil {
		h.logger.Debug("error info: " + err.Error())
		v1.HandleError(ctx, v1.ErrWeChatLoginCallbackCode, err.Error(), nil)
		return
	}

	// 完成登录流程
	//userInfo, token, err := h.userService.WeChatLogin(ctx, bindInfo)
	//if err != nil {
	//	h.logger.Debug("error info: " + err.Error())
	//	v1.HandleError(ctx, v1.ErrWeChatLoginCallbackCode, err.Error(), nil)
	//	return
	//}

	v1.HandleSuccess(
		ctx, model.LoginResponseData{
			BindInfo: bindInfo,
			//UserInfo:    userInfo,
			//AccessToken: token,
			Expired: common.ONE_DAY_SECONDS,
		},
	)
}

// WeChatMiniLogin godoc
// @Summary 微信小程序登录
// @Description 使用微信小程序临时登录凭证登录并获取平台token
// @Tags 账户
// @Accept json
// @Produce json
// @Param request body model.WeChatLoginRequest true "微信登录请求"
// @Success 200 {object} v1.Response{data=model.WeChatLoginResponse} "成功"
// @Failure 400 {object} v1.Response "参数错误"
// @Failure 500 {object} v1.Response "服务器内部错误"
// @Router /account/wechat_mini_login [post]
func (h *AccountHandler) WeChatMiniLogin(c *gin.Context) {
	var req model.WeChatLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		v1.HandleError(c, v1.ErrParamCode, "参数错误", err)
		return
	}

	// 调用服务层进行微信登录
	sessionKey, openID, err := h.userService.WeChatMiniLogin(c, &req)
	if err != nil {
		v1.HandleError(c, v1.ErrRegisterCode, "微信登录失败", err)
		return
	}

	// 返回token
	v1.HandleSuccess(c, model.WeChatLoginResponse{
		SessionKey: sessionKey,
		OpenID:     openID,
	})
}

// WeChatDecrypt godoc
// @Summary 解密微信加密数据
// @Description 使用session_key解密微信前端传来的加密数据，获取用户信息
// @Tags 账户
// @Accept json
// @Produce json
// @Param request body model.WeChatEncryptedDataRequest true "微信加密数据请求"
// @Success 200 {object} v1.Response{data=model.WeChatDecryptResponse} "成功"
// @Failure 400 {object} v1.Response "参数错误"
// @Failure 500 {object} v1.Response "服务器内部错误"
// @Router /account/wechat/decrypt [post]
func (h *AccountHandler) WeChatDecrypt(c *gin.Context) {
	var req model.WeChatEncryptedDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		v1.HandleError(c, v1.ErrParamCode, "参数错误", err)
		return
	}

	// 调用服务层解密数据
	response, err := h.userService.DecryptWeChatData(c, &req)
	if err != nil {
		v1.HandleError(c, v1.ErrRegisterCode, "[0]解密数据失败: "+err.Error(), nil)
		return
	}

	v1.HandleSuccess(c, response)
}
