package table

import (
	"time"

	"github.com/shopspring/decimal"
)

type FinancePayRecord struct {
	Id                     int             `gorm:"column:id" json:"id"`                                                                           // 充值表主键
	UserId                 int             `gorm:"column:user_id" json:"user_id"`                                                                 // 用户ID
	BillNo                 string          `gorm:"column:bill_no" json:"bill_no"`                                                                 // 系统账单号
	ThirdpartyOrderNo      string          `gorm:"column:thirdparty_order_no" json:"thirdparty_order_no"`                                         // 三方支付公司订单号
	OrderAmount            decimal.Decimal `gorm:"column:order_amount" json:"order_amount"`                                                       // 下单金额(法币)
	Amount                 decimal.Decimal `gorm:"column:amount" json:"amount"`                                                                   // 到账美金
	CompanyCode            string          `gorm:"column:company_code" json:"company_code"`                                                       // 支付公司编码（如TXZF）
	CountryCode            string          `gorm:"column:country_code" json:"country_code"`                                                       // 支付公司的国家编码(如BR)
	PayTypeCode            string          `gorm:"column:pay_type_code" json:"pay_type_code"`                                                     // 支付类型编码(如：ZFB)
	TransactionType        int             `gorm:"column:transaction_type" json:"transaction_type"`                                               // 交易类型1充值2上分
	Status                 int             `gorm:"column:status" json:"status"`                                                                   // 0待处理1处理中 2失败 3已成功 4已取消 5失效
	ExchangeRate           decimal.Decimal `gorm:"column:exchange_rate" json:"exchange_rate"`                                                     // 换算汇率
	UnionpayCode           string          `gorm:"column:unionpay_code" json:"unionpay_code"`                                                     // 银联编码(如选择的银联支付这个必须填写)
	EwalletCode            string          `gorm:"column:ewallet_code" json:"ewallet_code"`                                                       // 电子钱包编码（如选择的电子钱包支付这个必须填写）
	TransactionCompletedAt time.Time       `gorm:"column:transaction_completed_at;default:'1970-01-01 00:00:00'" json:"transaction_completed_at"` // 三方支付回调交易完成时间
	Remark                 string          `gorm:"column:remark" json:"remark"`                                                                   // 备注
	SiteId                 int             `gorm:"column:site_id" json:"site_id"`                                                                 // 站点id
	LanguageCode           string          `gorm:"column:language_code" json:"language_code"`                                                     // 下单时用户的语言
	CurrencyCode           string          `gorm:"column:currency_code" json:"currency_code"`                                                     // 币种编码
	UpdatedAt              time.Time       `gorm:"column:updated_at" json:"updated_at"`                                                           // 更新时间
	CreatedAt              time.Time       `gorm:"column:created_at" json:"created_at"`                                                           // 创建时间
}

func (FinancePayRecord) TableName() string {
	return "finance_pay_record"
}
