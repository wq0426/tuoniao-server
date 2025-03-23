package repository

import (
	"context"
	"time"

	"app/internal/model"

	"gorm.io/gorm"
)

type UserAddressRepository interface {
	AddAddress(ctx context.Context, userID string, address model.AddAddressRequest) error
	UpdateAddress(ctx context.Context, userID string, address model.UpdateAddressRequest) error
	GetUserAddresses(ctx context.Context, userID string) ([]model.UserAddress, error)
	DeleteAddress(ctx context.Context, userID string, addressID uint64) error
}

func NewUserAddressRepository(
	repository *Repository,
) UserAddressRepository {
	return &userAddressRepository{
		Repository: repository,
	}
}

type userAddressRepository struct {
	*Repository
}

func (r *userAddressRepository) AddAddress(ctx context.Context, userID string, req model.AddAddressRequest) error {
	// Begin transaction
	return r.DB(ctx).Transaction(func(tx *gorm.DB) error {
		// If this address is set as default, unset any existing default addresses
		if req.IsDefault == 1 {
			if err := tx.Model(&model.UserAddress{}).
				Where("user_id = ? AND is_default = ?", userID, 1).
				Update("is_default", 0).Error; err != nil {
				return err
			}
		}

		// Create the new address
		address := model.UserAddress{
			UserID:    userID,
			Name:      req.Name,
			Province:  req.Province,
			City:      req.City,
			District:  req.District,
			Street:    req.Street,
			IsDefault: req.IsDefault,
			Phone:     req.Phone,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		return tx.Create(&address).Error
	})
}

func (r *userAddressRepository) UpdateAddress(ctx context.Context, userID string, req model.UpdateAddressRequest) error {
	// Begin transaction
	return r.DB(ctx).Transaction(func(tx *gorm.DB) error {
		// Verify the address belongs to this user
		var address model.UserAddress
		if err := tx.Where("id = ? AND user_id = ?", req.ID, userID).First(&address).Error; err != nil {
			return err
		}

		// If this address is set as default, unset any existing default addresses
		if req.IsDefault == 1 && address.IsDefault != 1 {
			if err := tx.Model(&model.UserAddress{}).
				Where("user_id = ? AND is_default = ?", userID, 1).
				Update("is_default", 0).Error; err != nil {
				return err
			}
		}

		// Update the address
		address.Name = req.Name
		address.Province = req.Province
		address.City = req.City
		address.District = req.District
		address.Street = req.Street
		address.IsDefault = req.IsDefault
		address.Phone = req.Phone
		address.UpdatedAt = time.Now()

		return tx.Save(&address).Error
	})
}

func (r *userAddressRepository) GetUserAddresses(ctx context.Context, userID string) ([]model.UserAddress, error) {
	var addresses []model.UserAddress
	if err := r.DB(ctx).Where("user_id = ?", userID).Order("is_default DESC, updated_at DESC").Find(&addresses).Error; err != nil {
		return nil, err
	}
	return addresses, nil
}

func (r *userAddressRepository) DeleteAddress(ctx context.Context, userID string, addressID uint64) error {
	// Check if the address belongs to the user and delete it
	result := r.DB(ctx).Where("id = ? AND user_id = ?", addressID, userID).Delete(&model.UserAddress{})
	if result.Error != nil {
		return result.Error
	}

	// Check if any row was affected (to verify the address existed)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
