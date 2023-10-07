package main

import (
	"flag"
	"fmt"
	"github.com/JuliyaMS/service-metrics-alerting/internal/checks"
	"os"
)

var flagRunSerAddr string

func parseFlagsServer() {
	flag.StringVar(&flagRunSerAddr, "a", ":8080", "address and port to run server")
	if err := checks.CheckFlagsServer(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		flag.Usage()
		os.Exit(1)
	}
	if envRunAgAddr := os.Getenv("ADDRESS"); envRunAgAddr != "" {
		flagRunSerAddr = envRunAgAddr
	}
	flag.Parse()
}
