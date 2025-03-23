package service

import (
	"github.com/gin-gonic/gin"

	pb "app/internal/grpc"
	"app/internal/repository"
	"app/pkg/jwt"
	"app/pkg/log"
	"app/pkg/sid"
)

type Service struct {
	logger *log.Logger
	sid    *sid.Sid
	jwt    *jwt.JWT
	tm     repository.Transaction
}

func NewService(
	tm repository.Transaction,
	logger *log.Logger,
	sid *sid.Sid,
	jwt *jwt.JWT,
) *Service {
	return &Service{
		logger: logger,
		sid:    sid,
		jwt:    jwt,
		tm:     tm,
	}
}

func GetUserInfoFromCtx(ctx *gin.Context) *jwt.MyCustomClaims {
	v, exists := ctx.Get("claims")
	if !exists {
		return nil
	}
	return v.(*jwt.MyCustomClaims)
}

func GetUserIdFromCtx(ctx *gin.Context) string {
	v, exists := ctx.Get("claims")
	if !exists {
		return ""
	}
	return v.(*jwt.MyCustomClaims).UserId
}

func GetclaimsFromCtx(ctx *gin.Context) *jwt.MyCustomClaims {
	v, exists := ctx.Get("claims")
	if !exists {
		return nil
	}
	return v.(*jwt.MyCustomClaims)
}

func Transfer(module string, data any) *pb.PushMessageResponse {
	return nil
}
