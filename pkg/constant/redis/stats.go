package redis

const (
	LiveStats                    = "stats_%s"                  // 统计信息
	EventUserRegistrations       = "user_registrations"        // 用户注册数
	EventUniqueVisitors          = "unique_visitors"           // 唯一访客数（UV）
	EventActiveUsers             = "active_users"              // 活跃用户数
	EventRetentionRate           = "retention_rate"            // 留存率
	EventGameLaunchCount         = "game_launch_count"         // 游戏启动次数
	EventGameAwardCount          = "game_award_count"          // 游戏中奖次数
	EventHomepageBannerClicks    = "homepage_banner_clicks"    // 首页 Banner 点击次数
	EventRecommendedBannerClicks = "recommended_banner_clicks" // 推荐 Banner 点击次数
	EventPopularBannerClicks     = "popular_banner_clicks"     // 热门 Banner 点击次数
	EventPageStayTime            = "page_stay_time:%s"         // 页面停留时间
	StatsSyncTime                = "stats:last_updated_date"   // redis同步到mysql的最后时间
)
