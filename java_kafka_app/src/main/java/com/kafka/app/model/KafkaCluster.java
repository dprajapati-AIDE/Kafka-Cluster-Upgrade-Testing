package com.kafka.app.model;

import java.util.List;

public class KafkaCluster {
    public String name;
    public String version;
    public List<String> brokers;
    public List<Topic> topics;

    public KafkaCluster() {
    }

    public String getName() {
        return name;
    }

    public String getVersion() {
        return version;
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

    public void setVersion(String version) {
        this.version = version;
    }

    public void setBrokers(List<String> brokers) {
        this.brokers = brokers;
    }

    public void setTopics(List<Topic> topics) {
        this.topics = topics;
    }
}
