package v1

import (
	"app/internal/common"
)

// RegisterRequest godoc
// @Summary RegisterRequest 参数
// @Description RegisterRequest 参数描述
// @Tags 参数
// @Accept json
// @Produce json
// @Param phone body string true "手机号"
// @Param password body string true "密码"
// @Param re_password body string true "确认密码"
type RegisterRequest struct {
	Phone      string `form:"phone" json:"phone" binding:"required"`
	Password   string `form:"password" json:"password" binding:"required"`
	RePassword string `form:"re_password" json:"re_password" binding:"required"`
}

// 写一个validate方法来实现验证RegisterRequest每个字段
func (r *RegisterRequest) Validate() error {
	if !common.IsPhone(r.Phone) {
		return ErrFormatPhoneError
	}
	if len(r.Password) == 0 || len(r.RePassword) == 0 {
		return ErrParamFormateError
	}
	if r.Password != r.RePassword {
		return ErrPasswordNotMatchError
	}
	return nil
}

// CryptoRequest godoc
// @Summary CryptoRequest 参数
// @Description CryptoRequest 参数描述
// @Tags 参数
// @Accept json
// @Produce json
// @Param phone body string true "手机号"
type CryptoRequest struct {
	Phone string `json:"phone" binding:"required"`
}

func (r *CryptoRequest) Validate() error {
	if !common.IsPhone(r.Phone) {
		return ErrFormatPhoneError
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
// @Param code body string true "验证码"
type LoginRequest struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

// 写一个validate方法来实现验证LoginRequest每个字段
func (r *LoginRequest) Validate() error {
	if !common.IsPhone(r.Phone) {
		return ErrFormatPhoneError
	}
	if !common.IsValidCode(r.Code) {
		return ErrCodeFormatError
	}
	return nil
}

type LoginResponseData struct {
	AccessToken string `json:"accessToken"`
	Expired     int64  `json:"expired"`
}
type LoginResponse struct {
	Response
	Data LoginResponseData
}

type UpdateProfileRequest struct {
	Nickname string `json:"nickname" example:"alan"`
	Email    string `json:"email" binding:"required,email" example:"1234@gmail.com"`
}

type UserAsset struct {
	Points      int     `json:"points"`
	CouponCount int     `json:"coupon_count"`
	Balance     float64 `json:"balance"`
}

type GetProfileResponseData struct {
	UserId       string    `json:"userId"`
	Phone        string    `json:"phone"`
	Nickname     string    `json:"nickname" example:"alan"`
	Email        string    `json:"email"`
	Avatar       string    `json:"avatar"`
	Role         int       `json:"role"`
	AffCode      string    `json:"aff_code"`
	Point        int       `json:"point"`
	StepMaxPoint int       `json:"step_max_point"`
	UserAsset    UserAsset `json:"user_asset"`
	Gender       uint8     `json:"gender"`
	Birthday     string    `json:"birthday"`
	MemberLevel  int       `json:"member_level"`
	Address      string    `json:"address"`
}

type GetProfileResponse struct {
	Response
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
		return ErrFormatPhoneError
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
// @Param type body int true "类型"
type VerifyByPhoneAndCodeRequest struct {
	Code string `json:"code" binding:"required,len=6"`
	Type int    `json:"type" binding:"required"`
}

func (r *VerifyByPhoneAndCodeRequest) Validate() error {
	// 验证code
	if !common.IsValidCode(r.Code) {
		return ErrCodeFormatError
	}
	if r.Type == 0 {
		return ErrParamFormateError
	}
	return nil
}
