package api

import (
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"auth/internal/api/handler"
	"auth/internal/api/middleware"
	"auth/internal/config"
	"auth/internal/service"
	"auth/internal/utils"
)

const (
	maxHeaderBytes = 1 << 20
	readTimeout    = 10 * time.Second
	writeTimeout   = 10 * time.Second
)

type Server struct {
	httpServer *http.Server
	router     *mux.Router
}

func NewServer(ctx utils.MyContext, config config.Config) *Server {
	router := mux.NewRouter()

	allowedOrigins := handlers.AllowedOrigins([]string{"http://localhost:3000"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})

	corsRouter := handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)(router)

	recoveredRouter := middleware.RecoveryMiddleware(ctx, corsRouter)

	return &Server{
		httpServer: &http.Server{
			Addr:           config.ServerPort,
			MaxHeaderBytes: maxHeaderBytes,
			ReadTimeout:    readTimeout,
			WriteTimeout:   writeTimeout,
			Handler:        recoveredRouter,
		},
		router: router,
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx utils.MyContext) error {
	return s.httpServer.Shutdown(ctx.Ctx)
}

func (s *Server) HandleAuth(ctx utils.MyContext, service service.AuthorizationService) {
	s.router.HandleFunc("/api/register/", handler.Register(ctx, service)).Methods(http.MethodPost)
	s.router.HandleFunc("/api/login/", handler.Login(ctx, service)).Methods(http.MethodPost)
	s.router.HandleFunc("/api/refresh/", handler.Refresh(ctx, service)).Methods(http.MethodPost)
}
