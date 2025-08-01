package table

import "time"

type ChannelManage struct {
	Id                 int       `gorm:"column:id" json:"id"`                                     // 渠道管理表Id
	ChannelPartnerId   int       `gorm:"column:channel_partner_id" json:"channel_partner_id"`     // 渠道商Id，对应channel_partner表的id字段
	ChannelPartnerName string    `gorm:"column:channel_partner_name" json:"channel_partner_name"` // 渠道商名称，对应channel_partner表的渠道商名称
	ChannelCode        int       `gorm:"column:channel_code" json:"channel_code"`                 // 管理后台定义的唯一id（整数类型4位数字，如1001），每行唯一
	ChannelUrl         string    `gorm:"column:channel_url" json:"channel_url"`                   // 渠道链接
	DeductionType      int       `gorm:"column:deduction_type" json:"deduction_type"`             // 1-CPS充值笔数 2-CPA 安装量 3-注册量 4-访问量
	Threshold          int       `gorm:"column:threshold" json:"threshold"`                       // 阈值（比如：超过10笔）
	Per                int       `gorm:"column:per" json:"per"`                                   // 每隔（如：超过10笔后，每5笔）
	Deduction          int       `gorm:"column:deduction" json:"deduction"`                       // 扣量固定死1（比如：超过10笔，每5笔，扣量1笔）
	Status             int       `gorm:"column:status" json:"status"`                             // 0-禁用 1-启用
	Operator           string    `gorm:"column:operator" json:"operator"`                         // 操作人
	PixelId            string    `gorm:"column:pixel_id" json:"pixel_id"`                         // 渠道商对应的像素id
	PixelToken         string    `gorm:"column:pixel_token" json:"pixel_token"`                   // 渠道商对应的像素token
	CreatedAt          time.Time `gorm:"column:created_at" json:"created_at"`                     // 创建时间
	UpdatedAt          time.Time `gorm:"column:updated_at" json:"updated_at"`                     // 更新时间
}

func (ChannelManage) TableName() string {
	return "channel_manage"
}
