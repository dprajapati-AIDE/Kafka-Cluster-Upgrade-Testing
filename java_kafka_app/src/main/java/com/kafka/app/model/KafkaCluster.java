package com.kafka.app.model;

import java.util.List;

public class KafkaCluster {
    public String name;
    public List<String> brokers;

    public KafkaCluster() {
    }

    public String getName() {
        return name;
    }

    public List<String> getBrokers() {
        return brokers;
    }

    public void setName(String name) {
        this.name = name;
    }

    public void setBrokers(List<String> brokers) {
        this.brokers = brokers;
    }

}
