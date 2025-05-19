package consumer

import (
	"encoding/json"
	"go_producer_consumer/internal/logger"
	"go_producer_consumer/internal/utils"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"
)

type Consumer struct{}

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
			logger.Error("Failed to unmarshel message", zap.String("func", utils.GetFunctionName(1)), zap.Error(err))
		} else {
			logger.Info("Consumed Message", zap.Any("message", payload))
		}

		session.MarkMessage(msg, "")
	}

	return nil
}
