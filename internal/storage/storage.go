package storage

import (
	"github.com/JuliyaMS/service-metrics-alerting/internal/metrics"
	"net/http"
)

var Storage MemStorage

type Repositories interface {
	Init()
	Add(t, name, val string) int
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

func (s MemStorage) Add(t, name, val string) int {
	if !metrics.CheckType(t) {
		return http.StatusBadRequest
	}
	if t == "counter" {
		if s.MetricsCounter.Add(name, val) {
			return http.StatusOK
		}
	}
	if s.MetricsGauge.Add(name, val) {
		return http.StatusOK
	}

	return http.StatusBadRequest
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
