package handler

import (
	"github.com/gin-gonic/gin"

	v1 "app/api/v1"
	"app/internal/model"
	"app/internal/service"
)

type ProductEvaluateHandler struct {
	*Handler
	productEvaluateService service.ProductEvaluateService
}

func NewProductEvaluateHandler(
	handler *Handler,
	productEvaluateService service.ProductEvaluateService,
) *ProductEvaluateHandler {
	return &ProductEvaluateHandler{
		Handler:                handler,
		productEvaluateService: productEvaluateService,
	}
}

// CreateEvaluateReply godoc
// @Summary 创建评价回复
// @Description 回复商品评价
// @Tags 评价
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body model.CreateEvaluateReplyRequest true "回复信息"
// @Success 200 {object} v1.Response
// @Failure 400 {object} v1.Response "参数错误"
// @Failure 401 {object} v1.Response "未授权"
// @Failure 404 {object} v1.Response "评价不存在"
// @Failure 500 {object} v1.Response "服务器内部错误"
// @Router /review/reply [post]
func (h *ProductEvaluateHandler) CreateEvaluateReply(c *gin.Context) {
	var req model.CreateEvaluateReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		v1.HandleError(c, v1.ErrParamCode, "参数错误", err)
		return
	}

	// 校验内容不能为空
	if req.Content == "" {
		v1.HandleError(c, v1.ErrParamCode, "回复内容不能为空", nil)
		return
	}

	// 调用服务层创建评价回复
	err := h.productEvaluateService.CreateEvaluateReply(c, req)
	if err != nil {
		if err.Error() == "评价不存在" {
			v1.HandleError(c, v1.ErrNotFoundCode, err.Error(), nil)
			return
		}
		v1.HandleError(c, v1.ErrRegisterCode, "回复评价失败", err)
		return
	}

	v1.HandleSuccess(c, nil)
}

// CreateEvaluate godoc
// @Summary 创建主评论
// @Description 对商品评价进行评论（一级评论）
// @Tags 评价
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body model.CreateEvaluateRequest true "评论信息"
// @Success 200 {object} v1.Response
// @Failure 400 {object} v1.Response "参数错误"
// @Failure 401 {object} v1.Response "未授权"
// @Failure 404 {object} v1.Response "评价不存在"
// @Failure 500 {object} v1.Response "服务器内部错误"
// @Router /review/comment [post]
func (h *ProductEvaluateHandler) CreateEvaluate(c *gin.Context) {
	var req model.CreateEvaluateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		v1.HandleError(c, v1.ErrParamCode, "参数错误", err)
		return
	}

	// 校验内容不能为空
	if req.Content == "" {
		v1.HandleError(c, v1.ErrParamCode, "评论内容不能为空", nil)
		return
	}

	// 调用服务层创建评论
	err := h.productEvaluateService.CreateEvaluate(c, req)
	if err != nil {
		if err.Error() == "评价不存在" {
			v1.HandleError(c, v1.ErrNotFoundCode, err.Error(), nil)
			return
		}
		v1.HandleError(c, v1.ErrRegisterCode, "评论失败", err)
		return
	}

	v1.HandleSuccess(c, nil)
}

// UpdateEvaluateAnonymous godoc
// @Summary 更新评论匿名状态
// @Description 设置评论为匿名或公开状态
// @Tags 评价
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body model.UpdateEvaluateAnonymousRequest true "匿名状态信息"
// @Success 200 {object} v1.Response
// @Failure 400 {object} v1.Response "参数错误"
// @Failure 401 {object} v1.Response "未授权"
// @Failure 403 {object} v1.Response "评论不存在或您没有权限操作"
// @Failure 500 {object} v1.Response "服务器内部错误"
// @Router /review/update_anonymous [post]
func (h *ProductEvaluateHandler) UpdateEvaluateAnonymous(c *gin.Context) {
	var req model.UpdateEvaluateAnonymousRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		v1.HandleError(c, v1.ErrParamCode, "参数错误", err)
		return
	}

	// 调用服务层更新评论匿名状态
	err := h.productEvaluateService.UpdateEvaluateAnonymous(c, req)
	if err != nil {
		if err.Error() == "评论不存在或您没有权限操作" {
			v1.HandleError(c, v1.ErrForbiddenCode, err.Error(), nil)
			return
		}
		v1.HandleError(c, v1.ErrRegisterCode, "更新评论匿名状态失败", err)
		return
	}

	v1.HandleSuccess(c, nil)
}
