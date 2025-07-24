package table

import (
	"encoding/json"
	"time"
)

type LiveRoom struct {
	Id                  int             `gorm:"column:id" json:"id"`
	Sort                int             `gorm:"column:sort" json:"sort"`                                 // 排序 不能存在相同的排序，如果已存在先设置为0
	UserId              int             `gorm:"column:user_id" json:"user_id"`                           // 主播id
	ChatRoomId          string          `gorm:"column:chat_room_id" json:"chat_room_id"`                 // chat房间id(创建主播时调用rpc返回)
	CountryCode         string          `gorm:"column:country_code" json:"country_code"`                 // 国家code
	TagsJson            json.RawMessage `gorm:"column:tags_json" json:"tags_json"`                       // 标签 ["1", "2", "3"]  二级标签id
	Title               string          `gorm:"column:title" json:"title"`                               // 直播间标题
	Cover               string          `gorm:"column:cover" json:"cover"`                               // 封面图片
	Summary             string          `gorm:"column:summary" json:"summary"`                           // 简介
	VideoClarity        int             `gorm:"column:video_clarity" json:"video_clarity"`               // 视屏清晰度 1- 480p 2-540p 3-720p
	GiftRatio           int             `gorm:"column:gift_ratio" json:"gift_ratio"`                     // 主播抽成比例  实际*100
	PlatformRatio       int             `gorm:"column:platform_ratio" json:"platform_ratio"`             // 平台抽成比例  实际*100
	FamilyRatio         int             `gorm:"column:family_ratio" json:"family_ratio"`                 // 家族抽成比例  实际*100
	LevelId             int             `gorm:"column:level_id" json:"level_id"`                         // 等级ID
	SetLevelId          int             `gorm:"column:set_level_id" json:"set_level_id"`                 // 后台设置等级id 优先展示
	AccumulationDiamond int64           `gorm:"column:accumulation_diamond" json:"accumulation_diamond"` // 直播间累积收到的钻石(只有整数)
	Remark              string          `gorm:"column:remark" json:"remark"`                             // 备注
	Status              int             `gorm:"column:status" json:"status"`                             // 房间状态 1-正常  2-禁用 3-删除
	LiveStatus          int             `gorm:"column:live_status" json:"live_status"`                   // 直播状态 1-直播 2-下播
	PaidPurviewStatus   int             `gorm:"column:paid_purview_status" json:"paid_purview_status"`   // 是否有可以开启付费直播
	BottomSort          int             `gorm:"column:bottom_sort" json:"bottom_sort"`                   // 置底(默认0  比如 -1倒数第一 -2 倒数第二,每次设置值底需要重新排序)
	LastStartLiveTime   time.Time       `gorm:"column:last_start_live_time" json:"last_start_live_time"` // 最后开播时间
	LastEndLiveTime     time.Time       `gorm:"column:last_end_live_time" json:"last_end_live_time"`     // 最后下播时间
	CreatedAt           time.Time       `gorm:"column:created_at" json:"created_at"`                     // 创建时间
	UpdatedAt           time.Time       `gorm:"column:updated_at" json:"updated_at"`                     // 更新时间
}

func (LiveRoom) TableName() string {
	return "live_room"
}
