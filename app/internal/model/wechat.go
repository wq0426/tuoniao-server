// 如果文件不存在需创建
package model

// WeChatLoginRequest 微信登录请求
type WeChatLoginRequest struct {
	Code string `json:"code" binding:"required"` // 微信临时登录凭证
}

// WeChatSession 微信登录会话信息
type WeChatSession struct {
	OpenID     string `json:"openid"`      // 用户唯一标识
	SessionKey string `json:"session_key"` // 会话密钥
	UnionID    string `json:"unionid"`     // 用户在开放平台的唯一标识符
	ErrCode    int    `json:"errcode"`     // 错误码
	ErrMsg     string `json:"errmsg"`      // 错误信息
}

// WeChatLoginResponse 微信登录响应
type WeChatLoginResponse struct {
	OpenID     string `json:"openid"`      // 用户唯一标识
	Token      string `json:"token"`       // 平台token
	SessionKey string `json:"session_key"` // 会话密钥
}

// WeChatEncryptedDataRequest 微信加密数据解密请求
type WeChatEncryptedDataRequest struct {
	SessionKey    string `json:"sessionKey"`                       // 会话密钥
	OpenID        string `json:"openid"`                           // 用户唯一标识
	EncryptedData string `json:"encryptedData" binding:"required"` // 包括敏感数据在内的完整用户信息的加密数据
	IV            string `json:"iv" binding:"required"`            // 加密算法的初始向量
}

// WeChatDecryptedInfo 解密后的微信用户信息
type WeChatDecryptedInfo struct {
	NickName  string `json:"nickName"`  // 用户昵称
	Gender    int    `json:"gender"`    // 用户性别
	Language  string `json:"language"`  // 用户语言
	City      string `json:"city"`      // 用户所在城市
	Province  string `json:"province"`  // 用户所在省份
	Country   string `json:"country"`   // 用户所在国家
	AvatarURL string `json:"avatarUrl"` // 用户头像CountryCode     string `json:"countryCode,omitempty"`     // 区号
	Watermark struct {
		AppID     string `json:"appid"`     // 小程序appid
		Timestamp int64  `json:"timestamp"` // 时间戳
	} `json:"watermark"` // 数据水印
}

// WeChatDecryptResponse 微信数据解密响应
type WeChatDecryptResponse struct {
	UserInfo WeChatDecryptedInfo `json:"user_info"` // 用户信息
	Token    string              `json:"token"`     // 平台token，如果需要登录
}
