package model

import (
	"time"
)

// PointExchangeConfig 积分兑换配置表
type PointExchangeConfig struct {
	ID             uint64    `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	MinAmount      float64   `gorm:"column:min_amount;not null;comment:可使用的最小金额" json:"min_amount"`       // 可使用的最小金额
	ExchangeAmount float64   `gorm:"column:exchang_amount;not null;comment:兑换的金额" json:"exchange_amount"` // 兑换的金额
	Required       float64   `gorm:"column:required;not null;comment:消费的金额" json:"required"`              // 消费的金额
	Points         int       `gorm:"column:points;not null;comment:发放的积分" json:"points"`                  // 发放的积分
	Type           int8      `gorm:"column:type;not null;comment:类型（1:优惠券；2:兑换券）" json:"type"`            // 类型（1:优惠券；2:兑换券）
	Images         string    `gorm:"column:images;comment:积分兑换图片" json:"images"`                          // 积分兑换图片
	Title          string    `gorm:"column:title;comment:积分兑换标题" json:"title"`                            // 积分兑换标题
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`                                 // 创建时间
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updated_at"`                                 // 更新时间
}

// PointExchangeConfigResponse 积分兑换配置响应
type PointExchangeConfigResponse struct {
	PointExchangeConfig
	IsExchange bool `json:"is_exchange"` // 是否可兑换
}

// TableName 表名
func (p *PointExchangeConfig) TableName() string {
	return "point_exchange_config"
}

// PointExchangeConfigListResponse 积分兑换配置列表响应
type PointExchangeConfigListResponse struct {
	Total int                           `json:"total"`
	List  []PointExchangeConfigResponse `json:"list"`
}

// PointExchangeRequest 积分兑换请求
type PointExchangeRequest struct {
	ConfigID uint64 `json:"config_id" binding:"required"` // 兑换配置ID
}
