package rocketmq

import (
	"context"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go.uber.org/zap"

	"queueJob/pkg/common/config"
	"queueJob/pkg/zlogger"
)

type Consumer struct {
	group    string
	topic    string
	consumer rocketmq.PushConsumer
	hl       func(ctx context.Context, msgList ...*primitive.MessageExt) (consumer.ConsumeResult, error)
}

func NewConsumer(topic string, fun func(ctx context.Context, msgList ...*primitive.MessageExt) (consumer.ConsumeResult, error)) *Consumer {
	group := fmt.Sprintf(ConsumerGroupName, topic)
	cons, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer(config.Config.RocketMQ.RocketMQAddr),
		consumer.WithGroupName(group),
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromLastOffset),
	)
	if err != nil {
		zlogger.Errorw("new push consumer", zap.String("topic", topic), zap.Error(err))
		panic(err)
	}

	return &Consumer{group: group, topic: topic, consumer: cons, hl: fun}
}

func (c *Consumer) Start() {
	// 设置消息监听器
	if errInner := c.consumer.Subscribe(c.topic, consumer.MessageSelector{}, c.hl); errInner != nil {
		zlogger.Errorw("consumer subscribe error", zap.String("topic", c.topic), zap.Error(errInner))
		panic(errInner)
	}

	if errInner := c.consumer.Start(); errInner != nil {
		zlogger.Errorw("start consumer error", zap.Error(errInner))
		panic(errInner)
	}

	zlogger.Infow("RocketMQ consumer started",
		zap.Strings("nameServer", config.Config.RocketMQ.RocketMQAddr),
		zap.String("group", c.group),
		zap.String("topic", c.topic),
	)
	return
}

func (c *Consumer) Stop() {
	if err := c.consumer.Shutdown(); err != nil {
		zlogger.Errorw("Failed to shut down consumer", zap.Error(err),
			zap.Strings("nameServer", config.Config.RocketMQ.RocketMQAddr),
			zap.String("group", c.group),
			zap.String("topic", c.topic),
		)
		return
	}
	zlogger.Infow("RocketMQ consumer stopped",
		zap.Strings("nameServer", config.Config.RocketMQ.RocketMQAddr),
		zap.String("group", c.group),
		zap.String("topic", c.topic),
	)
}
