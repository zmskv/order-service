package kafka

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type Consumer struct {
	client  sarama.ConsumerGroup
	topic   string
	groupID string
	ctx     context.Context
	cancel  context.CancelFunc
	Handler func(data []byte) error
	logger  *zap.Logger
}

func NewConsumer(brokers []string, topic string, groupID string, handler func(data []byte) error, log *zap.Logger) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	client, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Consumer{
		client:  client,
		topic:   topic,
		groupID: groupID,
		ctx:     ctx,
		cancel:  cancel,
		Handler: handler,
		logger:  log.Named("kafka_consumer"),
	}, nil
}

func (c *Consumer) Start() error {
	c.logger.Info("starting Kafka consumer", zap.String("topic", c.topic), zap.String("group", c.groupID))

	go func() {
		handler := &consumerGroupHandler{
			handler: c.Handler,
			logger:  c.logger,
		}

		for {
			if err := c.client.Consume(c.ctx, []string{c.topic}, handler); err != nil {
				c.logger.Error("error from consumer", zap.Error(err))
			}

			if c.ctx.Err() != nil {
				c.logger.Info("consumer context cancelled, stopping consume loop")
				return
			}
		}
	}()

	return nil
}

func (c *Consumer) Stop() error {
	c.logger.Info("stopping Kafka consumer")

	c.cancel()

	if err := c.client.Close(); err != nil {
		c.logger.Error("failed to close Kafka consumer", zap.Error(err))
		return fmt.Errorf("failed to close Kafka consumer: %w", err)
	}

	c.logger.Info("Kafka consumer stopped successfully")
	return nil
}

type consumerGroupHandler struct {
	handler func(data []byte) error
	logger  *zap.Logger
}

func (h *consumerGroupHandler) Setup(sess sarama.ConsumerGroupSession) error {
	h.logger.Info("consumer group setup",
		zap.String("member_id", sess.MemberID()),
		zap.String("generation_id", fmt.Sprintf("%d", sess.GenerationID())),
	)
	return nil
}

func (h *consumerGroupHandler) Cleanup(sess sarama.ConsumerGroupSession) error {
	h.logger.Info("consumer group cleanup",
		zap.String("member_id", sess.MemberID()),
		zap.String("generation_id", fmt.Sprintf("%d", sess.GenerationID())),
	)
	return nil
}

func (h *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.logger.Debug("received message",
			zap.String("topic", msg.Topic),
			zap.Int32("partition", msg.Partition),
			zap.Int64("offset", msg.Offset),
			zap.ByteString("key", msg.Key),
		)

		if err := h.handler(msg.Value); err != nil {
			h.logger.Error("error processing message", zap.Error(err), zap.Int64("offset", msg.Offset))
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}
