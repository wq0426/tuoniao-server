package repository

import (
	"context"

	"app/internal/model"

	"gorm.io/gorm"
)

type UserCartRepository interface {
	AddToCart(ctx context.Context, userID string, productID uint64, quantity int, couponID uint64) error
	GetUserCartItems(ctx context.Context, userID string) (model.CartResponse, error)
	DeleteCartItems(ctx context.Context, userID string, cartIDs []uint) error
	GetCartList(ctx context.Context, userID string, productIds []uint64) ([]model.CartProductDTO, error)
}

func NewUserCartRepository(
	repository *Repository,
) UserCartRepository {
	return &userCartRepository{
		Repository: repository,
	}
}

type userCartRepository struct {
	*Repository
}

func (r *userCartRepository) AddToCart(ctx context.Context, userID string, productID uint64, quantity int, couponID uint64) error {
	// Check if the product already exists in the cart with status = 0
	var cart model.UserCart

	// Create a transaction to ensure data consistency
	err := r.DB(ctx).Transaction(func(tx *gorm.DB) error {
		// 根据product_id查询商品信息
		var product model.Product
		productResult := tx.Where("id = ?", productID).First(&product)
		if productResult.Error != nil {
			return productResult.Error
		}

		// Try to find the existing cart item
		cartResult := tx.Where("user_id = ? AND product_id = ? AND status = 0", userID, productID).First(&cart)
		if cartResult.Error != nil {
			var coupon model.ProductCoupon
			if cartResult.Error == gorm.ErrRecordNotFound {
				if couponID > 0 {
					// 查询优惠券信息
					couponResult := tx.Where("id = ?", couponID).First(&coupon)
					if couponResult.Error != nil {
						return couponResult.Error
					}
				}
				// Create a new cart item
				status := uint8(0) // Unprocessed status
				newCart := model.UserCart{
					UserID:       userID,
					ProductID:    productID,
					ProductName:  product.ProductName,
					CurrentPrice: product.CurrentPrice,
					Quantity:     quantity,
					Status:       &status,
					StoreID:      product.StoreID,
					StoreName:    product.StoreName,
					CouponID:     couponID,
					CouponPrice:  coupon.CouponPrice,
				}
				return tx.Create(&newCart).Error
			}
			return cartResult.Error
		}

		// If found, increment the product number
		if quantity > 0 {
			return tx.Model(&cart).
				Update("quantity", gorm.Expr("quantity + ?", quantity)).Error
		} else {
			return tx.Model(&cart).
				Update("quantity", gorm.Expr("quantity - ?", quantity)).Error
		}
	})

	return err
}

func (r *userCartRepository) GetUserCartItems(ctx context.Context, userID string) (model.CartResponse, error) {
	// Get all cart items for this user with status = 0 (unprocessed)
	var cartItems []model.UserCart
	if err := r.DB(ctx).Where("user_id = ? AND status = 0", userID).Find(&cartItems).Error; err != nil {
		return nil, err
	}

	// Group cart items by store
	storeMap := make(map[uint64]*model.CartStoreDTO)

	for _, item := range cartItems {
		// Skip if product has invalid data
		if item.StoreID == 0 || item.ProductID == 0 {
			continue
		}

		// Check if we already have this store in our map
		store, exists := storeMap[uint64(item.StoreID)]
		if !exists {
			// Create new store entry
			store = &model.CartStoreDTO{
				StoreID:   uint64(item.StoreID),
				StoreName: item.StoreName,
				StoreURL:  "", // You may need to get this from somewhere else
				List:      []model.CartProductDTO{},
			}
			storeMap[uint64(item.StoreID)] = store
		}

		product := model.CartProductDTO{
			CartID:         item.ID,
			ProductID:      item.ProductID,
			ProductName:    item.ProductName,
			Quantity:       item.Quantity,
			CurrentPrice:   item.CurrentPrice,
			CourierFeeMin:  item.CourierFeeMin,
			MemberDiscount: item.MemberDiscount,
			CouponID:       item.CouponID,
			CouponPrice:    item.CouponPrice,
		}

		store.List = append(store.List, product)
	}

	// Convert map to slice for response
	response := make(model.CartResponse, 0, len(storeMap))
	for _, store := range storeMap {
		response = append(response, *store)
	}

	return response, nil
}

func (r *userCartRepository) DeleteCartItems(ctx context.Context, userID string, cartIDs []uint) error {
	// Delete cart items that belong to the user with the specified IDs
	result := r.DB(ctx).Where("user_id = ? AND id IN (?)", userID, cartIDs).Delete(&model.UserCart{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *userCartRepository) GetCartList(ctx context.Context, userID string, productIds []uint64) ([]model.CartProductDTO, error) {
	// Get all cart items for this user with status = 0 (unprocessed)
	var cartItems []model.UserCart
	if err := r.DB(ctx).Where("user_id = ? AND status = 0", userID).Find(&cartItems).Error; err != nil {
		return nil, err
	}

	// 将cartItems转成map格式
	cartMap := make(map[uint64]model.CartProductDTO)
	for _, item := range cartItems {
		cartMap[item.ProductID] = model.CartProductDTO{
			CartID:         item.ID,
			ProductID:      item.ProductID,
			ProductName:    item.ProductName,
			Quantity:       item.Quantity,
			CurrentPrice:   item.CurrentPrice,
			CourierFeeMin:  item.CourierFeeMin,
			MemberDiscount: item.MemberDiscount,
			CouponID:       item.CouponID,
			CouponPrice:    item.CouponPrice,
		}
	}

	// 遍历productIds
	cartList := make([]model.CartProductDTO, 0, len(productIds))
	for _, productID := range productIds {
		if cart, ok := cartMap[productID]; ok {
			cartList = append(cartList, cart)
		}
	}

	return cartList, nil
}
