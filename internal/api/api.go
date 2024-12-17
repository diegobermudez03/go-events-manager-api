package api

import (
	"github.com/diegobermudez03/go-events-manager-api/pkg/storage"
)

type APIServer struct {
	storage 	*storage.Storage
}

func NewAPIServer(storage *storage.Storage) *APIServer {
	return &APIServer{
		storage: storage,
	}
}

func (s *APIServer) Run() error {
	return nil
}