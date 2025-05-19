package kafka

import (
	"fmt"
	"go_producer_consumer/internal/config"
	"go_producer_consumer/internal/logger"
	"strings"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"
)

type Client struct {
	config      *config.ClusterConfig
	adminClient sarama.ClusterAdmin
}

func NewClient(clusterConfig *config.ClusterConfig) (*Client, error) {

	saramaConfig := sarama.NewConfig()

	// setting up kafka client with id
	saramaConfig.ClientID = fmt.Sprintf("kafka-network-monitor-%s", clusterConfig.Name)

	// Configure version
	version, err := sarama.ParseKafkaVersion(clusterConfig.Version)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Kafka version: %w", err)
	}
	saramaConfig.Version = version

	// Create admin client
	adminClient, err := sarama.NewClusterAdmin(clusterConfig.Brokers, saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka admin client: %w", err)
	}

	logger.Info("Kafka client created",
		zap.String("cluster", clusterConfig.Name),
		zap.String("version", clusterConfig.Version),
		zap.String("brokers", strings.Join(clusterConfig.Brokers, ",")))

	return &Client{
		config:      clusterConfig,
		adminClient: adminClient,
	}, nil
}

// Close closes the Kafka client
func (c *Client) Close() error {
	if err := c.adminClient.Close(); err != nil {
		return fmt.Errorf("failed to close Kafka admin client: %w", err)
	}
	logger.Info("Kafka client closed", zap.String("cluster", c.config.Name))
	return nil
}

// GetBrokers returns the brokers for this client
func (c *Client) GetBrokers() []string {
	return c.config.Brokers
}

// GetConfig returns the cluster configuration
func (c *Client) GetConfig() *config.ClusterConfig {
	return c.config
}

// List Topics
func (c *Client) ListTopics() (map[string]sarama.TopicDetail, error) {
	return c.adminClient.ListTopics()
}

// Create Topic
func (c *Client) CreateTopic(name string, detail *sarama.TopicDetail) error {
	return c.adminClient.CreateTopic(name, detail, false)
}
