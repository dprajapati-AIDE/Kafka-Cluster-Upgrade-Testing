package utils

import (
	"encoding/json"
	"go_kafka_app/internal/logger"
	"runtime"
	"strings"
	"time"

	"github.com/go-faker/faker/v4"
	"go.uber.org/zap"
)

func GetFunctionName(skip int) string {
	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		return "unknown"
	}
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown"
	}
	parts := strings.Split(fn.Name(), "/")
	return parts[len(parts)-1]
}

func GenerateMessage(messageSizeKB int) ([]byte, error) {
	msg := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"messageID": faker.UUIDDigit(),
		"sourceIp":  faker.IPv4(),
	}

	baseMessageSize := EstimateBaseMessageSize(msg)
	targetMessageBytes := messageSizeKB * 1024

	payload := ""
	if messageSizeKB > 0 {
		targetPayloadBytes := targetMessageBytes - baseMessageSize
		if targetPayloadBytes < 0 {
			targetPayloadBytes = 0
		}

		var payloadBuilder strings.Builder
		for payloadBuilder.Len() < targetPayloadBytes {
			para := faker.Paragraph()
			if payloadBuilder.Len()+len(para) > targetPayloadBytes {
				break
			}
			payloadBuilder.WriteString(para + " ")
		}
		payload = payloadBuilder.String()
	} else {
		payload = faker.Sentence()
	}
	msg["data"] = payload

	jsonBytes, err := json.Marshal(msg)
	if err != nil {
		logger.Error("Failed to marshal message", zap.Error(err), zap.String("func", GetFunctionName(1)))
		return nil, err
	}

	return jsonBytes, nil
}

func EstimateBaseMessageSize(baseFields map[string]interface{}) int {
	tempData, hasData := baseFields["data"]
	delete(baseFields, "data")

	jsonBytes, err := json.Marshal(baseFields)
	if err != nil {
		logger.Warn("Could not estimate base message size", zap.Error(err), zap.String("func", GetFunctionName(1)))
		return 50
	}

	if hasData {
		baseFields["data"] = tempData
	}
	return len(jsonBytes)
}
