package rocketmq

const (
	ProducerGroupName = "LiveProducerGroup" // 直播生产者组名
	ConsumerGroupName = "Group_%s"          // 直播消费者组名

	LiveRoomSendGift        = "live_room_send_gift"          // 直播间发送礼物topic
	LiveRoomStart           = "live_room_start"              // 主播开播 topic
	LiveRoomStop            = "live_room_stop"               // 主播下播 topic
	SiteUserRegister        = "site_user_register"           // 用户注册 topic
	SiteMessage             = "site_message"                 // 站点通知消息 topic
	LiveRoomStartFeeLive    = "live_room_start_fee_delay"    // 主播开播付费-分钟延迟队列 topic
	VipLevelUp              = "vip_level_up"                 // vip升级 topic
	FinanceCancel           = "finance_order_cancel"         // 充值订单取消 topic
	LiveRoomTransferPayLive = "live_room_transfer_pay_delay" // 主播转付费 topic
	LiveRoomRobotDelay      = "live_room_robot_delay"        // 直播间机器人延迟队列 topic
	LiveRoomRobotEnter      = "live_room_robot_enter"        // 直播间机器人延迟进入队列 topic
	LiveRoomRobotLeave      = "live_room_robot_leave"        // 直播间机器人延迟退出队列 topic
	StatsEvent              = "live_stats_event"             // 统计事件 topic
	ActCheck                = "activity_involve_check"       // 活动参与检查 topic
	SiteBlackMonitor        = "site_black_monitor"           // 黑名单监控 topic
	SystemStats             = "live_system_stats"            // 系统统计事件 topic
)
