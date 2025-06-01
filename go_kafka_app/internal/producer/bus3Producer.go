package producer

import (
	"fmt"
	"go_kafka_app/internal/config"
	"go_kafka_app/internal/logger"
	"go_kafka_app/internal/utils"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/mistsys/mist_go_utils/bus3/producer3"
	"github.com/mistsys/mist_go_utils/cloud"
	"go.uber.org/zap"
)

type Bus3Producer struct {
	producer producer3.Producer
	cluster  *config.ClusterConfig
}

func NewBus3Producer(cluster *config.ClusterConfig) (*Bus3Producer, error) {
	cloud.KAFKA_BROKERS = cluster.Brokers
	bus3SaramaConfig := sarama.NewConfig()
	bus3SaramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	bus3SaramaConfig.Producer.Retry.Max = 5
	bus3SaramaConfig.Net.DialTimeout = 30 * time.Second
	bus3SaramaConfig.Net.ReadTimeout = 30 * time.Second
	bus3SaramaConfig.Net.WriteTimeout = 30 * time.Second
	bus3SaramaConfig.Metadata.RefreshFrequency = 2 * time.Minute
	bus3SaramaConfig.Metadata.Retry.Backoff = 2 * time.Second
	bus3SaramaConfig.Producer.MaxMessageBytes = 20_000_000
	bus3SaramaConfig.Producer.RequiredAcks = sarama.RequiredAcks(1)
	bus3SaramaConfig.Producer.Timeout = 10 * time.Second
	bus3SaramaConfig.Producer.Return.Successes = true
	bus3SaramaConfig.Producer.Flush.Frequency = 5 * time.Second
	bus3SaramaConfig.Producer.Flush.MaxMessages = 0
	bus3SaramaConfig.Version, _ = sarama.ParseKafkaVersion(cluster.Version)

	var err error
	bus3SaramaConfig.Version, err = sarama.ParseKafkaVersion(cluster.Version)
	if err != nil {
		logger.Error("Failed to parse Kafka version", zap.String("func", utils.GetFunctionName(1)), zap.Error(err))
		return nil, fmt.Errorf("invalid Kafka version: %w", err)
	}

	bus3Producer, err := producer3.NewProducer(bus3SaramaConfig)
	if err != nil {
		logger.Error("Failed to created producer", zap.String("func", utils.GetFunctionName(1)), zap.Error(err))
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	return &Bus3Producer{
		producer: bus3Producer,
		cluster:  cluster,
	}, nil
}

func (bus3Producer *Bus3Producer) ProduceMessage(topics []string, count int, messageSizeKB int) error {
	var wg sync.WaitGroup

	go func() {
		for err := range bus3Producer.producer.Errors() {
			logger.Error("Producer error",
				zap.String("func", utils.GetFunctionName(1)),
				zap.Error(err.Err),
				zap.Any("msg", err.Msg))
		}
	}()

	go func() {
		for msg := range bus3Producer.producer.Successes() {
			logger.Info("Message delivered",
				zap.String("cluster", bus3Producer.cluster.Name),
				zap.String("topic", msg.Topic),
				zap.Int32("partition", msg.Partition),
				zap.Int64("offset", msg.Offset))
			wg.Done()
		}
	}()

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

			wg.Add(1)
			bus3Producer.producer.Input() <- kafkaMessage

			logger.Info("Message queued for sending",
				zap.String("cluster", bus3Producer.cluster.Name),
				zap.String("topic", topicName),
				zap.Int("message length in bytes", len(jsonBytes)))
		}
	}

	// Wait until all messages are confirmed
	wg.Wait()
	return nil
}

func (bus3Producer *Bus3Producer) Close() error {
	return bus3Producer.producer.Close()
}
