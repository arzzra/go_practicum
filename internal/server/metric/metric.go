package metric

import (
	"fmt"
	"strconv"
)

type (
	MetricType string
	Gauge      float64
	Counter    int64
)

const (
	TypeGauge   MetricType = "gauge"
	TypeCounter MetricType = "counter"
)

type Metric struct {
	Name         string
	Type         MetricType
	ValueGauge   Gauge
	ValueCounter Counter
}

func (m MetricType) Validate() error {
	if m == TypeGauge || m == TypeCounter {
		return nil
	} else {
		return fmt.Errorf("unknown type of metric: %s", m)
	}
}

func MakeMetricStruct(name string, typeM MetricType, value string) (Metric, error) {
	switch typeM {
	case TypeGauge:
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return Metric{}, err
		}
		return Metric{
			Name:         name,
			Type:         typeM,
			ValueGauge:   Gauge(floatValue),
			ValueCounter: 0,
		}, nil
	case TypeCounter:
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return Metric{}, err
		}
		return Metric{
			Name:         name,
			Type:         TypeCounter,
			ValueGauge:   0,
			ValueCounter: Counter(intValue),
		}, nil
	default:
		return Metric{}, fmt.Errorf("unknown type of metric: %s", typeM)
	}
}

func (m *Metric) GetValueString() string {
	switch m.Type {
	case TypeGauge:
		return strconv.FormatFloat(float64(m.ValueGauge), 'f', -1, 64)
	case TypeCounter:
		return strconv.FormatInt(int64(m.ValueCounter), 10)
	default:
		return ""
	}
}
