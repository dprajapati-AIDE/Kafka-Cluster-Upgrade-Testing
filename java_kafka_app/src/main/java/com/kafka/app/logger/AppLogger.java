package com.kafka.app.logger;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Level;
import org.apache.logging.log4j.core.config.builder.api.*;
import org.apache.logging.log4j.core.config.builder.impl.BuiltConfiguration;
import org.apache.logging.log4j.core.config.builder.impl.DefaultConfigurationBuilder;
import org.apache.logging.log4j.core.LoggerContext;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.kafka.app.model.LoggingConfig;

import java.time.LocalDate;

public class AppLogger {

    private static boolean initialized = false;
    private static final String LOG_DIR = "./logs";

    public static synchronized void initializeLogger(LoggingConfig config) {
        if (initialized)
            return;

        Level logLevel = parseLogLevel(config != null ? config.level : null);
        String encoding = (config != null && config.encoding != null) ? config.encoding : "UTF-8";

        LoggerContext context = (LoggerContext) LogManager.getContext(false);
        ConfigurationBuilder<BuiltConfiguration> builder = new DefaultConfigurationBuilder<>();
        builder.setStatusLevel(Level.ERROR);
        builder.setConfigurationName("CustomConfig");

        String pattern = "%d{yyyy-MM-dd HH:mm:ss.SSS} [%t] %-5level %logger{36} - %msg%n";

        LayoutComponentBuilder layout = builder.newLayout("PatternLayout")
                .addAttribute("pattern", pattern)
                .addAttribute("charset", encoding);

        AppenderComponentBuilder console = builder.newAppender("Console", "CONSOLE")
                .add(layout);
        builder.add(console);

        String dateStr = LocalDate.now().toString();
        String fileName = LOG_DIR + "/" + dateStr + ".log";
        String filePattern = LOG_DIR + "/" + dateStr + "-%i.log.gz";

        ComponentBuilder<?> triggeringPolicy = builder.newComponent("Policies")
                .addComponent(builder.newComponent("SizeBasedTriggeringPolicy").addAttribute("size", "10MB"));

        AppenderComponentBuilder rollingFile = builder.newAppender("RollingFile", "RollingFile")
                .addAttribute("fileName", fileName)
                .addAttribute("filePattern", filePattern)
                .add(layout)
                .addComponent(triggeringPolicy);
        builder.add(rollingFile);

        builder.add(builder.newRootLogger(logLevel)
                .add(builder.newAppenderRef("Console"))
                .add(builder.newAppenderRef("RollingFile")));

        context.start(builder.build());

        LoggerFactory.getLogger(AppLogger.class).info(
                "Logger initialized with level: {} and encoding: {}",
                logLevel, encoding);

        initialized = true;
    }

    public static Logger getLogger(Class<?> clazz) {
        if (!initialized) {
            initializeLogger(null);
        }
        return LoggerFactory.getLogger(clazz);
    }

    private static Level parseLogLevel(String levelStr) {
        if (levelStr == null)
            return Level.INFO;
        switch (levelStr.toLowerCase()) {
            case "trace":
                return Level.TRACE;
            case "debug":
                return Level.DEBUG;
            case "info":
                return Level.INFO;
            case "warn":
                return Level.WARN;
            case "error":
                return Level.ERROR;
            case "off":
                return Level.OFF;
            default:
                System.err.println("Unknown log level: " + levelStr + ", defaulting to INFO.");
                return Level.INFO;
        }
    }
}
