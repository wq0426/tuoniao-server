package model

import (
	"time"
)

// UserOrder represents an order
type UserOrder struct {
	ID            uint64          `gorm:"primarykey" json:"id"`
	OrderNo       string          `gorm:"column:order_no" json:"order_no"`
	UserID        string          `gorm:"column:user_id" json:"user_id"`
	PaymentMethod int             `gorm:"column:payment_method;type:tinyint(4);not null;comment:支付方式" json:"payment_method"`           // 支付方式
	AddressID     int             `gorm:"column:address_id;type:int;not null;comment:地址ID" json:"address_id"`                          // 地址ID
	Name          string          `gorm:"column:name;type:varchar(30);not null;comment:收件人" json:"name"`                               // 收件人
	Phone         string          `gorm:"column:phone;type:varchar(20);not null;comment:手机号" json:"phone"`                             // 手机号
	Province      string          `gorm:"column:province;type:varchar(255);not null;comment:省" json:"province"`                        // 省
	City          string          `gorm:"column:city;type:varchar(255);not null;comment:城市" json:"city"`                               // 城市
	District      string          `gorm:"column:district;type:varchar(255);not null;comment:区/县" json:"district"`                      // 区/县
	Detail        string          `gorm:"column:detail;type:varchar(255);not null;comment:详细地址" json:"detail"`                         // 详细地址
	IsDefault     uint8           `gorm:"column:is_default;type:tinyint;not null;default:0;comment:是否默认地址（0:否；1:是）" json:"is_default"` // 是否默认地址
	TotalFee      float64         `gorm:"column:total_fee;type:int(13);not null;comment:总金额" json:"total_fee"`                         // 总金额
	CreatedAt     *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	OrderItems    []UserOrderItem `gorm:"foreignKey:OrderID;references:ID" json:"order_items"`
	OrderAddress  *UserAddress    `gorm:"foreignKey:ID;references:AddressID" json:"order_address"`
}

func (m *UserOrder) TableName() string {
	return "user_order"
}

type OrderItemRequest struct {
	CartID         uint    `json:"cart_id"`
	ProductID      uint64  `json:"product_id" binding:"required"`
	Quantity       int     `json:"quantity" binding:"required"`
	ProductName    string  `json:"product_name" binding:"required"`
	Image          string  `json:"image" binding:"required"`
	CurrentPrice   float64 `json:"current_price" binding:"required"`
	CourierFeeMin  float64 `json:"courier_fee_min"`
	MemberDiscount float64 `json:"member_discount"`
	StoreID        uint64  `json:"store_id"`
	StoreName      string  `json:"store_name"`
	Note           string  `json:"note"`
	CouponID       uint64  `json:"coupon_id"`
	CouponPrice    float64 `json:"coupon_price"`
}

type Address struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Province  string `json:"province"`
	City      string `json:"city"`
	District  string `json:"district"`
	Detail    string `json:"detail"`
	IsDefault bool   `json:"isDefault"`
}

// CreateOrderRequest represents the request to create multiple orders
type CreateOrderRequest struct {
	Items         []OrderItemRequest `json:"items" binding:"required,dive"`
	Address       Address            `json:"address" binding:"required"`
	PaymentMethod int                `json:"payment_method" binding:"required"`
	Status        int                `json:"status"`
}

// OrderQueryRequest 订单查询请求
type OrderQueryRequest struct {
	OrderNo  *string `form:"order_no" json:"order_no"`   // 订单号
	Status   *int    `form:"status" json:"status"`       // 订单状态，为空则查询所有状态
	Page     int     `form:"page" json:"page"`           // 页码
	PageSize int     `form:"page_size" json:"page_size"` // 每页条数
}

// OrderListResponse 订单列表响应
type OrderListResponse struct {
	Total    int64          `json:"total"`
	Orders   []OrderListDTO `json:"orders"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
}

// OrderListDTO 订单列表数据传输对象
type OrderListDTO struct {
	ID            uint64          `json:"id"`
	OrderNo       string          `json:"order_no"`
	UserID        string          `json:"user_id"`
	TotalFee      float64         `json:"total_fee"`
	Status        int             `json:"status"`
	StatusText    string          `json:"status_text"`
	PaymentMethod int             `json:"payment_method"`
	PayTime       *string         `json:"pay_time,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
	Address       string          `json:"address"`
	Product       OrderProductDTO `json:"product"`
	StoreName     string          `json:"store_name"` // 店铺名称
	StoreIcon     string          `json:"store_icon"` // 店铺图标
}

// OrderProductDTO 订单商品数据传输对象
type OrderProductDTO struct {
	ItemID      uint64  `json:"item_id"`
	ProductID   uint64  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	Image       string  `json:"image"`
}

// OrderProductsRequest 订单商品查询请求
type OrderProductsRequest struct {
	OrderItemID uint64 `form:"order_item_id" binding:"required"` // 订单商品ID
}

// OrderProductsResponse 订单商品查询响应
type OrderProductsResponse struct {
	Products []ProductListItemDTO `json:"products"` // 商品列表
}

// UpdateOrderStatusRequest 更新订单状态请求
type UpdateOrderStatusRequest struct {
	OrderItemID uint64 `json:"order_id" binding:"required"` // 订单ID
	Status      uint8  `json:"status"`                      // 订单状态
}

// OrderDetailResponse 订单详情响应
type OrderDetailResponse struct {
	// 订单基本信息
	OrderID       uint64  `json:"order_id"`        // 订单ID
	OrderNo       string  `json:"order_no"`        // 订单号
	UserID        string  `json:"user_id"`         // 用户ID
	OrderAmount   float64 `json:"order_amount"`    // 订单金额
	Status        int     `json:"status"`          // 订单状态
	StatusText    string  `json:"status_text"`     // 订单状态文本
	PayMethod     int     `json:"pay_method"`      // 支付方式
	PayMethodText string  `json:"pay_method_text"` // 支付方式文本
	PayTime       string  `json:"pay_time"`        // 支付时间
	ShippedAt     string  `json:"shipped_at"`      // 发货时间
	CompletedAt   string  `json:"completed_at"`    // 完成时间
	CreatedAt     string  `json:"created_at"`      // 创建时间
	UpdatedAt     string  `json:"updated_at"`      // 更新时间
	StoreID       uint64  `json:"store_id"`        // 店铺ID
	StoreName     string  `json:"store_name"`      // 店铺名称
	StoreLogo     string  `json:"store_logo"`      // 店铺Logo

	// 地址信息
	Address *AddressInfo `json:"address"` // 收货地址信息

	// 商品信息
	Products []ProductListItemDTO `json:"products"` // 商品列表

	// 订单汇总信息
	TotalPrice     float64 `json:"total_price"`     // 商品总价
	ShippingFee    float64 `json:"shipping_fee"`    // 运费
	CouponDiscount float64 `json:"coupon_discount"` // 优惠券折扣
	MemberDiscount float64 `json:"member_discount"` // 会员折扣
	ActualAmount   float64 `json:"actual_amount"`   // 实际支付金额
}

// AddressInfo 地址信息
type AddressInfo struct {
	ReceiverName  string `json:"receiver_name"`  // 收货人姓名
	ReceiverPhone string `json:"receiver_phone"` // 收货人电话
	Province      string `json:"province"`       // 省份
	City          string `json:"city"`           // 城市
	District      string `json:"district"`       // 区县
	DetailAddress string `json:"detail_address"` // 详细地址
}
