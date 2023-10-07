package main

import (
	"flag"
	"fmt"
	"github.com/JuliyaMS/service-metrics-alerting/internal/checks"
	"os"
	"strconv"
)

var flagRunAgAddr string
var reportInterval int
var pollInterval int

func getEnvValues() {

	if envRunAgAddr := os.Getenv("ADDRESS"); envRunAgAddr != "" {
		flagRunAgAddr = envRunAgAddr
	}
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		if interval, err := strconv.Atoi(string(envReportInterval[0])); err != nil {
			reportInterval = interval
		}

	}
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		if poll, err := strconv.Atoi(string(envPollInterval[0])); err != nil {
			pollInterval = poll
		}
	}
}

func parseFlagsAgent() {

	flag.StringVar(&flagRunAgAddr, "a", ":8080", "address and port to run server")
	flag.IntVar(&reportInterval, "r", 10, "time interval for generate metrics")
	flag.IntVar(&pollInterval, "p", 2, "time interval for send request to server")
	if err := checks.CheckFlagsAgent(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		flag.Usage()
		os.Exit(1)
	}
	flag.Parse()
	getEnvValues()
}
