package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/kotvytskyi/frontend/app/repository"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type CreatedResponse struct {
	Url string
}

type ShortService interface {
	Short(ctx context.Context, urlToShort string, short string) (string, error)
	CreateShortURL(r *http.Request, short string) string
}

type RestServer struct {
	Port         int
	ShortService ShortService
}

func NewRestServer(ctx context.Context) *RestServer {
	r := createRepository()
	a := createApi()

	return &RestServer{
		Port:         80,
		ShortService: NewShorter(r, a),
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
		respondWithError(w, http.StatusInternalServerError, "oops! an error occurred. please, try again later or contact the support")
		return
	}

	sURL := s.ShortService.CreateShortURL(r, short)
	respondCreated(w, sURL)
}

func (s *RestServer) shortProxyHandler(w http.ResponseWriter, r *http.Request) {

}

func respondCreated(w http.ResponseWriter, url string) {
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreatedResponse{Url: url})
}

func respondWithError(w http.ResponseWriter, status int, reason string) {
	w.WriteHeader(status)
	errResponse := &ErrorResponse{Error: reason}
	json.NewEncoder(w).Encode(errResponse)
}

func createRepository() UrlRepository {
	p := repository.NewParams(os.Getenv("MONGO"), os.Getenv("MONGO_USER"), os.Getenv("MONGO_PASS"))
	r, err := repository.NewShort(p)
	if err != nil {
		log.Fatalf("cannot create repository: %v", err)
	}

	return r
}

func createApi() ShortApi {
	api, err := NewShortApi(os.Getenv("SHORTSRV"))
	if err != nil {
		log.Fatalf("cannot create api: %v", err)
	}

	return api
}
