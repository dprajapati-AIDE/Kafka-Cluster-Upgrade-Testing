package com.kafka.app;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.dataformat.yaml.YAMLFactory;
import com.kafka.app.logger.AppLogger;
import com.kafka.app.model.AppConfig;
import com.kafka.app.model.KafkaCluster;
import com.kafka.app.producer.Producer;
import com.kafka.app.kafka.KafkaManager;
import com.kafka.app.topic.TopicManager;

import org.slf4j.Logger;

import java.io.InputStream;
import java.util.ArrayList;
import java.util.List;

public class App {

    private static final Logger logger = AppLogger.getLogger(App.class);

    public static void main(String[] args) {
        AppConfig config;

        try {
            // Load YAML config
            ObjectMapper mapper = new ObjectMapper(new YAMLFactory());
            InputStream inputStream = App.class.getClassLoader().getResourceAsStream("config.yaml");

            if (inputStream == null) {
                System.err.println("Error: config.yaml not found in resources!");
                return;
            }

            config = mapper.readValue(inputStream, AppConfig.class);

            // Initialize logger
            AppLogger.initializeLogger(config.getLogging());
            logger.info("--- Application Starting ---");

            // Confirm config load
            logger.info("Configuration loaded successfully.");

            List<KafkaCluster> clusterList = config.getKafka().getClusters();

            if (clusterList == null || clusterList.isEmpty()) {
                logger.error("No Kafka clusters found in configuration.");
                return;
            }

            // Check Kafka cluster availability
            List<KafkaCluster> reachableClusters = new ArrayList<>();

            for (KafkaCluster clusterConfig : clusterList) {
                KafkaManager kafkaManager = new KafkaManager(clusterConfig);
                boolean isKafkaUp = kafkaManager.connect();

                if (isKafkaUp) {
                    logger.info("Kafka cluster [{}] is reachable.", clusterConfig.getName());
                    reachableClusters.add(clusterConfig);
                } else {
                    logger.error("Kafka cluster [{}] is NOT reachable.", clusterConfig.getName());
                }
            }

            if (reachableClusters.isEmpty()) {
                logger.error("No reachable Kafka clusters found. Exiting.");
                return;
            }

            // Process topics on all reachable clusters
            TopicManager topicManager = new TopicManager();
            topicManager.processClusters(reachableClusters);

            Producer producerService = new Producer(config.getDeviceConfig());
            producerService.produceToClusters(reachableClusters, 10);

        } catch (Exception e) {
            logger.error("Failed to load configuration or initialize components: {}", e.getMessage(), e);
            return;
        }

        logger.info("--- Application Exiting ---");
    }
}
