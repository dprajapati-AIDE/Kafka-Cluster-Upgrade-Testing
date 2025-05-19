package consumer

import (
	"context"
	"fmt"
	"go_producer_consumer/internal/config"
	"go_producer_consumer/internal/logger"
	"go_producer_consumer/internal/utils"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"
)

func StartConsumerGroup(cluster *config.ClusterConfig, groupID string, topics []string) error {

	saramaConfig := sarama.NewConfig()
	saramaConfig.Version, _ = sarama.ParseKafkaVersion(cluster.Version)
	saramaConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(cluster.Brokers, groupID, saramaConfig)
	if err != nil {
		logger.Error("Failed to create consumer group", zap.String("func", utils.GetFunctionName(1)), zap.Error(err))
		return fmt.Errorf("failed to create consumer group: %w", err)
	}

	defer consumerGroup.Close()

	ctx := context.Background()
	handler := &Consumer{}

	for {
		if err := consumerGroup.Consume(ctx, topics, handler); err != nil {
			logger.Error("Error consuming from kafka", zap.String("func", utils.GetFunctionName(1)), zap.Error(err))
			break
		}
	}

	return nil
}
