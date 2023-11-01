package file

import (
	"encoding/json"
	"github.com/JuliyaMS/service-metrics-alerting/internal/storage"
	"os"
	"sync"
)

var fileMutex sync.Mutex

func WriteToFile(fileName string) error {

	fileMutex.Lock()
	defer fileMutex.Unlock()

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	if err = encoder.Encode(storage.Storage); err != nil {
		return err
	}
	file.Close()
	return nil
}

func ReadFromFile(fileName string) error {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	file, err := os.OpenFile(fileName, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	encoder := json.NewDecoder(file)
	if err = encoder.Decode(&storage.Storage); err != nil {
		return err
	}
	file.Close()
	return nil
}
