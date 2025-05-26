package com.kafka.app.model;

import java.util.List;

public class KafkaCluster {
    public String name;
    public List<String> brokers;
    public List<Topic> topics;

    public KafkaCluster() {
    }

    public String getName() {
        return name;
    }

    public List<String> getBrokers() {
        return brokers;
    }

    public List<Topic> getTopics() {
        return topics;
    }

    public void setName(String name) {
        this.name = name;
    }

    public void setBrokers(List<String> brokers) {
        this.brokers = brokers;
    }

    public void setTopics(List<Topic> topics) {
        this.topics = topics;
    }
}
