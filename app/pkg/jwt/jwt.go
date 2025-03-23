package jwt

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"

	"app/internal/common"
	"app/internal/model"
	"app/pkg/config"
)

type JWT struct {
	key []byte
}

type MyCustomClaims struct {
	UserId   string
	Role     int64
	Nickname string
	SiteNo   int
	Sign     string
	jwt.RegisteredClaims
}

func NewJwt(conf *viper.Viper) *JWT {
	return &JWT{key: []byte(conf.GetString("security.jwt.key"))}
}

func (j *JWT) GenToken(user *model.Account, expiresAt time.Time) (string, error) {
	security := config.ConfigInstance.GetString("security.api_sign.app_security")
	// 对security做md5加密
	securityMd5 := common.Md5Encrypt(user.UserId + user.Nickname + security)
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256, MyCustomClaims{
			UserId:   user.UserId,
			Nickname: user.Nickname,
			Role:     user.Role,
			Sign:     securityMd5,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expiresAt),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				Issuer:    user.UserId,
				Subject:   "wechat",
				ID:        user.UserId,
				Audience: []string{
					user.Nickname,
				},
			},
		},
	)
	tokenString, err := token.SignedString(j.key)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (j *JWT) ParseToken(tokenString string) (*MyCustomClaims, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	if strings.TrimSpace(tokenString) == "" {
		return nil, errors.New("token is empty")
	}
	token, err := jwt.ParseWithClaims(
		tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return j.key, nil
		},
	)
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
