package services

import (
	"log"
)

func pushToDownstream(metadata string) {
	log.Printf("Pushing metadata to downstream services: %s", metadata)
}
