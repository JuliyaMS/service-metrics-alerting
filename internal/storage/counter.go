package storage

import (
	"fmt"
	"strconv"
)

type Counter struct {
	Name  string
	Value int64
}

type CounterMetrics struct {
	Metrics map[string]int64
}

func (c *CounterMetrics) Init() {
	if c.Metrics == nil {
		c.Metrics = make(map[string]int64)
	}
}

func (c *CounterMetrics) Close() {
	if c.Metrics != nil {
		clear(c.Metrics)
	}
}

func (c *CounterMetrics) Add(name string, val string) bool {
	if value, err := strconv.ParseInt(val, 10, 64); err == nil {
		c.Metrics[name] += value
		c.Print()
		return true
	}
	return false
}

func (c *CounterMetrics) Get(name string) string {

	if val, ok := c.Metrics[name]; ok {
		return strconv.FormatInt(val, 10)
	} else {
		return "-1"
	}
}

func (c *CounterMetrics) Print() {
	fmt.Println(c.Metrics)
}

func (c *CounterMetrics) ReturnValues() map[string]int64 {
	return c.Metrics
}
