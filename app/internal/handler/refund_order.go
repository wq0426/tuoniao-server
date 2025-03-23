package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	v1 "app/api/v1"
	"app/internal/model"
	"app/internal/service"
)

type RefundOrderHandler struct {
	*Handler
	refundOrderService service.RefundOrderService
}

func NewRefundOrderHandler(
	handler *Handler,
	refundOrderService service.RefundOrderService,
) *RefundOrderHandler {
	return &RefundOrderHandler{
		Handler:            handler,
		refundOrderService: refundOrderService,
	}
}

// CreateRefund godoc
// @Summary 创建退款单
// @Description 创建退款申请
// @Tags 退款
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body model.RefundItemRequest true "创建退款请求"
// @Success 200 {object} v1.Response
// @Failure 400 {object} v1.Response "参数错误"
// @Failure 401 {object} v1.Response "未授权"
// @Failure 403 {object} v1.Response "权限不足"
// @Failure 500 {object} v1.Response "服务器内部错误"
// @Router /refund/create [post]
func (h *RefundOrderHandler) CreateRefund(c *gin.Context) {
	var req model.RefundItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		v1.HandleError(c, v1.ErrParamCode, "参数错误", err)
		return
	}

	// 调用服务层创建退款单
	err := h.refundOrderService.CreateRefund(c, req)
	if err != nil {
		v1.HandleError(c, v1.ErrRegisterCode, "创建退款单失败: "+err.Error(), nil)
		return
	}

	v1.HandleSuccess(c, nil)
}

// GetRefundList godoc
// @Summary 获取退款单列表
// @Description 获取用户退款单列表
// @Tags 退款
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param status query int false "退款状态：0-待处理，1-退款中，2-已退款，3-已拒绝；不传则查询所有状态"
// @Param page query int false "页码，默认1"
// @Param page_size query int false "每页条数，默认10"
// @Success 200 {object} model.RefundListResponse
// @Failure 401 {object} v1.Response "未授权"
// @Failure 500 {object} v1.Response "服务器内部错误"
// @Router /refund/list [get]
func (h *RefundOrderHandler) GetRefundList(c *gin.Context) {
	var req model.RefundQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		v1.HandleError(c, v1.ErrParamCode, "参数错误", err)
		return
	}

	// 设置默认值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	// 处理状态参数
	if statusStr := c.Query("status"); statusStr != "" {
		statusInt, err := strconv.Atoi(statusStr)
		if err == nil {
			statusUint8 := uint8(statusInt)
			req.Status = &statusUint8
		}
	}

	// 调用服务层获取退款单列表
	response, err := h.refundOrderService.GetRefundList(c, req)
	if err != nil {
		v1.HandleError(c, v1.ErrRegisterCode, "获取退款单列表失败", nil)
		return
	}

	v1.HandleSuccess(c, response)
}

// GetRefundDetail godoc
// @Summary 获取退款单详情
// @Description 获取退款单详情
// @Tags 退款
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param refund_id path uint64 true "退款单ID"
// @Success 200 {object} model.RefundDetailResponse
// @Failure 400 {object} v1.Response "参数错误"
// @Failure 401 {object} v1.Response "未授权"
// @Failure 404 {object} v1.Response "退款单不存在"
// @Failure 500 {object} v1.Response "服务器内部错误"
// @Router /refund/detail/{refund_id} [get]
func (h *RefundOrderHandler) GetRefundDetail(c *gin.Context) {
	refundIDStr := c.Param("refund_id")
	refundID, err := strconv.ParseUint(refundIDStr, 10, 64)
	if err != nil {
		v1.HandleError(c, v1.ErrParamCode, "参数错误", err)
		return
	}

	// 调用服务层获取退款单详情
	response, err := h.refundOrderService.GetRefundDetail(c, refundID)
	if err != nil {
		if err.Error() == "退款单不存在或无权限" {
			v1.HandleError(c, v1.ErrNotFoundCode, err.Error(), nil)
			return
		}
		v1.HandleError(c, v1.ErrRegisterCode, "获取退款单详情失败", err)
		return
	}

	v1.HandleSuccess(c, response)
}

// CancelRefund godoc
// @Summary 撤销退款申请
// @Description 撤销处理中的退款申请，恢复原订单状态
// @Tags 退款
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body model.CancelRefundRequest true "退款撤销信息"
// @Success 200 {object} v1.Response
// @Failure 400 {object} v1.Response "参数错误"
// @Failure 401 {object} v1.Response "未授权"
// @Failure 403 {object} v1.Response "退款单不存在或无权限"
// @Failure 422 {object} v1.Response "当前退款单状态不允许撤销"
// @Failure 500 {object} v1.Response "服务器内部错误"
// @Router /refund/cancel [post]
func (h *RefundOrderHandler) CancelRefund(c *gin.Context) {
	var req model.CancelRefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		v1.HandleError(c, v1.ErrParamCode, "参数错误", err)
		return
	}

	// 调用服务层撤销退款
	err := h.refundOrderService.CancelRefund(c, req.RefundID)
	if err != nil {
		if err.Error() == "退款单不存在或无权限" {
			v1.HandleError(c, v1.ErrForbiddenCode, err.Error(), nil)
			return
		} else if err.Error() == "当前退款单状态不允许撤销" {
			v1.HandleError(c, v1.ErrUnprocessableCode, err.Error(), nil)
			return
		}
		v1.HandleError(c, v1.ErrRegisterCode, "撤销退款失败", err)
		return
	}

	v1.HandleSuccess(c, nil)
}

// DeleteRefund godoc
// @Summary 删除退款记录
// @Description 删除用户的退款记录
// @Tags 退款
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param refund_id path uint64 true "退款单ID"
// @Success 200 {object} v1.Response
// @Failure 400 {object} v1.Response "参数错误"
// @Failure 401 {object} v1.Response "未授权"
// @Failure 403 {object} v1.Response "退款单不存在或无权限"
// @Failure 500 {object} v1.Response "服务器内部错误"
// @Router /refund/delete/{refund_id} [delete]
func (h *RefundOrderHandler) DeleteRefund(c *gin.Context) {
	refundIDStr := c.Param("refund_id")
	refundID, err := strconv.ParseUint(refundIDStr, 10, 64)
	if err != nil {
		v1.HandleError(c, v1.ErrParamCode, "参数错误", err)
		return
	}

	// 调用服务层删除退款记录
	err = h.refundOrderService.DeleteRefund(c, refundID)
	if err != nil {
		if err.Error() == "退款单不存在或无权限" {
			v1.HandleError(c, v1.ErrForbiddenCode, err.Error(), nil)
			return
		}
		v1.HandleError(c, v1.ErrRegisterCode, "删除退款记录失败", err)
		return
	}

	v1.HandleSuccess(c, nil)
}
