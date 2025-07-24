package kafka

import (
	"context"
	"fmt"
	"sync"

	"github.com/IBM/sarama"
	"go.uber.org/zap"

	"queueJob/pkg/common/config"
	"queueJob/pkg/safego"
	"queueJob/pkg/zlogger"
)

type Consumer struct {
	ready  chan bool
	ctx    context.Context
	cancel func()
	wg     sync.WaitGroup
	group  sarama.ConsumerGroup
	topic  string
	handle func(msg []byte)
}

func NewConsumer(topic string, hl func(msg []byte)) *Consumer {
	return &Consumer{
		ready:  make(chan bool),
		topic:  topic,
		handle: hl,
	}
}

func (c *Consumer) Start() {
	if len(config.Config.Kafka.KafkaAddr) == 0 {
		zlogger.Errorw("kafka consumer config empty", zap.String("prefix", "plat"))
		return
	}

	consumerConfig := sarama.NewConfig()
	consumerConfig.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
	consumerConfig.Consumer.Offsets.Initial = sarama.OffsetNewest

	group, err := sarama.NewConsumerGroup(config.Config.Kafka.KafkaAddr, fmt.Sprintf("group_%s", c.topic), consumerConfig)
	if err != nil {
		zlogger.Errorw("kafka consumer error", zap.String("topic", c.topic), zap.Error(err))
		panic(err)
	}

	c.group = group
	c.ctx, c.cancel = context.WithCancel(context.Background())

	c.wg.Add(1)
	safego.Go("kafka.consumer", func() {
		defer c.wg.Done()
		for {
			if innerErr := group.Consume(c.ctx, []string{c.topic}, c); innerErr != nil {
				zlogger.Errorw("error from consumer", zap.String("topic", c.topic), zap.Error(innerErr))
			}
			// Check if context was cancelled, signaling that the consumer should stop
			if c.ctx.Err() != nil {
				return
			}
			// Recreate ready channel for the next consumption cycle
			c.ready = make(chan bool)
		}
	})

	<-c.ready // wait till the consumer has been set up
	zlogger.Infow("kafka consumer group up and running.", zap.String("topic", c.topic))
	return
}

func (c *Consumer) Stop() {
	c.cancel()
	c.wg.Wait()
	if err := c.group.Close(); err != nil {
		zlogger.Errorw("error closing consumer group", zap.String("topic", c.topic), zap.Error(err))
		return
	}

	zlogger.Infow("kafka consumer stopped",
		zap.Strings("nameServer", config.Config.Kafka.KafkaAddr),
		zap.String("topic", c.topic),
	)
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *Consumer) Setup(_ sarama.ConsumerGroupSession) error {
	close(c.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *Consumer) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (c *Consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		// zlogger.Debugw("Message claimed",
		// 	zap.String("topic", message.Topic),
		// 	// zap.String("value", string(message.Value)),
		// 	zap.String("timestamp", message.Timestamp.String()),
		// 	zap.Int32("partition", message.Partition),
		// 	zap.Int64("offset", message.Offset),
		// )
		c.handle(message.Value)
		sess.MarkMessage(message, "")
	}
	return nil
}
