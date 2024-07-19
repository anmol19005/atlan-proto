package main

import (
	"log"

	"atlan-proto/handlers"
	"atlan-proto/services"
	"atlan-proto/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func setupOperationsConfig(redisClient *redis.Client) {
	configs := []utils.OperationConfig{
		{
			Source:     "jira",
			Operations: []string{"add_timestamp", "add_processing_info"},
		},
		{
			Source:     "github",
			Operations: []string{"add_timestamp"},
		},
		{
			Source:     "serviceX",
			Operations: []string{"add_processing_info"},
		},
	}

	for _, config := range configs {
		utils.SetOperationsConfig(redisClient, config)
	}
}

func main() {
	// Initialize Redis client
	redisClient := utils.InitRedis()
	defer redisClient.Close()

	// Setup operations configurations
	setupOperationsConfig(redisClient)

	// Set up HTTP server and routes
	router := gin.Default()
	router.Use(handlers.Authenticate)
	router.POST("/ingest", handlers.IngestHandler(redisClient))
	router.GET("/consume", handlers.ConsumeHandler(redisClient))

	// Start pre-process metadata goroutine
	go services.PreProcessMetadata(redisClient)
	go services.PostConsumeTransformations(redisClient)

	// Start HTTP server
	log.Println("Starting server on :8080")
	log.Fatal(router.Run(":8080"))
}
