package redis

const (
	UserCacheInfoKey      = "user_cache_info_:%d"          // 用户信息缓存key
	UserLoginToken        = "user_login_token_:%v"         // 用户登陆令牌Key
	UserLoginTokenVersion = "user_login_token_version_%:v" // 用户登陆令牌版本Key
	LevelDiamondConfig    = "g:level:diamond"              // 累计消费钻石对应vip等级表
	LevelStreamerConfig   = "g:level:streamer"             // 主播累计收入钻石对应vip等级表
	PaymentRecordKey      = "p:payment:record:%d"          // 用户充值缓存
	UserPay               = "pay_%d"                       // 修改用户钱包的分布式锁key
	UserWalletCache       = "u:wallet:%d"                  // 用户钱包缓存
	UserActWalletCache    = "u:act_wallet:%d"              // 用户活动钱包缓存, %d 用userId填充
	UserBonusWalletCache  = "u:bonus_wallet:%d"            // 用户彩金钱包缓存, %d 用userId填充
	GameProfitMonitor     = "game_profit_monitor_config"   // 游戏获利监控参数配置表缓存
	UserGameProfit        = "game:profit:%d:%s"            // 用户当天游戏赢利数据
	UserGameWager         = "game:wager:%d:%s"             // 用户当天游戏投注数据
	MonitorNotify         = "g:monitor"                    // 监控提醒
	UserDeviceId          = "u:device"                     // 用户设备id

	SameRelationIPHashKey                = "same_relation_ip_hash"      // 同IP关联IPKey
	SameRelationDeviceHashKey            = "same_relation_device_hash"  // 同设备关联设备Key
	SameRelationWithdrawalAccountHashKey = "same_relation_account_hash" // 同关联提现Key
	SameRelationUserHashKey              = "same_relation_user_hash"    // 同关联提现Key
)
