package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Logging LoggingConfig `yaml:"logging"`
	Kafka   KafkaConfig   `yaml:"kafka"`
	Devices DevicesConfig `yaml:"devices"`
}

type LoggingConfig struct {
	Level    string `yaml:"level"`
	Encoding string `yaml:"encoding"`
	Output   string `yaml:"output"`
}

type KafkaConfig struct {
	Clusters []ClusterConfig `yaml:"clusters"`
}

type ClusterConfig struct {
	Name    string        `yaml:"name"`
	Version string        `yaml:"version"`
	Brokers []string      `yaml:"brokers"`
	Topics  []TopicConfig `yaml:"topics"`
}

type TopicConfig struct {
	Device            string            `yaml:"device"`
	Name              string            `yaml:"name"`
	Partitions        int               `yaml:"partitions"`
	ReplicationFactor int               `yaml:"replication_factor"`
	Config            map[string]string `yaml:"config"`
}

type DevicesConfig struct {
	Vendor string              `yaml:"vendor"`
	Types  map[string][]string `yaml:"types"`
}

func LoadConfig(configPath string) (*Config, error) {

	if configPath == "" {
		configPath = "config/config.yaml"
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	config := &Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

func validateConfig(config *Config) error {
	// Validate logging
	if !isValidLogLevel(config.Logging.Level) {
		return fmt.Errorf("invalid log level: %s", config.Logging.Level)
	}

	// Validate Kafka clusters
	if len(config.Kafka.Clusters) == 0 {
		return fmt.Errorf("no Kafka clusters configured")
	}

	// Validate Kafka Cluster Configuration
	for _, cluster := range config.Kafka.Clusters {
		if len(cluster.Brokers) == 0 {
			return fmt.Errorf("no brokers defined for cluster %s", cluster.Name)
		}
		for _, topic := range cluster.Topics {
			if topic.Device == "" || topic.Name == "" {
				return fmt.Errorf("invalid topic config in cluster %s: device and name are required", cluster.Name)
			}
			if topic.Partitions <= 0 {
				return fmt.Errorf("invalid partitions for topic %s in cluster %s", topic.Name, cluster.Name)
			}
			if topic.ReplicationFactor <= 0 {
				return fmt.Errorf("invalid replication_factor for topic %s in cluster %s", topic.Name, cluster.Name)
			}
		}
	}

	// Validate devices
	if config.Devices.Vendor == "" {
		return fmt.Errorf("vendor name is required in devices config")
	}
	if len(config.Devices.Types) == 0 {
		return fmt.Errorf("no devices configured under vendor %s", config.Devices.Vendor)
	}

	return nil
}

func isValidLogLevel(level string) bool {
	level = strings.ToLower(level)
	validLevels := []string{"debug", "info", "warn", "error"}

	for _, validLevel := range validLevels {
		if level == validLevel {
			return true
		}
	}

	return false
}

// get cluster configuration by name
func (k *KafkaConfig) GetClusterByName(name string) (*ClusterConfig, error) {
	for _, cluster := range k.Clusters {
		if cluster.Name == name {
			return &cluster, nil
		}
	}
	return nil, fmt.Errorf("cluster not found: %s", name)
}
