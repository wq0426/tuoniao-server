package handler

import (
	"app/pkg/log"

	"github.com/gin-gonic/gin"

	"app/internal/service"
)

type Handler struct {
	logger *log.Logger
}

func NewHandler(
	logger *log.Logger,
) *Handler {
	return &Handler{
		logger: logger,
	}
}
func TriggerEvent(ctx *gin.Context, module string, account any) {
	homeHandler, isExist := ctx.Get("home_handler")
	if isExist {
		handler := homeHandler.(*HomeHandler)
		handler.TriggerEvent(ctx, service.Transfer(module, account))
	}
}
