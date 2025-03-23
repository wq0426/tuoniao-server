package model

import (
	"time"

	"gorm.io/gorm"
)

type FreeMarketMine struct {
	gorm.Model
	UserID   string     `gorm:"column:user_id;type:varchar(30);not null;comment:用户ID" json:"user_id"`           // 用户ID
	EggPrice float64    `gorm:"column:egg_price;type:decimal(10,2);comment:蛋价格" json:"egg_price"`               // 蛋价格
	EggNum   *uint      `gorm:"column:egg_num;type:int(13) unsigned;comment:蛋数量" json:"egg_num"`                // 蛋数量
	Status   *uint8     `gorm:"column:status;type:tinyint(4) unsigned;comment:出售状态（0:未出售;1:已出售）" json:"status"` // 出售状态
	Date     *time.Time `gorm:"column:date;type:date" json:"date"`                                              // 日期
}

func (m *FreeMarketMine) TableName() string {
	return "free_market_mine"
}

// Define a struct to receive parameters
type UpdateEggPriceRequest struct {
	ID    int     `form:"id" binding:"required"`
	Price float64 `form:"price" binding:"required"`
}
