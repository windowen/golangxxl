package table

import (
	"time"

	"github.com/shopspring/decimal"
)

type FinanceCoinExchange struct {
	Id           int             `gorm:"column:id" json:"id"`                         // 币种兑换表
	FromCoinCode string          `gorm:"column:from_coin_code" json:"from_coin_code"` // 要兑换的币种（固定死为美金）
	ToCoinCode   string          `gorm:"column:to_coin_code" json:"to_coin_code"`     // 被兑换的币种（如泰铢）
	ExchangeRate decimal.Decimal `gorm:"column:exchange_rate" json:"exchange_rate"`   // 兑换汇率（如：7.1400）
	Status       int             `gorm:"column:status" json:"status"`                 // 0启用1禁用
	CreatedAt    time.Time       `gorm:"column:created_at" json:"created_at"`         // 创建时间
	UpdatedAt    time.Time       `gorm:"column:updated_at" json:"updated_at"`         // 更新时间
}

func (FinanceCoinExchange) TableName() string {
	return "finance_coin_exchange"
}
