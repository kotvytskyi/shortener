package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/kotvytskyi/shortsrv/pkg/server"
)

const PORT = 80

func main() {
	config := server.Config{
		Mongo: server.MongoConfig{
			Address:  os.Getenv("MONGO"),
			User:     os.Getenv("MONGO_USER"),
			Password: os.Getenv("MONGO_PASS"),
		},
		Port: PORT,
	}

	httpServer, err := server.NewServer(config)
	if err != nil {
		panic(fmt.Sprintf("Cannot initiate http server %v", err))
	}

	err = httpServer.Run(context.Background())
	if err != nil {
		log.Printf("Http server was terminated: %v", err)
	}
}
