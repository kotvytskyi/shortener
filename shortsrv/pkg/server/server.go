package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kotvytskyi/shortsrv/pkg/mongo"
)

type HttpServer struct {
	Port        int
	DataService DataService
}

type DataService interface {
	ReserveKey(ctx context.Context) (string, error)
}

type MongoConfig struct {
	Address  string
	User     string
	Password string
}

type Config struct {
	Mongo MongoConfig
	Port  int
}

func NewServer(config Config) (*HttpServer, error) {
	mongoCfg := mongo.Config{
		Address:  config.Mongo.Address,
		User:     config.Mongo.User,
		Password: config.Mongo.Password,
	}

	serverMongo, err := mongo.NewKeyRepository(mongoCfg)
	if err != nil {
		return nil, err
	}

	srv := &HttpServer{
		Port:        config.Port,
		DataService: serverMongo,
	}

	return srv, nil
}

func (s *HttpServer) Run(ctx context.Context) error {
	log.Printf("[INFO] Starting REST server on port: %d", s.Port)

	router := s.router()
	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", s.Port),
		Handler:      router,
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

func (s *HttpServer) router() *mux.Router {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	r.HandleFunc("/api/keys", s.reserveKeyHandler).Methods("POST")

	return r
}

func (s *HttpServer) reserveKeyHandler(rw http.ResponseWriter, h *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	key, err := s.DataService.ReserveKey(context.Background())
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(rw).Encode(struct{ Err string }{Err: err.Error()})
		return
	}

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(key)
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func NewLoggingResponseWriter(rw http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{rw, http.StatusOK}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s", r.Method, r.URL.Path)

		lrw := NewLoggingResponseWriter(rw)
		next.ServeHTTP(lrw, r)

		log.Printf("[%s] %s - %d", r.Method, r.URL.Path, lrw.statusCode)
	})
}
