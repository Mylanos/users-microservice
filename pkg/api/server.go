package api

import (
	"log"
	"net/http"
	"time"
	"users-backend/pkg/config"
	"users-backend/pkg/services"
)

type APIResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Error     *APIError   `json:"error,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

type APIServer struct {
	listenAddr   string
	service      services.UserService
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type apiHandler func(w http.ResponseWriter, r *http.Request) error

func NewAPIServer(listenAddr string, service services.UserService, cfg *config.Config) *APIServer {
	return &APIServer{listenAddr: listenAddr, service: service, ReadTimeout: cfg.ReadTimeout, WriteTimeout: cfg.WriteTimeout, IdleTimeout: cfg.IdleTimeout}
}

func MakeHTTPHandleFunc(f apiHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		if err := f(w, r); err != nil {
			logError(r, err, time.Since(start))
			apiError := TranslateToAPIError(err)
			ConstructResponseWithError(w, apiError)
		} else {
			logSuccess(r, time.Since(start))
		}
	}
}

func methodCheckMiddleware(allowedMethod string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != allowedMethod {
			apiError := NewAPIError(http.StatusMethodNotAllowed,
				"HTTP method "+r.Method+" is not allowed for this endpoint")
			ConstructResponseWithError(w, *apiError)
			return
		}
		next(w, r)
	}
}

func (s *APIServer) Router() http.Handler {
	router := http.NewServeMux()

	getUserHandler := methodCheckMiddleware("GET", MakeHTTPHandleFunc(s.HandleGetUser))
	createUserHandler := methodCheckMiddleware("POST", MakeHTTPHandleFunc(s.HandleCreateUser))

	router.Handle("GET /{id}", getUserHandler)
	router.Handle("POST /save", createUserHandler)

	return router
}

func (s *APIServer) NewServer() *http.Server {
	return &http.Server{
		Addr:         s.listenAddr,
		Handler:      s.Router(),
		ReadTimeout:  s.ReadTimeout,
		WriteTimeout: s.WriteTimeout,
		IdleTimeout:  s.IdleTimeout,
	}
}

func (s *APIServer) Run() error {
	server := s.NewServer()
	log.Printf("Listening on %s", s.listenAddr)
	return server.ListenAndServe()
}
