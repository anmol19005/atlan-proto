package services

import (
	"log"
)

func pushToDownstream(payload []byte) {
	log.Printf("Pushing payload to downstream services: %s", string(payload))
}
