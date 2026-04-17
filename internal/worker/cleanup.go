package worker

import (
	"log"
	"time"
)

type CleanupService interface {
	CleanupExpired() error
}

func StartCleanup(service CleanupService, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			if err := service.CleanupExpired(); err != nil {
				log.Printf("cleanup failed: %v", err)
			}
		}
	}()
}
