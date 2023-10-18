package storage

import (
	"github.com/JuliyaMS/service-metrics-alerting/cmd/server/checks"
	"net/http"
	"strconv"
)

type MemStorage struct {
	MetricsGauge   map[string]float64
	MetricsCounter map[string]int64
}

func (s *MemStorage) Init() {
	s.MetricsGauge = make(map[string]float64)
	s.MetricsCounter = make(map[string]int64)
}

func (s MemStorage) addGauge(name string, val string) bool {
	if value, err := strconv.ParseFloat(val, 64); err == nil {
		s.MetricsGauge[name] = value
		return true
	}
	return false
}

func (s MemStorage) addCounter(name string, val string) bool {
	if value, err := strconv.ParseInt(val, 10, 64); err == nil {
		s.MetricsCounter[name] += value
		return true
	}
	return false
}

func (s MemStorage) Add(t string, name string, val string) int {
	if !checks.CheckType(t) {
		return http.StatusBadRequest
	} else {
		if t == "counter" {
			if s.addCounter(name, val) {
				return http.StatusOK
			}
		} else {
			if s.addGauge(name, val) {
				return http.StatusOK
			}
		}
	}
	return http.StatusBadRequest
}
