package main

import (
	"context"
	"log"

	"github.com/kotvytskyi/frontend/app"
)

func main() {
	server, err := app.NewRestServer(context.Background())
	if err != nil {
		panic(err)
	}

	err = server.Run(context.Background())
	if err != nil {
		log.Fatalf("An error occurred in http server: %v", err)
	}
}
