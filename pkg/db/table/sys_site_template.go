package table

import "time"

type SysSiteTemplate struct {
	Id              int       `gorm:"column:id" json:"id"`                             // 站内信模板表Id
	CountryCode     string    `gorm:"column:country_code" json:"country_code"`         // 国家编码
	LanguageCode    string    `gorm:"column:language_code" json:"language_code"`       // 语言 短信、邮件模版使用
	Category        int       `gorm:"column:category" json:"category"`                 // 通知类型1-充值成功通知2-充值失败通知3-提现成功通知4-提现失败通知5-兑换钻石6-管理后台通知7-公告维护 8-短信 9-邮件
	TemplateCode    string    `gorm:"column:template_code" json:"template_code"`       // 模版code  如短信登陆、注册模版
	Subject         string    `gorm:"column:subject" json:"subject"`                   // 邮件主题
	TemplateContent string    `gorm:"column:template_content" json:"template_content"` // 模板内容
	Status          int       `gorm:"column:status" json:"status"`                     // 状态0-停用1-启用
	CreatedAt       time.Time `gorm:"column:created_at" json:"created_at"`             // 创建日期
}

func (SysSiteTemplate) TableName() string {
	return "sys_site_template"
}

type SiteTransactionMsg struct {
	Id          int       `gorm:"column:id" json:"id"`                     // 交易消息表Id
	Category    int       `gorm:"column:category" json:"category"`         // 通知类型1-充值成功通知2-充值失败通知3-提现成功通知4-提现失败通知5-兑换钻石
	UserId      int       `gorm:"column:user_id" json:"user_id"`           // 接收通知的用户Id
	CountryCode string    `gorm:"column:country_code" json:"country_code"` // 国家编码
	Content     string    `gorm:"column:content" json:"content"`           // 通知内容
	ReadStatus  int       `gorm:"column:read_status" json:"read_status"`   // 0-未读1-已读
	Status      int       `gorm:"column:status" json:"status"`             // 0-未删除1-已删除
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`     // 创建时间
}

func (SiteTransactionMsg) TableName() string {
	return "site_transaction_msg"
}
