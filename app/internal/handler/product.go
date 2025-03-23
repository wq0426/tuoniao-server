package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	v1 "app/api/v1"
	"app/internal/service"
)

type ProductHandler struct {
	*Handler
	productService service.ProductService
}

func NewProductHandler(
	handler *Handler,
	productService service.ProductService,
) *ProductHandler {
	return &ProductHandler{
		Handler:        handler,
		productService: productService,
	}
}

// GetProductList godoc
// @Summary 获取商品列表
// @Description 获取商品列表，包含分类和详细信息
// @Tags 商品模块
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} []model.ProductListResponse
// @Router /product/list [get]
func (h *ProductHandler) GetProductList(c *gin.Context) {
	products, err := h.productService.GetProductList(c)
	if err != nil {
		v1.HandleError(c, v1.ErrRegisterCode, "获取商品列表失败", nil)
		return
	}
	v1.HandleSuccess(c, products)
}

// GetRecommendProductList godoc
// @Summary 获取推荐商品列表
// @Description 获取推荐商品列表
// @Tags 商品模块
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} []model.ProductListResponse
// @Router /product/recommend [get]
func (h *ProductHandler) GetRecommendProductList(c *gin.Context) {
	products, err := h.productService.GetRecommendProductList(c)
	if err != nil {
		v1.HandleError(c, v1.ErrRegisterCode, "获取推荐商品列表失败", nil)
		return
	}
	v1.HandleSuccess(c, products)
}

// GetProductByID godoc
// @Summary 获取商品详情
// @Description 获取商品详情
// @Tags 商品模块
// @Accept json
// @Produce json
// @Param product_id query int true "商品ID"
// @Success 200 {object} model.ProductListItemDTO
// @Router /product/detail [get]
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id := c.Query("product_id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		v1.HandleError(c, v1.ErrRegisterCode, "获取商品详情失败", nil)
		return
	}
	product, err := h.productService.GetProductByID(c, idInt)
	if err != nil {
		v1.HandleError(c, v1.ErrRegisterCode, "获取商品详情失败", nil)
		return
	}
	v1.HandleSuccess(c, product)
}

// GetProductDetailsByCartIDs godoc
// @Summary 通过购物车ID查询商品详情
// @Description 根据购物车ID列表查询对应的商品详情
// @Tags 商品
// @Accept json
// @Produce json
// @Param Authorization header string false "Bearer token"
// @Param cart_ids query string true "购物车ID列表"
// @Success 200 {object} []model.ProductListItemDTO "商品详情列表"
// @Router /product/details [get]
func (h *ProductHandler) GetProductDetailsByCartIDs(c *gin.Context) {
	ids := c.Query("cart_ids")
	// 调用服务层方法获取商品详情
	response, err := h.productService.GetProductDetailsByCartIDs(c, ids)
	if err != nil {
		v1.HandleError(c, v1.ErrRegisterCode, "获取商品详情失败", err)
		return
	}

	v1.HandleSuccess(c, response.Products)
}
