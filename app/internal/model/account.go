package model

import (
	"time"
)

type Account struct {
	Id           int       `gorm:"primaryKey;autoIncrement;column:id;type:int(11);NOT NULL;" json:"id"`
	UserId       string    `gorm:"column:user_id;type:varchar(30);NOT NULL;comment:用户ID" json:"user_id"`                                               // 用户ID
	Nickname     string    `gorm:"column:nickname;type:varchar(50);NOT NULL;comment:昵称" json:"nickname"`                                               // 昵称
	Avatar       string    `gorm:"column:avatar;type:varchar(255);NOT NULL;comment:头像" json:"avatar"`                                                  // 头像
	Phone        string    `gorm:"column:phone;type:varchar(255);NOT NULL;comment:手机号" json:"phone"`                                                   // 手机号
	CreatedAt    time.Time `gorm:"column:created_at;type:datetime;NOT NULL;comment:创建时间" json:"created_at"`                                            // 创建时间
	UpdatedAt    time.Time `gorm:"column:updated_at;type:datetime;NOT NULL;comment:更新时间" json:"updated_at"`                                            // 更新时间
	Username     string    `gorm:"column:username;type:varchar(100);NOT NULL;comment:用户名" json:"username"`                                             // 用户名
	DisplayName  string    `gorm:"column:display_name;type:varchar(100);NOT NULL;comment:展示名称" json:"display_name"`                                    // 展示名称
	Role         int64     `gorm:"column:role;type:bigint(20);default:1;NOT NULL;comment:角色（1:普通会员 2:高级会员 3:初级农场主 4:高级农场主 5:资深农场主 6:合伙人）" json:"role"` // 角色（1:普通 2:会员）
	Status       int64     `gorm:"column:status;type:bigint(20);default:1;NOT NULL;comment:状态（1: 正常 2: 注销）" json:"status"`                             // 状态（1: 正常 2: 注销）
	Email        string    `gorm:"column:email;type:varchar(191);NOT NULL;comment:邮箱" json:"email"`                                                    // 邮箱
	AuthStatus   int8      `gorm:"column:auth_status;type:tinyint(4);NOT NULL;comment:是否实名认证（0:未认证 1:已认证）" json:"auth_status"`                         // 是否实名认证（0:未认证 1:已认证）
	Platform     int8      `gorm:"column:platform;type:tinyint(4);NOT NULL;comment:注册平台(1:微信 2:抖音)" json:"platform"`                                   // 注册平台(1:微信 2:抖音)
	Gender       uint8     `gorm:"column:gender;comment:性别" json:"gender"`                                                                             // 性别
	Birthday     string    `gorm:"column:birthday;comment:出生日期" json:"birthday"`                                                                       // 出生日期
	MemberLevel  int       `gorm:"column:member_level;comment:会员等级" json:"member_level"`                                                               // 会员等级
	Address      string    `gorm:"column:address;comment:收获地址" json:"address"`                                                                         // 收获地址
	OpenID       string    `gorm:"column:open_id;type:varchar(64);comment:微信OpenID" json:"open_id"`
	UnionID      string    `gorm:"column:union_id;type:varchar(64);comment:微信UnionID" json:"union_id"`
	RegisterTime time.Time `gorm:"column:register_time;type:datetime;comment:注册时间" json:"register_time"`
	LoginTime    time.Time `gorm:"column:login_time;type:datetime;comment:登录时间" json:"login_time"`
}

type AccountWithRank struct {
	Account
	Rank int `gorm:"column:rank" json:"rank"`
}

func (m *Account) TableName() string {
	return "users"
}

type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
}
