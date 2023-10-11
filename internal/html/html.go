package html

import (
	"github.com/JuliyaMS/service-metrics-alerting/internal/metrics"
)

type DataGauge struct {
	NameMetric string
	Items      []metrics.Gauge
}

type DataCounter struct {
	NameMetric string
	Items      []metrics.Counter
}

func getHTMLStructGauge(metricsGauge metrics.GaugeMetrics) DataGauge {
	var htmlGauge []metrics.Gauge
	var htmlData DataGauge
	for key, value := range metricsGauge.ReturnValues() {
		htmlGauge = append(htmlGauge, metrics.Gauge{Name: key, Value: value})
	}
	htmlData.Items = htmlGauge
	htmlData.NameMetric = "Gauge metrics"
	return htmlData
}

func getHTMLStructCounter(metricsCounter metrics.CounterMetrics) DataCounter {
	var htmlCounter []metrics.Counter
	var htmlData DataCounter
	for key, value := range metricsCounter.ReturnValues() {
		htmlCounter = append(htmlCounter, metrics.Counter{Name: key, Value: value})
	}
	htmlData.Items = htmlCounter
	htmlData.NameMetric = "Counter metrics"
	return htmlData
}

func GetHTMLStructs(MetricsGauge metrics.GaugeMetrics, MetricsCounter metrics.CounterMetrics) (DataGauge, DataCounter) {
	return getHTMLStructGauge(MetricsGauge), getHTMLStructCounter(MetricsCounter)
}
