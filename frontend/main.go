package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	"github.com/kotvytskyi/frontend/app"
)

func main() {
	godotenv.Load()

	srv := app.NewRestServer(context.Background())
	err := srv.Run(context.Background())
	if err != nil {
		log.Fatalf("An error occurred in http server: %v", err)
	}
}
