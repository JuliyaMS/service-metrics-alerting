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
	logger.NewLogger()
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
			logger.Logger.Info("Open file")
			encode, err := file.NewStorageFileEncode(config.FileStoragePath)
			if err != nil {
				logger.Logger.Errorf(err.Error(), "Can't create NewStorageFileEncode")
				return
			}
			logger.Logger.Info("Writing...")
			err = encode.WriteToFile()
			if err != nil {
				logger.Logger.Errorf(err.Error(), "Can't write to file:", config.FileStoragePath)
				return
			}
			logger.Logger.Info("Close file")
			encode.Close()
		}
	}()

	waitGroup.Wait()

}
