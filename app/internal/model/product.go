package model

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	ProductName    string  `gorm:"column:product_name;type:varchar(255);not null;comment:商品名称" json:"product_name"`      // 商品名称
	Category1ID    *int    `gorm:"column:category1_id;type:int" json:"category1_id"`                                     // 一级分类ID
	Category2ID    *int    `gorm:"column:category2_id;type:int" json:"category2_id"`                                     // 二级分类ID
	Content        string  `gorm:"column:content;type:text;not null;comment:内容" json:"content"`                          // 内容
	HeaderImg      string  `gorm:"column:header_img;type:varchar(255);not null;comment:头部图片链接" json:"header_img"`        // 头部图片链接
	BannerImg      string  `gorm:"column:banner_img;type:text;not null;comment:banner图片链接" json:"banner_img"`            // banner图片链接
	CurrentPrice   float64 `gorm:"column:current_price;type:varchar(30);not null;comment:商品当前价格" json:"current_price"`   // 商品当前价格
	OriginPrice    float64 `gorm:"column:origin_price;type:varchar(30);not null;comment:商品原始价格" json:"origin_price"`     // 商品原始价格
	MemberDiscount float64 `gorm:"column:member_discount;type:varchar(30);not null;comment:会员折扣" json:"member_discount"` // 会员折扣
	IsSpecial      *int    `gorm:"column:is_special;type:int;comment:是否特享价（0:否,1:是）" json:"is_special"`                  // 是否特享价
	Sales          int     `gorm:"column:sales;type:int;default:0;comment:销量" json:"sales"`                              // 销量
	Specification  string  `gorm:"column:specification;type:varchar(50);not null;comment:规格" json:"specification"`       // 规格
	CourierFeeMin  float64 `gorm:"column:courier_fee_min;type:int;not null;comment:最低快递费" json:"courier_fee_min"`        // 最低快递费
	CourierFeeMax  float64 `gorm:"column:courier_fee_max;type:int;not null;comment:最高快递费" json:"courier_fee_max"`        // 最高快递费
	IsRecommend    *int    `gorm:"column:is_recommend;type:int;comment:是否推荐（0:否,1:是）" json:"is_recommend"`               // 是否推荐
	RecommendSort  int     `gorm:"column:recommend_sort;type:int;default:0;comment:推荐排序" json:"recommend_sort"`          // 推荐排序
	StoreID        int     `gorm:"column:store_id;type:int;not null;comment:店铺ID" json:"store_id"`
	StoreName      string  `gorm:"column:store_name;type:varchar(255);not null;comment:店铺名称" json:"store_name"`
	StoreIcon      string  `gorm:"column:store_icon;type:varchar(255);not null;comment:店铺图标" json:"store_icon"`
}

func (m *Product) TableName() string {
	return "product"
}

// ProductDetailsByCartRequest 通过购物车ID查询商品详情请求
type ProductDetailsByCartRequest struct {
	CartIDs []uint64 `json:"cart_ids" binding:"required"` // 购物车ID列表
}

// ProductDetailsResponse 商品详情响应
type ProductDetailsResponse struct {
	Products []ProductListItemDTO `json:"products"` // 商品列表
}
