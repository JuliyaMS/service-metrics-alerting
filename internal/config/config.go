package config

import (
	"errors"
	"flag"
	"fmt"
	"github.com/JuliyaMS/service-metrics-alerting/internal/logger"
	"os"
	"strconv"
	"strings"
	"time"
)

var FlagRunSerAddr string
var CountIteration int64
var FlagRunAgAddr string
var TimeInterval time.Duration
var reportInterval int
var pollInterval int

func getEnvConfig() {
	if envRunAgAddr := os.Getenv("ADDRESS"); envRunAgAddr != "" {
		FlagRunAgAddr = envRunAgAddr
	}
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		if inter, err := strconv.Atoi(string(envReportInterval[0])); err != nil {
			reportInterval = inter
		}

	}
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		if poll, err := strconv.Atoi(string(envPollInterval[0])); err != nil {
			pollInterval = poll
		}
	}
}

func checkFlagsServer() error {
	flags := os.Args[1:]
	if len(flags) > 0 {
		if len(flags) != 1 {
			return errors.New("incorrect count of command line arguments")
		} else {
			data := strings.Split(flags[0], "=")
			if data[0] != "-a" {
				return errors.New("incorrect flag's name")
			}
			if !checkFlagAddr(data[1]) {
				return errors.New("adress is not correct. Need address in a form host:port")
			}
		}

	}
	return nil
}

func GetAgentConfig() {

	flag.StringVar(&FlagRunAgAddr, "a", ":8080", "address and port to run server")
	flag.IntVar(&reportInterval, "r", 10, "time interval for generate metrics")
	flag.IntVar(&pollInterval, "p", 2, "time interval for send request to server")
	if err := checkFlagsAgent(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		flag.Usage()
		os.Exit(1)
	}

	flag.Parse()
	getEnvConfig()

	CountIteration = int64(reportInterval / pollInterval)
	TimeInterval = time.Duration(pollInterval) * time.Second
}

func GetServerConfig() {
	flag.StringVar(&FlagRunSerAddr, "a", ":8080", "address and port to run server")
	if err := checkFlagsServer(); err != nil {
		flag.Usage()
		logger.Logger.Fatalf(err.Error(), "event", "get server config")
	}
	if envRunAgAddr := os.Getenv("ADDRESS"); envRunAgAddr != "" {
		FlagRunSerAddr = envRunAgAddr
	}
	flag.Parse()
}
