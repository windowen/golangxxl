package rocketmq

import (
	"encoding/json"
	"errors"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"context"
	"sync/atomic"

	"queueJob/pkg/common/config"
	"queueJob/pkg/zlogger"
)

// errProducerClosed 定义一个全局错误，当生产者关闭时返回该错误
var errProducerClosed = errors.New("rocketmq producer is closed")

type Producer struct {
	producer rocketmq.Producer
	started  atomic.Bool
	ctx      context.Context
	cancel   context.CancelFunc
	update   func(producer *Producer)
}

// newProducer 返回一个新的 Producer 实例
func newProducer(update func(producer *Producer)) *Producer {
	rp := &Producer{
		update: update,
	}

	return rp
}

func (p *Producer) Start() {
	if len(config.Config.RocketMQ.RocketMQAddr) == 0 {
		// 未配置, 打印日志后忽略
		zlogger.Infow("rocketmq producer config empty")
		return
	}

	if p.started.Load() {
		zlogger.Info("producer already started")
		return
	}

	p.ctx, p.cancel = context.WithCancel(context.Background())

	prod, err := rocketmq.NewProducer(
		producer.WithNameServer(config.Config.RocketMQ.RocketMQAddr),
		producer.WithGroupName(ProducerGroupName),
		producer.WithRetry(2),
	)
	if err != nil {
		zlogger.Errorw("rocketmq producer failed", zap.Error(err))
		panic(err)
	}
	p.producer = prod

	if err = p.producer.Start(); err != nil {
		panic(err)
	}

	p.started.Store(true)
	zlogger.Infow("rocketMQ producer started", zap.Strings("nameServer", config.Config.RocketMQ.RocketMQAddr), zap.String("groupName", ProducerGroupName))

	p.update(p)
	return
}

func (p *Producer) Send(topic string, msg []byte) error {
	if !p.started.Load() {
		return errProducerClosed
	}

	message := &primitive.Message{
		Topic: topic,
		Body:  msg,
	}

	_, err := p.producer.SendSync(p.ctx, message)
	if err != nil {
		zlogger.Errorw("failed to send message", zap.Error(err))
	}
	return err
}

// SendAsync 异步发送消息
func (p *Producer) SendAsync(topic string, msg []byte) error {
	if !p.started.Load() {
		return errProducerClosed
	}

	message := &primitive.Message{
		Topic: topic,
		Body:  msg,
	}

	// 定义回调函数
	callback := func(ctx context.Context, result *primitive.SendResult, err error) {
		if err != nil {
			// 处理发送错误
			zlogger.Errorf("SendAsync | message=%v | Err=%v", message.String(), err.Error())
		}
	}

	err := p.producer.SendAsync(p.ctx, callback, message)
	if err != nil {
		zlogger.Errorf("SendAsync failed to send message | topic=%v,message=%v | Err=%v", topic, message.String(), zap.Error(err))
	}
	return err
}

// SendWithDelay 发送延时消息
// delayLevel 定义: 1s 5s 10s 30s 1m 2m 3m 4m 5m 6m 7m 8m 9m 10m 20m 30m 1h 2h
// 从1开始. 如果 level=1, 延时时间是 1s.
func (p *Producer) SendWithDelay(topic string, msg []byte, delayLevel int) error {
	if !p.started.Load() {
		return errProducerClosed
	}

	message := &primitive.Message{
		Topic: topic,
		Body:  msg,
	}
	// 设置消息的延时级别
	message.WithDelayTimeLevel(delayLevel)

	_, err := p.producer.SendSync(p.ctx, message)
	if err != nil {
		zlogger.Errorf("SendWithDelay failed | topic=%v,message=%v | Err=%v", topic, message.String(), zap.Error(err))
	}
	return err
}

// SendAsyncWithDelay 异步发送延时消息
func (p *Producer) SendAsyncWithDelay(topic string, msg []byte, delayLevel int) error {
	if !p.started.Load() {
		return errProducerClosed
	}

	message := &primitive.Message{
		Topic: topic,
		Body:  msg,
	}
	// 设置消息的延时级别
	message.WithDelayTimeLevel(delayLevel)

	// 定义回调函数
	callback := func(ctx context.Context, result *primitive.SendResult, err error) {
		if err != nil {
			// 处理发送错误
			zlogger.Errorf("SendAsyncWithDelay callback | message=%s | Err=%v", message.String(), err.Error())
		}
	}

	err := p.producer.SendAsync(p.ctx, callback, message)
	if err != nil {
		zlogger.Errorw("failed to send message", zap.Error(err))
	}
	return err
}

// SendPb 生产消息(protobuf)
func (p *Producer) SendPb(topic string, pb proto.Message) (err error) {
	data, err := proto.Marshal(pb)
	if err != nil {
		zlogger.Errorw("rocketmq producer marshal pb failed.", zap.String("topic", topic), zap.Int("len", len(data)),
			zap.Error(err))
		return err
	}
	return p.Send(topic, data)
}

// AsyncSendPb 生产消息(protobuf)
func (p *Producer) AsyncSendPb(topic string, pb proto.Message) (err error) {
	data, err := proto.Marshal(pb)
	if err != nil {
		zlogger.Errorw("AsyncSendPb producer marshal pb failed.", zap.String("topic", topic), zap.Int("len", len(data)),
			zap.Error(err))
		return err
	}
	return p.SendAsync(topic, data)
}

// AsyncSendJson 生产消息(json)
func (p *Producer) AsyncSendJson(topic string, st interface{}) (err error) {
	data, err := json.Marshal(st)
	if err != nil {
		zlogger.Errorw("rocketmq producer marshal json failed.", zap.String("topic", topic), zap.Int("len", len(data)),
			zap.Error(err))
		return err
	}
	return p.SendAsync(topic, data)
}

// SendPbWithDelay 生产延时消息(protobuf)
func (p *Producer) SendPbWithDelay(topic string, pb proto.Message, delayLevel int) (err error) {
	data, err := proto.Marshal(pb)
	if err != nil {
		zlogger.Errorw("rocketmq producer marshal pb failed.", zap.String("topic", topic), zap.Int("len", len(data)),
			zap.Error(err))
		return err
	}
	return p.SendWithDelay(topic, data, delayLevel)
}

// SendAsyncPbWithDelay 生产异步延时消息(protobuf)
func (p *Producer) SendAsyncPbWithDelay(topic string, pb proto.Message, delayLevel int) (err error) {
	data, err := proto.Marshal(pb)
	if err != nil {
		zlogger.Errorw("SendAsyncPbWithDelay async producer marshal pb failed.", zap.String("topic", topic), zap.Int("len", len(data)),
			zap.Error(err))
		return err
	}
	return p.SendAsyncWithDelay(topic, data, delayLevel)
}

// SendAsyncPbWithDelayJson 生产异步延时消息(json)
func (p *Producer) SendAsyncPbWithDelayJson(topic string, st interface{}, delayLevel int) (err error) {
	data, err := json.Marshal(st)
	if err != nil {
		zlogger.Errorw("SendAsyncPbWithDelayJson async producer marshal pb failed.", zap.String("topic", topic), zap.Int("len", len(data)),
			zap.Error(err))
		return err
	}
	return p.SendAsyncWithDelay(topic, data, delayLevel)
}

func (p *Producer) Stop() {
	if p.started.CompareAndSwap(true, false) {
		if err := p.producer.Shutdown(); err != nil {
			zlogger.Errorw("failed to shut down producer", zap.Error(err))
		}
		p.cancel()
	}
}
