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
                PreOperations:  []string{"add_timestamp"},
                PostOperations: []string{"add_processing_info"},
            },
            {
                Source:     "github",
                PreOperations:  []string{"add_timestamp"},
                PostOperations: []string{"add_processing_info"},
            },
            {
                Source:     "serviceX",
                PreOperations:  []string{"add_timestamp"},
                PostOperations: []string{"add_processing_info"},
            },
        }
    
        for _, config := range configs {
            utils.SetOperationsConfig(redisClient, config)
        }
    }

    <img width="606" alt="image" src="https://github.com/user-attachments/assets/ee391fd4-e6dc-420b-a044-d9c3778ae3f1">


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
    1. In-memory DB (Redis) stores partner-specific operations details for pre and post transformations.
    2. Client Systems send metadata change events to the HTTP Server.
    3. The HTTP Server authenticates the requests using Auth and forwards them to the Message Broker (Kafka), specifically to the pre_process_meta_data_events topic.
    4. Pre-Ingest Transformer (Go Routine) subscribes to the Channel (Kafka topic) and preprocesses the metadata events by reading partner configurations from Redis.
    5. Preprocessed events are stored in Redis (Meta data store DynamoBd) under processed:<id>.
    6. Post-Consume Transformations (Go Routine) functions read the processed data from Redis (Meta data store DynamoBd), apply post-transformations using configurations from Redis (for that partner), and publish events to the post_process_meta_data_events Kafka topic.
    7. Downstream Systems (Workers) consume events from Kafka and perform required actions, such as enforcing data access security.
    8. Observability tools like New Relic and Grafana monitor the entire workflow, providing real-time logging and performance metrics.
