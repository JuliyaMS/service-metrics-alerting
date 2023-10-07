package main

import (
	"github.com/JuliyaMS/service-metrics-alerting/internal/headers"
	"net/http"
)

func main() {
	parseFlagsServer()
	r := headers.Router()
	if err := http.ListenAndServe(flagRunSerAddr, r); err != nil {
		panic(err)
	}
}
