package model

import (
	"time"
)

// UserAsset represents the user_asset table
type UserAsset struct {
	ID          int       `gorm:"primarykey" json:"id"`
	UserID      string    `gorm:"column:user_id;type:varchar(30);not null;uniqueIndex;comment:用户ID" json:"user_id"` // 用户ID
	Points      int       `gorm:"column:points;type:int;not null;default:0;comment:用户积分" json:"points"`             // 用户积分
	Balance     float64   `gorm:"column:balance;type:int;not null;default:0;comment:用户余额" json:"balance"`           // 用户余额
	Consumption float64   `gorm:"column:consumption;type:int;not null;default:0;comment:用户消费" json:"consumption"`   // 用户消费
	CreatedAt   time.Time `gorm:"column:created_at;not null;comment:创建时间" json:"created_at"`                        // 创建时间
	UpdatedAt   time.Time `gorm:"column:updated_at;not null;comment:更新时间" json:"updated_at"`                        // 更新时间
}

// TableName specifies the table name for the UserAsset model
func (m *UserAsset) TableName() string {
	return "user_asset"
}

// UserAssetResponse represents the response for user asset queries
type UserAssetResponse struct {
	UserID      string  `json:"user_id"`      // 用户ID
	Points      int     `json:"points"`       // 用户积分
	Balance     float64 `json:"balance"`      // 用户余额
	CouponCount int     `json:"coupon_count"` // 用户优惠券数量
	Nickname    string  `json:"nickname"`     // 用户昵称
	Avatar      string  `json:"avatar"`       // 用户头像
}

// RechargeRequest 表示充值余额的请求
type RechargeRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"` // 充值金额，必须大于0
}

// WithdrawRequest 表示从用户余额提取的请求
type WithdrawRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"` // 提取金额，必须大于0
}
