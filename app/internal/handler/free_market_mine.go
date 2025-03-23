package handler

import (
	"github.com/gin-gonic/gin"

	v1 "app/api/v1"
	"app/internal/model"
	"app/internal/service"
)

type FreeMarketMineHandler struct {
	*Handler
	freeMarketMineService service.FreeMarketMineService
}

func NewFreeMarketMineHandler(
	handler *Handler,
	freeMarketMineService service.FreeMarketMineService,
) *FreeMarketMineHandler {
	return &FreeMarketMineHandler{
		Handler:               handler,
		freeMarketMineService: freeMarketMineService,
	}
}

// GetUserEggsSummary godoc
// @Summary 获取用户鸡蛋销售情况
// @Description 获取用户已售出和未售出的鸡蛋情况
// @Tags 自由市场
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} model.FreeMarketMineResponse
// @Router /market/mine [get]
func (h *FreeMarketMineHandler) GetUserEggsSummary(c *gin.Context) {
	// Get data from service
	summary, err := h.freeMarketMineService.GetUserEggsSummary(c)
	if err != nil {
		v1.HandleError(c, v1.ErrRegisterCode, "获取鸡蛋销售数据失败", nil)
		return
	}

	v1.HandleSuccess(c, summary)
}

// UpdateEggPrice godoc
// @Summary 更新用户鸡蛋价格
// @Description 更新用户未售出的鸡蛋价格
// @Tags 自由市场
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body model.UpdateEggPriceRequest true "更新鸡蛋价格请求"
// @Success 200 {object} v1.Response
// @Router /market/update-price [post]
func (h *FreeMarketMineHandler) UpdateEggPrice(c *gin.Context) {
	// Create a request instance
	var req model.UpdateEggPriceRequest

	// Bind query parameters to the struct
	if err := c.ShouldBind(&req); err != nil {
		v1.HandleError(c, v1.ErrParamCode, "参数错误", err)
		return
	}

	// Use the bound parameters
	id := req.ID
	// Update price
	err := h.freeMarketMineService.UpdateEggPrice(c, req.Price, id)
	if err != nil {
		v1.HandleError(c, v1.ErrRegisterCode, "更新价格失败", nil)
		return
	}

	v1.HandleSuccess(c, nil)
}
