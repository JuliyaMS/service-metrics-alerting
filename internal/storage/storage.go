package storage

import (
	"github.com/JuliyaMS/service-metrics-alerting/internal/metrics"
	"net/http"
)

type Repositories interface {
	Init()
	Add(t, name, val string) int
	Get(tp, name string) string
	GetAll() (metrics.GaugeMetrics, metrics.CounterMetrics)
}

type MemStorage struct {
	metricsGauge   metrics.GaugeMetrics
	metricsCounter metrics.CounterMetrics
}

func (s *MemStorage) Init() {
	s.metricsGauge.Init()
	s.metricsCounter.Init()
}

func (s MemStorage) Add(t, name, val string) int {
	if !metrics.CheckType(t) {
		return http.StatusBadRequest
	}
	if t == "counter" {
		if s.metricsCounter.Add(name, val) {
			return http.StatusOK
		}
	}
	if s.metricsGauge.Add(name, val) {
		return http.StatusOK
	}

	return http.StatusBadRequest
}

func (s MemStorage) Get(tp, name string) string {
	if metrics.CheckType(tp) {
		if tp == "gauge" {
			value := s.metricsGauge.Get(name)
			return value
		}
		value := s.metricsCounter.Get(name)
		return value

	}
	return "-1"
}

func (s *MemStorage) GetAll() (metrics.GaugeMetrics, metrics.CounterMetrics) {
	return s.metricsGauge, s.metricsCounter
}
