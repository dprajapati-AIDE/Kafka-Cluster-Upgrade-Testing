package com.kafka.app.producer;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.github.javafaker.Faker;
import com.kafka.app.logger.AppLogger;
import com.kafka.app.model.*;

import org.apache.kafka.clients.producer.*;
import org.apache.kafka.common.serialization.StringSerializer;
import org.slf4j.Logger;

import java.time.Instant;
import java.util.*;
import java.util.concurrent.*;

public class Producer {

    private final ObjectMapper objectMapper = new ObjectMapper();
    private final Faker faker = new Faker();
    private final DeviceConfig deviceConfig;
    private static final Logger logger = AppLogger.getLogger(Producer.class);

    public Producer(DeviceConfig deviceConfig) {
        this.deviceConfig = deviceConfig;
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
            for (Topic topic : cluster.getTopics()) {
                String deviceType = topic.getDevice();
                List<String> deviceModels = deviceConfig.getTypes().getByDeviceType(deviceType);

                if (deviceModels == null || deviceModels.isEmpty()) {
                    logger.error("No device models found for type [{}] in config.%n", deviceType);
                    continue;
                }

                for (int i = 0; i < messageCount; i++) {
                    DeviceMessageModel msg = new DeviceMessageModel();
                    msg.timestamp = Instant.now().toString();
                    msg.deviceId = UUID.randomUUID().toString();
                    msg.deviceIp = faker.internet().ipV4Address();
                    msg.deviceType = deviceType;
                    msg.deviceModel = deviceModels.get(0);
                    msg.vendor = deviceConfig.getVendor();
                    msg.message = faker.lorem().sentence();

                    String json = objectMapper.writeValueAsString(msg);
                    ProducerRecord<String, String> record = new ProducerRecord<>(topic.getName(), json);
                    RecordMetadata metadata = producer.send(record).get();

                    logger.info("Sent message to cluster [{}], topic [{}], partition [{}], offset [{}]",
                            cluster.getName(), topic.getName(), metadata.partition(), metadata.offset());

                }
            }
        } catch (Exception e) {
            logger.error("Error producing to cluster [{}]: [{}][{}]", cluster.getName(), e.getMessage());
            e.printStackTrace();
        }
    }
}
