package table

import (
	"time"

	"github.com/shopspring/decimal"
)

type ActRegisterWelfareConfig struct {
	Id                          int             `gorm:"column:id" json:"id"`                                                             // 注册福利配置表主键Id
	BonusFundsAmount            decimal.Decimal `gorm:"column:bonus_funds_amount" json:"bonus_funds_amount"`                             // 体验金
	WithdrawThresholdAmount     decimal.Decimal `gorm:"column:withdraw_threshold_amount" json:"withdraw_threshold_amount"`               // 达到额度出款
	InviteFriendsDeposit        int             `gorm:"column:invite_friends_deposit" json:"invite_friends_deposit"`                     // 邀请好友注册且充值
	Status                      int             `gorm:"column:status" json:"status"`                                                     // 状态1-启用2-停用
	GameBonusTypeConfigBatchNum string          `gorm:"column:game_bonus_type_config_batch_num" json:"game_bonus_type_config_batch_num"` // 游戏类型batch_num (对应game_bonus_type_config里的batch_num字段)
	RequiredTurnoverMultiple    int             `gorm:"column:required_turnover_multiple" json:"required_turnover_multiple"`             // 邀请人数充值后需要达成的流水倍数
	RechargeAmount              int             `gorm:"column:recharge_amount" json:"recharge_amount"`                                   // 充值金额
	CreatedAt                   time.Time       `gorm:"column:created_at" json:"created_at"`                                             // 创建时间
	UpdatedAt                   time.Time       `gorm:"column:updated_at" json:"updated_at"`                                             // 更新时间
}

func (ActRegisterWelfareConfig) TableName() string {
	return "act_register_welfare_config"
}
