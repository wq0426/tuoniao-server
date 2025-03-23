package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	v1 "app/api/v1"
	"app/internal/model"
	"app/internal/service"
)

type WithdrawOrderHandler struct {
	*Handler
	withdrawOrderService service.WithdrawOrderService
}

func NewWithdrawOrderHandler(
	handler *Handler,
	withdrawOrderService service.WithdrawOrderService,
) *WithdrawOrderHandler {
	return &WithdrawOrderHandler{
		Handler:              handler,
		withdrawOrderService: withdrawOrderService,
	}
}

// CreateWithdraw godoc
// @Summary 创建提现单
// @Description 创建提现申请
// @Tags 提现
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body model.CreateWithdrawRequest true "创建提现请求"
// @Success 200 {object} v1.Response
// @Failure 400 {object} v1.Response "参数错误"
// @Failure 401 {object} v1.Response "未授权"
// @Failure 403 {object} v1.Response "余额不足"
// @Failure 500 {object} v1.Response "服务器内部错误"
// @Router /withdraw/create [post]
func (h *WithdrawOrderHandler) CreateWithdraw(c *gin.Context) {
	var req model.CreateWithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		v1.HandleError(c, v1.ErrParamCode, "参数错误", err)
		return
	}

	// 调用服务层创建提现单
	err := h.withdrawOrderService.CreateWithdraw(c, req)
	if err != nil {
		if err.Error() == "余额不足" {
			v1.HandleError(c, v1.ErrOperateCode, err.Error(), nil)
			return
		}
		v1.HandleError(c, v1.ErrRegisterCode, "创建提现单失败", err)
		return
	}

	v1.HandleSuccess(c, nil)
}

// GetWithdrawList godoc
// @Summary 获取提现单列表
// @Description 获取提现单列表
// @Tags 提现
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param status query int false "提现状态(0:待审核;1:处理中;2:已完成;3:已拒绝)"
// @Param page query int false "页码，默认1"
// @Param page_size query int false "每页条数，默认10"
// @Success 200 {object} model.WithdrawListResponse
// @Failure 400 {object} v1.Response "参数错误"
// @Failure 401 {object} v1.Response "未授权"
// @Failure 500 {object} v1.Response "服务器内部错误"
// @Router /withdraw/list [get]
func (h *WithdrawOrderHandler) GetWithdrawList(c *gin.Context) {
	var req model.WithdrawQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		v1.HandleError(c, v1.ErrParamCode, "参数错误", err)
		return
	}

	// 处理状态参数
	if statusStr := c.Query("status"); statusStr != "" {
		statusInt, err := strconv.Atoi(statusStr)
		if err != nil {
			v1.HandleError(c, v1.ErrParamCode, "状态参数格式错误", err)
			return
		}
		statusUint8 := uint8(statusInt)
		req.Status = &statusUint8
	}

	// 调用服务层获取提现单列表
	response, err := h.withdrawOrderService.GetWithdrawList(c, req)
	if err != nil {
		v1.HandleError(c, v1.ErrRegisterCode, "获取提现单列表失败", err)
		return
	}

	v1.HandleSuccess(c, response)
}

// GetWithdrawDetail godoc
// @Summary 获取提现单详情
// @Description 获取提现单详情
// @Tags 提现
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param withdraw_id path uint64 true "提现单ID"
// @Success 200 {object} model.WithdrawDetailResponse
// @Failure 400 {object} v1.Response "参数错误"
// @Failure 401 {object} v1.Response "未授权"
// @Failure 404 {object} v1.Response "提现单不存在"
// @Failure 500 {object} v1.Response "服务器内部错误"
// @Router /withdraw/detail/{withdraw_id} [get]
func (h *WithdrawOrderHandler) GetWithdrawDetail(c *gin.Context) {
	withdrawIDStr := c.Param("withdraw_id")
	withdrawID, err := strconv.ParseUint(withdrawIDStr, 10, 64)
	if err != nil {
		v1.HandleError(c, v1.ErrParamCode, "参数错误", err)
		return
	}

	// 调用服务层获取提现单详情
	response, err := h.withdrawOrderService.GetWithdrawDetail(c, withdrawID)
	if err != nil {
		if err.Error() == "提现单不存在或无权限" {
			v1.HandleError(c, v1.ErrNotFoundCode, err.Error(), nil)
			return
		}
		v1.HandleError(c, v1.ErrRegisterCode, "获取提现单详情失败", err)
		return
	}

	v1.HandleSuccess(c, response)
}
