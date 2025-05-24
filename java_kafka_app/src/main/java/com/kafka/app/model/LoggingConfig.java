package com.kafka.app.model;

public class LoggingConfig {
    public String level;
    public String encoding;
    public String output;

    public LoggingConfig() {
    }

    public String getLevel() {
        return level;
    }

    public String getEncoding() {
        return encoding;
    }

    public String getOutput() {
        return output;
    }

    public void setLevel(String level) {
        this.level = level;
    }

    public void setEncoding(String encoding) {
        this.encoding = encoding;
    }

    public void setOutput(String output) {
        this.output = output;
    }
}
