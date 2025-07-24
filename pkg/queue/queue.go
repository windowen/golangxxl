package queue

import (
	"time"

	"github.com/shopspring/decimal"
)

type MoneyChange struct {
	UserId       int             `json:"user_id"`       // 用户id
	CountryCode  string          `json:"country_code"`  // 国家编码
	CountryName  string          `json:"country_name"`  // 国家名称
	TransNo      string          `json:"trans_no"`      // 业务流水号
	ChangeType   int             `json:"change_type"`   // 账变类型1存款2提款3余额兑换钻石 4游戏流水
	ChangeAmount decimal.Decimal `json:"change_amount"` // 账变金额(有正和付)
	BeforeAmount decimal.Decimal `json:"before_amount"` // 账变前金额
	AfterAmount  decimal.Decimal `json:"after_amount"`  // 账变后金额
	ExchangeRate decimal.Decimal `json:"exchange_rate"` // 汇率
	Remark       string          `json:"remark"`        // 备注
	GameProvider string          `json:"game_provider"` // 游戏提供商
	GameType     string          `json:"game_type"`     // 游戏类型
	GameName     string          `json:"game_name"`     // 游戏名字
	GameOrderNo  string          `json:"game_order_no"` // 游戏业务订单号
	WagerCode    string          `json:"wager_code"`    // 赌注id
	TradeType    int             `json:"trade_type"`    // 交易类型 1投注 2结算
	CreatedAt    time.Time       `json:"created_at"`    // 创建时间
}
