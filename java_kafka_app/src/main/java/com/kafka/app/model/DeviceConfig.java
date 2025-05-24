package com.kafka.app.model;

public class DeviceConfig {
    public String vendor;
    public DeviceTypes types;

    public DeviceConfig() {
    }

    public String getVendor() {
        return vendor;
    }

    public DeviceTypes getTypes() {
        return types;
    }

    public void setVendor(String vendor) {
        this.vendor = vendor;
    }

    public void setTypes(DeviceTypes types) {
        this.types = types;
    }
}
