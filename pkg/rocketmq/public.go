package rocketmq

import (
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"queueJob/pkg/service"
	"queueJob/pkg/zlogger"
)

type PublishProducer interface {
	Publish(topic string, st interface{})
	PublishWithDelay(topic string, msg interface{}, delayLevel int)
	PublishWithSeconds(topic string, msg interface{}, seconds int)
}

var PublicEvent PublishProducer = &emptyMethod{}

var Publish = PublicEvent.Publish

// PublishWithDelay 发送延时消息
// delayLevel 定义: 1s 5s 10s 30s 1m 2m 3m 4m 5m 6m 7m 8m 9m 10m 20m 30m 1h 2h
// 从1开始. 如果 level=1, 延时时间是 1s.
var PublishWithDelay = PublicEvent.PublishWithDelay

// PublishWithSeconds 发送延时秒数的消息
var PublishWithSeconds = PublicEvent.PublishWithSeconds

type emptyMethod struct{}

func (pe *emptyMethod) Publish(topic string, st interface{}) {}

func (pe *emptyMethod) PublicUserId(uid int, pb proto.Message) {}

func (pe *emptyMethod) PublishWithDelay(topic string, msg interface{}, delayLevel int) {}

func (pe *emptyMethod) PublishWithSeconds(topic string, msg interface{}, seconds int) {}

func Init() {
	rp := newProducer(func(producer *Producer) {
		// Replace implementation
		PublicEvent = &producerMethod{
			producer: producer,
		}
		Publish = PublicEvent.Publish
		PublishWithDelay = PublicEvent.PublishWithDelay
		PublishWithSeconds = PublicEvent.PublishWithSeconds
	})

	service.RegisterService(rp)
}

type producerMethod struct {
	producer *Producer
}

// Publish 发送异步队列消息(json)
func (pe *producerMethod) Publish(topic string, st interface{}) {
	if err := pe.producer.AsyncSendJson(topic, st); err != nil {
		zlogger.Errorw("Publish rocketmq send error", zap.Error(err))
		return
	}
}

// PublishWithDelay 发送延迟队列消息(json)
func (pe *producerMethod) PublishWithDelay(topic string, msg interface{}, delayLevel int) {
	if err := pe.producer.SendAsyncWithDelayJson(topic, msg, delayLevel); err != nil {
		zlogger.Errorw("PublishWithDelayJson rocketmq send error", zap.Error(err))
		return
	}
}

// PublishWithSeconds 发送延迟队列消息(json)
func (pe *producerMethod) PublishWithSeconds(topic string, msg interface{}, seconds int) {
	if err := pe.producer.SendAsyncJsonWithSeconds(topic, msg, seconds); err != nil {
		zlogger.Errorw("PublishWithSeconds rocketmq send error", zap.Error(err))
		return
	}
}
