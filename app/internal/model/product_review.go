package model

import (
	"time"
)

// ProductReview 商品评价模型
type ProductReview struct {
	ID              uint64     `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	UserID          string     `gorm:"column:user_id;type:varchar(255);not null;comment:用户ID" json:"user_id"`
	OrderID         uint64     `gorm:"column:order_id;not null;comment:订单ID" json:"order_id"`
	OrderNo         string     `gorm:"column:order_no;not null;comment:订单编号" json:"order_no"`
	OrderItemID     uint64     `gorm:"column:order_item_id;not null;comment:订单项ID" json:"order_item_id"`
	ProductID       uint64     `gorm:"column:product_id;not null;comment:商品ID" json:"product_id"`
	Rating          uint8      `gorm:"column:rating;not null;comment:商品评分(1-5)" json:"rating"`
	FreshnessRating uint8      `gorm:"column:freshness_rating;not null;comment:新鲜程度评分(1-5)" json:"freshness_rating"`
	PackagingRating uint8      `gorm:"column:packaging_rating;not null;comment:包装评分(1-5)" json:"packaging_rating"`
	DeliveryRating  uint8      `gorm:"column:delivery_rating;not null;comment:配送评分(1-5)" json:"delivery_rating"`
	ServiceRating   uint8      `gorm:"column:service_rating;not null;comment:服务态度评分(1-5)" json:"service_rating"`
	Content         string     `gorm:"column:content;type:text;comment:评价内容" json:"content"`
	Images          string     `gorm:"column:images;type:text;comment:评价图片,多个用逗号分隔" json:"images"`
	IsAnonymous     bool       `gorm:"column:is_anonymous;not null;default:false;comment:是否匿名评价" json:"is_anonymous"`
	Status          uint8      `gorm:"column:status;not null;default:1;comment:状态(0:待审核;1:已通过;2:已拒绝)" json:"status"`
	IsReturn        int8       `gorm:"column:is_return;type:tinyint;comment:是否是回头客（0:不是;1:是）" json:"is_return"`  // 是否是回头客
	ViewNums        uint       `gorm:"column:view_nums;type:int unsigned;comment:查看人数" json:"view_nums"`         // 查看人数
	EvaluateNums    uint       `gorm:"column:evaluate_nums;type:int unsigned;comment:评价人数" json:"evaluate_nums"` // 评价人数
	PraiseNums      uint       `gorm:"column:praise_nums;type:int unsigned;comment:点赞人数" json:"praise_nums"`     // 点赞人数
	CreatedAt       *time.Time `gorm:"column:created_at;comment:创建时间" json:"created_at"`
	UpdatedAt       *time.Time `gorm:"column:updated_at;comment:更新时间" json:"updated_at"`
}

func (m *ProductReview) TableName() string {
	return "product_review"
}

// CreateReviewRequest 创建评价请求
type CreateReviewRequest struct {
	OrderItemID     uint64   `json:"order_item_id" binding:"required"`                // 订单ID
	ProductID       uint64   `json:"product_id" binding:"required"`                   // 商品ID
	Rating          uint8    `json:"rating" binding:"required,min=1,max=5"`           // 商品评分
	FreshnessRating uint8    `json:"freshness_rating" binding:"required,min=1,max=5"` // 新鲜程度评分
	PackagingRating uint8    `json:"packaging_rating" binding:"required,min=1,max=5"` // 包装评分
	DeliveryRating  uint8    `json:"delivery_rating" binding:"required,min=1,max=5"`  // 配送评分
	ServiceRating   uint8    `json:"service_rating" binding:"required,min=1,max=5"`   // 服务态度评分
	Content         string   `json:"content"`                                         // 评价内容
	Images          []string `json:"images"`                                          // 评价图片数组
	IsAnonymous     bool     `json:"is_anonymous"`                                    // 是否匿名评价
}

// ReviewListResponse 评价列表响应
type ReviewListResponse struct {
	Total   int64           `json:"total"`   // 总数
	Reviews []ReviewListDTO `json:"reviews"` // 评价列表
}

// ReviewListDTO 评价列表数据传输对象
type ReviewListDTO struct {
	ID            uint64   `json:"id"`             // 评价ID
	OrderID       uint64   `json:"order_id"`       // 订单ID
	OrderNo       string   `json:"order_no"`       // 订单编号
	OrderItemID   uint64   `json:"order_item_id"`  // 订单项ID
	ProductID     uint64   `json:"product_id"`     // 商品ID
	StoreName     string   `json:"store_name"`     // 店铺名称
	StoreIcon     string   `json:"store_icon"`     // 店铺图标
	StatusText    string   `json:"status_text"`    // 订单状态文本
	ProductName   string   `json:"product_name"`   // 商品名称
	ProductImage  string   `json:"product_image"`  // 商品图片
	Price         float64  `json:"price"`          // 商品价格
	Quantity      int      `json:"quantity"`       // 购买数量
	ReviewStatus  uint8    `json:"review_status"`  // 评价状态
	OrderTime     string   `json:"order_time"`     // 订单时间
	ReviewContent string   `json:"review_content"` // 评价内容
	ReviewImages  []string `json:"review_images"`  // 评价图片
	ViewNums      uint     `json:"view_nums"`      // 查看人数
	EvaluateNums  uint     `json:"evaluate_nums"`  // 评价人数
	PraiseNums    uint     `json:"praise_nums"`    // 点赞人数
	IsAnonymous   bool     `json:"is_anonymous"`   // 是否匿名
	UserName      string   `json:"user_name"`      // 用户名
	UserAvatar    string   `json:"user_avatar"`    // 用户头像
}

type ReviewDoneListDTO struct {
	ID            uint64   `json:"id"`            // 评价ID
	Status        uint8    `json:"status"`        // 评价状态
	CreatedAt     string   `json:"createdAt"`     // 创建时间
	ProductID     uint64   `json:"productId"`     // 商品IDs
	ProductName   string   `json:"productName"`   // 商品名称
	ProductImage  string   `json:"productImage"`  // 商品图片
	Price         float64  `json:"price"`         // 商品价格
	Quantity      int      `json:"quantity"`      // 购买数量
	ReviewContent string   `json:"reviewContent"` // 评价内容
	ReviewImages  []string `json:"reviewImages"`  // 评价图片
	ReviewStatus  uint8    `json:"reviewStatus"`  // 评价状态
	StatusText    string   `json:"statusText"`    // 评价状态文本
	StoreName     string   `json:"storeName"`     // 店铺名称
	StoreIcon     string   `json:"storeIcon"`     // 店铺图标
	IsAnonymous   bool     `json:"isAnonymous"`   // 是否匿名
}

// ReviewTabType 评价标签类型
type ReviewTabType uint8

const (
	ReviewTabAll       ReviewTabType = 0 // 全部评价
	ReviewTabPending   ReviewTabType = 1 // 待评价
	ReviewTabCompleted ReviewTabType = 2 // 已评价
)

// UserReviewListRequest 用户评价列表请求
type UserReviewListRequest struct {
	Tab ReviewTabType `form:"tab" json:"tab"` // 评价标签类型
}

// PendingReviewItem 待评价项目
type PendingReviewItem struct {
	OrderID         uint64     `json:"order_id"`          // 订单ID
	OrderNo         string     `json:"order_no"`          // 订单编号
	ProductID       uint64     `json:"product_id"`        // 商品ID
	ProductName     string     `json:"product_name"`      // 商品名称
	ProductImage    string     `json:"product_image"`     // 商品图片
	ProductSpec     string     `json:"product_spec"`      // 商品规格
	Price           float64    `json:"price"`             // 商品价格
	Quantity        int        `json:"quantity"`          // 购买数量
	OrderCreatedAt  *time.Time `json:"order_created_at"`  // 订单创建时间
	OrderStatus     int        `json:"order_status"`      // 订单状态
	OrderStatusText string     `json:"order_status_text"` // 订单状态文本
	StoreName       string     `json:"store_name"`        // 店铺名称
	StoreIcon       string     `json:"store_icon"`        // 店铺图标
}

// PendingReviewListResponse 待评价列表响应
type PendingReviewListResponse struct {
	Total   int64               `json:"total"`   // 总数
	Reviews []PendingReviewItem `json:"reviews"` // 待评价列表
}

// UserReviewListResponse 用户评价列表响应
type UserReviewListResponse struct {
	Tab  ReviewTabType `json:"tab"`  // 当前标签
	List []any         `json:"list"` // 全部评价列表
}

// ReviewDetailResponse 评价详情响应
type ReviewDetailResponse struct {
	ID              uint64                `json:"id"`               // 评价ID
	UserID          string                `json:"user_id"`          // 用户ID
	UserName        string                `json:"user_name"`        // 用户名称
	UserAvatar      string                `json:"user_avatar"`      // 用户头像
	OrderID         uint64                `json:"order_id"`         // 订单ID
	OrderNo         string                `json:"order_no"`         // 订单编号
	OrderItemID     uint64                `json:"order_item_id"`    // 订单项ID
	ProductID       uint64                `json:"product_id"`       // 商品ID
	ProductName     string                `json:"product_name"`     // 商品名称
	ProductImage    string                `json:"product_image"`    // 商品图片
	ProductSpec     string                `json:"product_spec"`     // 商品规格
	Price           float64               `json:"price"`            // 商品价格
	Quantity        int                   `json:"quantity"`         // 购买数量
	StoreName       string                `json:"store_name"`       // 店铺名称
	StoreIcon       string                `json:"store_icon"`       // 店铺图标
	Rating          uint8                 `json:"rating"`           // 商品评分
	FreshnessRating uint8                 `json:"freshness_rating"` // 新鲜程度评分
	PackagingRating uint8                 `json:"packaging_rating"` // 包装评分
	DeliveryRating  uint8                 `json:"delivery_rating"`  // 配送评分
	ServiceRating   uint8                 `json:"service_rating"`   // 服务态度评分
	Content         string                `json:"content"`          // 评价内容
	Images          []string              `json:"images"`           // 评价图片
	IsAnonymous     bool                  `json:"is_anonymous"`     // 是否匿名评价
	Status          uint8                 `json:"status"`           // 状态
	StatusText      string                `json:"status_text"`      // 状态文本
	CreatedAt       string                `json:"created_at"`       // 创建时间
	CommentNums     uint                  `json:"comment_nums"`     // 评论数量
	EvaluateList    []ProductEvaluateList `json:"evaluate_list"`    // 评论和回复列表
}

// UpdateReviewCounterRequest 更新评价计数请求
type UpdateReviewCounterRequest struct {
	ReviewID uint64 `form:"review_id" json:"review_id" binding:"required"` // 评价ID
	Type     uint8  `form:"type" json:"type" binding:"oneof=0 1 2 3"`      // 计数类型：0-浏览，1-评论，2-点赞，3-取消点赞
}
