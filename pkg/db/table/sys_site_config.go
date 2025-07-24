package table

import "time"

type SysSiteConfig struct {
	Id         int       `gorm:"column:id" json:"id"`
	Category   int       `gorm:"column:category" json:"category"`       // 配置类型 1- App配置 2-直播配置 3-财务配置
	Name       string    `gorm:"column:name" json:"name"`               // 配置名称
	ConfigCode string    `gorm:"column:config_code" json:"config_code"` // 配置code 唯一
	Content    string    `gorm:"column:content" json:"content"`         // 配置内容
	Status     int       `gorm:"column:status" json:"status"`           // 状态0-停用1-启用
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`   // 创建日期
}

func (SysSiteConfig) TableName() string {
	return "sys_site_config"
}
