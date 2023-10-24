package main

import (
	"github.com/JuliyaMS/service-metrics-alerting/internal/config"
	"github.com/JuliyaMS/service-metrics-alerting/internal/logger"
	"github.com/JuliyaMS/service-metrics-alerting/internal/metrics"
	"github.com/JuliyaMS/service-metrics-alerting/internal/requests"
	"runtime"
	"time"
)

func main() {
	logger.NewLoggerAgent()
	config.GetAgentConfig()
	var rtm runtime.MemStats
	for {
		<-time.After(config.TimeInterval)
		metrics.ChangeMetrics(&rtm)
		if (metrics.PollCount % config.CountIteration) == 0 {
			err := requests.SendRequestJSON()
			if err != nil {
				panic(err)
			}
		}
	}

}
