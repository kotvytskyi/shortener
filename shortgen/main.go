package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/kotvytskyi/shortgen/app"
)

func main() {
	godotenv.Load()

	mongoaddr := os.Getenv("MONGO")
	mongousr := os.Getenv("MONGO_USER")
	mongopass := os.Getenv("MONGO_PASS")
	generatorMongo, err := app.NewMongo(app.MongoConfig{
		Endpoint: fmt.Sprintf("mongodb://%s:%s@%s:27017", mongousr, mongopass, mongoaddr),
	})
	if err != nil {
		panic(err)
	}

	scheduler := app.NewScheduler(generatorMongo)

	err = scheduler.Schedule(context.Background())
	if err != nil {
		panic(err)
	}
}
