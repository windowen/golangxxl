package table

import (
	"time"

	"github.com/shopspring/decimal"
)

type ActUserWallet struct {
	Id            int             `gorm:"column:id" json:"id"` // ID
	UserId        int             `gorm:"column:user_id" json:"user_id"`
	Balance       decimal.Decimal `gorm:"column:balance" json:"balance"`                                                   // 彩金数量
	Diamond       int             `gorm:"column:diamond" json:"diamond"`                                                   // 钻石(只有整数)
	FlowAmount    decimal.Decimal `gorm:"column:flow_amount" json:"flow_amount"`                                           // 活动彩金对应的累积流水
	GameBonusType string          `gorm:"column:game_bonus_type_config_batch_num" json:"game_bonus_type_config_batch_num"` // 游戏类型batch_num (对应game_bonus_type_config里的batch_num字段)
	CreatedAt     time.Time       `gorm:"column:created_at" json:"created_at"`                                             // 创建时间
	UpdatedAt     time.Time       `gorm:"column:updated_at" json:"updated_at"`                                             // 更新时间
}

func (ActUserWallet) TableName() string {
	return "act_user_wallet"
}
