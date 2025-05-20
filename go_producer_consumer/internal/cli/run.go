package cli

import (
	"context"
	"go_producer_consumer/internal/config"
	"go_producer_consumer/internal/consumer"
	"go_producer_consumer/internal/kafka"
	"go_producer_consumer/internal/logger"
	"go_producer_consumer/internal/producer"
	"go_producer_consumer/internal/topic"
	"go_producer_consumer/internal/utils"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func Run(role string, msgCount int, consumerGroupName string) error {
	// Load config
	cfg, err := config.LoadConfig("")
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	// Initialize logger
	if err := logger.Initialize(cfg.Logging.Level, cfg.Logging.Encoding, cfg.Logging.Output); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer logger.Sync()

	logger.Info("Application started", zap.String("role", role))

	// context to handle graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var clients []*kafka.Client

	// Connect to all Kafka clusters and store the clients
	for _, cluster := range cfg.Kafka.Clusters {
		logger.Info("Connecting to Kafka cluster",
			zap.String("name", cluster.Name),
			zap.String("version", cluster.Version))

		client, err := kafka.NewClient(&cluster)
		if err != nil {
			logger.Error("Failed to connect to Kafka cluster",
				zap.String("func", utils.GetFunctionName(1)),
				zap.String("cluster", cluster.Name),
				zap.Error(err))
			continue
		}

		logger.Info("Successfully connected to Kafka cluster",
			zap.String("cluster", cluster.Name),
			zap.Strings("brokers", client.GetBrokers()))

		// Ensure required topics are present
		err = topic.EnsureTopics(client)
		if err != nil {
			logger.Warn("Failed to ensure topics", zap.String("func", utils.GetFunctionName(1)), zap.String("cluster", cluster.Name), zap.Error(err))
		}

		clients = append(clients, client)
	}

	// launch producer or both role
	if role == "producer" || role == "both" {
		for _, cluster := range cfg.Kafka.Clusters {
			producer, err := producer.NewProducer(&cluster)
			if err != nil {
				logger.Error("Producer initialization failed", zap.String("func", utils.GetFunctionName(1)), zap.Error(err))
				continue
			}

			err = producer.ProduceMessage(cfg.Devices, msgCount)
			if err != nil {
				logger.Error("Failed to send messages", zap.String("func", utils.GetFunctionName(1)), zap.Error(err))
			}
			producer.Close()
		}
	}

	// launch consumer or both role
	if role == "consumer" || role == "both" {
		for _, cluster := range cfg.Kafka.Clusters {
			var topics []string
			for _, topic := range cluster.Topics {
				topics = append(topics, topic.Name)
			}
			go func(c config.ClusterConfig) {
				err := consumer.StartConsumerGroup(&c, consumerGroupName, topics)
				if err != nil {
					logger.Error("Consumer group failed", zap.String("cluster", c.Name), zap.Error(err))
				}
			}(cluster)
		}
	}

	// Wait for interrupt signal (Ctrl+C)
	<-ctx.Done()
	logger.Info("Shutdown signal received. Closing Kafka clients...")

	// Gracefully close all Kafka clients
	for _, client := range clients {
		if err := client.Close(); err != nil {
			logger.Error("Error closing Kafka client",
				zap.String("func", utils.GetFunctionName(1)),
				zap.String("cluster", client.GetConfig().Name),
				zap.Error(err))
		}
	}

	logger.Info("Shutdown complete")
	time.Sleep(500 * time.Millisecond)

	return nil

}
