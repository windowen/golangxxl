package table

import "time"

type User struct {
	Id                 int       `gorm:"column:id" json:"id"`                                                       // ID
	ChatUuid           string    `gorm:"column:chat_uuid" json:"chat_uuid"`                                         // 聊天室用户id
	CountryCode        string    `gorm:"column:country_code" json:"country_code"`                                   // 注册国家code
	Avatar             string    `gorm:"column:avatar" json:"avatar"`                                               // 用户头像
	Nickname           string    `gorm:"column:nickname" json:"nickname"`                                           // 昵称
	Sex                int       `gorm:"column:sex" json:"sex"`                                                     // 性别 1-男 2-女
	Sign               string    `gorm:"column:sign" json:"sign"`                                                   // 签名
	Birthday           string    `gorm:"column:birthday" json:"birthday"`                                           // 生日
	Feeling            int       `gorm:"column:feeling" json:"feeling"`                                             // 感情
	Country            string    `gorm:"column:country" json:"country"`                                             // 国家
	Area               string    `gorm:"column:area" json:"area"`                                                   // 地区
	Profession         int       `gorm:"column:profession" json:"profession"`                                       // 职业
	PayPassword        string    `gorm:"column:pay_password" json:"pay_password"`                                   // 支付密码
	Category           int       `gorm:"column:category" json:"category"`                                           // 类型 1-用户 2-主播 3-机器人
	SiteId             int       `gorm:"column:site_id" json:"site_id"`                                             // 站点ID
	InviteCode         string    `gorm:"column:invite_code" json:"invite_code"`                                     // 邀请码
	ParentId           int       `gorm:"column:parent_id" json:"parent_id"`                                         // 上级ID
	LevelId            int       `gorm:"column:level_id" json:"level_id"`                                           // 等级ID
	SetLevelId         int       `gorm:"column:set_level_id" json:"set_level_id"`                                   // 后台设置等级id 优先展示
	LoginErrorTimes    int       `gorm:"column:login_error_times" json:"login_error_times"`                         // 登陆错误次数
	Remark             string    `gorm:"column:remark" json:"remark"`                                               // 备注
	GmStatus           int       `gorm:"column:gm_status" json:"gm_status"`                                         // 超级管理员状态 1- 开启 2-关闭
	Status             int       `gorm:"column:status" json:"status"`                                               // 账号状态 1-正常 2-禁用 3-删除
	AgentRebateStatus  int       `gorm:"column:agent_rebate_status" json:"agent_rebate_status"`                     // 代理返点状态 1-正常 2-禁用
	InviteRebateStatus int       `gorm:"column:invite_rebate_status" json:"invite_rebate_status"`                   // 邀请返点状态 1-正常 2-禁用
	BetStatus          int       `gorm:"column:bet_status" json:"bet_status"`                                       // 投注状态 1-正常 2-禁用
	DrawStatus         int       `gorm:"column:draw_status" json:"draw_status"`                                     // 出款状态 1-正常 2-禁用
	LoginIp            string    `gorm:"column:login_ip" json:"login_ip"`                                           // 登录IP
	ChannelCode        int       `gorm:"column:channel_code" json:"channel_code"`                                   // 管理后台定义的唯一id（整数类型）
	ChannelPartnerId   int       `gorm:"column:channel_partner_id" json:"channel_partner_id"`                       // 渠道商Id
	ChannelStatus      int       `gorm:"column:channel_status" json:"channel_status"`                               // 渠道充值状态是否有效 1 有效 2 无效
	LastActiveAt       time.Time `gorm:"column:last_active_at;default:'1970-01-01 00:00:00'" json:"last_active_at"` // 最后活跃时间
	CreatedAt          time.Time `gorm:"column:created_at" json:"created_at"`                                       // 创建时间
	UpdatedAt          time.Time `gorm:"column:updated_at" json:"updated_at"`                                       // 更新时间
}

func (User) TableName() string {
	return "site_user"
}
