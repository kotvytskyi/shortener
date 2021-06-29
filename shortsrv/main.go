package main

import (
	"context"
	"log"
	"os"
	"strconv"

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
		Port:        getPort(),
		DataService: serverMongo,
	}

	err = httpServer.Run(context.Background())
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
