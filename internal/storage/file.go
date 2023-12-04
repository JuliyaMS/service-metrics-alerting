package storage

import (
	"encoding/json"
	"github.com/JuliyaMS/service-metrics-alerting/internal/config"
	"github.com/JuliyaMS/service-metrics-alerting/internal/logger"
	"os"
	"sync"
	"time"
)

var fileMutex sync.Mutex
var DBMutex sync.Mutex

func WriteToFile(fileName string, stor *Repositories) error {

	fileMutex.Lock()
	defer fileMutex.Unlock()

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	if err = encoder.Encode(stor); err != nil {
		return err
	}
	file.Close()
	return nil
}

func ReadFromFile(fileName string) (*MemStorage, error) {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	var stor MemStorage
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0666)
	if err != nil {
		return new(MemStorage), err
	}
	encoder := json.NewDecoder(file)
	if err = encoder.Decode(&stor); err != nil {
		return new(MemStorage), err
	}
	file.Close()
	return &stor, nil
}

func SaveToFile(storage *Repositories) {
	for {
		<-time.After(config.StoreInterval)
		logger.Logger.Info("Write data to file:", config.FileStoragePath)
		if err := WriteToFile(config.FileStoragePath, storage); err != nil {
			logger.Logger.Error("Function WriteToFile return error:", err.Error())
		}
	}
}
