package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go.uber.org/zap"

	redis2 "queueJob/pkg/constant/redis"
	"queueJob/pkg/db/redisdb/redis"
	"queueJob/pkg/queue"
	"queueJob/pkg/zlogger"
)

// 统计事件
var sEvent = &statsEvent{}

type statsEvent struct{}

func (o *statsEvent) handleMessages(ctx context.Context, msgList ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for _, msg := range msgList {
		o.statsEvent(ctx, msg)
	}
	return consumer.ConsumeSuccess, nil
}

func (o *statsEvent) statsEvent(ctx context.Context, msg *primitive.MessageExt) {
	zlogger.Debugw("statsEvent, receive message",
		zap.String("msgID", msg.MsgId), zap.String("topic", msg.Topic), zap.String("body", string(msg.Body)))

	pMessage := &queue.EventStats{}
	if err := json.Unmarshal(msg.Body, pMessage); err != nil {
		zlogger.Errorw("statsEvent, unmarshal msg fail", zap.String("msgID", msg.MsgId), zap.Error(err))
		return
	}

	currDate := time.Now().Format("2006-01-02")

	switch pMessage.EventType {
	case queue.EventUserRegistrations:
		redis.HIncr(fmt.Sprintf(redis2.LiveStats, currDate), redis2.EventUserRegistrations)
	case queue.EventHomepageBannerClicks:
		redis.HIncr(fmt.Sprintf(redis2.LiveStats, currDate), redis2.EventHomepageBannerClicks)
	case queue.EventRecommendedBannerClicks:
		redis.HIncr(fmt.Sprintf(redis2.LiveStats, currDate), redis2.EventRecommendedBannerClicks)
	case queue.EventPopularBannerClicks:
		redis.HIncr(fmt.Sprintf(redis2.LiveStats, currDate), redis2.EventPopularBannerClicks)
	case queue.EventGameLaunch:
		redis.HIncr(fmt.Sprintf(redis2.LiveStats, currDate), redis2.EventGameLaunchCount)
	case queue.EventGameAward:
		redis.HIncr(fmt.Sprintf(redis2.LiveStats, currDate), redis2.EventGameAwardCount)
	case queue.EventPageStay:
		by, err := redis.HIncrBy(fmt.Sprintf(redis2.LiveStats, currDate), fmt.Sprintf(redis2.EventPageStayTime, pMessage.PageName), pMessage.Timestamp)
		if err != nil {
			break
		}

		zlogger.Debugw("statsEvent, event page stay time", zap.Int64("by", by))
	default:
		break
	}
}
