package model

import "time"

// UserCoupon represents the user_coupon table
type UserCoupon struct {
	ID                uint64    `gorm:"primarykey" json:"id"`
	UserID            string    `gorm:"column:user_id;type:varchar(255);not null;comment:用户ID" json:"user_id"`                     // 用户ID
	ProductID         uint64    `gorm:"column:product_id;type:bigint unsigned;not null;comment:商品ID" json:"product_id"`            // 商品ID
	CouponID          uint64    `gorm:"column:coupon_id;type:bigint unsigned;not null;comment:优惠券ID" json:"coupon_id"`             // 优惠券ID (注意:列名在DB中有拼写错误)
	Type              uint8     `gorm:"column:type;type:tinyint;not null;default:0;comment:类型（1:优惠券；2:兑换券）" json:"type"`           // 类型
	Status            uint8     `gorm:"column:status;type:tinyint;not null;default:0;comment:状态（0:未处理:1:已使用;2:已过期）" json:"status"` // 状态
	CreatedAt         time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt         time.Time `gorm:"column:updated_at" json:"updated_at"`
	CouponName        string    `gorm:"column:coupon_name;type:varchar(100);not null;comment:优惠券名称" json:"coupon_name"`                 // 优惠券名称
	CouponPrice       float64   `gorm:"column:coupon_price;type:varchar(50);not null;comment:优惠券价格" json:"coupon_price"`                // 优惠券价格
	AvailableMinPrice float64   `gorm:"column:available_min_price;type:varchar(50);not null;comment:可用最低价格" json:"available_min_price"` // 可用最低价格
	Deadline          time.Time `gorm:"column:deadline;comment:截止时间" json:"deadline"`                                                   // 截止时间
}

// TableName specifies the table name for UserCoupon
func (m *UserCoupon) TableName() string {
	return "user_coupon"
}

// UserCouponResponse represents the response for user coupon queries
type UserCouponResponse struct {
	ID        uint64    `json:"id"`
	UserID    string    `json:"user_id"`
	ProductID uint64    `json:"product_id"`
	CouponID  uint64    `json:"coupon_id"`
	Status    uint8     `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ClaimCouponRequest 领取优惠券的请求
type ClaimCouponRequest struct {
	CouponID uint64 `json:"coupon_id" binding:"required"` // 优惠券ID
}

// CouponDetailResponse 优惠券详情响应
type CouponDetailResponse struct {
	ID                uint64    `json:"id"`
	CouponName        string    `json:"coupon_name"`         // 优惠券名称
	CouponPrice       float64   `json:"coupon_price"`        // 优惠券金额
	AvailableMinPrice float64   `json:"available_min_price"` // 可用最低金额
	Type              uint8     `json:"type"`                // 类型（1:优惠券；2:兑换券）
	Status            uint8     `json:"status"`              // 状态
	ProductID         uint64    `json:"product_id"`          // 商品ID
	Deadline          time.Time `json:"deadline"`            // 截止时间
	CreatedAt         time.Time `json:"created_at"`          // 创建时间
}
