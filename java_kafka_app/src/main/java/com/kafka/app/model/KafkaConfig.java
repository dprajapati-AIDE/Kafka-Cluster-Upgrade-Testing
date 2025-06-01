package com.kafka.app.model;

import java.util.List;

public class KafkaConfig {
    public List<KafkaCluster> clusters;
    public List<String> topics;

    public KafkaConfig() {
    }

    public List<KafkaCluster> getClusters() {
        return clusters;
    }

    public List<String> getTopics() {
        return topics;
    }

    public void setClusters(List<KafkaCluster> clusters) {
        this.clusters = clusters;
    }

    public void setTopics(List<String> topics) {
        this.topics = topics;
    }
}
