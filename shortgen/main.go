package main

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"github.com/kotvytskyi/shortgen/app"
)

func main() {
	godotenv.Load()

	generatorMongo, err := app.NewMongo(app.MongoConfig{
		Endpoint: os.Getenv("MONGO"),
	})
	if err != nil {
		panic(err)
	}

	scheduler := app.NewScheduler(generatorMongo)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		err = scheduler.Schedule(ctx)
		if err != nil {
			panic(err)
		}
	}()
}
