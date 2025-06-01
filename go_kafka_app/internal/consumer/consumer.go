package consumer

import (
	"encoding/json"
	"go_kafka_app/internal/config"
	"go_kafka_app/internal/logger"
	"go_kafka_app/internal/utils"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"
)

type Consumer struct {
	cluster *config.ClusterConfig
}

func (c *Consumer) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for msg := range claim.Messages() {
		var payload map[string]interface{}
		if err := json.Unmarshal(msg.Value, &payload); err != nil {
			logger.Error("Failed to unmarshal message",
				zap.String("func", utils.GetFunctionName(1)),
				zap.String("cluster", c.cluster.Name),
				zap.String("topic", msg.Topic),
				zap.Error(err))

		} else {
			logger.Info("Consumed Message",
				zap.String("cluster", c.cluster.Name),
				zap.String("topic", msg.Topic),
				zap.Int32("partition", msg.Partition),
				zap.Int64("offset", msg.Offset))

		}

		session.MarkMessage(msg, "")
	}

	return nil
}
