package model

import (
	"time"
)

type UserOrderItem struct {
	ID             uint64     `gorm:"primarykey"`
	OrderID        uint64     `gorm:"column:order_id;type:varchar(255);not null;comment:订单ID" json:"order_id"`         // 订单ID
	OrderNo        string     `gorm:"column:order_no;type:varchar(255);not null;comment:订单号" json:"order_no"`          // 订单号
	ProductID      uint64     `gorm:"column:product_id;type:bigint unsigned;not null;comment:商品ID" json:"product_id"`  // 商品ID
	Quantity       int        `gorm:"column:quantity;type:int;not null;default:0;comment:商品数量" json:"quantity"`        // 商品数量
	ProductName    string     `gorm:"column:product_name;type:varchar(255);not null;comment:商品名称" json:"product_name"` // 商品名称
	HeaderImg      string     `gorm:"column:header_img;type:varchar(255);not null;comment:商品头部图片" json:"header_img"`   // 商品头部图片
	StoreID        uint64     `gorm:"column:store_id;type:bigint unsigned;not null;comment:店铺ID" json:"store_id"`      // 店铺ID
	StoreName      string     `gorm:"column:store_name;type:varchar(255);not null;comment:店铺名称" json:"store_name"`     // 店铺名称
	StoreLogo      string     `gorm:"column:store_logo;comment:店铺LOGO" json:"store_logo"`                              // 店铺LOGO
	CreatedAt      time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at" json:"updated_at"`
	CurrentPrice   float64    `gorm:"column:current_price;type:varchar(30);not null;comment:商品当前价格" json:"current_price"`   // 商品当前价格
	CourierFeeMin  float64    `gorm:"column:courier_fee_min;type:varchar(30);not null;comment:运费价格" json:"courier_fee_min"` // 运费价格
	MemberDiscount float64    `gorm:"column:member_discount;type:varchar(30);not null;comment:会员价格" json:"member_discount"` // 会员价格
	Note           string     `gorm:"column:note;type:varchar(255);default:'';comment:备注" json:"note"`                      // 备注
	UserId         string     `gorm:"column:user_id;type:varchar(255);not null;comment:用户ID" json:"user_id"`                // 用户ID
	Category1Id    int        `gorm:"column:category1_id;type:varchar(255);not null;comment:一级分类ID" json:"category1_id"`    // 一级分类ID
	Category2Id    int        `gorm:"column:category2_id;type:varchar(255);not null;comment:二级分类ID" json:"category2_id"`    // 二级分类ID                   // 三级分类ID
	CouponID       uint64     `gorm:"column:coupon_id;type:bigint unsigned;not null;comment:优惠券ID" json:"coupon_id"`        // 优惠券ID
	CouponPrice    float64    `gorm:"column:coupon_price;comment:优惠券价格" json:"coupon_price"`                                // 优惠券价格
	TotalFee       float64    `gorm:"column:total_fee;comment:总金额" json:"total_fee"`
	Status         uint8      `gorm:"column:status;type:tinyint;not null;default:0;comment:购物车商品状态（0:待支付;1:待发货;2:待收货;3:待评价;4:已完成;5:已取消）" json:"status"` // 订单状态                                        // 总金额
	PayTime        *time.Time `gorm:"column:pay_time" json:"pay_time"`
	ShippedAt      *time.Time `gorm:"column:shipped_at" json:"shipped_at"`
	CompletedAt    *time.Time `gorm:"column:completed_at" json:"completed_at"`
}

func (m *UserOrderItem) TableName() string {
	return "user_order_item"
}
