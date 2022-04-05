package server

import (
	"errors"
	"github.com/arzzra/go_practicum/internal/server/storage/local"
)

type Server struct {
	Storage *local.MetricStorage
}

func MakeServer(storage *local.MetricStorage) (*Server, error) {
	if storage == nil {
		return nil, errors.New("storage is not initialized")
	}

	srv := &Server{
		Storage: storage,
	}

	return srv, nil
}
