package handler

import (
	"github.com/gin-gonic/gin"

	v1 "app/api/v1"
	"app/internal/service"
)

type SmsHandler struct {
	*Handler
	smsService service.SmsService
}

func NewSmsHandler(
	handler *Handler,
	smsService service.SmsService,
) *SmsHandler {
	return &SmsHandler{
		Handler:    handler,
		smsService: smsService,
	}
}

// SendCode godoc
// @Summary 发送验证码
// @Description 发送验证码
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param request body v1.CryptoRequest true "备注 (type 1:注册 2:登录 3:找回密码 4:修改手机号 5:重新绑定 6:交友密码)"
// @Success 0 {object} v1.Response{data=bool}
// @Router /account/code [post]
func (h *SmsHandler) SendCode(ctx *gin.Context) {
	var req v1.CryptoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, v1.ErrPhoneCodeSendErrorCode, v1.MsgPhoneCodeSendError, nil)
		return
	}
	// 使用validate方法验证字段
	if err := req.Validate(); err != nil {
		v1.HandleError(ctx, v1.ErrParamFormatCode, v1.MsgParamFormateError, nil)
		return
	}
	if err := h.smsService.SendSmsCodeToPhone(ctx, &req); err != nil {
		v1.HandleError(ctx, v1.ErrPhoneCodeSendErrorCode, v1.MsgPhoneCodeSendError, nil)
		return
	}
	v1.HandleSuccess(ctx, true)
}
