package model

import "gorm.io/gorm"

type Monitor struct {
	gorm.Model
	Img   string `gorm:"column:img;type:varchar(255);not null;comment:图片链接" json:"img"`   // 图片链接
	Url   string `gorm:"column:url;type:varchar(255);not null;comment:跳转链接" json:"url"`   // 跳转链接
	Title string `gorm:"column:title;type:varchar(255);not null;comment:标题" json:"title"` // 标题
}

func (m *Monitor) TableName() string {
	return "monitor"
}

type MonitorResponse struct {
	ID    uint   `json:"id"`
	Img   string `json:"img"`   // 图片链接
	Url   string `json:"url"`   // 跳转链接
	Title string `json:"title"` // 标题
}
