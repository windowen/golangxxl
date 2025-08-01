package queue

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type MoneyChange struct {
	Id                 int             `json:"id,omitempty"`                    // 表id
	UserId             int             `json:"user_id,omitempty"`               // 用户id
	CountryCode        string          `json:"country_code,omitempty"`          // 国家编码
	CountryName        string          `json:"country_name,omitempty"`          // 国家名称
	TransNo            string          `json:"trans_no,omitempty"`              // 业务流水号
	ChangeType         int             `json:"change_type,omitempty"`           // 账变类型1存款2提款3余额兑换钻石 4游戏流水
	ChangeAmount       decimal.Decimal `json:"change_amount,omitempty"`         // 账变金额(有正和付)
	BeforeAmount       decimal.Decimal `json:"before_amount,omitempty"`         // 账变前金额
	AfterAmount        decimal.Decimal `json:"after_amount,omitempty"`          // 账变后金额
	AfterFlowAmount    decimal.Decimal `json:"after_flow_amount,omitempty"`     // 账变后剩余打码金额
	ReqAmount          decimal.Decimal `json:"req_amount,omitempty"`            // 客户端发过来的请求金额
	ExchangeRate       decimal.Decimal `json:"exchange_rate,omitempty"`         // 汇率
	Currency           string          `json:"currency,omitempty"`              // 货币
	Remark             string          `json:"remark,omitempty"`                // 备注
	GameProvider       string          `json:"game_provider,omitempty"`         // 游戏提供商
	GameType           string          `json:"game_type,omitempty"`             // 游戏类型
	GameName           string          `json:"game_name,omitempty"`             // 游戏名字
	GameOrderNo        string          `json:"game_order_no,omitempty"`         // 游戏业务订单号
	WagerCode          string          `json:"wager_code,omitempty"`            // 赌注id
	TradeType          int             `json:"trade_type,omitempty"`            // 交易类型 1投注 2结算
	WalletKey          string          `json:"wallet_key,omitempty"`            // 钱包key
	BatchNum           string          `json:"batch_num,omitempty"`             // 游戏类型
	ActChangeAmount    decimal.Decimal `json:"act_change_amount,omitempty"`     // 活动账变金额(有正和付)
	ActBeforeAmount    decimal.Decimal `json:"act_before_amount,omitempty"`     // 活动账变前金额
	ActAfterAmount     decimal.Decimal `json:"act_after_amount,omitempty"`      // 活动账变后金额
	ActAfterFlowAmount decimal.Decimal `json:"act_after_flow_amount,omitempty"` // 活动账变后剩余打码金额
	CreatedAt          time.Time       `json:"created_at,omitempty"`            // 创建时间
}

func (mc MoneyChange) String() string {
	data, err := json.MarshalIndent(mc, "", "  ")
	if err != nil {
		return fmt.Sprintf("MoneyChange: error marshaling to JSON: %v", err)
	}
	return string(data)
}

// ActType 定义活动类型枚举
type ActType int

const (
	ActTypeNewRegistrationRecharge ActType = iota + 1 // 1: 新注册充值返现
	ActTypeDailyLossCashback                          // 2: 每日亏损返现
	ActTypeDailyBetCashback                           // 3: 每日投注返现
	ActTypeFreeBonusActivity                          // 4: 免费奖金活动
	ActTypeReferralActivity                           // 5: 推荐活动
	ActTypeRecharge                                   // 6: 新注册充值前七天充值活动
	ActTypeNewRegister                                // 7: 新用户注册
)

// ActCheck 活动触发检测
type ActCheck struct {
	ActType  ActType         `json:"actType"`            // 活动类型
	UserId   int             `json:"userId"`             // 用户id
	Amount   decimal.Decimal `json:"amount,omitempty"`   // 数值
	Ip       string          `json:"ip,omitempty"`       // ip
	DeviceId string          `json:"deviceId,omitempty"` // 设备id
}

// BlackMonitorReq 黑名单监控请求
type BlackMonitorReq struct {
	BlackType         string // 黑名单类型 Login Register Withdrawal
	UserId            int    // 用户Id
	Ip                string // 同IP
	DeviceIdentifier  string // 同设备
	WithdrawalAccount string // 同提现账户
}

const (
	BlackLogin      = "Login"      // 登录黑名单
	BlackRegister   = "Register"   // 注册黑名单
	BlackWithdrawal = "Withdrawal" // 提现黑名单
)

type LiveRoomUserJoinReq struct {
	RoomId     int `json:"roomId"`     // 直播间id
	UserId     int `json:"userId"`     // 用户id
	IsTraveler int `json:"isTraveler"` // 是否游客
}
