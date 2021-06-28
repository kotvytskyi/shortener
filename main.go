package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/kotvytskyi/shortener/generator"
	"github.com/kotvytskyi/shortener/server"
)

func main() {
	godotenv.Load()

	generatorMongo, err := generator.NewMongo(generator.MongoConfig{
		Endpoint: os.Getenv("MONGO"),
	})
	if err != nil {
		panic(err)
	}

	scheduler := generator.NewScheduler(generatorMongo)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		err = scheduler.Schedule(ctx)
		if err != nil {
			panic(err)
		}
	}()

	serverMongo, err := server.NewMongo(server.MongoConfig{
		Endpoint: os.Getenv("MONGO"),
	})
	if err != nil {
		panic(err)
	}

	httpServer := server.HttpServer{
		Port:        getPort(),
		DataService: serverMongo,
	}

	err = httpServer.Run(ctx)
	if err != nil {
		log.Printf("Http server was terminated: %v", err)
	}
}

func getPort() int {
	pEnv := os.Getenv("REST_PORT")
	if pEnv == "" {
		pEnv = "8080"
	}

	port, err := strconv.ParseInt(pEnv, 10, 64)
	if err != nil {
		panic(err)
	}

	return int(port)
}
