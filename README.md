## Atlan Prototype
This project is a prototype for ingesting and processing metadata in near real-time using Go, Redis, and Gin. It supports dynamic operations based on the source of the metadata and pushes the processed data to downstream services.

## Features
### Ingest Metadata:
    Ingest metadata from different sources and apply dynamic operations based on the source configuration.
### Pre-Process Metadata: 
    Apply operations such as adding a timestamp or processing information to the metadata.
### Store Processed Metadata: 
    Store the processed metadata in Redis.
### Post-Process Metadata: 
    Push the processed metadata to downstream services.
### Dynamic Operations: 
    Configure different operations for different sources.

## Requirements
### Go 1.16 or later
### Redis

## Installation

Clone the repository:

    git clone https://github.com/yourusername/atlan-prototype.git
    cd atlan-prototype
    
Install dependencies:
      
    go mod tidy
    
Start Redis server (ensure Redis is running on localhost:6379).

## Configuration
Set the secret key in config/config.go:

    const (
        SecretKey = "your-secret-key" // Replace with your actual secret key
    )
    
Configure operations for different sources in main.go:

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

## Usage

Start the server:
    
    go run main.go

Ingest metadata:


    curl -X POST http://localhost:8080/ingest -H "Authorization: Bearer your-secret-key" -H "Content-Type: application/json" -d '{
        "id": "12345",
        "source": "jira",
        "payload": {
            "title": "Fix login bug",
            "description": "Users are unable to log in due to a server error.",
            "status": "open",
            "assignee": "john.doe"
        },
        "metadata": {
            "created_at": "2024-07-20T12:34:56Z"
        }
    }'
    
Consume metadata:

    curl -X GET http://localhost:8080/consume -H "Authorization: Bearer your-secret-key"


# Flow
    # Ingest Metadata: Client sends a POST request to /ingest endpoint with metadata.
    # Pre-Process Metadata: The metadata is published to a Redis channel, where it is picked up by the pre-processing service.
    # Apply Operations: The pre-processing service fetches the operations configuration based on the source and applies the configured operations.
    # Store Processed Metadata: The processed metadata is stored in Redis.
    # Post-Process Metadata: The post-processing service retrieves the processed metadata from Redis, pushes it to downstream services, and deletes the entry from Redis.
    # Consume Metadata: Client sends a GET request to /consume endpoint to retrieve processed metadata.
