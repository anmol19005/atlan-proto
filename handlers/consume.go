package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"atlan-proto/services"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func ConsumeHandler(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		keys, err := redisClient.Keys(context.Background(), "processed:*").Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve payload keys"})
			return
		}

		var payloadList []services.GenericPayload
		for _, key := range keys {
			payloadStr, err := redisClient.Get(context.Background(), key).Result()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve payload"})
				return
			}

			var payload services.GenericPayload
			if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal payload"})
				return
			}
			payloadList = append(payloadList, payload)
		}

		c.JSON(http.StatusOK, gin.H{"payloads": payloadList})
	}
}
