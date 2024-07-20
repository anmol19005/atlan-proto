package utils

import (
	"context"
	"encoding/json"
	"log"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type OperationConfig struct {
	Source         string   `json:"source"`
	PreOperations  []string `json:"pre_operations"`
	PostOperations []string `json:"post_operations"`
}

func InitRedis() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	return redisClient
}

func SetOperationsConfig(redisClient *redis.Client, config OperationConfig) {
	configJSON, err := json.Marshal(config)
	if err != nil {
		log.Fatalf("Failed to marshal operations config: %v", err)
	}

	err = redisClient.Set(ctx, "config:"+config.Source, configJSON, 0).Err()
	if err != nil {
		log.Fatalf("Could not store operations config in Redis: %v", err)
	}

	log.Printf("Stored operations config for source: %s", config.Source)
}

func GetOperationsConfig(redisClient *redis.Client, source string) (OperationConfig, error) {
	var config OperationConfig

	configJSON, err := redisClient.Get(ctx, "config:"+source).Result()
	if err != nil {
		return config, err
	}

	err = json.Unmarshal([]byte(configJSON), &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
