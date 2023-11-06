package storage

import (
	"errors"
	"github.com/JuliyaMS/service-metrics-alerting/internal/metrics"
)

var Storage MemStorage

type Repositories interface {
	Init()
	Add(t, name, val string) error
	Get(tp, name string) string
	GetAll() (metrics.GaugeMetrics, metrics.CounterMetrics)
}

type MemStorage struct {
	MetricsGauge   metrics.GaugeMetrics
	MetricsCounter metrics.CounterMetrics
}

func (s *MemStorage) Init() {
	s.MetricsGauge.Init()
	s.MetricsCounter.Init()
}

func (s MemStorage) Add(t, name, val string) error {
	if !metrics.CheckType(t) {
		return errors.New("this type of metric doesn't exists")
	}
	if t == "counter" {
		if s.MetricsCounter.Add(name, val) {
			return nil
		}
	}
	if s.MetricsGauge.Add(name, val) {
		return nil
	}

	return errors.New("can't add metric")
}

func (s MemStorage) Get(tp, name string) string {
	if metrics.CheckType(tp) {
		if tp == "gauge" {
			value := s.MetricsGauge.Get(name)
			return value
		}
		value := s.MetricsCounter.Get(name)
		return value

	}
	return "-1"
}

func (s *MemStorage) GetAll() (metrics.GaugeMetrics, metrics.CounterMetrics) {
	return s.MetricsGauge, s.MetricsCounter
}
