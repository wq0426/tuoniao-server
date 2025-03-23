package model

import "gorm.io/gorm"

type Banner struct {
	gorm.Model
	Img  string `gorm:"column:img;type:varchar(255);not null;comment:图片链接" json:"img"`   // 图片链接
	Url  string `gorm:"column:url;type:varchar(255);not null;comment:跳转链接" json:"url"`   // 跳转链接
	Path string `gorm:"column:path;type:varchar(255);not null;comment:跳转路径" json:"path"` // 跳转路径
}

func (m *Banner) TableName() string {
	return "banner"
}

type BannerResponse struct {
	ID  uint   `json:"id"`
	Img string `json:"img"` // 图片链接
	Url string `json:"url"` // 跳转链接
}
