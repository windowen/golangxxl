package redis

const (
	UserCacheInfoKey      = "user_cache_info_%v"          // 用户信息缓存key
	UserLoginToken        = "user_login_token_%v"         // 用户登陆令牌Key
	UserLoginTokenVersion = "user_login_token_version_%v" // 用户登陆令牌版本Key
	LevelDiamondConfig    = "g:level:diamond"             // 累计消费钻石对应vip等级表
	LevelStreamerConfig   = "g:level:streamer"            // 主播累计收入钻石对应vip等级表
	PaymentRecordKey      = "p:payment:record:%d"         // 用户充值缓存
	UserPay               = "pay_%d"                      // 修改用户钱包的分布式锁key
	UserWalletCache       = "u:wallet:%d"                 // 用户钱包缓存
)
