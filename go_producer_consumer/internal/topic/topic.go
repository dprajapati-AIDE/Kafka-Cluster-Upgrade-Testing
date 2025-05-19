package topic

import (
	"go_producer_consumer/internal/kafka"
	"go_producer_consumer/internal/logger"
	"go_producer_consumer/internal/utils"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"
)

const (
	DefaultPartitions        = 3
	DefaultReplicationFactor = 2
)

func EnsureTopics(client *kafka.Client) error {
	clusterConfig := client.GetConfig()

	existingTopics, err := client.ListTopics()
	if err != nil {
		logger.Error("Failed to list kafka topics", zap.String("func", utils.GetFunctionName(1)), zap.Error(err))
		return err
	}

	for _, topicConfig := range clusterConfig.Topics {

		topicName := topicConfig.Name

		// check topic existance
		if _, exists := existingTopics[topicConfig.Name]; exists {
			logger.Info("Topic already exists", zap.String("topic", topicConfig.Name))
			continue
		}

		// create topic
		topicDetail := &sarama.TopicDetail{
			NumPartitions:     DefaultPartitions,
			ReplicationFactor: DefaultReplicationFactor,
			ConfigEntries:     map[string]*string{},
		}

		err := client.CreateTopic(topicName, topicDetail)
		if err != nil {
			logger.Error("Failed to create topic",
				zap.String("func", utils.GetFunctionName(1)),
				zap.String("cluster", clusterConfig.Name),
				zap.String("topic", topicName),
				zap.Error(err))
			continue
		}

		logger.Info("Created topic successfully",
			zap.String("cluster", clusterConfig.Name),
			zap.String("topic", topicName),
			zap.Int("partitions", DefaultPartitions),
			zap.Int("replication_factor", DefaultReplicationFactor))
	}

	return nil
}
