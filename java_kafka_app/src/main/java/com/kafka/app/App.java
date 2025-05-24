package com.kafka.app;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.dataformat.yaml.YAMLFactory;
import com.kafka.app.model.AppConfig;
import com.kafka.app.logger.AppLogger;

import org.slf4j.Logger;
import java.io.InputStream;

public class App {

    private static final Logger logger = AppLogger.getLogger(App.class);

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

            AppLogger.initializeLogger(config.getLogging());

            logger.info("--- Application Starting ---");

            if (config != null) {
                logger.info("Configuration loaded successfully.");
            }

        } catch (Exception e) {
            if (logger != null) {
                logger.error("Failed to load configuration or initialize logger: {}", e.getMessage(), e);
            } else {
                System.err.println("Failed to load configuration or initialize logger: " + e.getMessage());
                e.printStackTrace();
            }
            return;
        }

        logger.info("--- Application Exiting ---");
    }
}
