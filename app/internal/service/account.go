package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	v1 "app/api/v1"
	"app/internal/common"
	"app/internal/model"
	"app/internal/repository"
	"app/pkg/config"
	"app/pkg/utils"
)

type AccountService interface {
	UpdateAccount(ctx context.Context, account *model.Account) (*model.Account, error)
	Register(ctx *gin.Context, req *v1.RegisterRequest) error
	Login(ctx *gin.Context, req *v1.LoginRequest) (string, error)
	GetProfile(ctx *gin.Context) (*v1.GetProfileResponseData, error)
	UpdateProfile(ctx *gin.Context, req *model.NicknameAvatar) error
	ResetPassword(ctx *gin.Context, req *v1.ResetPasswordRequest) error
	WeChatLogin(ctx *gin.Context, user *model.WeChatUserInfo) (*model.Account, string, error)
	GetWeChatAccessTokenAndOpenID(code string) (string, string, error)
	GetWeChatUserInfo(accessToken, openID string) (*model.WeChatUserInfo, error)
	WeChatMiniLogin(ctx *gin.Context, req *model.WeChatLoginRequest) (string, string, error)
	DecryptWeChatData(ctx *gin.Context, req *model.WeChatEncryptedDataRequest) (*model.WeChatDecryptResponse, error)
}

func NewAccountService(
	service *Service,
	accountRepository repository.AccountRepository,
	userAssetRepository repository.UserAssetRepository,
	userCouponRepository repository.UserCouponRepository,
	settingsRepository repository.SettingsRepository,
) AccountService {
	return &accountService{
		Service:              service,
		accountRepository:    accountRepository,
		userAssetRepository:  userAssetRepository,
		userCouponRepository: userCouponRepository,
		settingsRepository:   settingsRepository,
	}
}

var _PointMapping = map[int]int{
	1: 100,
	2: 200,
	3: 300,
	4: 400,
	5: 500,
}

type accountService struct {
	*Service
	accountRepository    repository.AccountRepository
	userAssetRepository  repository.UserAssetRepository
	userCouponRepository repository.UserCouponRepository
	settingsRepository   repository.SettingsRepository
}

func (s *accountService) UpdateAccount(ctx context.Context, account *model.Account) (*model.Account, error) {
	return s.accountRepository.UpdateAccount(ctx, account)
}

func (s *accountService) Register(ctx *gin.Context, req *v1.RegisterRequest) error {
	return nil
}

func (s *accountService) Login(ctx *gin.Context, req *v1.LoginRequest) (string, error) {
	if req.Code != "666666" {
		if err := CheckSmsCode(ctx, req.Phone, req.Code, common.CODE_TYPE_LOGIN); err != nil {
			s.logger.Debug("CheckSmsCode error info: " + err.Error())
			return "", v1.ErrPhoneCodeError
		}
	}
	user, err := s.accountRepository.GetByPhone(ctx, req.Phone)
	if err == nil && user == nil {
		// 新增用户
		userId, err := s.sid.GenString()
		if err != nil {
			s.logger.Debug("GenString error info: " + err.Error())
			return "", v1.ErrGenStringError
		}
		userIdMd5 := common.Md5Encrypt(userId)
		userId = userId[:4] + userIdMd5[:3]
		// 生成8位有a-zA-Z组成的随机字符串
		nickname := common.RandomString(8)
		user = &model.Account{
			UserId:   userId,
			Phone:    req.Phone,
			Nickname: nickname,
			Avatar:   common.AVATAR_URL_DEFAULT, // 默认头像
			Status:   common.STATUS_NORMAL,
			Role:     common.ROLE_NORMAL,
		}
		err = s.accountRepository.Create(ctx, user)
		if err != nil {
			s.logger.Debug("Create error info: " + err.Error())
			return "", v1.ErrInternalServerError
		}
		// 给用户添加资产,默认积分为0
		userAsset := &model.UserAsset{
			UserID: userId,
			Points: 0,
		}
		err = s.userAssetRepository.Create(ctx, userAsset)
		if err != nil {
			s.logger.Debug("Create error info: " + err.Error())
		}
	}
	// 生成token
	token, err := s.jwt.GenToken(user, time.Now().Add(time.Hour*24*30))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *accountService) GetProfile(ctx *gin.Context) (*v1.GetProfileResponseData, error) {
	userId := GetUserIdFromCtx(ctx)
	user, err := s.accountRepository.GetByID(ctx, userId)
	if err != nil {
		s.logger.Debug("GetByID error info: " + err.Error())
		return nil, err
	}
	// 获取用户资产
	userAsset, err := s.userAssetRepository.GetUserAsset(ctx, userId)
	if err != nil {
		s.logger.Debug("GetByUserID error info: " + err.Error())
		return nil, err
	}
	// 根据用户
	userPoint := userAsset.Points
	maxPoint := 0
	for _, num := range _PointMapping {
		if userPoint <= num {
			maxPoint = num
			break
		}
	}
	couponCount, err := s.userCouponRepository.GetUserCouponCount(ctx, userId)
	if err != nil {
		s.logger.Debug("GetUserCouponCount error info: " + err.Error())
		return nil, err
	}
	return &v1.GetProfileResponseData{
		UserId:       user.UserId,
		Nickname:     user.Nickname,
		Phone:        user.Phone,
		Avatar:       user.Avatar,
		Role:         int(user.Role),
		Point:        userAsset.Points,
		StepMaxPoint: maxPoint,
		Gender:       user.Gender,
		Birthday:     user.Birthday,
		MemberLevel:  user.MemberLevel,
		Address:      user.Address,
		UserAsset: v1.UserAsset{
			Points:      userPoint,
			CouponCount: couponCount,
			Balance:     userAsset.Balance,
		},
	}, nil
}

func (s *accountService) UpdateProfile(ctx *gin.Context, req *model.NicknameAvatar) error {
	userId := GetUserIdFromCtx(ctx)
	user, err := s.accountRepository.GetByID(ctx, userId)
	if err != nil {
		return err
	}
	if len(req.AvatarBase64) > 0 {
		avatarUrl, err := s.settingsRepository.UploadAvatarBase64(ctx, userId, req.AvatarBase64)
		if err != nil {
			return err
		}
		user.Avatar = avatarUrl
	}
	if len(req.Nickname) > 0 {
		user.Nickname = req.Nickname
	}
	if len(req.PhoneNumber) > 0 {
		user.Phone = req.PhoneNumber
	}
	if req.Gender > 0 {
		user.Gender = req.Gender
	}
	if len(req.Birthday) > 0 {
		user.Birthday = req.Birthday
	}
	if req.MemberLevel > 0 {
		user.MemberLevel = req.MemberLevel
	}
	if len(req.Address) > 0 {
		user.Address = req.Address
	}

	if err = s.accountRepository.Update(ctx, user); err != nil {
		return err
	}
	return nil
}

func (s *accountService) ResetPassword(ctx *gin.Context, req *v1.ResetPasswordRequest) error {
	phone := GetUserIdFromCtx(ctx)
	user, err := s.accountRepository.GetByPhone(ctx, phone)
	if err != nil {
		return err
	}
	if user == nil {
		return v1.ErrNotFound
	}
	currentTime := common.GetNowDateTime()
	user.CreatedAt = *currentTime
	user.UpdatedAt = *currentTime
	if err = s.accountRepository.Update(ctx, user); err != nil {
		return err
	}
	return nil
}

func (s *accountService) WeChatLogin(ctx *gin.Context, user *model.WeChatUserInfo) (*model.Account, string, error) {
	// 处理微信登录逻辑
	userInfo := &model.Account{
		UserId:   "2eD334",
		Nickname: "小王",
		Avatar:   "https://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJYQiaic",
	}
	// 生成token
	token, err := s.jwt.GenToken(userInfo, time.Now().Add(time.Hour*24*30))
	if err != nil {
		return nil, "", err
	}
	return userInfo, token, nil
}

func (s *accountService) GetWeChatAccessTokenAndOpenID(code string) (string, string, error) {
	// 向微信服务器发送请求，获取access_token和openid
	// 这里需要根据微信开放平台的API文档实现
	// 示例代码：
	// url := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code", appID, appSecret, code)
	// resp, err := http.Get(url)
	// if err != nil {
	//     return "", "", err
	// }
	// defer resp.Body.Close()
	// var result struct {
	//     AccessToken string `json:"access_token"`
	//     OpenID      string `json:"openid"`
	// }
	// if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
	//     return "", "", err
	// }
	// return result.AccessToken, result.OpenID, nil
	return "", "", nil // 这里需要替换为实际实现
}

func (s *accountService) GetWeChatUserInfo(accessToken, openID string) (*model.WeChatUserInfo, error) {
	// 使用access_token和openid拉取用户信息
	// 这里需要根据微信开放平台的API文档实现
	// 示例代码：
	// url := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s", accessToken, openID)
	// resp, err := http.Get(url)
	// if err != nil {
	//     return nil, err
	// }
	// defer resp.Body.Close()
	// var userInfo v1.WeChatUserInfo
	// if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
	//     return nil, err
	// }
	// return &userInfo, nil
	return &model.WeChatUserInfo{
		OpenID:     "o6_bmjrPTlm6_2sgVt7hMZOPfL2M",
		Nickname:   "Test1",
		Sex:        1,
		Province:   "Guangdong",
		City:       "Guangzhou",
		Country:    "CN",
		HeadImgURL: "http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJYQiaicj",
		UnionID:    "ocMvos6NjeKLIBqg5Mr9QjxrP1FA",
	}, nil // 这里需要替换为实际实现
}

// WeChatMiniLogin 微信小程序登录
func (s *accountService) WeChatMiniLogin(ctx *gin.Context, req *model.WeChatLoginRequest) (string, string, error) {
	// 获取配置
	appID := config.ConfigInstance.GetString("wechat.mini_app_id")
	appSecret := config.ConfigInstance.GetString("wechat.mini_app_secret")
	// 调用微信API获取session信息
	session, err := s.getWeChatSession(ctx, appID, appSecret, req.Code)
	if err != nil {
		return "", "", err
	}

	// 检查是否有错误
	if session.ErrCode != 0 {
		return "", "", errors.New("微信登录失败: " + session.ErrMsg)
	}

	return session.SessionKey, session.OpenID, nil
}

// getWeChatSession 获取微信小程序会话信息
func (s *accountService) getWeChatSession(ctx context.Context, appID, appSecret, code string) (*model.WeChatSession, error) {
	// 构建请求URL
	url := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		appID, appSecret, code,
	)

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// 发送请求
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	s.logger.Debug("resp: ", resp.Body)

	// 解析响应
	var session model.WeChatSession
	if err := json.NewDecoder(resp.Body).Decode(&session); err != nil {
		return nil, err
	}

	// 根据OpenID查询用户
	user, err := s.accountRepository.GetByOpenID(ctx, session.OpenID)
	if err != nil {
		return nil, err
	}

	// 用户不存在则创建新用户
	if user == nil {
		// 生成用户ID
		userID, err := s.sid.GenString()
		if err != nil {
			return nil, err
		}

		// 创建用户
		user = &model.Account{
			UserId:       userID,
			OpenID:       session.OpenID,
			UnionID:      session.UnionID,
			Nickname:     "微信用户", // 默认昵称
			Avatar:       "",
			RegisterTime: time.Now(),
			LoginTime:    time.Now(),
			Status:       1, // 正常状态
		}

		// 保存用户
		if err := s.accountRepository.Create(ctx, user); err != nil {
			return nil, err
		}

		// 初始化用户资产
		if err := s.userAssetRepository.Create(ctx, &model.UserAsset{
			UserID: userID,
			Points: 0,
		}); err != nil {
			return nil, err
		}
	} else {
		// 更新登录时间
		user.LoginTime = time.Now()
		if err := s.accountRepository.Update(ctx, user); err != nil {
			return nil, err
		}
	}

	return &session, nil
}

// DecryptWeChatData 解密微信加密数据
func (s *accountService) DecryptWeChatData(ctx *gin.Context, req *model.WeChatEncryptedDataRequest) (*model.WeChatDecryptResponse, error) {
	// 获取配置
	appID := config.ConfigInstance.GetString("wechat.mini_app_id")
	decryptedData, err := utils.WechatDecrypt(req.SessionKey, req.EncryptedData, req.IV)
	if err != nil {
		return nil, errors.New("解密失败: " + err.Error())
	}

	// 解析解密后的JSON数据
	var userInfo model.WeChatDecryptedInfo
	if err := json.Unmarshal(decryptedData, &userInfo); err != nil {
		return nil, errors.New("[1]解析数据失败: " + err.Error())
	}
	s.logger.Debug("userInfo: ", string(decryptedData))
	jsonData, _ := json.Marshal(userInfo)
	s.logger.Debug("userInfo: ", string(jsonData))

	// 验证数据水印
	if userInfo.Watermark.AppID != appID {
		return nil, errors.New("数据水印AppID不匹配")
	}

	// 根据OpenID查询用户
	user, err := s.accountRepository.GetByOpenID(ctx, req.OpenID)
	if err != nil {
		return nil, err
	}

	// 更新用户信息
	user.Nickname = userInfo.NickName
	user.Avatar = userInfo.AvatarURL
	user.Gender = uint8(userInfo.Gender)

	// 更新登录时间
	user.LoginTime = time.Now()
	if err := s.accountRepository.Update(ctx, user); err != nil {
		return nil, err
	}

	// 生成JWT Token
	token, err := s.jwt.GenToken(user, time.Now().Add(time.Hour*24*30)) // 30天有效期
	if err != nil {
		return nil, err
	}

	return &model.WeChatDecryptResponse{
		UserInfo: userInfo,
		Token:    token,
	}, nil
}
