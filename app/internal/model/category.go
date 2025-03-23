package model

import (
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	CategoryName string `gorm:"column:category_name;type:varchar(255);not null;comment:分类名称" json:"category_name"` // 分类名称
	ParentID     *int8  `gorm:"column:parent_id;type:tinyint" json:"parent_id"`                                    // 父级ID
	Sort         *int   `gorm:"column:sort;type:int" json:"sort"`                                                  // 排序
}

func (m *Category) TableName() string {
	return "category"
}
