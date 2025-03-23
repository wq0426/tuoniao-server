package handler

import (
	"github.com/gin-gonic/gin"

	v1 "app/api/v1"
	"app/internal/model"
	"app/internal/service"
)

type UserAddressHandler struct {
	*Handler
	userAddressService service.UserAddressService
}

func NewUserAddressHandler(
	handler *Handler,
	userAddressService service.UserAddressService,
) *UserAddressHandler {
	return &UserAddressHandler{
		Handler:            handler,
		userAddressService: userAddressService,
	}
}

// AddAddress godoc
// @Summary 新增用户地址
// @Description 新增用户收货地址
// @Tags 用户地址
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body model.AddAddressRequest true "新增地址请求"
// @Success 200 {object} v1.Response
// @Router /address/add [post]
func (h *UserAddressHandler) AddAddress(c *gin.Context) {
	// Parse request body
	var req model.AddAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		v1.HandleError(c, v1.ErrParamCode, "参数错误", err)
		return
	}

	// Add address
	err := h.userAddressService.AddAddress(c, req)
	if err != nil {
		v1.HandleError(c, v1.ErrRegisterCode, "添加地址失败", nil)
		return
	}

	v1.HandleSuccess(c, nil)
}

// UpdateAddress godoc
// @Summary 更新用户地址
// @Description 更新用户收货地址
// @Tags 用户地址
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body model.UpdateAddressRequest true "更新地址请求"
// @Success 200 {object} v1.Response
// @Router /address/update [post]
func (h *UserAddressHandler) UpdateAddress(c *gin.Context) {
	// Parse request body
	var req model.UpdateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		v1.HandleError(c, v1.ErrParamCode, "参数错误", err)
		return
	}

	// Update address
	err := h.userAddressService.UpdateAddress(c, req)
	if err != nil {
		v1.HandleError(c, v1.ErrRegisterCode, "更新地址失败", nil)
		return
	}

	v1.HandleSuccess(c, nil)
}

// GetUserAddresses godoc
// @Summary 获取用户地址列表
// @Description 获取用户所有收货地址
// @Tags 用户地址
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {array} model.UserAddress
// @Router /address/list [get]
func (h *UserAddressHandler) GetUserAddresses(c *gin.Context) {
	// Get addresses
	addresses, err := h.userAddressService.GetUserAddresses(c)
	if err != nil {
		v1.HandleError(c, v1.ErrRegisterCode, "获取地址列表失败", nil)
		return
	}

	v1.HandleSuccess(c, addresses)
}

// DeleteAddress godoc
// @Summary 删除用户地址
// @Description 删除用户收货地址
// @Tags 用户地址
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body model.DeleteAddressRequest true "删除地址请求"
// @Success 200 {object} v1.Response
// @Router /address/delete [post]
func (h *UserAddressHandler) DeleteAddress(c *gin.Context) {
	// Parse request body
	var delRequest model.DeleteAddressRequest
	if err := c.ShouldBindJSON(&delRequest); err != nil {
		v1.HandleError(c, v1.ErrParamCode, "参数错误", err)
		return
	}

	// Delete address
	err := h.userAddressService.DeleteAddress(c, delRequest.ID)
	if err != nil {
		v1.HandleError(c, v1.ErrRegisterCode, "删除地址失败", nil)
		return
	}

	v1.HandleSuccess(c, nil)
}
