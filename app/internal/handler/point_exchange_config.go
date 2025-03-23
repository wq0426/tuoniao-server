package handler

import (
	"github.com/gin-gonic/gin"

	v1 "app/api/v1"
	"app/internal/model"
	"app/internal/service"
)

type PointExchangeConfigHandler struct {
	*Handler
	pointExchangeConfigService service.PointExchangeConfigService
}

func NewPointExchangeConfigHandler(
	handler *Handler,
	pointExchangeConfigService service.PointExchangeConfigService,
) *PointExchangeConfigHandler {
	return &PointExchangeConfigHandler{
		Handler:                    handler,
		pointExchangeConfigService: pointExchangeConfigService,
	}
}

// GetPointExchangeConfigList godoc
// @Summary 获取积分兑换配置列表
// @Description 获取所有积分兑换配置项
// @Tags 积分兑换
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} model.PointExchangeConfigListResponse
// @Router /point/exchange/list [get]
func (h *PointExchangeConfigHandler) GetPointExchangeConfigList(c *gin.Context) {
	configs, err := h.pointExchangeConfigService.GetPointExchangeConfigList(c)
	if err != nil {
		v1.HandleError(c, v1.ErrRegisterCode, "获取积分兑换配置列表失败", nil)
		return
	}

	v1.HandleSuccess(c, configs)
}

// ExchangePoints godoc
// @Summary 兑换积分
// @Description 使用积分兑换优惠券或兑换券
// @Tags 积分兑换
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body model.PointExchangeRequest true "兑换请求"
// @Success 200 {object} v1.Response
// @Router /point/exchange/exchange [post]
func (h *PointExchangeConfigHandler) ExchangePoints(c *gin.Context) {
	// 解析请求体
	var req model.PointExchangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		v1.HandleError(c, v1.ErrParamCode, "参数错误", err)
		return
	}

	// 调用服务层进行兑换
	err := h.pointExchangeConfigService.ExchangePoints(c, req)
	if err != nil {
		if err.Error() == "积分不足" {
			v1.HandleError(c, v1.ErrRegisterCode, "积分不足，无法完成兑换", nil)
			return
		}
		v1.HandleError(c, v1.ErrRegisterCode, "兑换失败", nil)
		return
	}

	v1.HandleSuccess(c, nil)
}
