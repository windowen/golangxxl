package rocketmq

import (
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"liveJob/pkg/service"
	"liveJob/pkg/zlogger"
)

type PublishProducer interface {
	Publish(topic string, msg []byte)
	PublishProto(topic string, pb proto.Message)
	PublishJson(topic string, st interface{})
	PublishWithDelay(topic string, msg []byte, delayLevel int)

	PublishWithDelayJson(topic string, msg interface{}, delayLevel int)
}

var PublicEvent PublishProducer = &emptyMethod{}

var Publish = PublicEvent.Publish
var PublishProto = PublicEvent.PublishProto
var PublishJson = PublicEvent.PublishJson

// PublishWithDelay 发送延时消息
// delayLevel 定义: 1s 5s 10s 30s 1m 2m 3m 4m 5m 6m 7m 8m 9m 10m 20m 30m 1h 2h
// 从1开始. 如果 level=1, 延时时间是 1s.
var PublishWithDelay = PublicEvent.PublishWithDelay

// PublishWithDelayJson 发送延时消息
// delayLevel 定义: 1s 5s 10s 30s 1m 2m 3m 4m 5m 6m 7m 8m 9m 10m 20m 30m 1h 2h
// 从1开始. 如果 level=1, 延时时间是 1s.
var PublishWithDelayJson = PublicEvent.PublishWithDelayJson

type emptyMethod struct{}

func (pe *emptyMethod) Publish(topic string, message []byte) {}

func (pe *emptyMethod) PublishProto(topic string, pb proto.Message) {}

func (pe *emptyMethod) PublishJson(topic string, st interface{}) {}

func (pe *emptyMethod) PublicUserId(uid int, pb proto.Message) {}

func (pe *emptyMethod) PublishWithDelay(topic string, msg []byte, delayLevel int) {}

func (pe *emptyMethod) PublishWithDelayJson(topic string, msg interface{}, delayLevel int) {}

func Init() {
	rp := newProducer(func(producer *Producer) {
		// Replace implementation
		PublicEvent = &producerMethod{
			producer: producer,
		}
		Publish = PublicEvent.Publish
		PublishProto = PublicEvent.PublishProto
		PublishJson = PublicEvent.PublishJson
		PublishWithDelay = PublicEvent.PublishWithDelay
		PublishWithDelayJson = PublicEvent.PublishWithDelayJson

	})

	service.RegisterService(rp)
}

type producerMethod struct {
	producer *Producer
}

func (pe *producerMethod) PublishProto(topic string, pb proto.Message) {
	if err := pe.producer.AsyncSendPb(topic, pb); err != nil {
		zlogger.Errorw("PublishProto rocketmq send error", zap.Error(err))
		return
	}
}

// Publish 发送异步队列消息
func (pe *producerMethod) Publish(topic string, msg []byte) {
	if err := pe.producer.SendAsync(topic, msg); err != nil {
		zlogger.Errorw("Publish rocketmq send error", zap.Error(err))
		return
	}
}

// PublishWithDelay 发送延迟队列消息
func (pe *producerMethod) PublishWithDelay(topic string, msg []byte, delayLevel int) {
	if err := pe.producer.SendAsyncWithDelay(topic, msg, delayLevel); err != nil {
		zlogger.Errorw("PublishWithDelay rocketmq send error", zap.Error(err))
		return
	}
}

// PublishJson 发送异步队列消息(json)
func (pe *producerMethod) PublishJson(topic string, st interface{}) {
	if err := pe.producer.AsyncSendJson(topic, st); err != nil {
		zlogger.Errorw("PublishJson rocketmq send error", zap.Error(err))
		return
	}
}

// PublishWithDelayJson 发送延迟队列消息(json)
func (pe *producerMethod) PublishWithDelayJson(topic string, msg interface{}, delayLevel int) {
	if err := pe.producer.SendAsyncPbWithDelayJson(topic, msg, delayLevel); err != nil {
		zlogger.Errorw("PublishWithDelayJson rocketmq send error", zap.Error(err))
		return
	}
}
