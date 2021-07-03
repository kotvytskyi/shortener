package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kotvytskyi/shortsrv/app"
)

func main() {
	godotenv.Load()

	mongoaddr := os.Getenv("MONGO")
	mongousr := os.Getenv("MONGO_USER")
	mongopass := os.Getenv("MONGO_PASS")
	serverMongo, err := app.NewMongo(app.MongoConfig{
		Endpoint: fmt.Sprintf("mongodb://%s:%s@%s:27017", mongousr, mongopass, mongoaddr),
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
