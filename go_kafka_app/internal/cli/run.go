package cli

import (
	"context"
	"go_kafka_app/internal/config"
	"go_kafka_app/internal/consumer"
	"go_kafka_app/internal/kafka"
	"go_kafka_app/internal/logger"
	"go_kafka_app/internal/producer"
	"go_kafka_app/internal/utils"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func Run(role string, msgCount int, consumerGroupName string, bus3 bool, messageSizeKB int) error {
	cfg, err := config.LoadConfig("")
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	if err := logger.Initialize(cfg.Logging.Level, cfg.Logging.Encoding, cfg.Logging.Output); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer logger.Sync()

	logger.Info("Application started",
		zap.String("role", role),
		zap.Bool("bus3", bus3),
		zap.Int("message size in KB", messageSizeKB))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var clients []*kafka.Client
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

		logger.Info("Connected to Kafka cluster",
			zap.String("cluster", cluster.Name),
			zap.Strings("brokers", client.GetBrokers()))

		clients = append(clients, client)
	}

	if bus3 {
		for _, cluster := range cfg.Kafka.Clusters {
			if role == "producer" || role == "both" {
				bus3Producer, err := producer.NewBus3Producer(&cluster)
				if err != nil {
					logger.Error("Bus3 producer initialization failed",
						zap.String("func", utils.GetFunctionName(1)), zap.Error(err))
					continue
				}

				err = bus3Producer.ProduceMessage(cfg.Kafka.Topics, msgCount, messageSizeKB)
				if err != nil {
					logger.Error("Bus3 message production failed",
						zap.String("func", utils.GetFunctionName(1)), zap.Error(err))
				}
				bus3Producer.Close()
			}

			if role == "consumer" || role == "both" {
				go func(c config.ClusterConfig) {
					err := consumer.StartBus3ConsumerGroup(&c, consumerGroupName, cfg.Kafka.Topics)
					if err != nil {
						logger.Error("Bus3 consumer group failed",
							zap.String("func", utils.GetFunctionName(1)),
							zap.String("cluster", c.Name),
							zap.Error(err))
					}
				}(cluster)
			}
		}

		<-ctx.Done()
		logger.Info("Shutdown signal received. Exiting Bus3 mode.")

		// Close Kafka clients
		for _, client := range clients {
			if err := client.Close(); err != nil {
				logger.Error("Error closing Kafka client",
					zap.String("func", utils.GetFunctionName(1)),
					zap.String("cluster", client.GetConfig().Name),
					zap.Error(err))
			}
		}

		time.Sleep(500 * time.Millisecond)
		return nil
	}

	// Standard mode (non-bus3)
	if role == "producer" || role == "both" {
		for _, cluster := range cfg.Kafka.Clusters {
			prod, err := producer.NewProducer(&cluster)
			if err != nil {
				logger.Error("Producer initialization failed", zap.String("func", utils.GetFunctionName(1)), zap.Error(err))
				continue
			}

			err = prod.ProduceMessage(cfg.Kafka.Topics, msgCount, messageSizeKB)
			if err != nil {
				logger.Error("Failed to send messages", zap.String("func", utils.GetFunctionName(1)), zap.Error(err))
			}
			prod.Close()
		}
	}

	if role == "consumer" || role == "both" {
		for _, cluster := range cfg.Kafka.Clusters {
			go func(c config.ClusterConfig) {
				err := consumer.StartConsumerGroup(&c, consumerGroupName, cfg.Kafka.Topics)
				if err != nil {
					logger.Error("Consumer group failed", zap.String("cluster", c.Name), zap.Error(err))
				}
			}(cluster)
		}
	}

	<-ctx.Done()
	logger.Info("Shutdown signal received. Closing Kafka clients...")

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
