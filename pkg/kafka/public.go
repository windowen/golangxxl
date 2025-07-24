package kafka

import (
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"queueJob/pkg/service"
	"queueJob/pkg/zlogger"
)

type publicProducerData interface {
	Public(topic string, st interface{})
	PublicUserId(uid int, topic string, st interface{})
	PublicKey(key, topic string, st interface{})
	PublicPb(topic string, pb proto.Message)
	PublicPbUserId(uid int, topic string, pb proto.Message)
}

var publicEvent publicProducerData = &emptyMethod{}

var Public = publicEvent.Public
var PublicUserId = publicEvent.PublicUserId
var PublicKey = publicEvent.PublicKey

// var PublicPb = publicEvent.PublicPb
// var PublicPbUserId = publicEvent.PublicPbUserId

type emptyMethod struct{}

func (pe *emptyMethod) Public(topic string, st interface{}) {}

func (pe *emptyMethod) PublicUserId(uid int, topic string, st interface{}) {}

func (pe *emptyMethod) PublicKey(key, topic string, st interface{}) {}

func (pe *emptyMethod) PublicPb(topic string, pb proto.Message) {
}

func (pe *emptyMethod) PublicPbUserId(uid int, topic string, pb proto.Message) {
}

func Init() {
	kp := newProducer(func(producer *Producer) {
		// 替换实现
		publicEvent = &producerMethod{
			producer: producer,
		}
		Public = publicEvent.Public
		PublicUserId = publicEvent.PublicUserId
		PublicKey = publicEvent.PublicKey
		// PublicPb = publicEvent.PublicPb
		// PublicPbUserId = publicEvent.PublicPbUserId

		zlogger.Infow("kafka init success")
	})

	service.RegisterService(kp)
}

type producerMethod struct {
	producer *Producer
}

func (pe *producerMethod) Public(topic string, st interface{}) {
	if err := pe.producer.SendJson(topic, st); err != nil {
		zlogger.Errorw("kafka send error",
			zap.Error(err),
		)
		return
	}
	return
}

func (pe *producerMethod) PublicUserId(uid int, topic string, st interface{}) {
	if err := pe.producer.SendJsonWithUserId(uid, topic, st); err != nil {
		zlogger.Errorw("kafka send error",
			zap.Error(err),
		)
		return
	}
	return
}

func (pe *producerMethod) PublicPb(topic string, pb proto.Message) {
	if err := pe.producer.SendPb(topic, pb); err != nil {
		zlogger.Errorw("kafka send error",
			zap.Error(err),
		)
		return
	}
	return
}

func (pe *producerMethod) PublicPbUserId(uid int, topic string, pb proto.Message) {
	if err := pe.producer.SendPbWithUserID(uid, topic, pb); err != nil {
		zlogger.Errorw("kafka send error",
			zap.Int("uid", uid),
			zap.Error(err),
		)
		return
	}
	return
}

func (pe *producerMethod) PublicKey(key, topic string, st interface{}) {
	if err := pe.producer.SendJsonWithKey(key, topic, st); err != nil {
		zlogger.Errorw("kafka send error",
			zap.String("key", key),
			zap.Error(err),
		)
		return
	}
	return
}
