package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kotvytskyi/frontend/pkg/mongo"
	"github.com/kotvytskyi/frontend/pkg/server/middleware"
	"github.com/kotvytskyi/frontend/pkg/shorter"
)

type errorResponse struct {
	Error string `json:"error"`
}

type createdResponse struct {
	Url string
}

type shortService interface {
	Short(ctx context.Context, urlToShort string, short string) (string, error)
	CreateShortURL(r *http.Request, short string) string
}

type RestServer struct {
	Port         int
	ShortService shortService
}

type MongoConfig struct {
	User     string
	Password string
	Address  string
}

type ShortServerConfig struct {
	Address string
}

type Config struct {
	Mongo MongoConfig
	Short ShortServerConfig
}

func NewRestServer(ctx context.Context, config Config) *RestServer {
	r := createRepository(config.Mongo)
	a := createApi(config.Short)

	return &RestServer{
		Port:         80,
		ShortService: shorter.NewShorter(r, a),
	}
}

func (s *RestServer) Run(ctx context.Context) error {
	log.Printf("[INFO] Starting REST server on port: %d", s.Port)

	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", s.Port),
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

	r.HandleFunc("/api/shorts", s.createShortHandler).Methods("POST")
	r.HandleFunc("/short/{short}", s.shortProxyHandler).Methods("GET")

	return r
}

func (s *RestServer) createShortHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json") // move to middleware

	req := &CreateShortRequest{}
	json.NewDecoder(r.Body).Decode(req)

	err := req.Validate()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	short, err := s.ShortService.Short(context.Background(), req.URL, req.Short)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error: , %v", err))
		return
	}

	sURL := s.ShortService.CreateShortURL(r, short)
	respondCreated(w, sURL)
}

func (s *RestServer) shortProxyHandler(w http.ResponseWriter, r *http.Request) {

}

func respondCreated(w http.ResponseWriter, url string) {
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdResponse{Url: url})
}

func respondWithError(w http.ResponseWriter, status int, reason string) {
	w.WriteHeader(status)
	errResponse := &errorResponse{Error: reason}
	json.NewEncoder(w).Encode(errResponse)
}

func createRepository(config MongoConfig) shorter.UrlRepository {
	p := mongo.NewParams(config.Address, config.User, config.Password)
	r, err := mongo.NewShort(p)
	if err != nil {
		log.Fatalf("cannot create repository: %v", err)
	}

	return r
}

func createApi(config ShortServerConfig) shorter.ShortApi {
	api, err := shorter.NewShortApi(config.Address)
	if err != nil {
		log.Fatalf("cannot create api: %v", err)
	}

	return api
}
