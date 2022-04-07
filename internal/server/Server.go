package server

import (
	"errors"
	"github.com/arzzra/go_practicum/internal/server/storage"
)

type Server struct {
	Storage storage.MetricStorage
}

//func (s *Server) GetFromStorage(metricType metric.MetricType) {
//	s.Storage.
//}
//
//func (s *Server) SaveToStorage(metricType metric.MetricType) {
//	s.Storage.
//}

func MakeServer(storage storage.MetricStorage) (*Server, error) {
	if storage == nil {
		return nil, errors.New("storage is not initialized")
	}

	srv := &Server{
		Storage: storage,
	}

	return srv, nil
}
