package main

import (
	"context"
	"log"
	"os"

	"github.com/kotvytskyi/frontend/pkg/config"
	"github.com/kotvytskyi/frontend/pkg/server"
)

func main() {
	mongoCfg := config.MongoConfig{
		Address:  os.Getenv("MONGO"),
		User:     os.Getenv("MONGO_USER"),
		Password: os.Getenv("MONGO_PASS"),
	}

	shortCgf := config.ShortServerConfig{
		Address: os.Getenv("SHORTSRV"),
	}

	srv := server.NewRestServer(context.Background(), mongoCfg, shortCgf)

	err := srv.Run(context.Background())
	if err != nil {
		log.Fatalf("An error occurred in http server: %v", err)
	}
}
