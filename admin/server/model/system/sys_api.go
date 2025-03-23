package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type SysApi struct {
	global.GVA_MODEL
	Path        string `json:"path" gorm:"comment:api路径"`             // api路径
	Description string `json:"description" gorm:"comment:api中文描述"`    // api中文描述
	ApiGroup    string `json:"apiGroup" gorm:"comment:api组"`          // api组
	Method      string `json:"method" gorm:"default:POST;comment:方法"` // 方法:创建POST(默认)|查看GET|更新PUT|删除DELETE
}

func (SysApi) TableName() string {
	return "sys_apis"
}

type SysIgnoreApi struct {
	global.GVA_MODEL
	Path   string `json:"path" gorm:"comment:api路径"`             // api路径
	Method string `json:"method" gorm:"default:POST;comment:方法"` // 方法:创建POST(默认)|查看GET|更新PUT|删除DELETE
	Flag   bool   `json:"flag" gorm:"-"`                         // 是否忽略
}

func (SysIgnoreApi) TableName() string {
	return "sys_ignore_apis"
}

type SysAgent struct {
	global.GVA_MODEL
	Site        string `json:"site" binding:"required" gorm:"comment:路径"` // 路径
	ShareId     string `json:"shareId" gorm:"comment:共享ID"`               // 共享ID
	Region      string `json:"region" gorm:"comment:地区"`                  // 地区
	Agent       string `json:"agent" gorm:"comment:代理人"`                  // 代理人
	Description string `json:"description" gorm:"comment:描述"`             // 描述
}

func (SysAgent) TableName() string {
	return "sys_agent"
}
