package redis

const (
	RoomCacheInfoKey          = "room_cache_info_:%d"          // 直播间信息缓存key
	RoomManageCacheKey        = "room_manage_cache_:%v"        // 直播间房管缓存key
	RoomCacheOnlineKey        = "room_cache_online_:%v"        // 直播间在线用户缓存key
	RoomCacheLiveKey          = "room_cache_like_:%v"          // 点赞
	RoomCacheStatKey          = "room_cache_stat_:%v"          // 房间数据统计
	RoomCacheChatHisKey       = "room_cache_chat_history_:%v"  // 聊天历史
	RoomCacheChatChatRooms    = "room_cache_chat_rooms"        // 房间缓存
	RoomCacheChatBarrageKey   = "room_cache_barrage_:%v"       // 弹幕消息
	RoomCacheKickOutKey       = "room_cache_kick_out_%:v"      // T出
	RoomCacheDefaultChatRooms = "room_cache_default_chat_room" // 房间缓存
	RoomStartLiveSet          = "room_start_live"              // 开播直播间缓存
	RoomRobotConfigHash       = "room_robot_config"            // 直播间机器人配置
	RoomRobotListZSet         = "room_robot_list_:%v"          // 直播间机器人列表
	RoomRobotSet              = "robot_list"                   // 平台机器人列表
	RoomRobotNum              = "room_robot_num_:%d"           // 直播间机器人数量
	UserWagerInfo             = "user_wager_info:%d:%s"        // 用户赌注缓存key
	NobleCfgCacheHash         = "noble_config_cache"           // 贵族配置缓存信息
	UserRecharge              = "u:recharge:%d"                // 用户是否充值过
	UserBitmapCheckRecharge   = "u:c:recharge:%d"              // 用户是否充值过
	BitmapUserRecharge        = "u:bitmap:recharge:%d"         // 用户连续充值记录
	DeviceKey                 = "g:device:%s"                  // 设备活跃统计
	PlayGameStats             = "u:play_game:%d"               // 用户玩游戏统计
	ExchangeDiamondStats      = "u:exchange_diamond:%d"        // 用户兑换钻石统计
)

// 场数据
const (
	SceneFollow        = "room_scene_follow_:%v"     // 直播间场关注数据缓存key  v1-房间id
	SceneTrial         = "room_scene_trial_:%v"      // 直播间试看信息缓存key   v1-房间id
	SceneIncome        = "room_scene_income_:%v"     // 直播间场收益数据缓存key  v1-房间id
	SceneMute          = "room_scene_mute_:%v"       // 直播间场禁言数据缓存key  v1-房间id
	SceneKickOut       = "room_scene_kick_out_:%v"   // 直播间场踢出数据缓存key  v1-房间id
	SceneBlock         = "room_scene_block_:%v"      // 直播间拉黑数据缓存key  v1-房间id
	ScenePayUsers      = "room_scene_pay_users_:%d"  // 直播间付费用户数据缓存key  v1-房间id
	SceneUsersZSet     = "room_scene_users_:%v"      // 直播间用户数据缓存key  v1-房间id
	SceneRealUsersZSet = "room_scene_real_users_:%v" // 直播间真实用户数据缓存key  v1-房间id
)

type StatField string

const (
	StatKeyMaxOnline     StatField = "maxOnline"
	StatKeyFollowerCount StatField = "followerCount"
	StatKeyGiftCount     StatField = "giftCount"
	StatKeyLikeCount     StatField = "LikeCount"
)

const (
	ChatHistoryStartIndex = 0    // 第一条从0开始
	ChatHistoryLimit      = 3    // 历史聊天记录限制条数 =
	LengthLimit           = 5000 // 20条消息的总长度(包装的结构大概有200字，100字留给内容，道具的话内容会x2)
	BarrageCost           = 1000 // 存的是分
)
