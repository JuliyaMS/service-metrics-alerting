package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

var metricsGauge = map[string]float64{"Alloc": 0, "BuckHashSys": 0, "Frees": 0,
	"GCCPUFraction": 0, "GCSys": 0, "HeapAlloc": 0, "HeapIdle": 0, "HeapInuse": 0,
	"HeapObjects": 0, "HeapReleased": 0, "HeapSys": 0, "LastGC": 0, "Lookups": 0,
	"MCacheInuse": 0, "MSpanInuse": 0, "MSpanSys": 0, "Mallocs": 0, "NextGC": 0,
	"NumForcedGC": 0, "NumGC": 0, "OtherSys": 0, "PauseTotalNs": 0, "StackInuse": 0,
	"StackSys": 0, "Sys": 0, "TotalAlloc": 0, "RandomValue": 0}

var PollCounter = uint64(0)

func randomValue() float64 {
	val1 := float64(rand.Intn(1000))
	val2 := float64(rand.Intn(100000))
	return val1 + (val2-val1)*rand.Float64()
}

func changeMetrics(rtm runtime.MemStats) {

	runtime.ReadMemStats(&rtm)
	metricsGauge["Alloc"] = float64(rtm.Alloc)
	metricsGauge["BuckHashSys"] = float64(rtm.BuckHashSys)
	metricsGauge["Frees"] = float64(rtm.Frees)
	metricsGauge["GCCPUFraction"] = rtm.GCCPUFraction
	metricsGauge["GCSys"] = float64(rtm.GCSys)
	metricsGauge["HeapAlloc"] = float64(rtm.HeapAlloc)
	metricsGauge["HeapIdle"] = float64(rtm.HeapIdle)
	metricsGauge["HeapInuse"] = float64(rtm.HeapInuse)
	metricsGauge["HeapObjects"] = float64(rtm.HeapObjects)
	metricsGauge["HeapReleased"] = float64(rtm.HeapReleased)
	metricsGauge["HeapSys"] = float64(rtm.HeapSys)
	metricsGauge["LastGC"] = float64(rtm.LastGC)
	metricsGauge["Lookups"] = float64(rtm.Lookups)
	metricsGauge["MCacheInuse"] = float64(rtm.MCacheInuse)
	metricsGauge["MSpanInuse"] = float64(rtm.MSpanInuse)
	metricsGauge["MSpanSys"] = float64(rtm.MSpanSys)
	metricsGauge["Mallocs"] = float64(rtm.Mallocs)
	metricsGauge["NextGC"] = float64(rtm.NextGC)
	metricsGauge["NumForcedGC"] = float64(rtm.NumForcedGC)
	metricsGauge["NumGC"] = float64(rtm.NumGC)
	metricsGauge["OtherSys"] = float64(rtm.OtherSys)
	metricsGauge["StackInuse"] = float64(rtm.StackInuse)
	metricsGauge["StackSys"] = float64(rtm.StackSys)
	metricsGauge["Sys"] = float64(rtm.Sys)
	metricsGauge["TotalAlloc"] = float64(rtm.TotalAlloc)
	metricsGauge["RandomValue"] = randomValue()
	PollCounter += 1

}

func sendRequest(serverPort int) error {
	for k, v := range metricsGauge {
		requestURL := fmt.Sprintf("http://localhost:%d/update/gauge/%s/%f", serverPort, k, v)
		_, err := http.Post(requestURL, "Content-Type: text/plain", nil)
		if err != nil {
			return err
		}
	}
	requestURL := fmt.Sprintf("http://localhost:%d/update/counter/PollCounter/%d", serverPort, PollCounter)
	_, err := http.Post(requestURL, "Content-Type: text/plain", nil)
	if err != nil {
		return err
	}
	time.Sleep(1 * time.Second)
	return nil
}

func main() {
	var rtm runtime.MemStats
	var duration = uint(2)
	var interval = time.Duration(duration) * time.Second
	for {
		<-time.After(interval)
		changeMetrics(rtm)
		if PollCounter%5 == 0 {
			err := sendRequest(8080)
			if err != nil {
				panic("Request failed")
			}
		}
	}

}
