package queue

import (
	"time"

	"github.com/shopspring/decimal"
)

// ChannelStatsType 定义渠道统计事件的类型
type ChannelStatsType int

const (
	// StatsUserRegistrations 用户注册事件
	StatsUserRegistrations ChannelStatsType = iota + 1
	// StatsUserFirstRecharge 用户首充
	StatsUserFirstRecharge
	// StatsUserRecharge 用户充值
	StatsUserRecharge
	// StatsWithdraw 用户提现
	StatsWithdraw
	// StatsBet 用户投注
	StatsBet
	DeviceActive
	UserActive
	StatsExchangeDiamond // 钻石兑换
	StatsPlayGame        // 打开游戏
)

// ChannelStats 统计事件
type ChannelStats struct {
	StatsType    ChannelStatsType `json:"statsType"`              // 事件类型
	UserId       int              `json:"userId"`                 // 用户id
	ChannelCode  int              `json:"channelCode"`            // 渠道唯一id
	RechargeNum  decimal.Decimal  `json:"rechargeNum,omitempty"`  // 充值金额
	WithdrawNum  decimal.Decimal  `json:"withdrawNum,omitempty"`  // 提现金额
	BetNum       decimal.Decimal  `json:"betNum,omitempty"`       // 下注金额
	RecordId     int              `json:"recordId,omitempty"`     // 订单号
	Email        string           `json:"email,omitempty"`        // 邮箱
	Ip           string           `json:"ip,omitempty"`           // ip地址
	Platform     string           `json:"platform,omitempty"`     // 平台码
	DeviceId     string           `json:"deviceId,omitempty"`     // 设备id
	RegisterTime time.Time        `json:"registerTime,omitempty"` // 注册时间
	Active1      int              `json:"active1,omitempty"`      // 是否1日活跃
	Active7      int              `json:"active7,omitempty"`      // 是否7日活跃
	Active30     int              `json:"active30,omitempty"`     // 是否30日活跃
	IsNewDevice  int              `json:"IsNewDevice,omitempty"`  // 是否是新用户
	DiamondNum   int              `json:"diamondNum,omitempty"`   // 是否是新用户
}
