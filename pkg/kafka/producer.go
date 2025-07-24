package kafka

import (
	"context"
	"encoding/json"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"queueJob/pkg/common/config"
	"queueJob/pkg/safego"
	"queueJob/pkg/zlogger"
)

// Producer kafka生产者
type Producer struct {
	ap     sarama.AsyncProducer
	start  atomic.Bool
	ctx    context.Context
	cancel func()
	update func(producer *Producer)
}

func newProducer(update func(producer *Producer)) *Producer {
	kp := &Producer{
		update: update,
	}

	return kp
}

func (kp *Producer) Start() {
	if len(config.Config.Kafka.KafkaAddr) == 0 {
		// 未配置, 打印日志后忽略
		zlogger.Errorw("kafka producer saramaConfig empty")
		return
	}

	if kp.start.Load() {
		zlogger.Errorw("kafka producer is start. can't modify saramaConfig.(need restart) ",
			zap.Strings("broker", config.Config.Kafka.KafkaAddr))
		return
	}

	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.RequiredAcks = sarama.WaitForLocal
	saramaConfig.Producer.Partitioner = sarama.NewHashPartitioner
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Return.Errors = true

	producer, err := sarama.NewAsyncProducer(config.Config.Kafka.KafkaAddr, saramaConfig)
	if err != nil {
		zlogger.Errorw("kafka sarama.NewAsyncProducer error", zap.Error(err))
		panic(err)
	}

	kp.ap = producer
	kp.ctx, kp.cancel = context.WithCancel(context.Background())

	safego.Go("kafka.handleError", kp.handleError)

	if saramaConfig.Producer.Return.Successes {
		safego.Go("kafka.handleSuccess", kp.handleSuccess)
	}

	kp.start.Store(true)
	zlogger.Infow("kafka producer open success",
		zap.Strings("brokers", config.Config.Kafka.KafkaAddr))

	kp.update(kp)
	return
}

func (kp *Producer) handleSuccess() {
	var (
		pm *sarama.ProducerMessage
	)
	for {
		select {
		case <-kp.ctx.Done():
			return
		case pm = <-kp.ap.Successes():
			if pm != nil && config.Config.App.Env == "dev" {
				encode, err := pm.Value.Encode()
				if err != nil {
					continue
				}
				zlogger.Debugw("kafka send success",
					zap.Int32("Partition", pm.Partition),
					zap.Int64("Offset", pm.Offset),
					zap.Any("Key", pm.Key),
					zap.String("Value", string(encode)),
				)
			}
		}
	}
}

// 处理kafka错误
func (kp *Producer) handleError() {
	var (
		err *sarama.ProducerError
	)
	for {
		select {
		case <-kp.ctx.Done():
			return
		case err = <-kp.ap.Errors():
			if err != nil {
				encode, ee := err.Msg.Value.Encode()
				if ee != nil {
					continue
				}
				zlogger.Errorw("kafka process failed.",
					zap.Int32("Partition", err.Msg.Partition),
					zap.Int64("Offset", err.Msg.Offset),
					zap.Any("Key", err.Msg.Key),
					zap.String("Value", string(encode)),
					zap.Error(err),
				)
			}
		}
	}
}

func (kp *Producer) Stop() {
	if kp.start.CompareAndSwap(true, false) {
		err := kp.ap.Close()
		if err != nil {
			return
		}
		kp.cancel()
		zlogger.Infow("kafka producer is stop.")
	}
}

// Send 生产消息
func (kp *Producer) Send(topic string, msg []byte) (err error) {
	nowTime := time.Now().Unix()
	proMsg := sarama.ProducerMessage{
		Topic: topic,
		// Key:   sarama.StringEncoder(strconv.Itoa(int(nowTime))),
		Value: sarama.ByteEncoder(msg),
	}

	zlogger.Debugw("kafka send begin", zap.String("topic", topic), zap.Int64("key", nowTime), zap.String("value", string(msg)))
	kp.ap.Input() <- &proMsg
	return
}

// SendWithKey 生产消息
func (kp *Producer) SendWithKey(key, topic string, msg []byte) (err error) {
	proMsg := sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(msg),
	}
	kp.ap.Input() <- &proMsg
	return
}

// SendWithUserID 生产消息
func (kp *Producer) SendWithUserID(uid int, topic string, msg []byte) (err error) {
	proMsg := sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(strconv.Itoa(uid)),
		Value: sarama.ByteEncoder(msg),
	}
	kp.ap.Input() <- &proMsg
	return
}

// SendPb 生产消息(protobuf)
func (kp *Producer) SendPb(topic string, pb proto.Message) (err error) {
	data, err := proto.Marshal(pb)
	if err != nil {
		zlogger.Errorw("kafka producer marshal pb failed.", zap.String("topic", topic), zap.Int("len", len(data)),
			zap.Error(err))
		return err
	}
	return kp.Send(topic, data)
}

// SendPbWithKey 生产消息(protobuf)
func (kp *Producer) SendPbWithKey(key, topic string, pb proto.Message) (err error) {
	data, err := proto.Marshal(pb)
	if err != nil {
		zlogger.Errorw("kafka producer marshal pb failed.", zap.String("topic", topic), zap.Int("len", len(data)),
			zap.String("key", key), zap.Error(err))
		return err
	}
	return kp.SendWithKey(key, topic, data)
}

// SendPbWithUserID 生产消息(protobuf)
func (kp *Producer) SendPbWithUserID(uid int, topic string, pb proto.Message) (err error) {
	data, err := proto.Marshal(pb)
	if err != nil {
		zlogger.Errorw("kafka producer marshal pb failed.", zap.String("topic", topic), zap.Int("len", len(data)),
			zap.Int("uid", uid), zap.Error(err))
		return err
	}
	return kp.SendWithUserID(uid, topic, data)
}

// SendJson 生产消息(json)
func (kp *Producer) SendJson(topic string, st interface{}) (err error) {
	data, err := json.Marshal(st)
	if err != nil {
		zlogger.Errorw("kafka producer marshal json failed.", zap.String("topic", topic), zap.Int("len", len(data)),
			zap.Error(err))
		return err
	}
	return kp.Send(topic, data)
}

// SendJsonWithUserId 生产消息(json)
func (kp *Producer) SendJsonWithUserId(uid int, topic string, st interface{}) (err error) {
	data, err := json.Marshal(st)
	if err != nil {
		zlogger.Errorw("kafka producer marshal json failed.", zap.String("topic", topic), zap.Int("len", len(data)),
			zap.Int("key", uid), zap.Error(err))
		return err
	}
	return kp.SendWithUserID(uid, topic, data)
}

// SendJsonWithKey 生产消息with key(json)
func (kp *Producer) SendJsonWithKey(key, topic string, st interface{}) (err error) {
	data, err := json.Marshal(st)
	if err != nil {
		zlogger.Errorw("kafka producer marshal json failed.", zap.String("topic", topic), zap.Int("len", len(data)),
			zap.Error(err))
		return err
	}
	return kp.SendWithKey(key, topic, data)
}
