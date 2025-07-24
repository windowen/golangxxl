package constant

const (
	RoomStatusNormal  = 1 // 正常
	RoomStatusDisable = 2 // 禁用
	RoomStatusDel     = 3 // 删除

	RoomOpCategoryMute    = 1 // 禁言
	RoomOpCategoryKickOut = 2 // 踢出

	RoomSceneStatusDo  = 1 // 直播中
	RoomSceneStatusEnd = 2 // 已结束

	RoomChargingRulesFree   = 1 // 免费
	RoomChargingRulesMinute = 2 // 分钟
	RoomChargingRulesScene  = 3 // 正场
)

const (
	RoomNotifyTypeLogin  = "login"
	RoomNotifyTypeSay    = "say"
	RoomNotifyTypeLike   = "like"
	RoomNotifyTypeMute   = "mute"
	RoomNotifyTypeLeave  = "leave"
	RoomNotifyTypeFollow = "follow"
	RoomNotifyTypeIncome = "income"
)

const (
	NormalUserCount = 10
)

const (
	RoomWelcomeMessage = "欢迎%s加入%s房间"
	LeaveMessage       = "%s离开%s房间"
	FollowMessage      = "%s关注了直播间"
)
