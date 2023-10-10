package storage

import (
	"fmt"
	"github.com/JuliyaMS/service-metrics-alerting/internal/checks"
	"net/http"
	"strconv"
)

type repositories interface {
	Init()
	Add()
	Get()
	Print()
	GetHtmlStructs()
}

type gauge struct {
	Name  string
	Value float64
}

type counter struct {
	Name  string
	Value int64
}

type HTMLDataGauge struct {
	NameMetric string
	Items      []gauge
}

type HTMLDataCounter struct {
	NameMetric string
	Items      []counter
}

type MemStorage struct {
	metricsGauge   map[string]float64
	metricsCounter map[string]int64
}

func (s *MemStorage) Init() {
	if s.metricsGauge == nil {
		s.metricsGauge = make(map[string]float64)
		s.metricsCounter = make(map[string]int64)
	}
}

func (s MemStorage) addGauge(name string, val string) bool {
	if value, err := strconv.ParseFloat(val, 64); err == nil {
		s.metricsGauge[name] = value
		return true
	}
	return false
}

func (s MemStorage) addCounter(name string, val string) bool {
	if value, err := strconv.ParseInt(val, 10, 64); err == nil {
		s.metricsCounter[name] += value
		return true
	}
	return false
}

func (s MemStorage) Add(t string, name string, val string) int {
	if !checks.CheckType(t) {
		return http.StatusBadRequest
	}
	if t == "counter" {
		if s.addCounter(name, val) {
			return http.StatusOK
		}
	}
	if s.addGauge(name, val) {
		return http.StatusOK
	}

	return http.StatusBadRequest
}

func (s MemStorage) getGaugeMetric(name string) string {

	if val, ok := s.metricsGauge[name]; ok {
		return strconv.FormatFloat(val, 'f', -1, 64)
	} else {
		return "-1"
	}
}

func (s MemStorage) getCounterMetric(name string) string {

	if val, ok := s.metricsCounter[name]; ok {
		return strconv.FormatInt(val, 10)
	} else {
		return "-1"
	}
}

func (s MemStorage) Get(tp, name string) string {
	if checks.CheckType(tp) {
		if tp == "gauge" {
			value := s.getGaugeMetric(name)
			return value
		}
		value := s.getCounterMetric(name)
		return value

	}
	return "-1"
}

func (s MemStorage) getHTMLStructGauge() HTMLDataGauge {
	var htmlGauge []gauge
	var htmlData HTMLDataGauge
	for key, value := range s.metricsGauge {
		htmlGauge = append(htmlGauge, gauge{Name: key, Value: value})
	}
	htmlData.Items = htmlGauge
	htmlData.NameMetric = "Gauge metrics"
	return htmlData
}

func (s MemStorage) getHTMLStructCounter() HTMLDataCounter {
	var htmlCounter []counter
	var htmlData HTMLDataCounter
	for key, value := range s.metricsCounter {
		htmlCounter = append(htmlCounter, counter{Name: key, Value: value})
	}
	htmlData.Items = htmlCounter
	htmlData.NameMetric = "Counter metrics"
	return htmlData
}

func (s MemStorage) GetHTMLStructs() (HTMLDataCounter, HTMLDataGauge) {
	return s.getHTMLStructCounter(), s.getHTMLStructGauge()
}

func (s MemStorage) Print() {
	fmt.Println(s.metricsCounter)
	fmt.Println(s.metricsGauge)
}
