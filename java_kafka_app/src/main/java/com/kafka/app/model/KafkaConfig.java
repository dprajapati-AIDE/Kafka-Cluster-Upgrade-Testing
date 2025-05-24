package com.kafka.app.model;

import java.util.List;

public class KafkaConfig {
    public List<KafkaCluster> clusters;

    public KafkaConfig() {
    }

    public List<KafkaCluster> getClusters() {
        return clusters;
    }

    public void setClusters(List<KafkaCluster> clusters) {
        this.clusters = clusters;
    }
}
