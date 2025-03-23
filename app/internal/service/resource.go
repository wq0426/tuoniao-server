package service

import (
	"github.com/gin-gonic/gin"

	"app/internal/repository"
)

type ResourceService interface {
	GetResource(ctx *gin.Context, url string) (string, []byte, error)
}

func NewResourceService(
	service *Service,
	resourceRepository repository.ResourceRepository,
) ResourceService {
	return &resourceService{
		Service:            service,
		resourceRepository: resourceRepository,
	}
}

type resourceService struct {
	*Service
	resourceRepository repository.ResourceRepository
}

func (s *resourceService) GetResource(ctx *gin.Context, url string) (string, []byte, error) {
	return s.resourceRepository.GetResource(ctx, url)
}
