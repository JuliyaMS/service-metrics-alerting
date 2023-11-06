package main

import (
	"fmt"
	"github.com/JuliyaMS/service-metrics-alerting/internal/agent"
	"github.com/JuliyaMS/service-metrics-alerting/internal/config"
	"github.com/JuliyaMS/service-metrics-alerting/internal/logger"
	"github.com/JuliyaMS/service-metrics-alerting/internal/metrics"
	"runtime"
	"time"
)

func main() {
	config.GetAgentConfig()

	requestURL := fmt.Sprintf("http://%s/update/", config.FlagRunAgAddr)
	agent := agent.NewAgent(requestURL, true)

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
				logger.Logger.Infow("Change metrics", "time", tm)
				metrics.ChangeMetrics(&rtm)
			case tm2 := <-ticker2.C:
				logger.Logger.Infow("Send metrics", "time", tm2)
				err := agent.SendRequestJSON()
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
