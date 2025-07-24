package redis

const (
	LiveStats                    = "stats_%s"                  // 统计信息
	EventUserRegistrations       = "user_registrations"        // 用户注册数
	EventGameLaunchCount         = "game_launch_count"         // 游戏启动次数
	EventGameAwardCount          = "game_award_count"          // 游戏中奖次数
	EventHomepageBannerClicks    = "homepage_banner_clicks"    // 首页 Banner 点击次数
	EventRecommendedBannerClicks = "recommended_banner_clicks" // 推荐 Banner 点击次数
	EventPopularBannerClicks     = "popular_banner_clicks"     // 热门 Banner 点击次数
	EventGameBannerClicks        = "game_banner_clicks"        // 游戏 Banner 点击次数
	EventUserH5Active            = "h5_active"                 // h5活跃数量
	EventPageStayTime            = "page_stay_time:%s"         // 页面停留时间
	StatsSyncTime                = "stats:last_updated_date"   // redis同步到mysql的最后时间
	CoinExchange                 = "g:exchange:%s"             // 币种转换缓存
)

// 系统统计

const (
	SystemStatsCountry         = "system_stats_country_%s"
	SystemStatsRegisterKey     = "system_stats_register_%s"
	SystemStatsReceiveAwardKey = "system_stats_receive_active_award_%s"
	SystemStatsPayKey          = "system_stats_pay_%s"
	SystemStatsPayTimeKey      = "system_stats_pay_time_%s"
	SystemStatsWithdrawKey     = "system_stats_withdraw_%s"
	SystemStatsRebatesKey      = "system_stats_rebates_%s"
	SystemStatsFirstPayKey     = "system_stats_first_pay_%s"
	SystemStatsGameKey         = "system_stats_game_%s"
	SystemStatsWagerKey        = "system_stats_wager_%s"
	SystemStatsSettleKey       = "system_stats_settle_%s"
	StatsSysSyncTime           = "system_stats:last_updated_date" // redis同步到mysql的最后时间
)

const (
	StatsAceLotteryTime = "ace_lottery:last_updated_date" // 15.查詢玩家遊戲記錄 Query/GameRecord
)
