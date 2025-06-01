package com.kafka.app.producer;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.github.javafaker.Faker;
import com.kafka.app.logger.AppLogger;
import com.kafka.app.model.KafkaCluster;

import org.apache.kafka.clients.producer.*;
import org.apache.kafka.common.serialization.StringSerializer;
import org.slf4j.Logger;

import java.time.Instant;
import java.util.*;
import java.util.concurrent.*;

public class Producer {

    private final ObjectMapper objectMapper = new ObjectMapper();
    private final Faker faker = new Faker();
    private final List<String> topicsToProduce;
    private final int messageSizeKB;
    private static final Logger logger = AppLogger.getLogger(Producer.class);

    public Producer(List<String> topicsToProduce, int messageSizeKB) {
        this.topicsToProduce = topicsToProduce;
        this.messageSizeKB = messageSizeKB;
    }

    public void produceToClusters(List<KafkaCluster> clusters, int messageCount) {
        ExecutorService executor = Executors.newFixedThreadPool(clusters.size());

        for (KafkaCluster cluster : clusters) {
            executor.submit(() -> produceToCluster(cluster, messageCount));
        }

        executor.shutdown();
        try {
            executor.awaitTermination(5, TimeUnit.MINUTES);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            logger.error("Producer execution interrupted.");
        }
    }

    private void produceToCluster(KafkaCluster cluster, int messageCount) {
        Properties props = new Properties();
        props.put(ProducerConfig.BOOTSTRAP_SERVERS_CONFIG, String.join(",", cluster.getBrokers()));
        props.put(ProducerConfig.KEY_SERIALIZER_CLASS_CONFIG, StringSerializer.class.getName());
        props.put(ProducerConfig.VALUE_SERIALIZER_CLASS_CONFIG, StringSerializer.class.getName());
        props.put(ProducerConfig.ACKS_CONFIG, "all");
        props.put(ProducerConfig.RETRIES_CONFIG, 5);
        props.put(ProducerConfig.LINGER_MS_CONFIG, 5);
        props.put(ProducerConfig.MAX_BLOCK_MS_CONFIG, 10000);
        props.put(ProducerConfig.REQUEST_TIMEOUT_MS_CONFIG, 15000);

        try (KafkaProducer<String, String> producer = new KafkaProducer<>(props)) {
            for (String topicName : topicsToProduce) {
                for (int i = 0; i < messageCount; i++) {
                    Map<String, String> msg = new HashMap<>();
                    msg.put("timestamp", Instant.now().toString());
                    msg.put("messageId", UUID.randomUUID().toString());
                    msg.put("sourceIp", faker.internet().ipV4Address());

                    String payload = "";
                    if (messageSizeKB > 0) {

                        int targetPayloadBytes = (messageSizeKB * 1024) - estimateBaseMessageSize(msg);
                        if (targetPayloadBytes < 0)
                            targetPayloadBytes = 0;

                        StringBuilder sb = new StringBuilder();
                        while (sb.length() < targetPayloadBytes) {
                            sb.append(faker.lorem().characters(1000, 2000, true, true));
                        }
                        payload = sb.substring(0, Math.min(sb.length(), targetPayloadBytes));
                    } else {
                        payload = faker.lorem().sentence();
                    }
                    msg.put("data", payload);

                    String json = objectMapper.writeValueAsString(msg);
                    ProducerRecord<String, String> record = new ProducerRecord<>(topicName, json);
                    RecordMetadata metadata = producer.send(record).get();

                    logger.info("Sent message to cluster [{}], topic [{}], partition [{}], offset [{}]",
                            cluster.getName(), topicName, metadata.partition(), metadata.offset());
                }
            }
        } catch (Exception e) {
            logger.error("Error producing to cluster [{}]: [{}][{}]", cluster.getName(), e.getMessage());
            e.printStackTrace();
        }
    }

    private int estimateBaseMessageSize(Map<String, String> baseFields) {
        try {
            return objectMapper.writeValueAsString(baseFields).getBytes("UTF-8").length;
        } catch (Exception e) {
            logger.warn("Could not estimate base message size: {}", e.getMessage());
            return 50;
        }
    }
}
