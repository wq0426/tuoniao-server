package model

import (
	"app/internal/common"
	"time"
)

type Settings struct {
	Id        int64  `gorm:"primaryKey;autoIncrement;column:id;type:bigint(20) unsigned;NOT NULL;" json:"id"`
	CreatedAt string `gorm:"column:created_at;type:datetime(3);NULL;" json:"created_at"`
	UpdatedAt string `gorm:"column:updated_at;type:datetime(3);NULL;" json:"updated_at"`
	DeletedAt string `gorm:"column:deleted_at;type:datetime(3);NULL;" json:"deleted_at"`
	Name      string `gorm:"column:name;type:varchar(191);NULL;comment:参数名称" json:"name"` // 参数名称
	Key       string `gorm:"column:key;type:varchar(191);NULL;comment:参数键" json:"key"`    // 参数键
	Value     string `gorm:"column:value;type:text;NULL;comment:参数值" json:"value"`        // 参数值
	Desc      string `gorm:"column:desc;type:varchar(191);NULL;comment:参数说明" json:"desc"` // 参数说明
}

func (m *Settings) TableName() string {
	return "sys_params"
}

type ProfileSettings struct {
	UserId      string `json:"user_id"`      // 用户ID
	PhoneNumber string `json:"phone_number"` // 用户的手机号
	AvatarURL   string `json:"avatar_url"`   // 头像
	Nickname    string `json:"nickname"`     // 昵称
}

// 昵称和头像
type NicknameAvatar struct {
	Nickname     string `json:"nickname"`      // 昵称
	AvatarBase64 string `json:"avatar_base64"` // 头像
	Gender       uint8  `json:"gender"`        // 性别
	Birthday     string `json:"birthday"`      // 出生日期
	MemberLevel  int    `json:"member_level"`  // 会员等级
	Address      string `json:"address"`       // 收获地址
	PhoneNumber  string `json:"phone"`         // 绑定手机号
}

func (r *NicknameAvatar) Validate() bool {
	// if r.Nickname == "" && r.AvatarURL == "" {
	// 	return false
	// }
	if len(r.PhoneNumber) > 0 {
		// 验证手机号
		if !common.IsMobile(r.PhoneNumber) {
			return false
		}
	}
	return true
}

type SetPasswordRequest struct {
	Password string `json:"password" binding:"required"`
}

func (r *SetPasswordRequest) Validate() bool {
	if len(r.Password) == 0 {
		return false
	}
	return true
}

type RecordsReponse struct {
	ConsumptionRecords []RecordItem `json:"consumption_records"`
	IncomeRecords      []RecordItem `json:"income_records"`
}

type RecordItem struct {
	ID           int64     `gorm:"primaryKey;autoIncrement;column:id;type:bigint(20) unsigned;NOT NULL;" json:"id"`
	BusinessType string    `gorm:"column:business_type;type:varchar(191);NULL;comment:业务类型名称" json:"business_type"`       // 业务类型名称
	ActionType   string    `gorm:"column:consumption_type;type:varchar(191);NULL;comment:消费类型名称" json:"consumption_type"` // 消费类型名称
	Quantity     float64   `gorm:"column:quantity;type:int;NULL;comment:消费数量" json:"quantity"`                            // 消费数量
	Balance      float64   `gorm:"column:balance;type:decimal(10,2);NULL;comment:余额" json:"balance"`                      // 余额
	CreatedAt    time.Time `gorm:"column:created_at;type:datetime(3);NULL;" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:datetime(3);NULL;" json:"updated_at"`
}

type VerifyResponseBody struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Ta      interface{} `json:"ta"`
}

type VerifyRequestParams struct {
	Identification string `json:"identification"`
	Name           string `json:"name"`
}

type VerifyRequest struct {
	RequestURL    string              `json:"request_url"`
	RequestParams VerifyRequestParams `json:"request_params"`
	UserId        string              `json:"UserId"`
	ResponseBody  VerifyResponseBody  `json:"response_body"`
	ErrorCode     int                 `json:"error_code"`
	Reason        string              `json:"reason"`
	Result        CheckResult         `json:"result"`
}

type CheckResult struct {
	Realname    string    `json:"realname"`
	Idcard      string    `json:"idcard"`
	Isok        bool      `json:"isok"`
	IdCardInfor *UserInfo `json:"IdCardInfor"`
}

type SwitchStatusRequest struct {
	Status int8 `json:"status"` // 活动推送状态
}

// validate
func (r *SwitchStatusRequest) Validate() bool {
	if r.Status != 0 && r.Status != 1 {
		return false
	}
	return true
}

type AvatarRequest struct {
	File string `json:"file"`
}

// validate
func (r *AvatarRequest) Validate() bool {
	if r.File == "" {
		return false
	}
	return true
}

type TurntableReward struct {
	Id        int `json:"id"`
	Percent   int `json:"percent"`
	AssetType int `json:"asset_type"`
	RewardNum int `json:"reward_num"`
}
