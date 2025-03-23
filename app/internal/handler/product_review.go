package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	v1 "app/api/v1"
	"app/internal/model"
	"app/internal/service"
)

type ProductReviewHandler struct {
	*Handler
	productReviewService service.ProductReviewService
}

func NewProductReviewHandler(
	handler *Handler,
	productReviewService service.ProductReviewService,
) *ProductReviewHandler {
	return &ProductReviewHandler{
		Handler:              handler,
		productReviewService: productReviewService,
	}
}

// CreateReview godoc
// @Summary 创建商品评价
// @Description 对已购买的商品进行评价
// @Tags 评价
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body model.CreateReviewRequest true "评价信息"
// @Success 200 {object} v1.Response
// @Failure 400 {object} v1.Response "参数错误"
// @Failure 401 {object} v1.Response "未授权"
// @Failure 403 {object} v1.Response "已评价过该商品或订单不存在"
// @Failure 500 {object} v1.Response "服务器内部错误"
// @Router /review/create [post]
func (h *ProductReviewHandler) CreateReview(c *gin.Context) {
	var req model.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		v1.HandleError(c, v1.ErrParamCode, "参数错误", err)
		return
	}

	// 调用服务层创建评价
	err := h.productReviewService.CreateReview(c, req)
	if err != nil {
		v1.HandleError(c, v1.ErrRegisterCode, err.Error(), err)
		return
	}

	v1.HandleSuccess(c, nil)
}

// GetProductReviews godoc
// @Summary 获取商品评价列表
// @Description 获取指定商品的评价列表
// @Tags 评价
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param product_id path uint64 true "商品ID"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} model.ReviewListResponse
// @Failure 400 {object} v1.Response "参数错误"
// @Failure 401 {object} v1.Response "未授权"
// @Failure 500 {object} v1.Response "服务器内部错误"
// @Router /review/product/{product_id} [get]
func (h *ProductReviewHandler) GetProductReviews(c *gin.Context) {
	productIDStr := c.Param("product_id")
	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		v1.HandleError(c, v1.ErrParamCode, "参数错误", err)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// 调用服务层获取商品评价
	response, err := h.productReviewService.GetProductReviews(c, productID, page, pageSize)
	if err != nil {
		v1.HandleError(c, v1.ErrRegisterCode, "获取商品评价失败", err)
		return
	}

	v1.HandleSuccess(c, response)
}

// GetUserReviews godoc
// @Summary 获取用户评价列表
// @Description 获取当前用户的评价列表
// @Tags 评价
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} model.ReviewListResponse
// @Failure 400 {object} v1.Response "参数错误"
// @Failure 401 {object} v1.Response "未授权"
// @Failure 500 {object} v1.Response "服务器内部错误"
// @Router /review/user [get]
func (h *ProductReviewHandler) GetUserReviews(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// 调用服务层获取用户评价
	response, err := h.productReviewService.GetUserReviews(c, page, pageSize)
	if err != nil {
		v1.HandleError(c, v1.ErrRegisterCode, "获取用户评价失败", err)
		return
	}

	v1.HandleSuccess(c, response)
}

// GetUserReviewsByTab godoc
// @Summary 按标签获取用户评价列表
// @Description 获取当前用户的待评价、已评价或全部评价列表
// @Tags 评价
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param tab query string false "标签类型(all:全部;pending:待评价;completed:已评价)" default(all)
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} model.UserReviewListResponse
// @Failure 400 {object} v1.Response "参数错误"
// @Failure 401 {object} v1.Response "未授权"
// @Failure 500 {object} v1.Response "服务器内部错误"
// @Router /review/user/tab [get]
func (h *ProductReviewHandler) GetUserReviewsByTab(c *gin.Context) {
	var req model.UserReviewListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		v1.HandleError(c, v1.ErrParamCode, "参数错误", err)
		return
	}

	// 调用服务层获取用户评价
	response, err := h.productReviewService.GetUserReviewsByTab(c, req)
	if err != nil {
		if err.Error() == "无效的标签类型" {
			v1.HandleError(c, v1.ErrParamCode, err.Error(), nil)
			return
		}
		v1.HandleError(c, v1.ErrRegisterCode, "获取用户评价失败", err)
		return
	}

	v1.HandleSuccess(c, response)
}

// DeleteReview godoc
// @Summary 删除商品评价
// @Description 删除用户自己的商品评价
// @Tags 评价
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param review_id path uint64 true "评价ID"
// @Success 200 {object} v1.Response
// @Failure 400 {object} v1.Response "参数错误"
// @Failure 401 {object} v1.Response "未授权"
// @Failure 403 {object} v1.Response "无权限删除该评价"
// @Failure 404 {object} v1.Response "评价不存在"
// @Failure 500 {object} v1.Response "服务器内部错误"
// @Router /review/delete/{review_id} [delete]
func (h *ProductReviewHandler) DeleteReview(c *gin.Context) {
	reviewIDStr := c.Param("review_id")
	reviewID, err := strconv.ParseUint(reviewIDStr, 10, 64)
	if err != nil {
		v1.HandleError(c, v1.ErrParamCode, "参数错误", err)
		return
	}

	// 调用服务层删除评价
	err = h.productReviewService.DeleteReview(c, reviewID)
	if err != nil {
		if err.Error() == "评价不存在" {
			v1.HandleError(c, v1.ErrNotFoundCode, err.Error(), nil)
			return
		} else if err.Error() == "无权限删除该评价" {
			v1.HandleError(c, v1.ErrForbiddenCode, err.Error(), nil)
			return
		}
		v1.HandleError(c, v1.ErrRegisterCode, "删除评价失败", err)
		return
	}

	v1.HandleSuccess(c, nil)
}

// GetReviewDetail godoc
// @Summary 获取评价详情
// @Description 根据评价ID获取评价详情
// @Tags 评价
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param review_id path uint64 true "评价ID"
// @Success 200 {object} model.ReviewDetailResponse
// @Failure 400 {object} v1.Response "参数错误"
// @Failure 401 {object} v1.Response "未授权"
// @Failure 404 {object} v1.Response "评价不存在"
// @Failure 500 {object} v1.Response "服务器内部错误"
// @Router /review/detail/{review_id} [get]
func (h *ProductReviewHandler) GetReviewDetail(c *gin.Context) {
	reviewIDStr := c.Param("review_id")
	reviewID, err := strconv.ParseUint(reviewIDStr, 10, 64)
	if err != nil {
		v1.HandleError(c, v1.ErrParamCode, "参数错误", err)
		return
	}

	// 调用服务层获取评价详情
	response, err := h.productReviewService.GetReviewDetail(c, reviewID)
	if err != nil {
		if err.Error() == "评价不存在" {
			v1.HandleError(c, v1.ErrNotFoundCode, err.Error(), nil)
			return
		}
		v1.HandleError(c, v1.ErrRegisterCode, "获取评价详情失败", err)
		return
	}

	v1.HandleSuccess(c, response)
}

// IncrementReviewCounter godoc
// @Summary 更新评价计数器
// @Description 更新评价的浏览人数、评价人数或点赞人数
// @Tags 评价
// @Accept json
// @Produce json
// @Param Authorization header string false "Bearer token"
// @Param request body model.UpdateReviewCounterRequest true "计数器更新信息"
// @Success 200 {object} v1.Response
// @Failure 400 {object} v1.Response "参数错误"
// @Failure 404 {object} v1.Response "评价不存在"
// @Failure 500 {object} v1.Response "服务器内部错误"
// @Router /review/increment_counter [post]
func (h *ProductReviewHandler) IncrementReviewCounter(c *gin.Context) {
	var req model.UpdateReviewCounterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		v1.HandleError(c, v1.ErrParamCode, "参数错误", err)
		return
	}

	// 验证计数器类型
	if req.Type > 3 {
		v1.HandleError(c, v1.ErrParamCode, "无效的计数器类型", nil)
		return
	}

	// 调用服务层更新计数器
	err := h.productReviewService.IncrementReviewCounter(c, req)
	if err != nil {
		if err.Error() == "评价不存在" {
			v1.HandleError(c, v1.ErrNotFoundCode, err.Error(), nil)
			return
		}
		v1.HandleError(c, v1.ErrRegisterCode, "更新评价计数器失败", err)
		return
	}

	v1.HandleSuccess(c, nil)
}
