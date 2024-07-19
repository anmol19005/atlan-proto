package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"context"

	"atlan-proto/services"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func IngestHandler(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload services.GenericPayload
		if err := c.BindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		payload.Metadata["processed_at"] = time.Now().Format(time.RFC3339)

		jsonData, err := json.Marshal(payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal payload"})
			return
		}

		// Publish payload to the channel
		redisClient.Publish(context.Background(), "preProcessChannel", jsonData)

		c.JSON(http.StatusOK, gin.H{"message": "Payload received"})
	}
}
