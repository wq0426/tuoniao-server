package handler

import (
	"github.com/gin-gonic/gin"

	v1 "app/api/v1"
	"app/internal/service"
)

type MonitorHandler struct {
	*Handler
	monitorService service.MonitorService
}

func NewMonitorHandler(
	handler *Handler,
	monitorService service.MonitorService,
) *MonitorHandler {
	return &MonitorHandler{
		Handler:        handler,
		monitorService: monitorService,
	}
}

// GetMonitorList godoc
// @Summary 获取Monitor列表
// @Description 获取Monitor列表
// @Tags Monitor模块
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} []model.MonitorResponse
// @Router /monitor/list [get]
func (h *MonitorHandler) GetMonitorList(c *gin.Context) {
	monitors, err := h.monitorService.GetMonitorList(c)
	if err != nil {
		v1.HandleError(c, v1.ErrRegisterCode, v1.MsgMonitorListError, nil)
		return
	}
	v1.HandleSuccess(c, monitors)
}
