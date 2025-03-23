package model

import (
	"time"

	"gorm.io/gorm"
)

type ProductCoupon struct {
	gorm.Model
	ProductID         uint64     `gorm:"column:product_id;type:bigint unsigned;not null;comment:商品ID" json:"product_id"`                 // 商品ID
	CouponName        string     `gorm:"column:coupon_name;type:varchar(100);not null;comment:优惠券名称" json:"coupon_name"`                 // 优惠券名称
	CouponPrice       float64    `gorm:"column:coupon_price;type:varchar(50);not null;comment:优惠券价格" json:"coupon_price"`                // 优惠券价格
	AvailableMinPrice float64    `gorm:"column:available_min_price;type:varchar(50);not null;comment:可用最低价格" json:"available_min_price"` // 可用最低价格
	Deadline          *time.Time `gorm:"column:deadline;type:datetime;comment:截止时间" json:"deadline"`                                     // 截止时间
	Status            uint8      `gorm:"column:status;type:tinyint unsigned;not null;default:0;comment:状态(0:未领取;1:已领取)" json:"status"`   // 状态
}

func (m *ProductCoupon) TableName() string {
	return "product_coupon"
}
