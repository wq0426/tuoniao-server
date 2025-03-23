package model

import (
	"time"

	"gorm.io/gorm"
)

type News struct {
	gorm.Model
	Title   string    `gorm:"column:title;type:varchar(255);not null;comment:标题" json:"title"` // 标题
	Content string    `gorm:"column:content;type:text;not null;comment:内容" json:"content"`     // 内容
	Img     string    `gorm:"column:img;type:varchar(255);not null;comment:图片链接" json:"img"`   // 图片链接
	Url     string    `gorm:"column:url;type:varchar(255);not null;comment:跳转链接" json:"url"`   // 跳转链接
	Date    time.Time `gorm:"column:date;type:date;not null" json:"date"`                      // 日期
	Type    int8      `gorm:"column:type;comment:类型(1:鸵鸟信息;2:农场信息;3:市场咨询;4:供需发布)" json:"type"` // 类型
}

func (m *News) TableName() string {
	return "news"
}

type NewsResponse struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`   // 标题
	Content string `json:"content"` // 内容
	Img     string `json:"img"`     // 图片链接
	Url     string `json:"url"`     // 跳转链接
	Date    string `json:"date"`    // 日期
}

// NewsType represents a group of news with the same type
type NewsType struct {
	Type  int8            `json:"type"`  // 类型
	Title string          `json:"title"` // 类型标题
	List  []*NewsResponse `json:"list"`  // 新闻列表
}
