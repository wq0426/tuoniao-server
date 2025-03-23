package model

import (
	"time"
)

// WithdrawOrder 提现单模型
type WithdrawOrder struct {
	ID           uint64     `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	WithdrawNo   string     `gorm:"column:withdraw_no;type:varchar(255);not null;comment:提现单号" json:"withdraw_no"`                     // 提现单号
	UserID       string     `gorm:"column:user_id;type:varchar(255);not null;comment:用户ID" json:"user_id"`                             // 用户ID
	Amount       float64    `gorm:"column:amount;not null;comment:提现金额" json:"amount"`                                                 // 提现金额
	Fee          float64    `gorm:"column:fee;not null;default:0;comment:手续费" json:"fee"`                                              // 手续费
	ActualAmount float64    `gorm:"column:actual_amount;not null;comment:实际到账金额" json:"actual_amount"`                                 // 实际到账金额
	Status       uint8      `gorm:"column:status;type:tinyint;not null;default:0;comment:提现状态(0:待审核;1:处理中;2:已完成;3:已拒绝)" json:"status"` // 提现状态
	RejectReason string     `gorm:"column:reject_reason;type:varchar(500);comment:拒绝原因" json:"reject_reason"`                          // 拒绝原因
	BankName     string     `gorm:"column:bank_name;type:varchar(100);comment:银行名称" json:"bank_name"`                                  // 银行名称
	AccountName  string     `gorm:"column:account_name;type:varchar(100);comment:账户名" json:"account_name"`                             // 账户名
	AccountNo    string     `gorm:"column:account_no;type:varchar(100);comment:账号" json:"account_no"`                                  // 账号
	Remark       string     `gorm:"column:remark;type:varchar(500);comment:备注" json:"remark"`                                          // 备注
	AuditTime    *time.Time `gorm:"column:audit_time;comment:审核时间" json:"audit_time"`                                                  // 审核时间
	CompleteTime *time.Time `gorm:"column:complete_time;comment:完成时间" json:"complete_time"`                                            // 完成时间
	CreatedAt    *time.Time `gorm:"column:created_at;comment:创建时间" json:"created_at"`                                                  // 创建时间
	UpdatedAt    *time.Time `gorm:"column:updated_at;comment:更新时间" json:"updated_at"`                                                  // 更新时间
}

func (m *WithdrawOrder) TableName() string {
	return "withdraw_order"
}

// CreateWithdrawRequest 创建提现请求
type CreateWithdrawRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"` // 提现金额
	// BankName    string  `json:"bank_name" binding:"required,max=100"`    // 银行名称
	// AccountName string  `json:"account_name" binding:"required,max=100"` // 账户名
	// AccountNo   string  `json:"account_no" binding:"required,max=100"`   // 账号
	// Remark      string  `json:"remark" binding:"max=500"`                // 备注
}

// WithdrawQueryRequest 提现查询请求
type WithdrawQueryRequest struct {
	Status   *uint8 `form:"status"`                     // 提现状态
	Page     int    `form:"page" json:"page"`           // 页码
	PageSize int    `form:"page_size" json:"page_size"` // 每页条数
}

// WithdrawListItem 提现列表项
type WithdrawListItem struct {
	ID           uint64     `json:"id"`            // 提现单ID
	Title        string     `json:"title"`         // 标题
	WithdrawNo   string     `json:"withdraw_no"`   // 提现单号
	Amount       float64    `json:"amount"`        // 提现金额
	Fee          float64    `json:"fee"`           // 手续费
	ActualAmount float64    `json:"actual_amount"` // 实际到账金额
	Status       uint8      `json:"status"`        // 提现状态
	StatusText   string     `json:"status_text"`   // 提现状态文本
	BankName     string     `json:"bank_name"`     // 银行名称
	AccountName  string     `json:"account_name"`  // 账户名
	AccountNo    string     `json:"account_no"`    // 账号
	CreatedAt    *time.Time `json:"created_at"`    // 创建时间
}

// WithdrawListResponse 提现列表响应
type WithdrawListResponse struct {
	Total int64              `json:"total"` // 总数
	List  []WithdrawListItem `json:"list"`  // 提现列表
	Page  int                `json:"page"`  // 页码
	Size  int                `json:"size"`  // 每页条数
}

// WithdrawDetailResponse 提现详情响应
type WithdrawDetailResponse struct {
	ID           uint64     `json:"id"`            // 提现单ID
	WithdrawNo   string     `json:"withdraw_no"`   // 提现单号
	UserID       string     `json:"user_id"`       // 用户ID
	Amount       float64    `json:"amount"`        // 提现金额
	Fee          float64    `json:"fee"`           // 手续费
	ActualAmount float64    `json:"actual_amount"` // 实际到账金额
	Status       uint8      `json:"status"`        // 提现状态
	StatusText   string     `json:"status_text"`   // 提现状态文本
	RejectReason string     `json:"reject_reason"` // 拒绝原因
	BankName     string     `json:"bank_name"`     // 银行名称
	AccountName  string     `json:"account_name"`  // 账户名
	AccountNo    string     `json:"account_no"`    // 账号
	Remark       string     `json:"remark"`        // 备注
	AuditTime    *time.Time `json:"audit_time"`    // 审核时间
	CompleteTime *time.Time `json:"complete_time"` // 完成时间
	CreatedAt    *time.Time `json:"created_at"`    // 创建时间
	UpdatedAt    *time.Time `json:"updated_at"`    // 更新时间
}
