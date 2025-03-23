package service

import (
	"app/internal/model"
	"app/internal/repository"

	"github.com/gin-gonic/gin"
)

type ProductEvaluateService interface {
	CreateEvaluateReply(ctx *gin.Context, req model.CreateEvaluateReplyRequest) error
	CreateEvaluate(ctx *gin.Context, req model.CreateEvaluateRequest) error
	UpdateEvaluateAnonymous(ctx *gin.Context, req model.UpdateEvaluateAnonymousRequest) error
}

type productEvaluateService struct {
	*Service
	productEvaluateRepository repository.ProductEvaluateRepository
}

func NewProductEvaluateService(
	service *Service,
	productEvaluateRepository repository.ProductEvaluateRepository,
) ProductEvaluateService {
	return &productEvaluateService{
		Service:                   service,
		productEvaluateRepository: productEvaluateRepository,
	}
}

// CreateEvaluateReply 创建评价回复
func (s *productEvaluateService) CreateEvaluateReply(ctx *gin.Context, req model.CreateEvaluateReplyRequest) error {
	userID := GetUserIdFromCtx(ctx)

	// 调用仓储层创建评价回复
	return s.productEvaluateRepository.CreateEvaluateReply(ctx, userID, req)
}

// CreateEvaluate 创建主评论
func (s *productEvaluateService) CreateEvaluate(ctx *gin.Context, req model.CreateEvaluateRequest) error {
	userID := GetUserIdFromCtx(ctx)

	// 调用仓储层创建评论
	return s.productEvaluateRepository.CreateEvaluate(ctx, userID, req)
}

// UpdateEvaluateAnonymous 更新评论匿名状态
func (s *productEvaluateService) UpdateEvaluateAnonymous(ctx *gin.Context, req model.UpdateEvaluateAnonymousRequest) error {
	userID := GetUserIdFromCtx(ctx)

	// 调用仓储层更新评论匿名状态
	return s.productEvaluateRepository.UpdateEvaluateAnonymous(ctx, userID, req)
}
