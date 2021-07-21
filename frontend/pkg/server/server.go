package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kotvytskyi/frontend/pkg/config"
	"github.com/kotvytskyi/frontend/pkg/mongo"
	"github.com/kotvytskyi/frontend/pkg/server/controller"
	"github.com/kotvytskyi/frontend/pkg/server/middleware"
	"github.com/kotvytskyi/frontend/pkg/shorter"
)

type RestServer struct {
	port       int
	controller *controller.Short
}

func NewRestServer(ctx context.Context, mongoCfg config.MongoConfig, shortCfg config.ShortServerConfig) *RestServer {
	return &RestServer{
		port:       8081,
		controller: createShortController(mongoCfg, shortCfg),
	}
}

func (s *RestServer) Run(ctx context.Context) error {
	log.Printf("[INFO] Starting REST server on port: %d", s.port)

	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      s.router(),
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	go func() {
		<-ctx.Done()
		if e := srv.Close(); e != nil {
			log.Printf("[WARN] failed to close http server, %v", e)
		}
	}()

	return srv.ListenAndServe()
}

func (s *RestServer) router() *mux.Router {
	r := mux.NewRouter()
	r.Use(middleware.LoggingMiddleware)

	r.HandleFunc("/api/shorts", s.controller.CreateShort).Methods("POST")
	r.HandleFunc("/short/{short}", s.controller.ProxyShort).Methods("GET")

	return r
}

func createShortController(mongoCfg config.MongoConfig, shortCfg config.ShortServerConfig) *controller.Short {
	api, err := shorter.NewShortApi(shortCfg)
	if err != nil {
		panic(fmt.Sprintf("Cannot initialize the short api. %v", err))
	}

	mongo, err := mongo.NewShort(mongoCfg)
	if err != nil {
		panic(fmt.Sprintf("Cannot initialize the mongo. %v", err))
	}

	service := shorter.NewShorter(mongo, api)

	return controller.NewShort(service)
}
