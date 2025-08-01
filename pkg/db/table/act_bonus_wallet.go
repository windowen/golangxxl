package table

import (
	"time"

	"github.com/shopspring/decimal"
)

type ActBonusWallet struct {
	Id                          int             `gorm:"column:id" json:"id"`                                                             // ID
	UserId                      int             `gorm:"column:user_id" json:"user_id"`                                                   // 用户id
	Balance                     decimal.Decimal `gorm:"column:balance" json:"balance"`                                                   // 彩金数量
	ReachAmount                 decimal.Decimal `gorm:"column:reach_amount" json:"reach_amount"`                                         // 彩金可以提现需要达到的数量
	FlowAmount                  decimal.Decimal `gorm:"column:flow_amount" json:"flow_amount"`                                           // 彩金提现需要的流水
	ReachUsers                  int             `gorm:"column:reach_users" json:"reach_users"`                                           // 彩金可以提现需要达到的邀请人数
	GameBonusTypeConfigBatchNum string          `gorm:"column:game_bonus_type_config_batch_num" json:"game_bonus_type_config_batch_num"` // 游戏类型batch_num (对应game_bonus_type_config里的batch_num字段)
	Ip                          string          `gorm:"column:ip" json:"ip"`                                                             // IP
	DeviceId                    string          `gorm:"column:device_id" json:"device_id"`                                               // 设备ID
	CreatedAt                   time.Time       `gorm:"column:created_at" json:"created_at"`                                             // 创建时间
	UpdatedAt                   time.Time       `gorm:"column:updated_at" json:"updated_at"`                                             // 更新时间
}

func (ActBonusWallet) TableName() string {
	return "act_bonus_wallet"
}
