package handler

import (
	"github.com/gin-gonic/gin"

	v1 "app/api/v1"
	"app/internal/service"
)

type BannerHandler struct {
	*Handler
	bannerService service.BannerService
}

func NewBannerHandler(
	handler *Handler,
	bannerService service.BannerService,
) *BannerHandler {
	return &BannerHandler{
		Handler:       handler,
		bannerService: bannerService,
	}
}

// GetBannerList godoc
// @Summary 获取首页Banner列表
// @Description 获取首页Banner列表
// @Tags 首页模块
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} []model.BannerResponse
// @Router /banner/list [get]
func (h *BannerHandler) GetBannerList(c *gin.Context) {
	banners, err := h.bannerService.GetBannerList(c)
	if err != nil {
		v1.HandleError(c, v1.ErrRegisterCode, v1.MsgBannerListError, nil)
		return
	}
	v1.HandleSuccess(c, banners)
}
