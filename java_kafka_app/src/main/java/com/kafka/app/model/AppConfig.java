package com.kafka.app.model;

import com.fasterxml.jackson.annotation.JsonProperty;

public class AppConfig {

    public LoggingConfig logging;
    public KafkaConfig kafka;
    @JsonProperty("devices")
    public DeviceConfig deviceConfig;

    public AppConfig() {
    }

    public LoggingConfig getLogging() {
        return logging;
    }

    public KafkaConfig getKafka() {
        return kafka;
    }

    public DeviceConfig getDeviceConfig() {
        return deviceConfig;
    }

    public void setLogging(LoggingConfig logging) {
        this.logging = logging;
    }

    public void setKafka(KafkaConfig kafka) {
        this.kafka = kafka;
    }

    public void setDeviceConfig(DeviceConfig deviceConfig) {
        this.deviceConfig = deviceConfig;
    }
}
