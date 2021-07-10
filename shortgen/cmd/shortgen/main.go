package main

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"github.com/kotvytskyi/shortgen/pkg/scheduler"
)

func main() {
	godotenv.Load()

	cfg := scheduler.Config{
		Mongo: scheduler.MongoConfig{
			Address:  os.Getenv("MONGO"),
			Username: os.Getenv("MONGO_USER"),
			Password: os.Getenv("MONGO_PASS"),
		},
	}

	scheduler, err := scheduler.NewScheduler(cfg)

	if err != nil {
		panic(err)
	}

	err = scheduler.Schedule(context.Background())
	if err != nil {
		panic(err)
	}
}
