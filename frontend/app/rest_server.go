package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type CreateShortRequest struct {
	URL   string `json:"url"`
	Short string `json:"short"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type CreatedResponse struct {
	Url string
}

type ShortService interface {
	Create(ctx context.Context, urlToShort string, short string) error
	CreateShortURL(r *http.Request, short string) string
}

type RestServer struct {
	Port         int
	ShortService ShortService
}

func NewRestServer(ctx context.Context) (*RestServer, error) {
	service := &AppShortService{}

	mongoaddr := os.Getenv("MONGO")
	mongousr := os.Getenv("MONGO_USER")
	mongopass := os.Getenv("MONGO_PASS")
	repository, err := NewMongoShortRepository(MongoParams{Endpoint: fmt.Sprintf("mongodb://%s:%s@%s:27017", mongousr, mongopass, mongoaddr)})
	if err != nil {
		return nil, err
	}
	service.Repository = repository

	server := &RestServer{
		Port:         getPort(),
		ShortService: service,
	}

	return server, nil
}

func (s *RestServer) Run(ctx context.Context) error {
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

func (s *RestServer) router() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/shorts", s.createShortHandler).Methods("POST")
	r.HandleFunc("/short/{short}", s.createShortHandler).Methods("GET")

	return r
}

func (s *RestServer) createShortHandler(w http.ResponseWriter, r *http.Request) {
	request := &CreateShortRequest{}
	json.NewDecoder(r.Body).Decode(request)

	w.Header().Add("Content-Type", "application/json")

	if request.URL == "" {
		RespondWithError(w, http.StatusBadRequest, "url field is missing.")
		return
	}

	_, err := url.ParseRequestURI(request.URL)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "url field is malformed.")
		return
	}

	err = s.ShortService.Create(context.Background(), request.URL, request.Short)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "oops! an error occurred. please, try again later or contact the support")
		return
	}

	shortURL := s.ShortService.CreateShortURL(r, request.Short)

	RespondCreated(w, shortURL)
}

func (s *RestServer) shortProxyHandler(w http.ResponseWriter, r *http.Request) {

}

func RespondCreated(w http.ResponseWriter, url string) {
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreatedResponse{Url: url})
}

func RespondWithError(w http.ResponseWriter, status int, reason string) {
	w.WriteHeader(status)
	errResponse := &ErrorResponse{Error: reason}
	json.NewEncoder(w).Encode(errResponse)
}

func getPort() int {
	pEnv := os.Getenv("REST_PORT")
	if pEnv == "" {
		pEnv = "8080"
	}

	port, err := strconv.ParseInt(pEnv, 10, 64)
	if err != nil {
		panic(err)
	}

	return int(port)
}
