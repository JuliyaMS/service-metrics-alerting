package main

import (
	"github.com/JuliyaMS/service-metrics-alerting/internal/headers"
	"net/http"
)

func main() {
	r := headers.Router()
	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}
