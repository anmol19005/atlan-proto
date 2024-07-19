package services

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

func PostConsumeTransformations(redisClient *redis.Client) {
	for {
		keys, err := redisClient.Keys(context.Background(), "payload:*").Result()
		if err != nil {
			log.Printf("Could not retrieve payload keys: %v", err)
			continue
		}

		for _, key := range keys {
			payloadStr, err := redisClient.Get(context.Background(), key).Result()
			if err != nil {
				log.Printf("Could not retrieve payload for key %s: %v", key, err)
				continue
			}

			var payload GenericPayload
			if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
				log.Printf("Failed to unmarshal payload: %v", err)
				continue
			}

			log.Printf("Post-processing payload: %s", payload.ID)
			pushToDownstream(payloadStr)
			//deletes the key once pushing data
			// err = redisClient.Del(context.Background(), key).Err()
			// if err != nil {
			// 	log.Printf("Could not delete payload for key %s: %v", key, err)
			// }
		}
		time.Sleep(10 * time.Second) // Interval for processing payloads
	}
}
