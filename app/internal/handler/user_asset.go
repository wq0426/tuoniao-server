package handler

import (
	"github.com/gin-gonic/gin"

	v1 "app/api/v1"
	"app/internal/model"
	"app/internal/service"
)

type UserAssetHandler struct {
	*Handler
	userAssetService service.UserAssetService
}

func NewUserAssetHandler(
	handler *Handler,
	userAssetService service.UserAssetService,
) *UserAssetHandler {
	return &UserAssetHandler{
		Handler:          handler,
		userAssetService: userAssetService,
	}
}

// GetUserAsset godoc
// @Summary 获取用户资产
// @Description 获取用户的积分等资产信息
// @Tags 用户资产
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} model.UserAssetResponse
// @Router /asset/info [get]
func (h *UserAssetHandler) GetUserAsset(c *gin.Context) {
	// Get user asset
	userAsset, err := h.userAssetService.GetUserAsset(c)
	if err != nil {
		v1.HandleError(c, v1.ErrRegisterCode, "获取用户资产失败", nil)
		return
	}

	v1.HandleSuccess(c, userAsset)
}

// RechargeBalance godoc
// @Summary 充值用户余额
// @Description 给用户账户充值指定金额
// @Tags 用户资产
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body model.RechargeRequest true "充值请求"
// @Success 200 {object} v1.Response
// @Router /asset/recharge [post]
func (h *UserAssetHandler) RechargeBalance(c *gin.Context) {
	// 解析请求体
	var req model.RechargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		v1.HandleError(c, v1.ErrParamCode, "参数错误", err)
		return
	}

	// 调用服务层进行充值
	err := h.userAssetService.RechargeBalance(c, req)
	if err != nil {
		v1.HandleError(c, v1.ErrRegisterCode, "充值失败", nil)
		return
	}

	v1.HandleSuccess(c, nil)
}

// WithdrawBalance godoc
// @Summary 从用户余额提现
// @Description 从用户账户提取指定金额
// @Tags 用户资产
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body model.WithdrawRequest true "提现请求"
// @Success 200 {object} v1.Response
// @Router /asset/withdraw [post]
func (h *UserAssetHandler) WithdrawBalance(c *gin.Context) {
	// 解析请求体
	var req model.WithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		v1.HandleError(c, v1.ErrParamCode, "参数错误", err)
		return
	}

	// 调用服务层进行提现
	err := h.userAssetService.WithdrawBalance(c, req)
	if err != nil {
		if err.Error() == "余额不足" {
			v1.HandleError(c, v1.ErrRegisterCode, "余额不足，无法完成提现", nil)
			return
		}
		v1.HandleError(c, v1.ErrRegisterCode, "提现失败", nil)
		return
	}

	v1.HandleSuccess(c, nil)
}

// GetBalanceRecords godoc
// @Summary 查询用户余额变动记录
// @Description 分页查询用户余额的所有变动记录
// @Tags 用户资产
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param page query int false "页码，默认1"
// @Param page_size query int false "每页条数，默认10"
// @Success 200 {object} model.BalanceRecordResponse
// @Router /asset/balance/records [get]
func (h *UserAssetHandler) GetBalanceRecords(c *gin.Context) {
	var req model.BalanceRecordQueryRequest
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

	// 调用服务层查询记录
	records, err := h.userAssetService.GetBalanceRecords(c, req)
	if err != nil {
		v1.HandleError(c, v1.ErrRegisterCode, "查询余额记录失败", nil)
		return
	}

	v1.HandleSuccess(c, records)
}

// GetWithdrawRecords godoc
// @Summary 查询用户提现记录
// @Description 分页查询用户的提现记录
// @Tags 用户资产
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param page query int false "页码，默认1"
// @Param page_size query int false "每页条数，默认10"
// @Success 200 {object} model.BalanceRecordResponse
// @Router /asset/withdraw/records [get]
func (h *UserAssetHandler) GetWithdrawRecords(c *gin.Context) {
	var req model.BalanceRecordQueryRequest
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

	// 调用服务层查询记录
	records, err := h.userAssetService.GetWithdrawRecords(c, req)
	if err != nil {
		v1.HandleError(c, v1.ErrRegisterCode, "查询提现记录失败", nil)
		return
	}

	v1.HandleSuccess(c, records)
}

// GetExchangeRecords godoc
// @Summary 查询用户积分兑换记录
// @Description 分页查询用户的积分兑换记录
// @Tags 用户资产
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param page query int false "页码，默认1"
// @Param page_size query int false "每页条数，默认10"
// @Success 200 {object} model.BalanceRecordResponse
// @Router /asset/exchange/records [get]
func (h *UserAssetHandler) GetExchangeRecords(c *gin.Context) {
	var req model.BalanceRecordQueryRequest
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

	// 调用服务层查询记录
	records, err := h.userAssetService.GetExchangeRecords(c, req)
	if err != nil {
		v1.HandleError(c, v1.ErrRegisterCode, "查询积分兑换记录失败", nil)
		return
	}

	v1.HandleSuccess(c, records)
}
