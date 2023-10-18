package checks

import "strconv"

var metricsGauge = []string{"Alloc", "BuckHashSys", "Frees",
	"GCCPUFraction", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse",
	"HeapObjects", "HeapReleased", "HeapSys", "LastGC", "Lookups",
	"MCacheInuse", "MSpanInuse", "MSpanSys", "Mallocs", "NextGC",
	"NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse",
	"StackSys", "Sys", "TotalAlloc", "RandomValue"}

func CheckName(name string) bool {
	for _, nm := range metricsGauge {
		if nm == name {
			return true
		}
	}
	return false
}
func CheckType(value string) bool {
	types := []string{"gauge", "counter"}
	for _, tp := range types {
		if tp == value {
			return true
		}
	}
	return false
}

func CheckDigit(t string, val string) bool {
	if t == "gauge" {
		if _, err := strconv.ParseFloat(val, 64); err == nil {
			return true
		}
	}
	if t == "counter" {
		if _, err := strconv.ParseInt(val, 10, 64); err == nil {
			return true
		}
	}
	return false
}
