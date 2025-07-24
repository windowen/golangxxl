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
	RechargeAmount      decimal.Decimal `gorm:"column:recharge_amount" json:"recharge_amount"`           // 充值金额
	RebateAmount        decimal.Decimal `gorm:"column:rebate_amount" json:"rebate_amount"`               // 返点金额
	Diamond             int             `gorm:"column:diamond" json:"diamond"`                           // 钻石(只有整数)
	SettlementDiamond   int             `gorm:"column:settlement_diamond" json:"settlement_diamond"`     // 已经结算的钻石(只有整数)
	AccumulationDiamond int             `gorm:"column:accumulation_diamond" json:"accumulation_diamond"` // 累积消费的钻石(只有整数)
	CreatedAt           time.Time       `gorm:"column:created_at" json:"created_at"`                     // 创建时间
	UpdatedAt           time.Time       `gorm:"column:updated_at" json:"updated_at"`                     // 更新时间
}

func (SiteUserWallet) TableName() string {
	return "site_user_wallet"
}
