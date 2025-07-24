package kafkaconsumer

import (
	"queueJob/pkg/kafka"
	"queueJob/pkg/service"
)

// Init 初始化消费者,还有一个game_record_topic游戏使用，不在此项目
func Init() {
	consumerMap := map[string]func(msg []byte){
		kafka.FinanceMoneyChange:     bChange.handleMessages,           // 美元账变记录
		kafka.UserJoinRoomTopic:      ujr.handleMessages,               // 用户加入直播间
		kafka.StreamerReceiveDiamond: receiveDiamond.handleMessages,    // 主播累积收到钻石
		kafka.ChannelStatsTopic:      cStats.handleMessages,            // 渠道数据统计
		kafka.CheckRegisterBonus:     checkBonusMessage.handleMessages, // 检测注册福利是否可提现
	}

	// 启动消费者, 每个topic一个消费者组
	for topic, fun := range consumerMap {
		service.RegisterService(kafka.NewConsumer(topic, fun))
	}

	return
}
