package com.kafka.app;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.dataformat.yaml.YAMLFactory;
import com.kafka.app.consumer.ConsumerGroup;
import com.kafka.app.logger.AppLogger;
import com.kafka.app.model.AppConfig;
import com.kafka.app.model.KafkaCluster;
import com.kafka.app.producer.Producer;
import com.kafka.app.kafka.KafkaManager;
import com.kafka.app.topic.TopicManager;
import org.slf4j.Logger;

import java.io.InputStream;
import java.util.*;
import java.util.stream.Collectors;

public class Run {
    private static final Logger logger = AppLogger.getLogger(Run.class);

    public static void run(String role, int msgCount, String consumerGroupName) {
        Runtime.getRuntime().addShutdownHook(new Thread(() -> {
            logger.info("Shutdown signal received. Exiting...");
        }));

        AppConfig config;

        try {
            ObjectMapper mapper = new ObjectMapper(new YAMLFactory());
            InputStream inputStream = Run.class.getClassLoader().getResourceAsStream("config.yaml");

            if (inputStream == null) {
                System.err.println("config.yaml not found in resources!");
                return;
            }

            config = mapper.readValue(inputStream, AppConfig.class);
            AppLogger.initializeLogger(config.getLogging());
            logger.info("Configuration loaded");

            List<KafkaCluster> clusterList = config.getKafka().getClusters();
            List<KafkaCluster> reachableClusters = new ArrayList<>();

            for (KafkaCluster cluster : clusterList) {
                KafkaManager kafkaManager = new KafkaManager(cluster);
                if (kafkaManager.connect()) {
                    logger.info("Connected to Kafka cluster: {}", cluster.getName());
                    reachableClusters.add(cluster);
                } else {
                    logger.warn("Failed to connect to cluster: {}", cluster.getName());
                }
            }

            if (reachableClusters.isEmpty()) {
                logger.error("No reachable Kafka clusters");
                return;
            }

            TopicManager topicManager = new TopicManager();
            topicManager.processClusters(reachableClusters);

            if ("producer".equalsIgnoreCase(role) || "both".equalsIgnoreCase(role)) {
                Producer producer = new Producer(config.getDeviceConfig());
                producer.produceToClusters(reachableClusters, msgCount);
            }

            if ("consumer".equalsIgnoreCase(role) || "both".equalsIgnoreCase(role)) {
                for (KafkaCluster cluster : reachableClusters) {
                    List<String> topics = cluster.getTopics().stream()
                            .map(t -> t.getName())
                            .collect(Collectors.toList());

                    Thread consumerThread = new Thread(() -> {
                        ConsumerGroup.start(cluster, consumerGroupName, topics);
                    });
                    consumerThread.setDaemon(true);
                    consumerThread.start();
                }

                // Keep main thread alive for consumers
                Thread.currentThread().join();
            }

        } catch (Exception e) {
            logger.error("Error during app run: {}", e.getMessage(), e);
        }
    }
}
