package table

import (
	"time"

	"github.com/shopspring/decimal"
)

type ChannelDailyUserStats struct {
	Id                      int             `gorm:"column:id" json:"id"`                                               // 自增主键，唯一标识每条记录
	ReportDate              time.Time       `gorm:"column:report_date" json:"report_date"`                             // 统计日期，用于标识每天的统计数据
	UserRegistrations       int             `gorm:"column:user_registrations" json:"user_registrations"`               // 用户注册数
	UniqueVisitors          int             `gorm:"column:unique_visitors" json:"unique_visitors"`                     // 唯一访客数（UV）
	ActiveUsers             int             `gorm:"column:active_users" json:"active_users"`                           // 每日活跃用户数
	RetentionRate           decimal.Decimal `gorm:"column:retention_rate" json:"retention_rate"`                       // 留存率，百分比，保留两位小数
	PageStayTime            int             `gorm:"column:page_stay_time" json:"page_stay_time"`                       // 游戏页面停留时间，单位为秒
	GameLaunchCount         int             `gorm:"column:game_launch_count" json:"game_launch_count"`                 // 游戏启动次数
	GameAwardCount          int             `gorm:"column:game_award_count" json:"game_award_count"`                   // 游戏中奖次数
	HomepageBannerClicks    int             `gorm:"column:homepage_banner_clicks" json:"homepage_banner_clicks"`       // 首页 banner 点击次数
	RecommendedBannerClicks int             `gorm:"column:recommended_banner_clicks" json:"recommended_banner_clicks"` // 推荐 banner 点击次数
	PopularBannerClicks     int             `gorm:"column:popular_banner_clicks" json:"popular_banner_clicks"`         // 热门 banner 点击次数
}

func (ChannelDailyUserStats) TableName() string {
	return "channel_daily_user_stats"
}
