package model

import (
	"time"
)

// PointExchangeRecord 积分兑换记录表
type PointExchangeRecord struct {
	ID             uint64    `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	UserID         string    `gorm:"column:user_id;not null;comment:用户ID" json:"user_id"`                 // 用户ID
	ConfigID       int       `gorm:"column:config_id;not null;comment:积分兑换项ID" json:"config_id"`          // 积分兑换项ID
	MinAmount      float64   `gorm:"column:min_amount;not null;comment:可使用的最小金额" json:"min_amount"`       // 可使用的最小金额
	ExchangeAmount float64   `gorm:"column:exchang_amount;not null;comment:兑换的金额" json:"exchange_amount"` // 兑换的金额
	Required       float64   `gorm:"column:required;not null;comment:消费的金额" json:"required"`              // 消费的金额
	Points         int       `gorm:"column:points;not null;comment:发放的积分" json:"points"`                  // 发放的积分
	Type           int8      `gorm:"column:type;not null;comment:类型（1:优惠券；2:兑换券）" json:"type"`            // 类型（1:优惠券；2:兑换券）
	Images         string    `gorm:"column:images;comment:积分兑换图片" json:"images"`                          // 积分兑换图片
	Title          string    `gorm:"column:title;comment:积分兑换标题" json:"title"`                            // 积分兑换标题
	CreatedAt      time.Time `gorm:"column:created_at;not null" json:"created_at"`                        // 创建时间
	UpdatedAt      time.Time `gorm:"column:updated_at;not null" json:"updated_at"`
	Deadline       time.Time `gorm:"column:deadline;not null" json:"deadline"` // 截止时间
}

// TableName 表名
func (p *PointExchangeRecord) TableName() string {
	return "point_exchange_record"
}

// PointExchangeRecordQueryRequest 积分兑换记录查询请求
type PointExchangeRecordQueryRequest struct {
	Page     int `form:"page" json:"page"`           // 页码
	PageSize int `form:"page_size" json:"page_size"` // 每页条数
}

// PointExchangeRecordDTO 积分兑换记录DTO
type PointExchangeRecordDTO struct {
	ID             uint64  `json:"id"`
	UserID         string  `json:"user_id"`
	ConfigID       int     `json:"config_id"`
	MinAmount      float64 `json:"min_amount"`
	ExchangeAmount float64 `json:"exchange_amount"`
	Required       float64 `json:"required"`
	Points         int     `json:"points"`
	Type           int8    `json:"type"`
	Images         string  `json:"images"`
	Title          string  `json:"title"`
	CreatedAt      string  `json:"created_at"`
	Deadline       string  `json:"deadline"`
}

// PointExchangeRecordResponse 积分兑换记录响应
type PointExchangeRecordResponse struct {
	Total    int64                    `json:"total"`
	Records  []PointExchangeRecordDTO `json:"records"`
	Page     int                      `json:"page"`
	PageSize int                      `json:"page_size"`
}
