package com.kafka.app.model;

import java.util.List;

public class DeviceTypes {
    public List<String> firewall;
    public List<String> router;

    public List<String> switchNW;

    public DeviceTypes() {
    }

    public List<String> getFirewall() {
        return firewall;
    }

    public List<String> getRouter() {
        return router;
    }

    public List<String> getSwitch() {
        return switchNW;
    }

    public void setFirewall(List<String> firewall) {
        this.firewall = firewall;
    }

    public void setRouter(List<String> router) {
        this.router = router;
    }

    public void setSwitch(List<String> switchNW) {
        this.switchNW = switchNW;
    }

    public List<String> getByDeviceType(String deviceType) {
        if (deviceType == null)
            return null;
        switch (deviceType.toLowerCase()) {
            case "firewall":
                return firewall;
            case "router":
                return router;
            case "switch":
                return switchNW;
            default:
                return null;
        }
    }
}
