package v1

import (
	"errors"
	"fmt"
)

const (
	SuccessCode                    = 0
	ErrFormatPhoneCode             = 10001
	ErrPhoneCodeSendErrorCode      = 10002
	ErrParamFormatCode             = 10003
	ErrBadRequestCode              = 10004
	ErrUnauthorizedCode            = 10005
	ErrNotFoundCode                = 10006
	ErrTooManyRequestsCode         = 10007
	ErrInternalServerErrorCode     = 10008
	ErrWeChatLoginCallbackCode     = 10009
	ErrHeaderParamsFormatErrorCode = 10010
	ErrAuthorizationFormatCode     = 10011
	ErrAuthorizationCheckCode      = 10012
	ErrAuthorizationExpireCode     = 10013
	ErrGetExchangeCode             = 10039
	ErrProductInfoCode             = 10045
	ErrGetMyDevicesCode            = 10046
	ErrCodeFormatErrorCode         = 10048
	ErrPhoneAlreadyUseCode         = 10049
	ErrResetPasswordCode           = 10050
	ErrVerifyByPhoneAndCode        = 10051
	ErrGetProfileCode              = 10052
	ErrLoginCode                   = 10053
	ErrRegisterCode                = 10054
	ErrUpdateEggPriceCode          = 10055
	ErrParamCode                   = 10056
	ErrInviteCodeNotExistCode      = 10057
	ErrInputExploreCode            = 10058
	ErrUserProductBuyCode          = 10073
	ErrSignOutCode                 = 10081
	ErrUploadAvatarCode            = 10091
	ErrPasswordNotMatchErrorCode   = 10096
	ErrGetProfileSettingsCode      = 10097
	ErrUpdateProfileSettingsCode   = 10098
	ErrLogoutCode                  = 10099
	ErrGetTurntableCountCode       = 10106
	ErrGetTurntableResultCode      = 10107
	ErrGetPunchingNumCode          = 10108
	ErrPunchingCode                = 10109
	ErrActionNotFoundCode          = 10110
	ErrRoomFullCode                = 10111
	ErrUserInRoomCode              = 10112
	ErrRoomNotFoundCode            = 10113
	ErrRoomStatusErrorCode         = 10114
	ErrPhoneCodeErrorCode          = 10115
	ErrGenStringErrorCode          = 10116
	ErrWithdrawCode                = 10117
	ErrWithdrawBalanceCode         = 10118
	ErrOperateCode                 = 10119
	ErrForbiddenCode               = 10120
	ErrUnprocessableCode           = 10121
)

var (
	MsgSuccess               = "ok"
	MsgFormatPhoneError      = "手机号格式错误"
	MsgPhoneCodeSendError    = "验证码发送失败"
	MsgPhoneCodeError        = "验证码错误"
	MsgParamFormateError     = "参数格式错误"
	MsgBadRequest            = "错误的请求"
	MsgUnauthorized          = "未授权"
	MsgNotFound              = "未找到"
	MsgTooManyRequests       = "请求过于频繁"
	MsgInternalServerError   = "服务器内部错误"
	MsgWeChatLoginCallback   = "微信登录回调错误"
	MsgHeaderParamsFormat    = "请求头参数格式错误"
	MsgAuthorizationFormat   = "授权格式错误"
	MsgAuthorizationCheck    = "授权验证错误"
	MsgAuthorizationExpire   = "授权已过期"
	MsgCountryCodeError      = "国家代码错误"
	MsgPlunderError          = "抢夺错误"
	MsgGetPutdownError       = "获取放置信息错误"
	MsgGetRankingListError   = "获取排行榜错误"
	MsgAddMemberInviteError  = "添加成员邀请错误"
	MsgGetExchange           = "获取兑换信息错误"
	MsgGetMyBuyDevices       = "获取我购买的设备错误"
	MsgGetMyPublishDevices   = "获取我发布的设备错误"
	MsgProductInfo           = "产品信息错误"
	MsgGetMyDevices          = "获取我的设备错误"
	MsgPasswordNotMatch      = "密码不匹配"
	MsgCodeFormatError       = "验证码格式错误"
	MsgPhoneAlreadyUse       = "手机号已被使用"
	MsgResetPassword         = "重置密码错误"
	MsgVerifyByPhoneAndCode  = "手机号和验证码验证错误"
	MsgGetProfile            = "获取个人资料错误"
	MsgLogin                 = "登录错误"
	MsgRegister              = "注册错误"
	MsgInviteCodeNotExist    = "邀请码不存在"
	MsgInputExplore          = "输入探索错误"
	MsgSignOut               = "退出登录错误"
	MsgUploadAvatar          = "上传头像错误"
	MsgUpdateProfileSettings = "更新个人资料设置错误"
	MsgLogout                = "退出登录错误"
	MessageActionNotFound    = "操作未找到"
	MsgRoomFull              = "房间已满"
	MsgUserInRoom            = "用户已在房间中"
	MsgRoomNotFound          = "房间不存在"
	MsgRoomStatusError       = "房间状态错误"
	MsgGenStringError        = "生成用户ID错误"
	MsgBannerListError       = "获取Banner列表错误"
	MsgMonitorListError      = "获取Monitor列表失败"
	MsgNewsListError         = "获取新闻资讯列表失败"
)

var (
	ErrFormatPhoneError      = NewCustomError(ErrFormatPhoneCode, MsgFormatPhoneError, nil)
	ErrParamFormateError     = NewCustomError(ErrParamFormatCode, MsgParamFormateError, nil)
	ErrPhoneCodeSendError    = NewCustomError(ErrPhoneCodeSendErrorCode, MsgPhoneCodeSendError, nil)
	ErrPhoneCodeError        = NewCustomError(ErrPhoneCodeErrorCode, MsgPhoneCodeError, nil)
	ErrBadRequest            = NewCustomError(ErrBadRequestCode, MsgBadRequest, nil)
	ErrUnauthorized          = NewCustomError(ErrUnauthorizedCode, MsgUnauthorized, nil)
	ErrNotFound              = NewCustomError(ErrNotFoundCode, MsgNotFound, nil)
	ErrTooManyRequests       = NewCustomError(ErrTooManyRequestsCode, MsgTooManyRequests, nil)
	ErrInternalServerError   = NewCustomError(ErrInternalServerErrorCode, MsgInternalServerError, nil)
	ErrPasswordNotMatchError = NewCustomError(ErrPasswordNotMatchErrorCode, MsgPasswordNotMatch, nil)
	ErrCodeFormatError       = NewCustomError(ErrCodeFormatErrorCode, MsgCodeFormatError, nil)
	ErrPhoneAlreadyUseError  = NewCustomError(ErrPhoneAlreadyUseCode, MsgPhoneAlreadyUse, nil)
	ErrAuthorizationFormat   = NewCustomError(ErrAuthorizationFormatCode, MsgAuthorizationFormat, nil)
	ErrAuthorizationCheck    = NewCustomError(ErrAuthorizationCheckCode, MsgAuthorizationCheck, nil)
	ErrInviteCodeNotExist    = NewCustomError(ErrInviteCodeNotExistCode, MsgInviteCodeNotExist, nil)
	ErrInputExplore          = NewCustomError(ErrInputExploreCode, MsgInputExplore, nil)
	ErrAuthorizationExpire   = NewCustomError(ErrAuthorizationExpireCode, MsgAuthorizationExpire, nil)
	ErrActionNotFound        = NewCustomError(ErrActionNotFoundCode, MessageActionNotFound, nil)
	ErrGenStringError        = NewCustomError(ErrGenStringErrorCode, MsgGenStringError, nil)
)

// CustomError 定义一个自定义错误类型
type CustomError struct {
	error
	Code    int         // 错误码
	Message string      // 错误信息
	Data    interface{} // 附带的数据
}

// 实现 Error 接口
func (e *CustomError) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s, Data: %v", e.Code, e.Message, e.Data)
}

// 获取错误码的方法
func (e *CustomError) GetCode() int {
	return e.Code
}

// 获取错误信息的方法
func (e *CustomError) GetMessage() string {
	return e.Message
}

// 获取附带数据的方法
func (e *CustomError) GetData() interface{} {
	return e.Data
}

// 构造函数，用于创建新的 CustomError
func NewCustomError(code int, message string, data interface{}) error {
	return errors.New(message)
	// return &CustomError{
	// 	error: err,
	// 	Code:  code,
	// 	Message: message,
	// 	Data:    data,
	// }
}
