package metrics

import (
	"math/rand"
	"runtime"
)

var PollCount = int64(0)

var GaugeAgent = GaugeMetrics{metrics: map[string]float64{"Alloc": 0, "BuckHashSys": 0, "Frees": 0,
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
	GaugeAgent.metrics["Alloc"] = float64(rtm.Alloc)
	GaugeAgent.metrics["BuckHashSys"] = float64(rtm.BuckHashSys)
	GaugeAgent.metrics["Frees"] = float64(rtm.Frees)
	GaugeAgent.metrics["GCCPUFraction"] = rtm.GCCPUFraction
	GaugeAgent.metrics["GCSys"] = float64(rtm.GCSys)
	GaugeAgent.metrics["HeapAlloc"] = float64(rtm.HeapAlloc)
	GaugeAgent.metrics["HeapIdle"] = float64(rtm.HeapIdle)
	GaugeAgent.metrics["HeapInuse"] = float64(rtm.HeapInuse)
	GaugeAgent.metrics["HeapObjects"] = float64(rtm.HeapObjects)
	GaugeAgent.metrics["HeapReleased"] = float64(rtm.HeapReleased)
	GaugeAgent.metrics["HeapSys"] = float64(rtm.HeapSys)
	GaugeAgent.metrics["LastGC"] = float64(rtm.LastGC)
	GaugeAgent.metrics["Lookups"] = float64(rtm.Lookups)
	GaugeAgent.metrics["MCacheInuse"] = float64(rtm.MCacheInuse)
	GaugeAgent.metrics["MSpanInuse"] = float64(rtm.MSpanInuse)
	GaugeAgent.metrics["MSpanSys"] = float64(rtm.MSpanSys)
	GaugeAgent.metrics["Mallocs"] = float64(rtm.Mallocs)
	GaugeAgent.metrics["NextGC"] = float64(rtm.NextGC)
	GaugeAgent.metrics["NumForcedGC"] = float64(rtm.NumForcedGC)
	GaugeAgent.metrics["NumGC"] = float64(rtm.NumGC)
	GaugeAgent.metrics["OtherSys"] = float64(rtm.OtherSys)
	GaugeAgent.metrics["StackInuse"] = float64(rtm.StackInuse)
	GaugeAgent.metrics["StackSys"] = float64(rtm.StackSys)
	GaugeAgent.metrics["Sys"] = float64(rtm.Sys)
	GaugeAgent.metrics["TotalAlloc"] = float64(rtm.TotalAlloc)
	GaugeAgent.metrics["RandomValue"] = randomValue()
	PollCount += 1

}
