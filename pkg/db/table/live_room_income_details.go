package table

import "time"

type LiveRoomIncomeDetails struct {
	Id               int       `gorm:"column:id" json:"id"`
	BillNo           string    `gorm:"column:bill_no" json:"bill_no"`                     // 直播礼物订单号
	UserId           int       `gorm:"column:user_id" json:"user_id"`                     // 消费人用户id
	RoomId           int       `gorm:"column:room_id" json:"room_id"`                     // 直播间id
	FamilyId         int       `gorm:"column:family_id" json:"family_id"`                 // 家族id
	AnchorId         int       `gorm:"column:anchor_id" json:"anchor_id"`                 // 主播id
	SceneId          int       `gorm:"column:scene_id" json:"scene_id"`                   // 直播场次id
	FamilyMasterId   int       `gorm:"column:family_master_id" json:"family_master_id"`   // 家族长id
	CountryName      string    `gorm:"column:country_name" json:"country_name"`           // 国家名字( 如巴西)
	CountryCode      string    `gorm:"column:country_code" json:"country_code"`           // 国家编码(eg:BL)
	Category         int       `gorm:"column:category" json:"category"`                   // 类别 1-礼物 2-弹幕 3-游戏 4-直播收费 5-道具
	BillingMode      int       `gorm:"column:billing_mode" json:"billing_mode"`           // 直播类型收费模式1-按时收费2-按场收费
	ProjectId        int       `gorm:"column:project_id" json:"project_id"`               // 收入项目id(礼物id、游戏id)
	ProjectNum       int       `gorm:"column:project_num" json:"project_num"`             // 数量
	UnitPrice        int       `gorm:"column:unit_price" json:"unit_price"`               // 单价
	ProjectTotal     int       `gorm:"column:project_total" json:"project_total"`         // 合计(如: 礼物单价 * 数量)
	GiftRatio        int       `gorm:"column:gift_ratio" json:"gift_ratio"`               // 主播抽成比例  实际*100
	AnchorIncome     int64     `gorm:"column:anchor_income" json:"anchor_income"`         // 主播收入
	PlatformRatio    int       `gorm:"column:platform_ratio" json:"platform_ratio"`       // 平台抽成比例  实际*100
	PlatformIncome   int64     `gorm:"column:platform_income" json:"platform_income"`     // 平台收入
	FamilyRatio      int       `gorm:"column:family_ratio" json:"family_ratio"`           // 家族抽成比例  实际*100
	FamilyIncome     int64     `gorm:"column:family_income" json:"family_income"`         // 家族收入
	IsDivide         int       `gorm:"column:is_divide" json:"is_divide"`                 // 是否需要分成
	SettlementStatus int       `gorm:"column:settlement_status" json:"settlement_status"` // 结算状态 0-未结算, 1-已结算，2-冻结，3-结算失败
	BillSerial       string    `gorm:"column:bill_serial" json:"bill_serial"`             // 直播消费流水号
	CreatedAt        time.Time `gorm:"column:created_at" json:"created_at"`               // 创建时间
	UpdatedAt        time.Time `gorm:"column:updated_at" json:"updated_at"`               // 更新时间
}

func (LiveRoomIncomeDetails) TableName() string {
	return "live_room_income_details"
}
