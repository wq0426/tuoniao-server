package model

import (
	"time"
)

type UserAssetRecord struct {
	Id            int       `gorm:"primaryKey;autoIncrement;column:id;type:int(11);NOT NULL;" json:"id"`
	UserId        string    `gorm:"column:user_id;type:varchar(255);NOT NULL;comment:用户ID" json:"user_id"` // 用户ID
	BusinessType  int8      `gorm:"column:business_type;type:tinyint(4);default:0;NOT NULL;comment:业务类型" json:"business_type"`
	ActionType    int8      `gorm:"column:action_type;type:tinyint(4);default:0;NOT NULL;comment:消费类型(1:使用 2:奖励 3:购买 5:售出所得 6:收益 7:使用)" json:"action_type"` // 消费类型(1:使用 2:奖励 3:购买 5:售出所得 6:收益 7:使用)
	AssetType     int8      `gorm:"column:asset_type;type:tinyint(4);default:0;NOT NULL;comment:资产类型(1:金币 2:命数 3:扑克 4:转盘次数 5:广告加倍)" json:"asset_type"`      // 资产类型(1:金币 2:命数 3:扑克 4:转盘次数 5:广告加倍)
	ActionNum     float32   `gorm:"column:action_num;type:double;default:0;NOT NULL;comment:使用数量" json:"action_num"`                                        // 使用数量
	LeftNum       float32   `gorm:"column:left_num;type:double;default:0;NOT NULL;comment:剩余数量" json:"left_num"`                                            // 剩余数量
	CreatedAt     time.Time `gorm:"column:created_at;type:datetime;NOT NULL;comment:创建时间" json:"created_at"`                                                // 创建时间
	UpdatedAt     time.Time `gorm:"column:updated_at;type:datetime;NOT NULL;comment:更新时间" json:"updated_at"`                                                // 更新时间
	RelationId    int       `gorm:"column:relation_id;comment:关联ID" json:"relation_id"`                                                                     // 关联ID                                                      // 关联图片
	RelationTitle string    `gorm:"column:relation_title;comment:关联标题" json:"relation_title"`                                                               // 关联标题
}

func (m *UserAssetRecord) TableName() string {
	return "user_asset_record"
}

// BalanceRecordQueryRequest 表示查询余额记录的请求
type BalanceRecordQueryRequest struct {
	Page     int `form:"page" json:"page"`           // 页码，从1开始
	PageSize int `form:"page_size" json:"page_size"` // 每页条数
}

// BalanceRecordResponse 余额记录响应
type BalanceRecordResponse struct {
	Total    int64                `json:"total"`     // 总记录数
	Records  []UserAssetRecordDTO `json:"records"`   // 记录列表
	Page     int                  `json:"page"`      // 当前页码
	PageSize int                  `json:"page_size"` // 每页条数
}

// UserAssetRecordDTO 用户资产记录DTO
type UserAssetRecordDTO struct {
	ID           uint64  `json:"id"`            // 记录ID
	UserId       string  `json:"user_id"`       // 用户ID
	Title        string  `json:"title"`         // 标题
	BusinessType int8    `json:"business_type"` // 业务类型
	ActionType   int8    `json:"action_type"`   // 操作类型：1-增加，2-减少
	AssetType    int8    `json:"asset_type"`    // 资产类型：1-积分，2-余额
	ActionNum    float32 `json:"action_num"`    // 操作数量
	LeftNum      float32 `json:"left_num"`      // 剩余数量
	Remark       string  `json:"remark"`        // 备注
	CreatedAt    string  `json:"created_at"`    // 创建时间
}
