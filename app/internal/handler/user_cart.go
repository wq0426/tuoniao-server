package handler

import (
	"github.com/gin-gonic/gin"

	v1 "app/api/v1"
	"app/internal/model"
	"app/internal/service"
)

type UserCartHandler struct {
	*Handler
	userCartService service.UserCartService
}

func NewUserCartHandler(
	handler *Handler,
	userCartService service.UserCartService,
) *UserCartHandler {
	return &UserCartHandler{
		Handler:         handler,
		userCartService: userCartService,
	}
}

// AddToCart godoc
// @Summary 添加商品到购物车
// @Description 添加商品到购物车，如果已存在则数量加1
// @Tags 购物车
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body model.AddToCartRequest true "添加购物车请求"
// @Success 200 {object} v1.Response
// @Router /cart/add [post]
func (h *UserCartHandler) AddToCart(c *gin.Context) {
	// Parse request body
	var req model.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		v1.HandleError(c, v1.ErrParamCode, "参数错误", err)
		return
	}

	// Add to cart
	err := h.userCartService.AddToCart(c, req)
	if err != nil {
		v1.HandleError(c, v1.ErrRegisterCode, "添加购物车失败", nil)
		return
	}

	v1.HandleSuccess(c, nil)
}

// GetUserCartItems godoc
// @Summary 获取用户购物车列表
// @Description 获取用户购物车中未处理的商品列表，按店铺分组
// @Tags 购物车
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {array} model.CartStoreDTO
// @Router /cart/list [get]
func (h *UserCartHandler) GetUserCartItems(c *gin.Context) {
	// Get cart items from service
	cartItems, err := h.userCartService.GetUserCartItems(c)
	if err != nil {
		v1.HandleError(c, v1.ErrRegisterCode, "获取购物车失败", nil)
		return
	}

	v1.HandleSuccess(c, cartItems)
}

// DeleteCartItems godoc
// @Summary 删除购物车商品
// @Description 删除购物车中的商品，支持批量删除
// @Tags 购物车
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body model.DeleteCartItemsRequest true "删除购物车请求"
// @Success 200 {object} v1.Response
// @Router /cart/delete [post]
func (h *UserCartHandler) DeleteCartItems(c *gin.Context) {
	// Parse request body
	var req model.DeleteCartItemsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		v1.HandleError(c, v1.ErrParamCode, "参数错误", err)
		return
	}

	// Delete cart items
	err := h.userCartService.DeleteCartItems(c, req)
	if err != nil {
		v1.HandleError(c, v1.ErrRegisterCode, "删除购物车商品失败", nil)
		return
	}

	v1.HandleSuccess(c, nil)
}
