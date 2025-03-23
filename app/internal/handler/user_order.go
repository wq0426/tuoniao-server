package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	v1 "app/api/v1"
	"app/internal/model"
	"app/internal/service"
)

type UserOrderHandler struct {
	*Handler
	userOrderService service.UserOrderService
}

func NewUserOrderHandler(
	handler *Handler,
	userOrderService service.UserOrderService,
) *UserOrderHandler {
	return &UserOrderHandler{
		Handler:          handler,
		userOrderService: userOrderService,
	}
}

// CreateOrders godoc
// @Summary 创建用户订单
// @Description 创建用户订单，支持多个商品一次性提交
// @Tags 订单
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body model.CreateOrderRequest true "创建订单请求"
// @Success 200 {object} v1.Response
// @Router /order/create [post]
func (h *UserOrderHandler) CreateOrders(c *gin.Context) {
	// Parse request body
	var req model.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		v1.HandleError(c, v1.ErrParamCode, "参数错误", err)
		return
	}

	// Validate request has items
	if len(req.Items) == 0 {
		v1.HandleError(c, v1.ErrParamCode, "订单项不能为空", nil)
		return
	}

	// Create orders
	err := h.userOrderService.CreateOrders(c, req)
	if err != nil {
		v1.HandleError(c, v1.ErrRegisterCode, "创建订单失败", nil)
		return
	}

	v1.HandleSuccess(c, nil)
}

// GetOrderList godoc
// @Summary 获取订单列表
// @Description 获取用户订单列表，可根据状态筛选
// @Tags 订单
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param status query int false "订单状态：0-待付款，1-已付款，2-已发货，3-已完成，4-已取消；不传则查询所有状态"
// @Param page query int false "页码，默认1"
// @Param page_size query int false "每页条数，默认10"
// @Success 200 {object} model.OrderListResponse
// @Router /order/list [get]
func (h *UserOrderHandler) GetOrderList(c *gin.Context) {
	var req model.OrderQueryRequest
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

	// 调用服务层查询订单列表
	orders, err := h.userOrderService.GetOrderList(c, req)
	if err != nil {
		v1.HandleError(c, v1.ErrRegisterCode, "获取订单列表失败", nil)
		return
	}

	v1.HandleSuccess(c, orders)
}

// GetOrderProductDetails godoc
// @Summary 获取订单商品详情
// @Description 根据订单ID查询该订单下的所有商品和优惠券信息
// @Tags 订单
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param order_id query uint64 true "订单ID"
// @Success 200 {object} model.OrderProductsResponse "订单商品详情"
// @Router /order/products [get]
func (h *UserOrderHandler) GetOrderProductDetails(c *gin.Context) {
	var req model.OrderProductsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		v1.HandleError(c, v1.ErrParamCode, "参数错误", err)
		return
	}

	// 调用服务层方法获取订单商品详情
	response, err := h.userOrderService.GetOrderProductDetails(c, req)
	if err != nil {
		v1.HandleError(c, v1.ErrRegisterCode, "获取订单商品详情失败", err)
		return
	}

	v1.HandleSuccess(c, response.Products)
}

// UpdateOrderStatus godoc
// @Summary 更新订单状态
// @Description 根据订单ID更新订单状态
// @Tags 订单
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body model.UpdateOrderStatusRequest true "更新订单状态请求"
// @Success 200 {object} v1.Response "操作成功"
// @Failure 400 {object} v1.Response "参数错误或非法的状态转换"
// @Failure 401 {object} v1.Response "未授权"
// @Failure 404 {object} v1.Response "订单不存在或无权限"
// @Router /order/status [post]
func (h *UserOrderHandler) UpdateOrderStatus(c *gin.Context) {
	var req model.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		v1.HandleError(c, v1.ErrParamCode, "参数错误", err)
		return
	}

	// 调用服务层更新订单状态
	err := h.userOrderService.UpdateOrderStatus(c, req)
	if err != nil {
		if err.Error() == "订单不存在或无权限" {
			v1.HandleError(c, v1.ErrNotFoundCode, err.Error(), nil)
			return
		} else if err.Error() == "非法的状态转换" {
			v1.HandleError(c, v1.ErrParamCode, err.Error(), nil)
			return
		}
		v1.HandleError(c, v1.ErrRegisterCode, "更新订单状态失败", err)
		return
	}

	v1.HandleSuccess(c, nil)
}

// GetOrderDetail godoc
// @Summary 获取订单详情
// @Description 获取订单详情
// @Tags 订单
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param order_id path uint64 true "订单ID"
// @Success 200 {object} model.OrderDetailResponse
// @Failure 400 {object} v1.Response "参数错误"
// @Failure 401 {object} v1.Response "未授权"
// @Failure 404 {object} v1.Response "订单不存在"
// @Failure 500 {object} v1.Response "服务器内部错误"
// @Router /order/detail/{order_item_id} [get]
func (h *UserOrderHandler) GetOrderDetail(c *gin.Context) {
	orderItemIDStr := c.Param("order_item_id")
	orderItemID, err := strconv.ParseUint(orderItemIDStr, 10, 64)
	if err != nil {
		v1.HandleError(c, v1.ErrParamCode, "参数错误", err)
		return
	}

	// 调用服务层获取订单详情
	response, err := h.userOrderService.GetOrderDetail(c, orderItemID)
	if err != nil {
		if err.Error() == "订单不存在或无权限" {
			v1.HandleError(c, v1.ErrNotFoundCode, err.Error(), nil)
			return
		}
		v1.HandleError(c, v1.ErrRegisterCode, "获取订单详情失败", err)
		return
	}

	v1.HandleSuccess(c, response)
}
