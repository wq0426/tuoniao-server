package service

import (
	"app/internal/model"
	"app/internal/repository"

	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserCartService interface {
	AddToCart(ctx *gin.Context, req model.AddToCartRequest) error
	GetUserCartItems(ctx *gin.Context) (model.CartResponse, error)
	DeleteCartItems(ctx *gin.Context, req model.DeleteCartItemsRequest) error
}

func NewUserCartService(
	service *Service,
	userCartRepository repository.UserCartRepository,
) UserCartService {
	return &userCartService{
		Service:            service,
		userCartRepository: userCartRepository,
	}
}

type userCartService struct {
	*Service
	userCartRepository repository.UserCartRepository
}

func (s *userCartService) AddToCart(ctx *gin.Context, req model.AddToCartRequest) error {
	// Get user ID from context (assuming it's set by auth middleware)
	userID := GetUserIdFromCtx(ctx)
	return s.userCartRepository.AddToCart(ctx, userID, req.ProductID, req.Quantity, req.CouponID)
}

func (s *userCartService) GetUserCartItems(ctx *gin.Context) (model.CartResponse, error) {
	userID := GetUserIdFromCtx(ctx)
	return s.userCartRepository.GetUserCartItems(ctx, userID)
}

func (s *userCartService) DeleteCartItems(ctx *gin.Context, req model.DeleteCartItemsRequest) error {
	userID := GetUserIdFromCtx(ctx)

	// Convert comma-separated IDs to a slice of uint
	idStrings := strings.Split(req.IDs, ",")
	cartIDs := make([]uint, 0, len(idStrings))

	for _, idStr := range idStrings {
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			return err
		}
		cartIDs = append(cartIDs, uint(id))
	}

	return s.userCartRepository.DeleteCartItems(ctx, userID, cartIDs)
}
