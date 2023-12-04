package main

import (
	"fmt"
	"github.com/JuliyaMS/service-metrics-alerting/internal/agent"
	"github.com/JuliyaMS/service-metrics-alerting/internal/config"
	"github.com/JuliyaMS/service-metrics-alerting/internal/logger"
	"github.com/JuliyaMS/service-metrics-alerting/internal/storage"
	"runtime"
	"time"
)

func main() {
	config.GetAgentConfig()

	requestURL := fmt.Sprintf("http://%s/updates/", config.FlagRunAgAddr)
	client := agent.NewAgent(requestURL, true)

	var rtm runtime.MemStats
	PollTicker := time.NewTicker(config.PollInterval)
	ReportTicker := time.NewTicker(config.ReportInterval)

	tickerChan := make(chan bool)

	go func() {
		in := make(chan storage.GaugeMetrics, 2)
		for {
			select {
			case <-tickerChan:
				close(in)
				return
			case tm := <-PollTicker.C:
				logger.Logger.Infow("Change metrics", "time", tm)
				agent.ChangeMetrics(&rtm, in)
			case tm2 := <-ReportTicker.C:
				logger.Logger.Infow("Send metrics", "time", tm2)
				for i := 1; i <= config.RateLimit; i++ {
					go client.Worker(i, in)
				}
			}
		}
	}()

	time.Sleep(100 * time.Second)
	PollTicker.Stop()
	ReportTicker.Stop()
	tickerChan <- true

}
