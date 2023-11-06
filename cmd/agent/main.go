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
	ticker := time.NewTicker(config.TimeInterval)
	ticker2 := time.NewTicker(config.TimeInterval2)

	tickerChan := make(chan bool)

	go func() {
		for {
			select {
			case <-tickerChan:
				return
			case tm := <-ticker.C:
				logger.Agent.Infow("Change metrics", "time", tm)
				metrics.ChangeMetrics(&rtm)
			case tm2 := <-ticker2.C:
				logger.Agent.Infow("Send metrics", "time", tm2)
				err := requests.SendRequestJSON()
				if err != nil {
					panic(err)
				}
			}
		}
	}()

	time.Sleep(100 * time.Second)
	ticker.Stop()
	ticker2.Stop()
	tickerChan <- true

}
