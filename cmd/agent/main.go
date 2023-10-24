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
	config.GetAgentConfig()
	logger.NewLogger()
	var rtm runtime.MemStats
	time.Sleep(5 * time.Second)
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
