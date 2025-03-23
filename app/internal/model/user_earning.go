package model

import "time"

// EarningType 收益类型
type EarningType uint8

const (
	EarningTypeEgg      EarningType = 1 // 鸟蛋
	EarningTypeBird     EarningType = 2 // 商品鸟
	EarningTypeBreeding EarningType = 3 // 种鸟
)

// UserEarning 用户收益模型
type UserEarning struct {
	ID          uint `gorm:"primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	UserID      string      `gorm:"column:user_id;type:varchar(255);not null;comment:用户ID" json:"user_id"`
	EarningType EarningType `gorm:"column:earning_type;type:tinyint;not null;comment:收益类型:1-鸟蛋,2-商品鸟,3-种鸟" json:"earning_type"`
	TypeName    string      `gorm:"column:type_name;type:varchar(255);comment:类型名称" json:"type_name"`
	Image       string      `gorm:"column:image;type:varchar(255);comment:图片" json:"image"`
	Amount      int64       `gorm:"column:amount;type:bigint;not null;comment:收益数量" json:"amount"`
	EarningDate string      `gorm:"column:earning_date;type:varchar(10);comment:收益日期" json:"earning_date"`
	Year        int         `gorm:"column:year;type:int;comment:所属年份" json:"year"`
	Month       int         `gorm:"column:month;type:int;comment:所属月份" json:"month"`
}

// AddEarningRequest 添加收益请求
type AddEarningRequest struct {
	EarningType EarningType `json:"earning_type" binding:"required"` // 收益类型
	Amount      int64       `json:"amount" binding:"required"`       // 收益数量
}

// QueryEarningRequest 查询收益请求
type QueryEarningRequest struct {
	Date     string `form:"date"`                                // 日期，格式：YYYY-MM-DD
	Page     int    `form:"page" binding:"omitempty,min=1"`      // 页码
	PageSize int    `form:"page_size" binding:"omitempty,min=1"` // 每页条数
}

// EarningListResponse 收益列表响应
type EarningListResponse struct {
	List []EarningListItem `json:"list"` // 收益列表
}

type EarningListItem struct {
	ID       EarningType           `json:"id"`        // 收益类型
	TypeName string                `json:"type_name"` // 分类名称
	DateList []EarningListDateItem `json:"date_list"` // 日期列表
}

type EarningListDateItem struct {
	Date       string                      `json:"date"`        // 日期，格式：YYYY-MM
	Amount     int64                       `json:"amount"`      // 收益数量
	AmountItem []EarningListDateAmountItem `json:"amount_item"` // 收益数量
}

type EarningListDateAmountItem struct {
	Name   string `json:"name"`   // 名称
	Image  string `json:"image"`  // 图片
	Date   string `json:"date"`   // 日期，格式：YYYY-MM
	Amount int64  `json:"amount"` // 收益数量
}
