package middleware

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	v1 "app/api/v1"
	"app/internal/common"
	"app/pkg/config"
	"app/pkg/jwt"
	"app/pkg/log"
)

func SignMiddleware(logger *log.Logger, conf *viper.Viper, handler ...any) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 判断配置文件中debug是否为true
		//if config.ConfigInstance.GetBool("debug") {
		//	claims := &jwt.MyCustomClaims{
		//		UserId:   "2eD334",
		//		Level:    1,
		//		Nickname: "小王",
		//		Sign:     "",
		//		RegisteredClaims: v5jwt.RegisteredClaims{
		//			ExpiresAt: v5jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
		//			IssuedAt:  v5jwt.NewNumericDate(time.Now()),
		//			NotBefore: v5jwt.NewNumericDate(time.Now()),
		//			Issuer:    "2eD334",
		//			Subject:   "wechat",
		//			ID:        "2eD334",
		//			Audience:  []string{"小王"},
		//		},
		//	}
		//	ctx.Set("claims", claims)
		//	ctx.Next()
		//	return
		//}
		// 验证Authorization
		authStr := ctx.Request.Header.Get("Authorization")
		if !strings.Contains(authStr, "Bearer") || len(authStr) < 7 {
			v1.HandleError(ctx, v1.ErrAuthorizationFormatCode, v1.MsgAuthorizationFormat, nil)
			ctx.Abort()
			return
		}
		claims, err := jwt.NewJwt(conf).ParseToken(authStr[7:])
		if err != nil {
			v1.HandleError(ctx, v1.ErrAuthorizationCheckCode, v1.MsgAuthorizationCheck, nil)
			ctx.Abort()
			return
		}
		// 验证expireAt是否过期
		if claims.ExpiresAt.Unix() < time.Now().Unix() {
			v1.HandleError(ctx, v1.ErrAuthorizationExpireCode, v1.MsgAuthorizationExpire, nil)
			ctx.Abort()
			return
		}
		// 验证token有效性
		security := config.ConfigInstance.GetString("security.api_sign.app_security")
		securityMd5 := common.Md5Encrypt(claims.UserId + claims.Nickname + security)
		if claims.Sign != securityMd5 {
			v1.HandleError(ctx, v1.ErrAuthorizationCheckCode, v1.MsgAuthorizationCheck, nil)
			ctx.Abort()
			return
		}
		fmt.Println("claims", claims)
		ctx.Set("claims", claims)
		if len(handler) > 0 {
			ctx.Set("home_handler", handler[0])
		}
		recoveryLoggerFunc(ctx, logger)
		ctx.Next()
	}
}

func GetTokenInfo(ctx *gin.Context, token string) *jwt.MyCustomClaims {
	if !strings.Contains(token, "Bearer") || len(token) < 7 {
		v1.HandleError(ctx, v1.ErrAuthorizationFormatCode, v1.MsgAuthorizationFormat, nil)
		return nil
	}
	conf := config.ConfigInstance
	claims, err := jwt.NewJwt(conf).ParseToken(token[7:])
	if err != nil {
		v1.HandleError(ctx, v1.ErrAuthorizationCheckCode, v1.MsgAuthorizationCheck, nil)
		return nil
	}
	// 验证expireAt是否过期
	if claims.ExpiresAt.Unix() < time.Now().Unix() {
		v1.HandleError(ctx, v1.ErrAuthorizationExpireCode, v1.MsgAuthorizationExpire, nil)
		return nil
	}
	// 验证token有效性
	security := config.ConfigInstance.GetString("security.api_sign.app_security")
	securityMd5 := common.Md5Encrypt(claims.UserId + claims.Nickname + strconv.Itoa(int(claims.Role)) + security)
	if claims.Sign != securityMd5 {
		v1.HandleError(ctx, v1.ErrAuthorizationCheckCode, v1.MsgAuthorizationCheck, nil)
		return nil
	}
	ctx.Set("claims", claims)
	return claims
}
