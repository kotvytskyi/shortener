package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	"github.com/kotvytskyi/frontend/pkg/server"
)

func main() {
	godotenv.Load()

	srv := server.NewRestServer(context.Background())
	err := srv.Run(context.Background())
	if err != nil {
		log.Fatalf("An error occurred in http server: %v", err)
	}
}
