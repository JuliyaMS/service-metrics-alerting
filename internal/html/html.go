package html

import (
	"github.com/JuliyaMS/service-metrics-alerting/internal/storage"
)

type DataGauge struct {
	NameMetric string
	Items      []storage.Gauge
}

type DataCounter struct {
	NameMetric string
	Items      []storage.Counter
}

func getHTMLStructGauge(metricsGauge storage.GaugeMetrics) DataGauge {
	var htmlGauge []storage.Gauge
	var htmlData DataGauge
	for key, value := range metricsGauge.ReturnValues() {
		htmlGauge = append(htmlGauge, storage.Gauge{Name: key, Value: value})
	}
	htmlData.Items = htmlGauge
	htmlData.NameMetric = "Gauge metrics"
	return htmlData
}

func getHTMLStructCounter(metricsCounter storage.CounterMetrics) DataCounter {
	var htmlCounter []storage.Counter
	var htmlData DataCounter
	for key, value := range metricsCounter.ReturnValues() {
		htmlCounter = append(htmlCounter, storage.Counter{Name: key, Value: value})
	}
	htmlData.Items = htmlCounter
	htmlData.NameMetric = "Counter metrics"
	return htmlData
}

func GetHTMLStructs(MetricsGauge storage.GaugeMetrics, MetricsCounter storage.CounterMetrics) (DataGauge, DataCounter) {
	return getHTMLStructGauge(MetricsGauge), getHTMLStructCounter(MetricsCounter)
}
