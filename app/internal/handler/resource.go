package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	v1 "app/api/v1"
	"app/internal/service"
)

type ResourceHandler struct {
	*Handler
	resourceService service.ResourceService
}

func NewResourceHandler(
	handler *Handler,
	resourceService service.ResourceService,
) *ResourceHandler {
	return &ResourceHandler{
		Handler:         handler,
		resourceService: resourceService,
	}
}

func (h *ResourceHandler) GetResource(ctx *gin.Context) {
	url := ctx.Request.URL.Path
	contentType, data, err := h.resourceService.GetResource(ctx, url)
	if err != nil {
		v1.HandleError(ctx, 400, err.Error(), nil)
		return
	}
	h.logger.Debug("contentType:", contentType)
	// 设置响应头
	ctx.Header("Content-Type", contentType)
	// 返回图片的二进制数据
	ctx.Data(http.StatusOK, contentType, data)
}
