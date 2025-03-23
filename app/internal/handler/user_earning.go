package handler

import (
	"github.com/gin-gonic/gin"

	v1 "app/api/v1"
	"app/internal/model"
	"app/internal/service"
)

type UserEarningHandler struct {
	*Handler
	userEarningService service.UserEarningService
}

func NewUserEarningHandler(
	handler *Handler,
	userEarningService service.UserEarningService,
) *UserEarningHandler {
	return &UserEarningHandler{
		Handler:            handler,
		userEarningService: userEarningService,
	}
}

// AddEarning godoc
// @Summary 添加用户收益
// @Description 添加用户的收益记录
// @Tags 收益
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body model.AddEarningRequest true "收益信息"
// @Success 200 {object} v1.Response
// @Failure 400 {object} v1.Response "参数错误"
// @Failure 401 {object} v1.Response "未授权"
// @Failure 500 {object} v1.Response "服务器内部错误"
// @Router /earning/add [post]
func (h *UserEarningHandler) AddEarning(c *gin.Context) {
	var req model.AddEarningRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		v1.HandleError(c, v1.ErrParamCode, "参数错误", err)
		return
	}

	// 调用服务层添加收益
	err := h.userEarningService.AddEarning(c, req)
	if err != nil {
		v1.HandleError(c, v1.ErrRegisterCode, err.Error(), err)
		return
	}

	v1.HandleSuccess(c, nil)
}

// GetEarningList godoc
// @Summary 获取用户收益列表
// @Description 获取用户的收益记录列表，可按日期筛选
// @Tags 收益
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param date query string false "日期，格式：YYYY-MM-DD"
// @Param page query int false "页码，默认1"
// @Param page_size query int false "每页条数，默认10"
// @Success 200 {object} model.EarningListResponse
// @Failure 400 {object} v1.Response "参数错误"
// @Failure 401 {object} v1.Response "未授权"
// @Failure 500 {object} v1.Response "服务器内部错误"
// @Router /earning/list [get]
func (h *UserEarningHandler) GetEarningList(c *gin.Context) {
	var req model.QueryEarningRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		v1.HandleError(c, v1.ErrParamCode, "参数错误", err)
		return
	}

	// 调用服务层获取收益列表
	response, err := h.userEarningService.GetEarningList(c, req)
	if err != nil {
		v1.HandleError(c, v1.ErrRegisterCode, "获取收益列表失败", err)
		return
	}

	v1.HandleSuccess(c, response)
}
