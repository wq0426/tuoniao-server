package model

import (
	"time"
)

// RefundOrder 退款单模型
type RefundOrder struct {
	ID             uint64     `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	RefundNo       string     `gorm:"column:refund_no;type:varchar(255);not null;comment:退款单号" json:"refund_no"`                   // 退款单号
	UserID         string     `gorm:"column:user_id;type:varchar(255);not null;comment:用户ID" json:"user_id"`                       // 用户ID
	OrderID        uint64     `gorm:"column:order_id;type:bigint unsigned;not null;comment:订单ID" json:"order_id"`                  // 订单ID
	OrderItemID    uint64     `gorm:"column:order_item_id;type:bigint unsigned;not null;comment:订单项ID" json:"order_item_id"`       // 订单项ID
	OrderNo        string     `gorm:"column:order_no;type:varchar(255);not null;comment:订单号" json:"order_no"`                      // 订单号
	RefundAmount   float64    `gorm:"column:refund_amount;not null;comment:退款金额" json:"refund_amount"`                             // 退款金额
	RefundReason   string     `gorm:"column:refund_reason;type:varchar(500);comment:退款原因" json:"refund_reason"`                    // 退款原因
	RefundType     uint8      `gorm:"column:refund_type;type:tinyint;not null;comment:退款类型(1:仅退款;2:退货退款)" json:"refund_type"`      // 退款类型
	Status         uint8      `gorm:"column:status;type:tinyint;not null;default:0;comment:退款状态(0:退款中;1:已退款;2:已拒绝)" json:"status"` // 退款状态
	OriginStatus   uint8      `gorm:"column:origin_status;type:tinyint;not null;default:0;comment:原始状态" json:"origin_status"`      // 原始状态
	RejectReason   string     `gorm:"column:reject_reason;type:varchar(500);comment:拒绝原因" json:"reject_reason"`                    // 拒绝原因
	Images         string     `gorm:"column:images;type:text;comment:图片凭证，多个图片用逗号分隔" json:"images"`                                // 图片凭证
	ApplyTime      string     `gorm:"column:apply_time;comment:申请时间" json:"apply_time"`                                            // 申请时间
	ProcessTime    string     `gorm:"column:process_time;comment:处理时间" json:"process_time"`                                        // 处理时间
	CompletionTime string     `gorm:"column:completion_time;comment:完成时间" json:"completion_time"`                                  // 完成时间
	CreatedAt      *time.Time `gorm:"column:created_at;comment:创建时间" json:"created_at"`                                            // 创建时间
	UpdatedAt      *time.Time `gorm:"column:updated_at;comment:更新时间" json:"updated_at"`                                            // 更新时间
	ProductID      uint64     `gorm:"column:product_id;type:bigint unsigned;not null;comment:商品ID" json:"product_id"`              // 商品ID
	ProductName    string     `gorm:"column:product_name;type:varchar(255);not null;comment:商品名称" json:"product_name"`             // 商品名称
	HeaderImg      string     `gorm:"column:header_img;type:varchar(255);comment:商品图片" json:"header_img"`                          // 商品图片
	Quantity       int        `gorm:"column:quantity;type:int;not null;comment:退款数量" json:"quantity"`                              // 退款数量
	Price          float64    `gorm:"column:price;not null;comment:商品单价" json:"price"`                                             // 商品单价
	StoreName      string     `gorm:"column:store_name;type:varchar(255);not null;comment:店铺名称" json:"store_name"`                 // 店铺名称
	StoreIcon      string     `gorm:"column:store_icon;type:varchar(255);not null;comment:店铺图标" json:"store_icon"`                 // 店铺图标
}

func (m *RefundOrder) TableName() string {
	return "refund_order"
}

// RefundItem 退款商品项模型
type RefundItem struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	RefundID     uint64    `gorm:"column:refund_id;type:bigint unsigned;not null;comment:退款单ID" json:"refund_id"`         // 退款单ID
	OrderItemID  uint64    `gorm:"column:order_item_id;type:bigint unsigned;not null;comment:订单项ID" json:"order_item_id"` // 订单项ID
	ProductID    uint64    `gorm:"column:product_id;type:bigint unsigned;not null;comment:商品ID" json:"product_id"`        // 商品ID
	ProductName  string    `gorm:"column:product_name;type:varchar(255);not null;comment:商品名称" json:"product_name"`       // 商品名称
	HeaderImg    string    `gorm:"column:header_img;type:varchar(255);comment:商品图片" json:"header_img"`                    // 商品图片
	Quantity     int       `gorm:"column:quantity;type:int;not null;comment:退款数量" json:"quantity"`                        // 退款数量
	Price        float64   `gorm:"column:price;not null;comment:商品单价" json:"price"`                                       // 商品单价
	RefundAmount float64   `gorm:"column:refund_amount;not null;comment:退款金额" json:"refund_amount"`                       // 退款金额
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`                                                   // 创建时间
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`                                                   // 更新时间
}

func (m *RefundItem) TableName() string {
	return "refund_item"
}

// CreateRefundRequest 创建退款请求
// type CreateRefundRequest struct {
// 	OrderID      uint64              `json:"order_id" binding:"required"`              // 订单ID
// 	RefundAmount float64             `json:"refund_amount" binding:"required"`         // 退款金额
// 	RefundReason string              `json:"refund_reason" binding:"required"`         // 退款原因
// 	RefundType   uint8               `json:"refund_type" binding:"required,oneof=1 2"` // 退款类型(1:仅退款;2:退货退款)
// 	Images       string              `json:"images"`                                   // 图片凭证
// 	RefundItems  []RefundItemRequest `json:"refund_items" binding:"required,min=1"`    // 退款商品项
// }

// RefundItemRequest 退款商品项请求
type RefundItemRequest struct {
	OrderItemID  uint64 `json:"order_item_id" binding:"required"`         // 订单项ID
	RefundReason string `json:"refund_reason" binding:"required"`         // 退款原因
	RefundType   uint8  `json:"refund_type" binding:"required,oneof=1 2"` // 退款类型(1:仅退款;2:退货退款)
	Images       string `json:"images"`                                   // 图片凭证
}

// RefundQueryRequest 退款查询请求
type RefundQueryRequest struct {
	Status   *uint8 `form:"status"`                     // 退款状态
	Page     int    `form:"page" json:"page"`           // 页码
	PageSize int    `form:"page_size" json:"page_size"` // 每页条数
}

// RefundDetailResponse 退款详情响应
type RefundDetailResponse struct {
	Refund      RefundOrder  `json:"refund"`       // 退款单信息
	OrderDetail OrderDetails `json:"order_detail"` // 订单详情
}

type OrderDetails struct {
	OrderID      uint64     `json:"order_id"`      // 订单ID
	OrderNo      string     `json:"order_no"`      // 订单号
	OrderStatus  uint8      `json:"order_status"`  // 订单状态
	OrderCreated *time.Time `json:"order_created"` // 订单创建时间
}

// RefundListItem 退款列表项
type RefundListItem struct {
	ID           uint64     `json:"id"`            // 退款单ID
	RefundNo     string     `json:"refund_no"`     // 退款单号
	OrderID      uint64     `json:"order_id"`      // 订单ID
	OrderNo      string     `json:"order_no"`      // 订单号
	RefundAmount float64    `json:"refund_amount"` // 退款金额
	RefundType   uint8      `json:"refund_type"`   // 退款类型
	Status       uint8      `json:"status"`        // 退款状态
	StatusText   string     `json:"status_text"`   // 退款状态文本
	ApplyTime    string     `json:"apply_time"`    // 申请时间
	CreatedAt    *time.Time `json:"created_at"`    // 创建时间
	StoreName    string     `json:"store_name"`    // 店铺名称
	StoreIcon    string     `json:"store_icon"`    // 店铺图标
}

// RefundListResponse 退款列表响应
type RefundListResponse struct {
	Total int64             `json:"total"` // 总数
	List  []RefundOrderItem `json:"list"`  // 退款列表
	Page  int               `json:"page"`  // 页码
	Size  int               `json:"size"`  // 每页条数
}

type RefundOrderItem struct {
	RefundListItem
	ProductName string  `json:"product_name"` // 商品名称
	HeaderImg   string  `json:"header_img"`   // 商品图片
	Price       float64 `json:"price"`        // 商品单价
	Quantity    int     `json:"quantity"`     // 退款数量
}

// CancelRefundRequest 取消退款请求
type CancelRefundRequest struct {
	RefundID uint64 `json:"refund_id" binding:"required"` // 退款单ID
}
