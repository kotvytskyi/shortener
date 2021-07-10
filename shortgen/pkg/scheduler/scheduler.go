package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/kotvytskyi/shortgen/pkg/generator"
	"github.com/kotvytskyi/shortgen/pkg/mongo"
)

type MongoConfig struct {
	Address  string
	Username string
	Password string
}

type Config struct {
	Mongo MongoConfig
}

type Scheduler struct {
	generator *generator.Generator
}

func NewScheduler(config Config) (*Scheduler, error) {
	scheduler := &Scheduler{}

	repo, err := mongo.NewKeyRepository(mongo.Config{
		Username: config.Mongo.Username,
		Password: config.Mongo.Password,
		Address:  config.Mongo.Address,
	})

	if err != nil {
		return nil, err
	}

	scheduler.generator = generator.NewGenerator(repo)

	return scheduler, nil
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
			err := s.generator.Generate(ctx)

			if err != nil {
				log.Printf("[ERROR] An error occurred while saving the key: %s", err)
			} else {
				log.Print("[INFO] The key was generated")
			}
		}
	}
}
