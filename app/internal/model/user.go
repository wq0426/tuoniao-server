package model

import (
	v1 "app/api/v1"
)

// RegisterRequest godoc
// @Summary RegisterRequest 参数
// @Description RegisterRequest 参数描述
// @Tags 参数
// @Accept json
// @Produce json
// @Param phone body string true "手机号"
// @Param code body string true "验证码"
// @Param password body string true "密码"
// @Param invite_code body string false "邀请码"
type RegisterRequest struct {
	Phone      string `form:"phone" json:"phone" binding:"required"`
	Code       string `form:"code" json:"code" binding:"required,len=6"`
	Password   string `form:"password" json:"password" binding:"required"`
	InviteCode string `form:"invite_code" json:"invite_code"`
}

// 写一个validate方法来实现验证RegisterRequest每个字段
func (r *RegisterRequest) Validate() error {
	return nil
}

// CryptoRequest godoc
// @Summary CryptoRequest 参数
// @Description CryptoRequest 参数描述
// @Tags 参数
// @Accept json
// @Produce json
// @Param phone body string true "手机号"
// @Param type body int true "类型 (1:注册 2:登录 3:找回密码 4:修改手机号 5:重新绑定 6:交友密码)"
type CryptoRequest struct {
	Phone string `json:"phone" binding:"required"`
	Type  int    `json:"type" binding:"required" example:"1" description:"类型 (1:注册 2:登录 3:找回密码 4:修改手机号 5:重新绑定 6:交友密码)"`
}

func (r *CryptoRequest) Validate() error {
	if r.Type != 1 && r.Type != 2 && r.Type != 3 {
		return v1.ErrParamFormateError
	}
	return nil
}

// LoginRequest godoc
// @Summary LoginRequest 参数
// @Description LoginRequest 参数描述
// @Tags 参数
// @Accept json
// @Produce json
// @Param phone body string true "手机号"
// @Param password body string true "密码"
// @Param code body string true "验证码"
// @Param type body int true "类型 (1:验证码登录 2:密码登录)"
type LoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password"`
	Code     string `form:"code"`
	Type     int    `json:"type" binding:"required"`
}

// 写一个validate方法来实现验证LoginRequest每个字段
func (r *LoginRequest) Validate() error {
	if r.Type != 1 && r.Type != 2 {
		return v1.ErrParamFormateError
	}
	if r.Type == 1 && len(r.Code) == 0 {
		return v1.ErrParamFormateError
	}
	if r.Type == 2 && len(r.Password) == 0 {
		return v1.ErrParamFormateError
	}
	return nil
}

type LoginResponseData struct {
	BindInfo    *WeChatUserInfo `json:"bind_info"`
	UserInfo    *UserInfo       `json:"user_info"`
	AccessToken string          `json:"access_token"`
	Expired     int64           `json:"expired"`
}

type LoginResponse struct {
	v1.Response
	Data LoginResponseData
}

type UserInfo struct {
	UserId   string `json:"user_id"`
	Nickname string `json:"nickname"`
	Level    int    `json:"level"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
}

type GetProfileResponseData struct {
	UserId   string `json:"userId"`
	Phone    string `json:"phone"`
	Nickname string `json:"nickname" example:"alan"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Role     string `json:"role"`
	AffCode  string `json:"aff_code"`
}

type GetProfileResponse struct {
	v1.Response
	Data GetProfileResponseData
}

// ResetPasswordRequest godoc
// @Summary ResetPasswordRequest 参数
// @Description ResetPasswordRequest 参数描述
// @Tags 参数
// @Accept json
// @Produce json
// @Param password body string true "密码"
type ResetPasswordRequest struct {
	Password string `json:"password" binding:"required"`
}

func (r *ResetPasswordRequest) Validate() error {
	if len(r.Password) == 0 {
		return v1.ErrFormatPhoneError
	}
	return nil
}

// VerifyByPhoneAndCodeRequest godoc
// @Summary VerifyByPhoneAndCodeRequest 参数
// @Description VerifyByPhoneAndCodeRequest 参数描述
// @Tags 参数
// @Accept json
// @Produce json
// @Param code body string true "验证码"
// @Param type body int true "类型 (1:重置密码 2:修改手机号 3:交友密码未设置验证 4:交友密码已设置验证)"
type VerifyByPhoneAndCodeRequest struct {
	Code string `json:"code" binding:"required,len=6"`
	Type int    `json:"type" binding:"required"`
}

func (r *VerifyByPhoneAndCodeRequest) Validate() error {
	// 验证code
	if r.Type != 1 && r.Type != 2 && r.Type != 3 && r.Type != 4 {
		return v1.ErrParamFormateError
	}
	return nil
}

type WeChatUserInfo struct {
	OpenID     string `json:"openid"`
	Nickname   string `json:"nickname"`
	Sex        int    `json:"sex"`
	Province   string `json:"province"`
	City       string `json:"city"`
	Country    string `json:"country"`
	HeadImgURL string `json:"headimgurl"`
	UnionID    string `json:"unionid"`
}
