package table

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type MoneyChange struct {
	Id                 int             `gorm:"column:id" json:"id"`                                       // 账变记录表id
	UserId             int             `gorm:"column:user_id" json:"user_id"`                             // 用户id
	CountryCode        string          `gorm:"column:country_code" json:"country_code"`                   // 国家编码
	CountryName        string          `gorm:"column:country_name" json:"country_name"`                   // 国家名称
	TransNo            string          `gorm:"column:trans_no" json:"trans_no"`                           // 业务流水号
	ChangeType         int             `gorm:"column:change_type" json:"change_type"`                     // 账变类型 1存款 2提款 3余额兑换钻石 4游戏流水 5彩金领取 6彩金转余额
	ChangeAmount       decimal.Decimal `gorm:"column:change_amount" json:"change_amount"`                 // 账变金额(有正和付)
	BeforeAmount       decimal.Decimal `gorm:"column:before_amount" json:"before_amount"`                 // 账变前金额
	AfterAmount        decimal.Decimal `gorm:"column:after_amount" json:"after_amount"`                   // 账变后金额
	AfterFlowAmount    decimal.Decimal `gorm:"column:after_flow_amount" json:"after_flow_amount"`         // 账变后剩余打码金额
	ReqAmount          decimal.Decimal `gorm:"column:req_amount" json:"req_amount"`                       // 客户端发过来的请求金额(当地币种）
	ExchangeRate       decimal.Decimal `gorm:"column:exchange_rate" json:"exchange_rate"`                 // 汇率
	Currency           string          `gorm:"column:currency" json:"currency"`                           // 货币类型
	Remark             string          `gorm:"column:remark" json:"remark"`                               // 备注
	GameProvider       string          `gorm:"column:game_provider" json:"game_provider"`                 // 游戏提供商
	GameType           string          `gorm:"column:game_type" json:"game_type"`                         // 游戏类型
	GameName           string          `gorm:"column:game_name" json:"game_name"`                         // 游戏名字
	GameOrderNo        string          `gorm:"column:game_order_no" json:"game_order_no"`                 // 游戏业务订单号
	WagerCode          string          `gorm:"column:wager_code" json:"wager_code"`                       // 赌注id
	TradeType          int             `gorm:"column:trade_type" json:"trade_type"`                       // 交易类型 1投注 2结算
	BatchNum           string          `gorm:"column:batch_num" json:"batch_num"`                         // 游戏类型
	ActChangeAmount    decimal.Decimal `gorm:"column:act_change_amount" json:"act_change_amount"`         // 活动彩金账变金额(有正和付)
	ActBeforeAmount    decimal.Decimal `gorm:"column:act_before_amount" json:"act_before_amount"`         // 活动彩金账变前金额
	ActAfterAmount     decimal.Decimal `gorm:"column:act_after_amount" json:"act_after_amount"`           // 活动账变后金额
	ActAfterFlowAmount decimal.Decimal `gorm:"column:act_after_flow_amount" json:"act_after_flow_amount"` // 活动账变后剩余打码金额
	SiteId             int             `gorm:"column:site_id" json:"site_id"`                             // 站点id
	CreatedAt          time.Time       `gorm:"column:created_at" json:"created_at"`                       // 创建时间
	CreatedDay         string          `gorm:"column:created_day" json:"created_day"`                     // 创建日期天
}

func (m MoneyChange) TableName() string {
	return fmt.Sprintf("finance_money_change_%d", m.UserId%10)
}
