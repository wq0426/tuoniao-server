package model

import "gorm.io/gorm"

type Resource struct {
	gorm.Model
}

func (m *Resource) TableName() string {
    return "resource"
}
