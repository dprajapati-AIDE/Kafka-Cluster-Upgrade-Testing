package producer

import (
	"encoding/json"
	"fmt"
	"go_producer_consumer/internal/config"
	"go_producer_consumer/internal/logger"
	"go_producer_consumer/internal/utils"
	"time"

	"github.com/Shopify/sarama"
	"github.com/go-faker/faker/v4"
	"go.uber.org/zap"
)

type Producer struct {
	producer sarama.SyncProducer
	cluster  *config.ClusterConfig
}

func NewProducer(cluster *config.ClusterConfig) (*Producer, error) {

	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Retry.Max = 5
	saramaConfig.Version, _ = sarama.ParseKafkaVersion(cluster.Version)

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

func (p *Producer) ProduceMessage(topic string, count int) error {

	for i := 0; i < count; i++ {

		msg := DeviceMessageModel{
			Timestamp:  time.Now().Format(time.RFC3339),
			DeviceID:   faker.UUIDDigit(),
			DeviceIP:   faker.IPv4(),
			DeviceType: "firewall",
			Vendor:     "Juniper",
			Message:    faker.Sentence(),
		}

		jsonBytes, err := json.Marshal(msg)
		if err != nil {
			logger.Error("Failed to marshal message", zap.String("func", utils.GetFunctionName(1)), zap.Error(err))
			continue
		}

		kafkaMessage := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.ByteEncoder(jsonBytes),
		}

		partition, offset, err := p.producer.SendMessage(kafkaMessage)
		if err != nil {
			logger.Error("Failed to send message", zap.String("func", utils.GetFunctionName(1)), zap.Error(err))
			continue
		}

		logger.Info("Message sent",
			zap.String("topic", topic),
			zap.Int32("partition", partition),
			zap.Int64("offset", offset))
	}

	return nil
}

func (p *Producer) Close() error {
	return p.producer.Close()
}
