package storage

import (
	"github.com/arzzra/go_practicum/internal/server/metric"
)

type MetricStorage interface {
	SaveMetric(m metric.Metric)
	UpdateMetric(m metric.Metric) error
	GetMetricFromStorage(typeM metric.MetricType, name string) (*metric.Metric, error)
	GetAllMetricFromStorage() (*[]metric.Metric, error)
}
