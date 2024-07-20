package services

import (
	"context"
	"encoding/json"
	"log"

	"atlan-proto/utils"

	"github.com/go-redis/redis/v8"
)

func PostConsumeTransformations(redisClient *redis.Client) {
	subscriber := redisClient.Subscribe(context.Background(), "postProcessChannel")
	channel := subscriber.Channel()

	for msg := range channel {
		var payload GenericPayload
		if err := json.Unmarshal([]byte(msg.Payload), &payload); err != nil {
			log.Printf("Failed to unmarshal payload: %v", err)
			continue
		}

		log.Printf("Post-processing payload: %s", payload.ID)

		config, err := utils.GetOperationsConfig(redisClient, payload.Source)
		if err != nil {
			log.Printf("Failed to get operations config for source %s: %v", payload.Source, err)
			continue
		}

		payload = executeOperations(payload, config.PostOperations)

		updatedPayload, err := json.Marshal(payload)
		if err != nil {
			log.Printf("Failed to marshal payload: %v", err)
			continue
		}

		// Push processed payload to downstream services
		pushToDownstream(updatedPayload)

		err = redisClient.Del(context.Background(), "processed:"+payload.ID).Err()
		if err != nil {
			log.Printf("Could not delete payload for key %s: %v", payload.ID, err)
		}
	}
}
