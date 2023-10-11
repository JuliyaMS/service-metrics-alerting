package main

import (
	"github.com/JuliyaMS/service-metrics-alerting/internal/config"
	"github.com/JuliyaMS/service-metrics-alerting/internal/metrics"
	"github.com/JuliyaMS/service-metrics-alerting/internal/requests"
	"runtime"
	"time"
)

func main() {
	config.GetAgentConfig()
	var rtm runtime.MemStats
	for {
		<-time.After(config.TimeInterval)
		metrics.ChangeMetrics(&rtm)
		if metrics.PollCounter == config.CountIteration {
			err := requests.SendRequest()
			if err != nil {
				panic(err)
			}
		}
	}

}
