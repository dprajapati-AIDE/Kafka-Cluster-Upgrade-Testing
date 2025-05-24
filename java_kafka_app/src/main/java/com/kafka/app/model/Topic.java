package com.kafka.app.model;

public class Topic {
    public String device;
    public String name;

    public Topic() {
    }

    public String getDevice() {
        return device;
    }

    public String getName() {
        return name;
    }

    public void setDevice(String device) {
        this.device = device;
    }

    public void setName(String name) {
        this.name = name;
    }
}
