package com.kafka.app.topic;

import com.kafka.app.logger.AppLogger;
import com.kafka.app.model.KafkaCluster;
import com.kafka.app.model.Topic;
import org.apache.kafka.clients.admin.AdminClient;
import org.apache.kafka.clients.admin.CreateTopicsResult;
import org.apache.kafka.clients.admin.NewTopic;
import org.slf4j.Logger;

import java.util.*;
import java.util.concurrent.*;

public class TopicManager {

    int DEFAULT_PARTIION_COUNT = 3;
    int DEFAULT_REPLICATION_FACTOR = 3;

    private static final Logger logger = AppLogger.getLogger(TopicManager.class);

    public void processClusters(List<KafkaCluster> clusters) {
        ExecutorService executor = Executors.newFixedThreadPool(clusters.size());

        for (KafkaCluster cluster : clusters) {
            executor.submit(() -> handleClusterTopics(cluster));
        }

        executor.shutdown();
        try {
            if (!executor.awaitTermination(5, TimeUnit.MINUTES)) {
                logger.error("Timeout: Not all cluster tasks completed.");
            }
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            logger.error("Execution interrupted: " + e.getMessage());
        }
    }

    private void handleClusterTopics(KafkaCluster cluster) {
        logger.info("Checking topics for cluster: " + cluster.getName());

        Properties config = new Properties();
        config.put("bootstrap.servers", String.join(",", cluster.getBrokers()));

        try (AdminClient adminClient = AdminClient.create(config)) {
            // Get existing topics
            Set<String> existingTopics = adminClient.listTopics().names().get();

            // Prepare topics to create
            List<NewTopic> topicsToCreate = new ArrayList<>();

            for (Topic topic : cluster.getTopics()) {
                if (!existingTopics.contains(topic.getName())) {
                    logger.info("Topic missing: " + topic.getName() + " — scheduling for creation.");
                    topicsToCreate.add(
                            new NewTopic(topic.getName(), DEFAULT_PARTIION_COUNT, (short) DEFAULT_REPLICATION_FACTOR));
                } else {
                    logger.info("Topic already exists: " + topic.getName());
                }
            }

            // Create missing topics
            if (!topicsToCreate.isEmpty()) {
                CreateTopicsResult result = adminClient.createTopics(topicsToCreate);
                result.all().get();
                logger.info("Created topics on cluster: " + cluster.getName());
            }
        } catch (Exception e) {
            logger.error("Error managing topics on cluster: " + cluster.getName());
            e.printStackTrace();
        }
    }
}
