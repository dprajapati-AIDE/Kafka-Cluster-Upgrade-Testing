package consumer

import (
	"fmt"
	"go_kafka_app/internal/config"
	"go_kafka_app/internal/logger"
	"go_kafka_app/internal/utils"
	"sync"
	"time"

	"github.com/mistsys/mist_go_utils/bus3/consumer3"
	"github.com/mistsys/mist_go_utils/cloud"
	consumer "github.com/mistsys/sarama-consumer"
	"github.com/mistsys/sarama-consumer/offsets"
	"github.com/mistsys/sarama-consumer/stable"
	"go.uber.org/zap"
)

func StartBus3ConsumerGroup(cluster *config.ClusterConfig, groupID string, topics []string) error {
	cloud.KAFKA_BROKERS = cluster.Brokers
	bus3Config := consumer.NewConfig()
	bus3Config.Heartbeat.Interval = 3 * time.Second
	bus3Config.Session.Timeout = 15 * time.Second
	bus3Config.Rebalance.Timeout = 10 * time.Second
	bus3Config.Partitioner = stable.New(false)
	bus3Config.StartingOffset, bus3Config.OffsetOutOfRange = offsets.NoOlderThan(time.Minute * 10)
	bus3Config.SidechannelTopic = ""

	sclient, client, err := consumer3.NewClient("bus3_consumer", groupID, nil, bus3Config)
	if err != nil {
		logger.Error("Failed to create consumer group", zap.String("func", utils.GetFunctionName(1)), zap.Error(err))
		return fmt.Errorf("failed to create consumer group: %w", err)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range client.Errors() {
			logger.Fatal("Error", zap.String("func", utils.GetFunctionName(1)), zap.Error(err))
		}
	}()

	for _, topic := range topics {
		client, err := client.Consume(topic)
		if err != nil {
			logger.Error("Error consuming from kafka", zap.String("func", utils.GetFunctionName(1)), zap.String("topic", topic), zap.Error(err))
			continue
		}
		go func(c consumer.Consumer, t string) {
			for msg := range c.Messages() {
				logger.Info("Consumed message", zap.String("cluster", cluster.Name), zap.String("topic", t), zap.Int32("partition", msg.Partition), zap.Int64("offset", msg.Offset))
				c.Done(msg)
			}
		}(client, topic)
	}

	defer sclient.Close()
	defer client.Close()
	wg.Wait()

	return nil
}
