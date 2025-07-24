package message

type MessageType int

const (
	RechargeSuccess         MessageType = 1 // 充值成功通知
	RechargeFailure         MessageType = 2 // 充值失败通知
	WithdrawalSuccess       MessageType = 3 // 提现成功通知
	WithdrawalFailure       MessageType = 4 // 提现失败通知
	DiamondExchange         MessageType = 5 // 兑换钻石
	Admin                   MessageType = 6 // 管理后台通知
	AnnouncementMaintenance MessageType = 7 // 公告维护
	SMS                     MessageType = 8 // 短信
	Email                   MessageType = 9 // 邮件
)

// PackageMessage 消息包
type PackageMessage struct {
	MessageType MessageType      `json:"message_type"` // 消息类型 对应sys_site_template里的category字段
	UserId      int              `json:"user_id"`      // 用户id
	Language    string           `json:"language"`     // 用户语言
	RechargeMsg *RechargeMessage `json:"recharge_msg"` // 充值信息
}

// RechargeMessage 充值成功消息参数
type RechargeMessage struct {
	CoinCode string `json:"coin_code"` // 货币类型
	Amount   int    `json:"amount"`    // 货币数量
	USD      int    `json:"usd"`       // USD数量
	Diamond  int    `json:"diamond"`   // 钻石数量
}

// AccumulationDiamond 钻石消费的总数量
type AccumulationDiamond struct {
	UserId int `json:"user_id"` // 用户id
	RoomId int `json:"room_id"` // 房间id
}
