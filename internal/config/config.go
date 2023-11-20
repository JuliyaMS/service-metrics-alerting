package config

import (
	"flag"
	"fmt"
	"github.com/JuliyaMS/service-metrics-alerting/internal/logger"
	"os"
	"strconv"
	"time"
)

var FlagRunSerAddr string
var FlagRunAgAddr string
var HashKeyAgent string
var HashKeyServer string
var TimeInterval time.Duration
var TimeInterval2 time.Duration
var StoreInterval time.Duration
var reportInterval int
var pollInterval int
var saveInterval int
var FileStoragePath string
var Restore bool
var DatabaseDsn string

func getEnvConfigServer() {
	if envRunSerAddr := os.Getenv("ADDRESS"); envRunSerAddr != "" {
		FlagRunSerAddr = envRunSerAddr
	}
	if envSaveInterval := os.Getenv("STORE_INTERVAL"); envSaveInterval != "" {
		if inter, err := strconv.Atoi(envSaveInterval); err == nil {
			saveInterval = inter
			fmt.Println(inter)
		}

	}
	if envPath := os.Getenv("FILE_STORAGE_PATH"); envPath != "" {
		FileStoragePath = envPath
	}
	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		if fl, err := strconv.ParseBool(envRestore); err == nil {
			Restore = fl
		}

	}
	if BDAddr := os.Getenv("DATABASE_DSN"); BDAddr != "" {
		DatabaseDsn = BDAddr
	}

	if HashKey := os.Getenv("KEY"); HashKey != "" {
		HashKeyServer = HashKey
	}
}

func getEnvConfigAgent() {
	if envRunAgAddr := os.Getenv("ADDRESS"); envRunAgAddr != "" {
		FlagRunAgAddr = envRunAgAddr
	}
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		if inter, err := strconv.Atoi(envReportInterval); err == nil {
			reportInterval = inter
		}

	}
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		if poll, err := strconv.Atoi(envPollInterval); err == nil {
			pollInterval = poll
		}
	}
	if HashKey := os.Getenv("KEY"); HashKey != "" {
		HashKeyAgent = HashKey
	}
}

func GetAgentConfig() {

	logger.Logger.Infow("Parse Agent config")

	flag.StringVar(&FlagRunAgAddr, "a", ":8080", "address and port to run server")
	flag.IntVar(&reportInterval, "r", 10, "time interval for generate metrics")
	flag.IntVar(&pollInterval, "p", 2, "time interval for send request to server")
	flag.StringVar(&HashKeyAgent, "k", "", "secret hash key for agent")

	flag.Parse()
	getEnvConfigAgent()

	TimeInterval = time.Duration(pollInterval) * time.Second
	TimeInterval2 = time.Duration(reportInterval) * time.Second
}

func GetServerConfig() {

	logger.Logger.Infow("Parse Server config")

	flag.StringVar(&FlagRunSerAddr, "a", ":8080", "address and port to run server")
	flag.IntVar(&saveInterval, "i", 300, "time interval to save metrics in file")
	flag.StringVar(&FileStoragePath, "f", "/tmp/metrics-db.json", "path to save file")
	flag.BoolVar(&Restore, "r", true, "restore data from file or not")
	flag.StringVar(&DatabaseDsn, "d", "", "database address")
	flag.StringVar(&HashKeyServer, "k", "", "secret hash key for server")

	flag.Parse()
	getEnvConfigServer()

	StoreInterval = time.Duration(saveInterval) * time.Second
}
