package table

import (
	"time"

	"github.com/shopspring/decimal"
)

type SiteUserWallet struct {
	Id                  int             `gorm:"column:id" json:"id"` // ID
	UserId              int             `gorm:"column:user_id" json:"user_id"`
	Balance             decimal.Decimal `gorm:"column:balance" json:"balance"`                           // 总余额(美金)
	Freeze              decimal.Decimal `gorm:"column:freeze" json:"freeze"`                             // 冻结金额
	Diamond             int             `gorm:"column:diamond" json:"diamond"`                           // 钻石(只有整数)
	SettlementDiamond   int             `gorm:"column:settlement_diamond" json:"settlement_diamond"`     // 结算的钻石(只有整数)
	WithdrawDiamond     int             `gorm:"column:withdraw_diamond" json:"withdraw_diamond"`         // 可提现钻石(只有整数)
	WithdrawDate        time.Time       `gorm:"column:withdraw_date" json:"withdraw_date"`               // 主播下次可提现钻石的日期（如：2025-02-11）
	FreezeDiamond       int             `gorm:"column:freeze_diamond" json:"freeze_diamond"`             // 冻结的提现钻石(只有整数)
	FlowAmount          decimal.Decimal `gorm:"column:flow_amount" json:"flow_amount"`                   // 充值本金提现还剩余的流水
	AccumulationDiamond int             `gorm:"column:accumulation_diamond" json:"accumulation_diamond"` // 累积消费钻石(只有整数)
	LevelUpExp          int             `gorm:"column:level_up_exp" json:"level_up_exp"`                 // 当前的经验
	TotalRecharge       decimal.Decimal `gorm:"column:total_recharge" json:"total_recharge"`             // 累积充值金额
	TotalBet            decimal.Decimal `gorm:"column:total_bet" json:"total_bet"`                       // 累积下注金额
	FirstAmount         decimal.Decimal `gorm:"column:first_amount" json:"first_amount"`                 // 首充金额
	Trc20Address        string          `gorm:"column:trc20_address" json:"trc20_address"`               // trc20支付地址
	Trc20PrivateKey     string          `gorm:"column:trc20_private_key" json:"trc20_private_key"`       // trc20支付私钥
	CreatedAt           time.Time       `gorm:"column:created_at" json:"created_at"`                     // 创建时间
	UpdatedAt           time.Time       `gorm:"column:updated_at" json:"updated_at"`                     // 更新时间
}

func (SiteUserWallet) TableName() string {
	return "site_user_wallet"
}
