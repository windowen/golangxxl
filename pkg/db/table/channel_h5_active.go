package table

import "time"

type ChannelH5Active struct {
	Id           int       `gorm:"column:id" json:"id"`                       // 自增id
	DeviceId     string    `gorm:"column:device_id" json:"device_id"`         // 设备id
	ChannelCode  int       `gorm:"column:channel_code" json:"channel_code"`   // 渠道号（整数类型，如1001）
	ActiveTime   time.Time `gorm:"column:active_time" json:"active_time"`     // 活跃时间
	RegisterTime time.Time `gorm:"column:register_time" json:"register_time"` // 注册时间
}

func (ChannelH5Active) TableName() string {
	return "channel_h5_active"
}
