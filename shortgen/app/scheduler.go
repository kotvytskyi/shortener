package app

import (
	"context"
	"log"
	"time"
)

type Key struct {
	Value string `json:"val" bson:"val"`
}

type DataService interface {
	create(ctx context.Context, key Key) error
}

type Scheduler struct {
	dataService DataService
}

func NewScheduler(dataService DataService) *Scheduler {
	scheduler := &Scheduler{}
	scheduler.dataService = dataService
	return scheduler
}

func (s *Scheduler) Schedule(ctx context.Context) error {
	log.Print("The scheduler has been started.")

	ticker := time.NewTicker(time.Second * 5)

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			log.Print("The scheduler has been stopped.")
			return nil
		case <-ticker.C:
			key := generate(6)
			err := s.dataService.create(ctx, Key{
				Value: key,
			})

			if err != nil {
				log.Printf("[ERROR] An error occurred while saving the key: %s", err)
			} else {
				log.Print("[INFO] The key was generated")
			}
		}
	}
}
