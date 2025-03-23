package handler

import (
	"github.com/gin-gonic/gin"

	v1 "app/api/v1"
	"app/internal/service"
)

type NewsHandler struct {
	*Handler
	newsService service.NewsService
}

func NewNewsHandler(
	handler *Handler,
	newsService service.NewsService,
) *NewsHandler {
	return &NewsHandler{
		Handler:     handler,
		newsService: newsService,
	}
}

// GetNewsList godoc
// @Summary 获取新闻资讯列表
// @Description 获取新闻资讯列表
// @Tags 新闻模块
// @Accept json
// @Produce json
// @Param keyword query string false "关键词"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} []model.NewsType
// @Router /news/list [get]
func (h *NewsHandler) GetNewsList(c *gin.Context) {
	// 获取请求参数keyword
	keyword := c.Query("keyword")
	news, err := h.newsService.GetNewsList(c, keyword)
	if err != nil {
		v1.HandleError(c, v1.ErrRegisterCode, v1.MsgNewsListError, nil)
		return
	}
	v1.HandleSuccess(c, news)
}
