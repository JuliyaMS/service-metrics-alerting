package main

import (
	"flag"
	"fmt"
	"github.com/JuliyaMS/service-metrics-alerting/internal/checks"
	"os"
)

var flagRunAgAddr string
var reportInterval int
var pollInterval int

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
}
