package com.kafka.app;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.dataformat.yaml.YAMLFactory;
import com.kafka.app.model.AppConfig;
import com.kafka.app.logger.LogManager;

import org.slf4j.Logger;
import java.io.InputStream;

public class App {

    private static final Logger logger = LogManager.getLogger(App.class);

    public static void main(String[] args) {
        AppConfig config = null;
        try {
            ObjectMapper mapper = new ObjectMapper(new YAMLFactory());
            InputStream inputStream = App.class.getClassLoader().getResourceAsStream("config.yaml");

            if (inputStream == null) {
                System.err.println("Error: config.yaml not found in resources!");
                return;
            }

            config = mapper.readValue(inputStream, AppConfig.class);

            LogManager.initializeLogger(config.getLogging());

            logger.info("--- Application Starting ---");
            logger.debug("Configuration loaded successfully.");

            logger.info("Logging Level from config: {}", config.getLogging().getLevel());
            logger.info("First Kafka Cluster Name: {}", config.getKafka().getClusters().get(0).getName());

        } catch (Exception e) {
            if (logger != null) {
                logger.error("Failed to load configuration or initialize logger: {}", e.getMessage(), e);
            } else {
                System.err.println("Failed to load configuration or initialize logger: " + e.getMessage());
                e.printStackTrace();
            }
            return;
        }

        if (config != null) {
            logger.info("Configuration loaded successfully. Proceeding with application logic...");

            config.getKafka().getClusters().forEach(cluster -> {
                logger.info("  Processing Cluster: {} (Version: {})", cluster.getName(), cluster.getVersion());
                cluster.getBrokers().forEach(broker -> logger.debug("    Broker: {}", broker));
                cluster.getTopics().forEach(
                        topic -> logger.info("    Topic: {} (Device: {})", topic.getName(), topic.getDevice()));
            });

            logger.info("Application finished processing clusters.");
            logger.debug("This is a debug message.");
            logger.trace("This is a trace message.");
            logger.warn("This is a warning message.");
            logger.error("This is an error message. An example exception might follow.",
                    new RuntimeException("Simulated error"));
        }

        logger.info("--- Application Exiting ---");
    }
}
