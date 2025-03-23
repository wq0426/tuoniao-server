package model

import (
	"time"
)

// ProductEvaluate 商品评价模型
type ProductEvaluate struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	ParentID  int8   `gorm:"column:parent_id;type:tinyint;comment:父级ID" json:"parent_id"`
	ProductID uint64 `gorm:"column:product_id;type:bigint unsigned;not null;comment:商品ID" json:"product_id"` // 商品ID
	ReviewID  uint64 `gorm:"column:review_id;type:bigint unsigned;not null;comment:评价ID" json:"review_id"`   // 评价ID
	UserID    string `gorm:"column:user_id;type:varchar(30);not null;comment:用户ID" json:"user_id"`           // 用户ID
	Nickname  string `gorm:"column:nickname;type:varchar(50);comment:用户昵称" json:"nickname"`                  // 用户昵称
	Avatar    string `gorm:"column:avatar;type:varchar(255);comment:用户头像" json:"avatar"`                     // 用户头像
	Content   string `gorm:"column:content;type:text;comment:评价内容" json:"content"`                           // 评价内容
}

func (m *ProductEvaluate) TableName() string {
	return "product_evaluate"
}

type ProductEvaluateList struct {
	ProductEvaluate
	CreatedAtStr string                `gorm:"-" json:"created_at_str"`
	Children     []ProductEvaluateList `json:"comment_list"`
}

// CreateEvaluateReplyRequest 创建评价回复请求
type CreateEvaluateReplyRequest struct {
	EvaluateID uint64 `json:"evaluate_id" binding:"required"` // 被回复的评价ID
	Content    string `json:"content" binding:"required"`     // 回复内容
}

// CreateEvaluateRequest 创建主评论请求
type CreateEvaluateRequest struct {
	ReviewID uint64 `json:"review_id" binding:"required"` // 评价ID
	Content  string `json:"content" binding:"required"`   // 评论内容
}

// UpdateEvaluateAnonymousRequest 更新评论匿名状态请求
type UpdateEvaluateAnonymousRequest struct {
	ReviewId    uint64 `json:"review_id" binding:"required"` // 评论ID
	IsAnonymous bool   `json:"is_anonymous"`                 // 是否匿名：true-匿名，false-公开
}
