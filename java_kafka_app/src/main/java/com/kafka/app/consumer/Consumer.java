package com.kafka.app.consumer;

import com.kafka.app.logger.AppLogger;

import org.apache.kafka.clients.consumer.ConsumerRecord;
import org.apache.kafka.clients.consumer.ConsumerRecords;
import org.slf4j.Logger;

public class Consumer {

    private static final Logger logger = AppLogger.getLogger(Consumer.class);

    public static void process(String clusterName, ConsumerRecords<String, String> records) {
        for (ConsumerRecord<String, String> record : records) {
            try {
                logger.info(
                        "Consumed Message from cluster [{}], topic [{}], partition [{}], offset [{}]",
                        clusterName, record.topic(), record.partition(), record.offset());
            } catch (Exception e) {
                logger.error("Error processing message from cluster [{}], topic [{}]: {}",
                        clusterName, record.topic(), e.getMessage());
                e.printStackTrace();
            }
        }
    }
}
