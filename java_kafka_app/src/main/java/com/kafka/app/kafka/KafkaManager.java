package com.kafka.app.kafka;

import com.kafka.app.logger.AppLogger;
import com.kafka.app.model.KafkaCluster;
import org.slf4j.Logger;

import java.io.IOException;
import java.net.InetSocketAddress;
import java.net.Socket;
import java.util.List;

public class KafkaManager {

    private final KafkaCluster clusterConfig;
    private static final Logger logger = AppLogger.getLogger(KafkaManager.class);

    public KafkaManager(KafkaCluster clusterConfig) {
        this.clusterConfig = clusterConfig;
    }

    public boolean connect() {
        List<String> brokers = clusterConfig.getBrokers();

        if (brokers == null || brokers.isEmpty()) {
            logger.error("No brokers specified in KafkaCluster config");
            throw new IllegalArgumentException("No brokers specified in KafkaCluster config");
        }

        for (String broker : brokers) {
            try {
                String[] parts = broker.split(":");
                String host = parts[0];
                int port = Integer.parseInt(parts[1]);

                try (Socket socket = new Socket()) {
                    socket.connect(new InetSocketAddress(host, port), 2000);
                    logger.info("Successfully connected to Kafka broker: {}:{}", host, port);
                    return true;
                }
            } catch (IOException | NumberFormatException e) {
                logger.warn("Failed to connect to Kafka broker {}: {}",
                        broker, e.getMessage());
            }
        }

        logger.error("Unable to connect to any Kafka brokers");
        return false;
    }
}
