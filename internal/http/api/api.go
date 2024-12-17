package api

import (
	"context"
	"net/http"

	"github.com/diegobermudez03/go-events-manager-api/internal/http/handlers"
	"github.com/diegobermudez03/go-events-manager-api/pkg/app"
	"github.com/diegobermudez03/go-events-manager-api/pkg/storage"
	"github.com/go-chi/chi/v5"
)

type APIServer struct {
	address 	string
	storage 	*storage.Storage
	server 		*http.Server
}

func NewAPIServer(address string, storage *storage.Storage) *APIServer {
	return &APIServer{
		address: address,
		storage: storage,
	}
}

func (s *APIServer) Run() error {
	router := chi.NewMux()
	r := chi.NewMux()
	router.Mount("/v1", r)

	//create services
	authService := app.NewAuthService(s.storage.UsersRepo)

	//create handlers
	authHandler := handlers.NewAuthHandler(authService)

	//mount routes
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	authHandler.MountRoutes(r)

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