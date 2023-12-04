package storage

import (
	"github.com/JuliyaMS/service-metrics-alerting/internal/config"
	"github.com/JuliyaMS/service-metrics-alerting/internal/logger"
	"github.com/JuliyaMS/service-metrics-alerting/internal/metrics"
)

type Repositories interface {
	Init()
	Add(t, name, val string) error
	Get(tp, name string) string
	GetAll() (GaugeMetrics, CounterMetrics)
	CheckConnection() error
	AddAnyData(req []metrics.Metrics) error
	Close() error
}

func NewStorage() Repositories {
	logger.Logger.Info("Create new storage")

	if config.DatabaseDsn != "" {
		return NewConnectionDB()
	}

	var storage Repositories

	if config.FileStoragePath != "" {
		go SaveToFile(&storage)
	}

	if config.Restore && config.FileStoragePath != "" {
		logger.Logger.Info("restore data from file:", config.FileStoragePath)

		var err error
		storage, err = ReadFromFile(config.FileStoragePath)
		if err != nil {
			logger.Logger.Errorf(err.Error(), "can't read data from file:", config.FileStoragePath)
		}

		return storage
	}

	storage = new(MemStorage)
	return storage
}
