package producer

import (
	"fmt"
	"go_kafka_app/internal/config"
	"go_kafka_app/internal/logger"
	"go_kafka_app/internal/utils"
	"time"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"
)

type Producer struct {
	producer sarama.SyncProducer
	cluster  *config.ClusterConfig
}

func NewProducer(cluster *config.ClusterConfig) (*Producer, error) {

	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Retry.Max = 5
	saramaConfig.Net.DialTimeout = 30 * time.Second
	saramaConfig.Net.ReadTimeout = 30 * time.Second
	saramaConfig.Net.WriteTimeout = 30 * time.Second
	saramaConfig.Metadata.RefreshFrequency = 2 * time.Minute
	saramaConfig.Metadata.Retry.Backoff = 2 * time.Second
	saramaConfig.Producer.MaxMessageBytes = 20_000_000
	saramaConfig.Producer.RequiredAcks = sarama.RequiredAcks(1)
	saramaConfig.Producer.Timeout = 10 * time.Second
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Flush.Frequency = 5 * time.Second
	saramaConfig.Producer.Flush.MaxMessages = 0
	saramaConfig.Version, _ = sarama.ParseKafkaVersion(cluster.Version)

	var err error
	saramaConfig.Version, err = sarama.ParseKafkaVersion(cluster.Version)
	if err != nil {
		logger.Error("Failed to parse Kafka version", zap.String("func", utils.GetFunctionName(1)), zap.Error(err))
		return nil, fmt.Errorf("invalid Kafka version: %w", err)
	}

	producer, err := sarama.NewSyncProducer(cluster.Brokers, saramaConfig)
	if err != nil {
		logger.Error("Failed to created producer", zap.String("func", utils.GetFunctionName(1)), zap.Error(err))
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	return &Producer{
		producer: producer,
		cluster:  cluster,
	}, nil
}

func (p *Producer) ProduceMessage(topics []string, count int, messageSizeKB int) error {

	for _, topicName := range topics {

		for i := 0; i < count; i++ {

			jsonBytes, err := utils.GenerateMessage(messageSizeKB)
			if err != nil {
				logger.Error("Failed to generate message",
					zap.String("func", utils.GetFunctionName(1)),
					zap.Error(err))
				continue
			}

			kafkaMessage := &sarama.ProducerMessage{
				Topic: topicName,
				Value: sarama.ByteEncoder(jsonBytes),
			}

			partition, offset, err := p.producer.SendMessage(kafkaMessage)
			if err != nil {
				logger.Error("Failed to send message", zap.String("func", utils.GetFunctionName(1)), zap.Error(err))
				continue
			}

			logger.Info("Message sent",
				zap.String("cluster", p.cluster.Name),
				zap.String("topic", topicName),
				zap.Int32("partition", partition),
				zap.Int64("offset", offset),
				zap.Int("message length in bytes", len(jsonBytes)))
		}
	}

	return nil
}

func (p *Producer) Close() error {
	return p.producer.Close()
}
