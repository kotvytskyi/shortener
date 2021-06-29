package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kotvytskyi/shortsrv/app"
)

func main() {
	godotenv.Load()
	serverMongo, err := app.NewMongo(app.MongoConfig{
		Endpoint: os.Getenv("MONGO"),
	})
	if err != nil {
		panic(err)
	}

	httpServer := app.HttpServer{
		Port:        8080,
		DataService: serverMongo,
	}

	err = httpServer.Run(context.Background())
	if err != nil {
		log.Printf("Http server was terminated: %v", err)
	}
}
