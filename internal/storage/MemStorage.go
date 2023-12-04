package storage

import (
	"errors"
	"github.com/JuliyaMS/service-metrics-alerting/internal/metrics"
	"strconv"
)

type MemStorage struct {
	MetricsGauge   GaugeMetrics
	MetricsCounter CounterMetrics
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

func (s *MemStorage) GetAll() (GaugeMetrics, CounterMetrics) {
	return s.MetricsGauge, s.MetricsCounter
}

func (s *MemStorage) CheckConnection() error {
	if s.MetricsGauge.Metrics != nil && s.MetricsCounter.Metrics != nil {
		return nil
	}
	return errors.New("storage for metrics isn`t initialize")
}

func (s *MemStorage) AddAnyData(req []metrics.Metrics) error {

	for _, r := range req {
		if r.MType == "gauge" {
			value := strconv.FormatFloat(*r.Value, 'g', -1, 64)
			if err := s.Add(r.MType, r.ID, value); err != nil {
				return err
			}
		}
		if r.MType == "counter" {
			value := strconv.FormatInt(*r.Delta, 10)
			if err := s.Add(r.MType, r.ID, value); err != nil {
				return err
			}
		}

	}
	return nil
}

func (s MemStorage) Close() error {
	s.MetricsGauge.Close()
	s.MetricsCounter.Close()
	return nil
}
