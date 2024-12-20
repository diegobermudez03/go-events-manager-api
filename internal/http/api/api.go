package api

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/diegobermudez03/go-events-manager-api/internal/config"
	"github.com/diegobermudez03/go-events-manager-api/internal/http/handlers"
	"github.com/diegobermudez03/go-events-manager-api/internal/http/middlewares"
	"github.com/diegobermudez03/go-events-manager-api/pkg/app"
	"github.com/diegobermudez03/go-events-manager-api/pkg/domain"
	"github.com/diegobermudez03/go-events-manager-api/pkg/storage"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type APIServer struct {
	address 	string
	storage 	*storage.Storage
	server 		*http.Server
	config 		*config.Config
}

func NewAPIServer(address string, storage *storage.Storage, config *config.Config) *APIServer {
	return &APIServer{
		address: 	address,
		storage: 	storage,
		config:		config,
	}
}

func (s *APIServer) Run() error {
	router := chi.NewMux()

	//CORS, this basic CORS has no filters, accepts all domains, methods, and ages
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	  }))

	//	LOG MIDDLEWARE
	//router.Use(logMiddleware)

	// A good base middleware stack
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	r := chi.NewMux()
	router.Mount("/v1", r)

	//	inject dependencies and suscribe routes
	initializer := s.injectDependencies(r)
	if err := initializer.RegisterRoles(); err != nil{
		return fmt.Errorf("unable to initialize app %s", err.Error())
	}

	//	health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	server := http.Server{
		Handler: router,
		Addr: s.address,
	}
	s.server = &server
	return server.ListenAndServe()
}

func (s *APIServer) Shutdown() error {
	return s.server.Shutdown(context.TODO())
}


func (s *APIServer) injectDependencies(router *chi.Mux) domain.InitializeSvc{
	// create middlewares
	middlerares := middlewares.NewMiddlewares(s.config.AuthConfig.JWTSecret)

	//create services
	authService := app.NewAuthService(
		s.storage.AuthRepo,
		s.storage.UsersRepo, 
		s.storage.SessionsRepo, 
		s.config.AuthConfig.SecondsLife, 
		s.config.AuthConfig.AccessTokenExpiration,
		s.config.AuthConfig.JWTSecret,
	)

	eventsService := app.NewEventsService()

	//create handlers
	authHandler := handlers.NewAuthHandler(authService)
	eventsHandler := handlers.NewEventsHandler(eventsService, middlerares)

	//mount routes
	authHandler.MountRoutes(router)
	eventsHandler.MountRoutes(router)

	return app.NewInitializeService(s.storage.RolesRepo)
}

func logMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%s request on %s from %s", r.Method, r.URL, r.Host)
			next.ServeHTTP(w, r)
		},
	)
}