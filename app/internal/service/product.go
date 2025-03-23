package service

import (
	"context"
	"strconv"
	"strings"

	"app/internal/model"
	"app/internal/repository"

	"github.com/gin-gonic/gin"
)

type ProductService interface {
	GetProductList(ctx *gin.Context) ([]*model.ProductListResponse, error)
	GetRecommendProductList(ctx context.Context) ([]*model.ProductListItemDTO, error)
	GetProductByID(ctx *gin.Context, id int) (*model.ProductListItemDTO, error)
	GetProductDetailsByCartIDs(ctx *gin.Context, cart_ids string) (*model.ProductDetailsResponse, error)
}

func NewProductService(
	service *Service,
	productRepository repository.ProductRepository,
) ProductService {
	return &productService{
		Service:           service,
		productRepository: productRepository,
	}
}

type productService struct {
	*Service
	productRepository repository.ProductRepository
}

func (s *productService) GetProductList(ctx *gin.Context) ([]*model.ProductListResponse, error) {
	userId := GetUserIdFromCtx(ctx)
	return s.productRepository.GetProductList(ctx, userId)
}

func (s *productService) GetRecommendProductList(ctx context.Context) ([]*model.ProductListItemDTO, error) {
	return s.productRepository.GetRecommendProductList(ctx)
}

func (s *productService) GetProductByID(ctx *gin.Context, id int) (*model.ProductListItemDTO, error) {
	userId := GetUserIdFromCtx(ctx)
	return s.productRepository.GetProductByID(ctx, id, userId)
}

func (s *productService) GetProductDetailsByCartIDs(ctx *gin.Context, cart_ids string) (*model.ProductDetailsResponse, error) {
	// 转[]uint64
	cart_ids_list := strings.Split(cart_ids, ",")
	cart_ids_list_uint64 := make([]uint64, 0)
	for _, cart_id := range cart_ids_list {
		cart_id_uint64, err := strconv.ParseUint(cart_id, 10, 64)
		if err != nil || cart_id_uint64 == 0 {
			continue
		}
		cart_ids_list_uint64 = append(cart_ids_list_uint64, cart_id_uint64)
	}
	products, err := s.productRepository.GetProductDetailsByCartIDs(ctx, cart_ids_list_uint64)
	if err != nil {
		return nil, err
	}

	return &model.ProductDetailsResponse{
		Products: products,
	}, nil
}
