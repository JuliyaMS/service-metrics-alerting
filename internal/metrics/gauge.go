package metrics

import (
	"fmt"
	"strconv"
)

type Gauge struct {
	Name  string
	Value float64
}

type GaugeMetrics struct {
	Metrics map[string]float64
}

func (g *GaugeMetrics) Init() {
	if g.Metrics == nil {
		g.Metrics = make(map[string]float64)
	}
}

func (g *GaugeMetrics) Add(name string, val string) bool {
	if value, err := strconv.ParseFloat(val, 64); err == nil {
		g.Metrics[name] = value
		g.Print()
		return true
	}
	return false
}

func (g *GaugeMetrics) Get(name string) string {

	if val, ok := g.Metrics[name]; ok {
		return strconv.FormatFloat(val, 'f', -1, 64)
	} else {
		return "-1"
	}
}

func (g *GaugeMetrics) Print() {
	fmt.Println(g.Metrics)
}

func (g *GaugeMetrics) ReturnValues() map[string]float64 {
	return g.Metrics
}
