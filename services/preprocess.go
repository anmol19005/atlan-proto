package services

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"atlan-proto/utils"

	"github.com/go-redis/redis/v8"
)

type GenericPayload struct {
	ID       string                 `json:"id"`
	Source   string                 `json:"source"`
	Payload  map[string]interface{} `json:"payload"`
	Metadata map[string]string      `json:"metadata"`
}

func addTimestamp(payload map[string]interface{}) map[string]interface{} {
	payload["timestamp"] = time.Now().Format(time.RFC3339)
	return payload
}

func addProcessingInfo(payload map[string]interface{}, source string) map[string]interface{} {
	payload["processed_by"] = source
	payload["processed_at"] = time.Now().Format(time.RFC3339)
	return payload
}

func executeOperations(payload GenericPayload, operations []string) GenericPayload {
	for _, operation := range operations {
		switch operation {
		case "add_timestamp":
			payload.Payload = addTimestamp(payload.Payload)
		case "add_processing_info":
			payload.Payload = addProcessingInfo(payload.Payload, payload.Source)
		default:
			log.Printf("Unknown operation: %s", operation)
		}
	}
	return payload
}

func PreProcessMetadata(redisClient *redis.Client) {
	subscriber := redisClient.Subscribe(context.Background(), "preProcessChannel")
	channel := subscriber.Channel()

	for msg := range channel {
		var payload GenericPayload
		if err := json.Unmarshal([]byte(msg.Payload), &payload); err != nil {
			log.Printf("Failed to unmarshal payload: %v", err)
			continue
		}

		log.Printf("Pre-processing payload: %s", payload.ID)

		config, err := utils.GetOperationsConfig(redisClient, payload.Source)
		if err != nil {
			log.Printf("Failed to get operations config for source %s: %v", payload.Source, err)
			continue
		}

		payload = executeOperations(payload, config.Operations)

		updatedPayload, err := json.Marshal(payload)
		if err != nil {
			log.Printf("Failed to marshal payload: %v", err)
			continue
		}

		err = redisClient.Set(context.Background(), "payload:"+payload.ID, updatedPayload, 0).Err()
		if err != nil {
			log.Printf("Could not store payload in Redis: %v", err)
		}
	}
}
