package metrics

import (
	"math/rand"
	"runtime"
)

var PollCount = int64(0)

var GaugeAgent = GaugeMetrics{Metrics: map[string]float64{"Alloc": 0, "BuckHashSys": 0, "Frees": 0,
	"GCCPUFraction": 0, "GCSys": 0, "HeapAlloc": 0, "HeapIdle": 0, "HeapInuse": 0,
	"HeapObjects": 0, "HeapReleased": 0, "HeapSys": 0, "LastGC": 0, "Lookups": 0,
	"MCacheInuse": 0, "MSpanInuse": 0, "MSpanSys": 0, "Mallocs": 0, "NextGC": 0,
	"NumForcedGC": 0, "NumGC": 0, "OtherSys": 0, "PauseTotalNs": 0, "StackInuse": 0,
	"StackSys": 0, "Sys": 0, "TotalAlloc": 0, "RandomValue": 0}}

func randomValue() float64 {
	val1 := float64(rand.Intn(1000))
	val2 := float64(rand.Intn(100000))
	return val1 + (val2-val1)*rand.Float64()
}

func ChangeMetrics(rtm *runtime.MemStats) {

	runtime.ReadMemStats(rtm)
	GaugeAgent.Metrics["Alloc"] = float64(rtm.Alloc)
	GaugeAgent.Metrics["BuckHashSys"] = float64(rtm.BuckHashSys)
	GaugeAgent.Metrics["Frees"] = float64(rtm.Frees)
	GaugeAgent.Metrics["GCCPUFraction"] = rtm.GCCPUFraction
	GaugeAgent.Metrics["GCSys"] = float64(rtm.GCSys)
	GaugeAgent.Metrics["HeapAlloc"] = float64(rtm.HeapAlloc)
	GaugeAgent.Metrics["HeapIdle"] = float64(rtm.HeapIdle)
	GaugeAgent.Metrics["HeapInuse"] = float64(rtm.HeapInuse)
	GaugeAgent.Metrics["HeapObjects"] = float64(rtm.HeapObjects)
	GaugeAgent.Metrics["HeapReleased"] = float64(rtm.HeapReleased)
	GaugeAgent.Metrics["HeapSys"] = float64(rtm.HeapSys)
	GaugeAgent.Metrics["LastGC"] = float64(rtm.LastGC)
	GaugeAgent.Metrics["Lookups"] = float64(rtm.Lookups)
	GaugeAgent.Metrics["MCacheInuse"] = float64(rtm.MCacheInuse)
	GaugeAgent.Metrics["MSpanInuse"] = float64(rtm.MSpanInuse)
	GaugeAgent.Metrics["MSpanSys"] = float64(rtm.MSpanSys)
	GaugeAgent.Metrics["Mallocs"] = float64(rtm.Mallocs)
	GaugeAgent.Metrics["NextGC"] = float64(rtm.NextGC)
	GaugeAgent.Metrics["NumForcedGC"] = float64(rtm.NumForcedGC)
	GaugeAgent.Metrics["NumGC"] = float64(rtm.NumGC)
	GaugeAgent.Metrics["OtherSys"] = float64(rtm.OtherSys)
	GaugeAgent.Metrics["StackInuse"] = float64(rtm.StackInuse)
	GaugeAgent.Metrics["StackSys"] = float64(rtm.StackSys)
	GaugeAgent.Metrics["Sys"] = float64(rtm.Sys)
	GaugeAgent.Metrics["TotalAlloc"] = float64(rtm.TotalAlloc)
	GaugeAgent.Metrics["RandomValue"] = randomValue()
	PollCount += 1

}
