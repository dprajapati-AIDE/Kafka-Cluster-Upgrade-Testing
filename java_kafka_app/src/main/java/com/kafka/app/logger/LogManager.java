package com.kafka.app.logger;

import com.kafka.app.model.LoggingConfig;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import ch.qos.logback.classic.Level;
import ch.qos.logback.classic.LoggerContext;
import ch.qos.logback.classic.encoder.PatternLayoutEncoder;
import ch.qos.logback.classic.spi.ILoggingEvent;
import ch.qos.logback.core.ConsoleAppender;
import ch.qos.logback.core.rolling.RollingFileAppender;
import ch.qos.logback.core.rolling.TimeBasedRollingPolicy;

import java.nio.charset.Charset;
import java.nio.charset.StandardCharsets;
import java.io.File;

public class LogManager {

    private static boolean initialized = false;
    private static final String LOG_DIR = "./logs";

    public static synchronized void initializeLogger(LoggingConfig config) {
        if (initialized) {
            System.err.println("Logger already initialized. Skipping re-initialization.");
            return;
        }

        if (config == null) {
            System.err.println("Logging configuration is null. Using default INFO level and UTF-8 encoding.");
        }

        LoggerContext loggerContext = (LoggerContext) LoggerFactory.getILoggerFactory();
        loggerContext.reset();

        Level logLevel = parseLogLevel(config != null ? config.level : null);
        if (logLevel == null) {
            System.err.println("Invalid or missing log level. Defaulting to INFO.");
            logLevel = Level.INFO;
        }

        Charset charset = StandardCharsets.UTF_8;
        if (config != null && config.encoding != null && !config.encoding.equalsIgnoreCase("json")) {
            Charset parsedCharset = parseCharset(config.encoding);
            if (parsedCharset != null) {
                charset = parsedCharset;
            } else {
                System.err.println("Invalid encoding specified: " + config.encoding + ". Defaulting to UTF-8.");
            }
        }

        PatternLayoutEncoder encoder = new PatternLayoutEncoder();
        encoder.setContext(loggerContext);
        encoder.setPattern("%d{yyyy-MM-dd HH:mm:ss.SSS} [%thread] %-5level %logger{36} - %msg%n");
        encoder.setCharset(charset);
        encoder.start();

        addConsoleAppender(loggerContext, encoder);
        addFileAppender(loggerContext, encoder);

        ch.qos.logback.classic.Logger rootLogger = loggerContext.getLogger(Logger.ROOT_LOGGER_NAME);
        rootLogger.setLevel(logLevel);

        initialized = true;
        LoggerFactory.getLogger(LogManager.class).info(
                "Logger initialized successfully with level: {} and encoding: {}. Output to console and daily file.",
                logLevel, charset.name());
    }

    public static Logger getLogger(Class<?> clazz) {
        if (!initialized) {
            initializeLogger(null);
        }
        return LoggerFactory.getLogger(clazz);
    }

    private static Level parseLogLevel(String levelStr) {
        if (levelStr == null) {
            return null;
        }
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
                return null;
        }
    }

    private static Charset parseCharset(String encodingStr) {
        if (encodingStr == null) {
            return null;
        }
        try {
            return Charset.forName(encodingStr);
        } catch (Exception e) {
            return null; // Invalid charset name
        }
    }

    private static void addConsoleAppender(LoggerContext loggerContext, PatternLayoutEncoder encoder) {
        ConsoleAppender<ILoggingEvent> consoleAppender = new ConsoleAppender<>();
        consoleAppender.setContext(loggerContext);
        consoleAppender.setName("STDOUT");
        consoleAppender.setEncoder(encoder);
        consoleAppender.start();
        loggerContext.getLogger(Logger.ROOT_LOGGER_NAME).addAppender(consoleAppender);
    }

    private static void addFileAppender(LoggerContext loggerContext, PatternLayoutEncoder encoder) {
        File logDirectory = new File(LOG_DIR);
        if (!logDirectory.exists()) {
            if (!logDirectory.mkdirs()) {
                System.err.println("Failed to create log directory: " + LOG_DIR);
                return;
            }
        }

        RollingFileAppender<ILoggingEvent> fileAppender = new RollingFileAppender<>();
        fileAppender.setContext(loggerContext);
        fileAppender.setName("FILE");
        fileAppender.setFile(LOG_DIR + "/app.log");
        fileAppender.setEncoder(encoder);

        TimeBasedRollingPolicy<ILoggingEvent> rollingPolicy = new TimeBasedRollingPolicy<>();
        rollingPolicy.setContext(loggerContext);
        rollingPolicy.setParent(fileAppender);
        rollingPolicy.setFileNamePattern(LOG_DIR + "/app-%d{yyyy-MM-dd}.log");
        rollingPolicy.start();

        fileAppender.setRollingPolicy(rollingPolicy);
        fileAppender.start();

        loggerContext.getLogger(Logger.ROOT_LOGGER_NAME).addAppender(fileAppender);
    }
}
