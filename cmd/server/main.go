package main

import (
	"github.com/JuliyaMS/service-metrics-alerting/internal/config"
	"github.com/JuliyaMS/service-metrics-alerting/internal/handlers"
	"net/http"
)

func main() {
	config.GetServerConfig()
	r := handlers.NewRouter()
	if err := http.ListenAndServe(config.FlagRunSerAddr, r); err != nil {
		panic(err)
	}
}
