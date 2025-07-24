package consumer

import (
	"context"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"

	"queueJob/pkg/rocketmq"
	"queueJob/pkg/service"
)

// Init 初始化消费者
func Init() {
	consumerMap := map[string]func(ctx context.Context, msgList ...*primitive.MessageExt) (consumer.ConsumeResult, error){
		rocketmq.LiveRoomSendGift:        gift.handleMessages,                // 添加gift消费主题
		rocketmq.LiveRoomStop:            live.handleMessages,                // 添加下播消费
		rocketmq.SiteUserRegister:        userReg.handleMessages,             // 添加用户注册消费
		rocketmq.SiteMessage:             notify.handleMessages,              // 站点消息通知
		rocketmq.LiveRoomStartFeeLive:    liveMPaid.handleMessages,           // 直播间分钟扣费延迟队列
		rocketmq.VipLevelUp:              vipLevel.handleMessages,            // 处理vip升级队列
		rocketmq.FinanceCancel:           fCancel.handleMessages,             // 用户未按时支付订单
		rocketmq.LiveRoomTransferPayLive: liveRoomTransferPay.handleMessages, // 主播转付费
		rocketmq.StreamerReceiveDiamond:  receiveDiamond.handleMessages,      // 主播累积收到钻石
		rocketmq.LiveRoomRobotDelay:      liveRoomRobot.handleMessages,       // 直播间机器人延迟队列
		rocketmq.FinanceMoneyChange:      bChange.handleMessages,             // 美元账变记录
		rocketmq.StatsEvent:              sEvent.handleMessages,              // 统计事件
	}

	// 启动消费者, 每个topic一个消费者组
	for topic, fun := range consumerMap {
		service.RegisterService(rocketmq.NewConsumer(topic, fun))
	}

	return
}
