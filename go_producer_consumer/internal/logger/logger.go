package logger

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Global logger instance
var log *zap.Logger

func Initialize(level, encoding, output string) error {
	var zapLevel zapcore.Level
	switch strings.ToLower(level) {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		return fmt.Errorf("invalid log level: %s", level)
	}

	// Configure how log messages are encoded (e.g., time format, keys)
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Choose encoder: JSON or console format
	var encoder zapcore.Encoder
	switch strings.ToLower(encoding) {
	case "json":
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	case "console":
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	default:
		return fmt.Errorf("invalid log encoding: %s", encoding)
	}

	// Set up the log output destination
	var writeSyncer zapcore.WriteSyncer
	switch strings.ToLower(output) {
	case "stdout":
		writeSyncer = zapcore.AddSync(os.Stdout)
	default:
		file, err := os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		writeSyncer = zapcore.AddSync(file)
	}

	// Create the zap core using encoder, output, and log level
	core := zapcore.NewCore(encoder, writeSyncer, zapLevel)

	// Initialize the global logger
	log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	Info("Logger initialized", zap.String("level", level), zap.String("encoding", encoding), zap.String("output", output))
	return nil
}

// debug-level message
func Debug(msg string, fields ...zap.Field) {
	log.Debug(msg, fields...)
}

// info-level message
func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

// warning-level message
func Warn(msg string, fields ...zap.Field) {
	log.Warn(msg, fields...)
}

// error-level message
func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

// fatal-level message and exist the app
func Fatal(msg string, fields ...zap.Field) {
	log.Fatal(msg, fields...)
}

func With(fields ...zap.Field) *zap.Logger {
	return log.With(fields...)
}

// Sync flushes any buffered log entries to the outpu
func Sync() error {
	return log.Sync()
}
