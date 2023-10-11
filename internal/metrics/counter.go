package metrics

import (
	"fmt"
	"strconv"
)

type Counter struct {
	Name  string
	Value int64
}

type CounterMetrics struct {
	metrics map[string]int64
}

func (c *CounterMetrics) Init() {
	if c.metrics == nil {
		c.metrics = make(map[string]int64)
	}
}

func (c *CounterMetrics) Add(name string, val string) bool {
	if value, err := strconv.ParseInt(val, 10, 64); err == nil {
		c.metrics[name] += value
		return true
	}
	return false
}

func (c *CounterMetrics) Get(name string) string {

	if val, ok := c.metrics[name]; ok {
		return strconv.FormatInt(val, 10)
	} else {
		return "-1"
	}
}

func (c *CounterMetrics) Print() {
	fmt.Println(c.metrics)
}

func (c *CounterMetrics) ReturnValues() map[string]int64 {
	return c.metrics
}
