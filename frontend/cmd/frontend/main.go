package main

import (
	"context"
	"log"
	"os"

	"github.com/kotvytskyi/frontend/pkg/server"
)

func main() {
	config := server.Config{
		Mongo: server.MongoConfig{
			Address:  os.Getenv("MONGO"),
			User:     os.Getenv("MONGO_USER"),
			Password: os.Getenv("MONGO_PASS"),
		},
		Short: server.ShortServerConfig{
			Address: os.Getenv("SHORTSRV"),
		},
	}

	srv := server.NewRestServer(context.Background(), config)
	err := srv.Run(context.Background())
	if err != nil {
		log.Fatalf("An error occurred in http server: %v", err)
	}
}
