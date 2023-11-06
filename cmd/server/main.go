package main

import (
	"github.com/JuliyaMS/service-metrics-alerting/internal/config"
	"github.com/JuliyaMS/service-metrics-alerting/internal/file"
	"github.com/JuliyaMS/service-metrics-alerting/internal/handlers"
	"github.com/JuliyaMS/service-metrics-alerting/internal/logger"
	"net/http"
	"sync"
	"time"
)

func main() {
	config.GetServerConfig()
	r := handlers.NewRouter()

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	go func() {
		defer waitGroup.Done()

		logger.Logger.Infow("Start server", "addr", config.FlagRunSerAddr)
		if err := http.ListenAndServe(config.FlagRunSerAddr, r); err != nil {
			logger.Logger.Fatalf(err.Error(), "event", "start server")
			return
		}

	}()

	waitGroup.Add(1)

	go func() {
		defer waitGroup.Done()
		for {
			<-time.After(config.StoreInterval)
			logger.Logger.Info("Write data to file:", config.FileStoragePath)
			if err := file.WriteToFile(config.FileStoragePath); err != nil {
				logger.Logger.Error("Function WriteToFile return error:", err.Error())
			}
		}
	}()

	waitGroup.Wait()

}
