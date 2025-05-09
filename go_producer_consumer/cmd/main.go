package main

import (
	"context"
	"go_producer_consumer/internal/config"
	"go_producer_consumer/internal/kafka"
	"go_producer_consumer/internal/logger"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
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

	logger.Info("Application started")

	// Create a context to handle graceful shutdown
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
				zap.String("cluster", cluster.Name),
				zap.Error(err))
			continue
		}

		logger.Info("Successfully connected to Kafka cluster",
			zap.String("cluster", cluster.Name),
			zap.Strings("brokers", client.GetBrokers()))

		clients = append(clients, client)
	}

	// Wait for interrupt signal (Ctrl+C)
	<-ctx.Done()
	logger.Info("Shutdown signal received. Closing Kafka clients...")

	// Gracefully close all Kafka clients
	for _, client := range clients {
		if err := client.Close(); err != nil {
			logger.Warn("Error closing Kafka client",
				zap.String("cluster", client.GetConfig().Name),
				zap.Error(err))
		}
	}

	logger.Info("Shutdown complete")
	time.Sleep(500 * time.Millisecond)
}
