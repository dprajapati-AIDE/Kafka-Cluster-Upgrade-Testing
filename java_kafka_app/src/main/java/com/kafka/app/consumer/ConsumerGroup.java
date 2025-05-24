package com.kafka.app.consumer;

import com.kafka.app.logger.AppLogger;
import com.kafka.app.model.KafkaCluster;
import org.apache.kafka.clients.consumer.*;
import org.apache.kafka.common.serialization.StringDeserializer;
import org.slf4j.Logger;

import java.util.*;
import java.util.concurrent.*;

public class ConsumerGroup {

    private static final Logger logger = AppLogger.getLogger(Consumer.class);

    public static void start(KafkaCluster cluster, String groupId, List<String> topics) {
        Properties props = new Properties();
        props.put(ConsumerConfig.BOOTSTRAP_SERVERS_CONFIG, String.join(",", cluster.getBrokers()));
        props.put(ConsumerConfig.GROUP_ID_CONFIG, groupId);
        props.put(ConsumerConfig.KEY_DESERIALIZER_CLASS_CONFIG, StringDeserializer.class.getName());
        props.put(ConsumerConfig.VALUE_DESERIALIZER_CLASS_CONFIG, StringDeserializer.class.getName());
        props.put(ConsumerConfig.AUTO_OFFSET_RESET_CONFIG, "earliest");
        props.put(ConsumerConfig.ENABLE_AUTO_COMMIT_CONFIG, "false");

        KafkaConsumer<String, String> consumer = new KafkaConsumer<>(props);
        consumer.subscribe(topics);

        ExecutorService executor = Executors.newSingleThreadExecutor();
        executor.submit(() -> {
            try {
                while (true) {
                    ConsumerRecords<String, String> records = consumer.poll(1000);
                    Consumer.process(cluster.getName(), records);
                    consumer.commitAsync();
                }
            } catch (Exception e) {
                logger.error("Error consuming from cluster [{}]: [{}][{}]", cluster.getName(), e.getMessage());
                e.printStackTrace();
            } finally {
                consumer.close();
            }
        });
    }
}
