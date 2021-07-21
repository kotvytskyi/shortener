package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kotvytskyi/frontend/pkg/config"
	"github.com/kotvytskyi/frontend/pkg/server/controller"
	"github.com/kotvytskyi/frontend/pkg/server/middleware"
)

type RestServer struct {
	port       int
	controller *controller.Short
}

func NewRestServer(ctx context.Context, mongoCfg config.MongoConfig, shortCfg config.ShortServerConfig) *RestServer {
	controller, err := controller.NewShort(mongoCfg, shortCfg)
	if err != nil {
		panic(fmt.Sprintf("Cannot create the server. Error: %v", err))
	}

	return &RestServer{
		port:       80,
		controller: controller,
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
	r.HandleFunc("/short/{short}", s.shortProxyHandler).Methods("GET")

	return r
}

func (s *RestServer) shortProxyHandler(w http.ResponseWriter, r *http.Request) {

}
