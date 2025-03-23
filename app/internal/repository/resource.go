package repository

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"app/internal/common"
)

type ResourceRepository interface {
	GetResource(ctx *gin.Context, url string) (string, []byte, error)
}

func NewResourceRepository(
	repository *Repository,
) ResourceRepository {
	return &resourceRepository{
		Repository: repository,
	}
}

type resourceRepository struct {
	*Repository
}

func (r *resourceRepository) GetResource(ctx *gin.Context, url string) (string, []byte, error) {
	url = common.HTTPS_PREFIX + common.BUCKET_NAME + "." + common.OSS_ENDPOINT +
		ctx.Request.URL.Path[len("/resource/avatar"):]
	response, err := http.Get(url)
	if err != nil {
		r.logger.Debug("Error:", err)
		return "", nil, err
	}
	defer response.Body.Close()
	// 获取response中的内容
	body, err := io.ReadAll(response.Body)
	if err != nil {
		r.logger.Debug("Error reading response body:", err)
		return "", nil, errors.New("Error reading response body")
	}
	contentType := http.DetectContentType(body)

	return contentType, body, nil
}
