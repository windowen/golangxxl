package table

import "time"

type LiveRoomIncomeSettlement struct {
	Id               int       `gorm:"column:id" json:"id"`                               // 唯一标识符，自增ID
	OwnerId          int       `gorm:"column:owner_id" json:"owner_id"`                   // 房主ID
	SettlementCycle  int       `gorm:"column:settlement_cycle" json:"settlement_cycle"`   // 结算周期(202401 表示24年第一周）
	FamilyId         int       `gorm:"column:family_id" json:"family_id"`                 // 家族id
	FamilyMasterId   int       `gorm:"column:family_master_id" json:"family_master_id"`   // 家族长id
	RoomId           int       `gorm:"column:room_id" json:"room_id"`                     // 房间ID
	SceneId          int       `gorm:"column:scene_id" json:"scene_id"`                   // 直播场次id
	CountryName      string    `gorm:"column:country_name" json:"country_name"`           // 国家名字
	StreamerIncome   int64     `gorm:"column:streamer_income" json:"streamer_income"`     // 主播上周总流水
	PlatformIncome   int64     `gorm:"column:platform_income" json:"platform_income"`     // 平台收入
	FamilyIncome     int64     `gorm:"column:family_income" json:"family_income"`         // 家族长收入
	SettlementTime   time.Time `gorm:"column:settlement_time" json:"settlement_time"`     // 上周结算时间
	SettlementStatus int       `gorm:"column:settlement_status" json:"settlement_status"` // 结算状态 0-未结算, 1-已结算，2-冻结，3-结算失败
}

func (LiveRoomIncomeSettlement) TableName() string {
	return "live_room_income_settlement"
}
