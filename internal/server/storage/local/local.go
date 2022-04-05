package local

import (
	"fmt"
	metric "github.com/arzzra/go_practicum/internal/server/metric"
	"sync"
)

type MetricStorage struct {
	mutex   sync.RWMutex
	metrics map[string]metric.Metric
}

func MakeMetricStorage() *MetricStorage {
	var storage = MetricStorage{
		mutex:   sync.RWMutex{},
		metrics: make(map[string]metric.Metric),
	}
	return &storage
}

func (s *MetricStorage) SaveMetric(m metric.Metric) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.metrics[m.Name] = m
}

func (s *MetricStorage) UpdateMetric(m metric.Metric) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if dataM, ok := s.metrics[m.Name]; ok {
		if dataM.Type != m.Type {
			return fmt.Errorf("different types: %s and %s", dataM.Type, m.Type)
		}
		m.ValueGauge += dataM.ValueGauge
		m.ValueCounter += dataM.ValueCounter
	}

	s.metrics[m.Name] = m
	return nil
}

func (s *MetricStorage) GetMetricFromStorage(typeM metric.MetricType, name string) (*metric.Metric, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if m, ok := s.metrics[name]; ok && m.Type == typeM {
		return &m, nil
	}

	return nil, fmt.Errorf("metric not found")
}

func (s *MetricStorage) GetAllMetricFromStorage() (*[]metric.Metric, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	metrics := make([]metric.Metric, 0, len(s.metrics))
	for _, metric := range s.metrics {
		metrics = append(metrics, metric)
	}

	return &metrics, nil
}
