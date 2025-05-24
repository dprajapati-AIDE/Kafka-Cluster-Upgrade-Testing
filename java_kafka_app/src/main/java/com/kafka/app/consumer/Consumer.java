package com.kafka.app.consumer;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.kafka.app.logger.AppLogger;

import org.apache.kafka.clients.consumer.ConsumerRecord;
import org.apache.kafka.clients.consumer.ConsumerRecords;
import org.slf4j.Logger;

import java.util.Map;

public class Consumer {

    private static final Logger logger = AppLogger.getLogger(Consumer.class);
    private static final ObjectMapper objectMapper = new ObjectMapper();

    public static void process(String clusterName, ConsumerRecords<String, String> records) {
        for (ConsumerRecord<String, String> record : records) {
            try {
                Map<String, Object> payload = objectMapper.readValue(record.value(), Map.class);

                logger.info(
                        "Consumed Message from cluster [{}], topic [{}], partition [{}], offset [{}]: [{}][{}]",
                        clusterName, record.topic(), record.partition(), record.offset(), payload);
            } catch (Exception e) {
                logger.error("Failed to deserialize message from cluster [{}], topic [{}]: [{}][{}]",
                        clusterName, record.topic(), e.getMessage());
                e.printStackTrace();
            }
        }
    }
}
