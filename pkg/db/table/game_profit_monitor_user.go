package table

import (
	"time"

	"github.com/shopspring/decimal"
)

type GameProfitMonitorUser struct {
	Id            int             `gorm:"column:id" json:"id"`                                                       // 游戏获利监控用户列表Id
	UserId        int             `gorm:"column:user_id" json:"user_id"`                                             // 用户Id
	UserAccount   string          `gorm:"column:user_account" json:"user_account"`                                   // 用户账户
	UserTag       string          `gorm:"column:user_tag" json:"user_tag"`                                           // 用户标签逗号隔开(比如：禁止提现, 风险用户)
	LoginIp       string          `gorm:"column:login_ip" json:"login_ip"`                                           // 登录Ip
	RegSource     string          `gorm:"column:reg_source" json:"reg_source"`                                       // 注册来源
	MonitorType   string          `gorm:"column:monitor_type" json:"monitor_type"`                                   // 监控类型（高倍爆奖、大额中奖、会员当天获利比）
	ActualValue   decimal.Decimal `gorm:"column:actual_value" json:"actual_value"`                                   // 实际值
	BetOrderNo    string          `gorm:"column:bet_order_no" json:"bet_order_no"`                                   // 注单编号(与游戏ES唯一注单对应)
	Status        int             `gorm:"column:status" json:"status"`                                               // 状态0-待处理 1-已忽略 2-已冻结 3-禁止领取奖励
	Operator      string          `gorm:"column:operator" json:"operator"`                                           // 操作人
	OperationTime time.Time       `gorm:"column:operation_time;default:'1970-01-01 00:00:00'" json:"operation_time"` // 操作时间
	Remark        string          `gorm:"column:remark" json:"remark"`                                               // 备注
	CreatedDate   time.Time       `gorm:"column:created_date" json:"created_date"`                                   // 创建日期（没时间只显示日期部分比如2024-12-26）
	CreatedAt     time.Time       `gorm:"column:created_at" json:"created_at"`                                       // 创建时间
	UpdatedAt     time.Time       `gorm:"column:updated_at" json:"updated_at"`                                       // 更新时间
}

func (GameProfitMonitorUser) TableName() string {
	return "game_profit_monitor_user"
}
