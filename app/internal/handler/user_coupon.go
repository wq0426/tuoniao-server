package handler

import (
	"github.com/gin-gonic/gin"

	v1 "app/api/v1"
	"app/internal/model"
	"app/internal/service"
	"strconv"
)

type UserCouponHandler struct {
	*Handler
	userCouponService service.UserCouponService
}

func NewUserCouponHandler(
	handler *Handler,
	userCouponService service.UserCouponService,
) *UserCouponHandler {
	return &UserCouponHandler{
		Handler:           handler,
		userCouponService: userCouponService,
	}
}

// GetUserCoupons godoc
// @Summary 获取用户特定商品的优惠券
// @Description 获取用户针对特定商品的优惠券列表
// @Tags 优惠券
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param product_id query uint64 true "商品ID"
// @Success 200 {array} model.UserCoupon
// @Router /coupon/product [get]
func (h *UserCouponHandler) GetUserCoupons(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Query("product_id"), 10, 64)
	if err != nil {
		v1.HandleError(c, v1.ErrParamCode, "商品ID格式错误", err)
		return
	}

	coupons, err := h.userCouponService.GetUserCoupons(c, productID)
	if err != nil {
		v1.HandleError(c, v1.ErrRegisterCode, "获取优惠券失败", nil)
		return
	}

	v1.HandleSuccess(c, coupons)
}

// GetAllUserCoupons godoc
// @Summary 获取用户所有优惠券
// @Description 获取用户所有优惠券列表，可根据状态过滤
// @Tags 优惠券
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {array} model.ProductCouponDTO
// @Router /coupon/list [get]
func (h *UserCouponHandler) GetAllUserCoupons(c *gin.Context) {
	coupons, err := h.userCouponService.GetAllUserCoupons(c)
	if err != nil {
		v1.HandleError(c, v1.ErrRegisterCode, "获取优惠券失败", nil)
		return
	}

	v1.HandleSuccess(c, coupons)
}

// ClaimCoupon godoc
// @Summary 领取优惠券
// @Description 根据优惠券ID领取优惠券，如果已领取会提示
// @Tags 优惠券
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body model.ClaimCouponRequest true "领取优惠券请求"
// @Success 200 {object} v1.Response
// @Router /coupon/claim [post]
func (h *UserCouponHandler) ClaimCoupon(c *gin.Context) {
	// 解析请求体
	var req model.ClaimCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		v1.HandleError(c, v1.ErrParamCode, "参数错误", err)
		return
	}

	// 领取优惠券
	err := h.userCouponService.ClaimCoupon(c, req)
	if err != nil {
		// 检查是否是已领取的错误
		if err.Error() == "coupon already claimed" {
			v1.HandleError(c, v1.ErrRegisterCode, "您已经领取过该优惠券", nil)
			return
		}
		v1.HandleError(c, v1.ErrRegisterCode, "领取优惠券失败", nil)
		return
	}

	v1.HandleSuccess(c, nil)
}

// GetUserCouponDetail godoc
// @Summary 获取优惠券详情
// @Description 根据优惠券ID获取优惠券详情
// @Tags 优惠券
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param coupon_id query uint64 true "优惠券ID"
// @Success 200 {object} model.CouponDetailResponse
// @Router /coupon/detail [get]
func (h *UserCouponHandler) GetUserCouponDetail(c *gin.Context) {
	couponIDStr := c.Query("coupon_id")
	if couponIDStr == "" {
		v1.HandleError(c, v1.ErrParamCode, "优惠券ID不能为空", nil)
		return
	}

	couponID, err := strconv.ParseUint(couponIDStr, 10, 64)
	if err != nil {
		v1.HandleError(c, v1.ErrParamCode, "优惠券ID格式错误", err)
		return
	}

	coupon, err := h.userCouponService.GetUserCouponByID(c, couponID)
	if err != nil {
		v1.HandleError(c, v1.ErrRegisterCode, "获取优惠券详情失败", err)
		return
	}

	v1.HandleSuccess(c, coupon)
}
