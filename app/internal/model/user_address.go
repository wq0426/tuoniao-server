package model

import (
	"time"
)

type UserAddress struct {
	ID        uint64    `gorm:"primarykey" json:"id"`
	UserID    string    `gorm:"column:user_id;type:varchar(255);not null;comment:用户ID" json:"user_id"`              // 用户ID
	Name      string    `gorm:"column:name;type:varchar(255);not null;comment:收件人姓名" json:"name"`                   // 收件人姓名
	Province  string    `gorm:"column:province;type:varchar(255);not null;comment:省份" json:"province"`              // 省份
	City      string    `gorm:"column:city;type:varchar(255);not null;comment:城市" json:"city"`                      // 城市
	District  string    `gorm:"column:district;type:varchar(255);not null;comment:区/县" json:"district"`             // 区/县
	Street    string    `gorm:"column:street;type:varchar(255);not null;comment:详细地址" json:"street"`                // 详细地址
	IsDefault uint8     `gorm:"column:is_default;type:tinyint;not null;default:0;comment:是否默认地址" json:"is_default"` // 是否默认地址
	Phone     string    `gorm:"column:phone;type:varchar(20);not null;comment:联系电话" json:"phone"`                   // 联系电话
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (m *UserAddress) TableName() string {
	return "user_address"
}

// AddAddressRequest represents the request to add a new address
type AddAddressRequest struct {
	Name      string `json:"name" binding:"required"`
	Province  string `json:"province" binding:"required"`
	City      string `json:"city" binding:"required"`
	District  string `json:"district" binding:"required"`
	Street    string `json:"street" binding:"required"`
	IsDefault uint8  `json:"is_default"`
	Phone     string `json:"phone" binding:"required"`
}

// UpdateAddressRequest represents the request to update an existing address
type UpdateAddressRequest struct {
	ID        uint64 `json:"id" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Province  string `json:"province" binding:"required"`
	City      string `json:"city" binding:"required"`
	District  string `json:"district" binding:"required"`
	Street    string `json:"street" binding:"required"`
	IsDefault uint8  `json:"is_default"`
	Phone     string `json:"phone" binding:"required"`
}

// Parse request body
type DeleteAddressRequest struct {
	ID uint64 `json:"id" binding:"required"`
}
