package model

import "time"

type UserCart struct {
	ID             uint    `gorm:"primarykey"`
	UserID         string  `gorm:"column:user_id;type:varchar(255);not null;comment:用户ID" json:"user_id"`                     // 用户ID
	ProductID      uint64  `gorm:"column:product_id;type:bigint unsigned;not null;comment:商品ID" json:"product_id"`            // 商品ID
	ProductName    string  `gorm:"column:product_name;type:varchar(255);not null;comment:商品名称" json:"product_name"`           // 商品名称
	Quantity       int     `gorm:"column:quantity;type:int;not null;default:0;comment:商品数量" json:"quantity"`                  // 商品数量
	CurrentPrice   float64 `gorm:"column:current_price;type:varchar(255);not null;comment:商品当前价格" json:"current_price"`       // 商品当前价格
	Status         *uint8  `gorm:"column:status;type:tinyint;not null;default:0;comment:购物车商品状态（0:未处理； 1:已处理）" json:"status"` // 购物车商品状态
	StoreID        int     `gorm:"column:store_id;type:int;not null;comment:店铺ID" json:"store_id"`
	StoreName      string  `gorm:"column:store_name;type:varchar(255);not null;comment:店铺名称" json:"store_name"`
	CourierFeeMin  float64 `gorm:"column:courier_fee_min;type:int;not null;default:0;comment:快递费最小值" json:"courier_fee_min"`
	MemberDiscount float64 `gorm:"column:member_discount;type:varchar(255);not null;default:0;comment:会员折扣" json:"member_discount"`
	CouponID       uint64  `gorm:"column:coupon_id;type:bigint unsigned;not null;comment:优惠券ID" json:"coupon_id"`
	CouponPrice    float64 `gorm:"column:coupon_price;comment:优惠券价格" json:"coupon_price"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (m *UserCart) TableName() string {
	return "user_carts"
}

// AddToCartRequest defines the request structure for adding to cart
type AddToCartRequest struct {
	ProductID uint64 `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required"`
	CouponID  uint64 `json:"coupon_id"`
}

// DeleteCartItemsRequest represents the request to delete items from cart
type DeleteCartItemsRequest struct {
	IDs string `json:"ids" binding:"required"` // Comma-separated list of cart IDs
}
