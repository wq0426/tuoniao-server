package service

import (
	"app/internal/model"
	"app/internal/repository"

	"github.com/gin-gonic/gin"
)

type UserAddressService interface {
	AddAddress(ctx *gin.Context, req model.AddAddressRequest) error
	UpdateAddress(ctx *gin.Context, req model.UpdateAddressRequest) error
	GetUserAddresses(ctx *gin.Context) ([]model.UserAddress, error)
	DeleteAddress(ctx *gin.Context, addressID uint64) error
}

func NewUserAddressService(
	service *Service,
	userAddressRepository repository.UserAddressRepository,
) UserAddressService {
	return &userAddressService{
		Service:               service,
		userAddressRepository: userAddressRepository,
	}
}

type userAddressService struct {
	*Service
	userAddressRepository repository.UserAddressRepository
}

func (s *userAddressService) AddAddress(ctx *gin.Context, req model.AddAddressRequest) error {
	userID := GetUserIdFromCtx(ctx)
	return s.userAddressRepository.AddAddress(ctx, userID, req)
}

func (s *userAddressService) UpdateAddress(ctx *gin.Context, req model.UpdateAddressRequest) error {
	userID := GetUserIdFromCtx(ctx)
	return s.userAddressRepository.UpdateAddress(ctx, userID, req)
}

func (s *userAddressService) GetUserAddresses(ctx *gin.Context) ([]model.UserAddress, error) {
	userID := GetUserIdFromCtx(ctx)
	return s.userAddressRepository.GetUserAddresses(ctx, userID)
}

func (s *userAddressService) DeleteAddress(ctx *gin.Context, addressID uint64) error {
	userID := GetUserIdFromCtx(ctx)
	return s.userAddressRepository.DeleteAddress(ctx, userID, addressID)
}
