package main

import (
	"github.com/JuliyaMS/service-metrics-alerting/internal/config"
	"github.com/JuliyaMS/service-metrics-alerting/internal/handlers"
	"github.com/JuliyaMS/service-metrics-alerting/internal/logger"
	"net/http"
)

func main() {
	logger.NewLogger()
	config.GetServerConfig()
	r := handlers.NewRouter()

	logger.Logger.Infow("Start server", "addr", config.FlagRunSerAddr)

	if err := http.ListenAndServe(config.FlagRunSerAddr, r); err != nil {
		logger.Logger.Fatalf(err.Error(), "event", "start server")
	}
}
