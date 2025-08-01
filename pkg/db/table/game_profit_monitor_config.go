package table

import "github.com/shopspring/decimal"

type GameProfitMonitorConfig struct {
	Id                            int             `gorm:"column:id" json:"id"`                                                             // 游戏获利监控参数配置
	HighMultiplierJackpot         int             `gorm:"column:high_multiplier_jackpot" json:"high_multiplier_jackpot"`                   // 高倍爆奖开关0-关闭 1-开启
	HighMultiplierJackpotMultiple int             `gorm:"column:high_multiplier_jackpot_multiple" json:"high_multiplier_jackpot_multiple"` // 高倍爆奖倍数
	HighMultiplierJackpotMoney    decimal.Decimal `gorm:"column:high_multiplier_jackpot_money" json:"high_multiplier_jackpot_money"`       // 高倍爆奖中奖金额
	LargePrize                    decimal.Decimal `gorm:"column:large_prize" json:"large_prize"`                                           // 大额中奖
	ProfitMargin                  int             `gorm:"column:profit_margin" json:"profit_margin"`                                       // 当日会员获利比开关1-开启 0-关闭
	DailyProfitMargin             int             `gorm:"column:daily_profit_margin" json:"daily_profit_margin"`                           // 当日会员获利比
	DailyProfitMarginQuota        decimal.Decimal `gorm:"column:daily_profit_margin_quota" json:"daily_profit_margin_quota"`               // 当日获利比触发额度值
	RiskWarning                   int             `gorm:"column:risk_warning" json:"risk_warning"`                                         // 风控提醒开关1-开启 0关闭
	RiskWarningPeriod             int             `gorm:"column:risk_warning_period" json:"risk_warning_period"`                           // 风控提醒间隔(秒)
}

func (GameProfitMonitorConfig) TableName() string {
	return "game_profit_monitor_config"
}
