package main

import (
	"fmt"
	"github.com/JuliyaMS/service-metrics-alerting/cmd/server/checks"
	"github.com/JuliyaMS/service-metrics-alerting/cmd/server/storage"
	"net/http"
	"regexp"
)

var memStor = storage.MemStorage{}

func paths(url string) []string {

	var reg = regexp.MustCompile(`\w+\.?\w*`)
	p := reg.FindAllString(url, -1)
	return p
}

func control(p []string) int {
	switch count := len(p); count {
	case 4:
		return memStor.Add(p[1], p[2], p[3])
	case 3:
		if checks.CheckDigit(p[1], p[2]) || checks.CheckType(p[1]) {
			return http.StatusNotFound
		}
	case 2:
		if checks.CheckType(p[1]) {
			return http.StatusNotFound
		}
	}

	return http.StatusBadRequest
}

func request(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	p := paths(r.URL.Path)
	fmt.Println(p)
	w.WriteHeader(control(p))
	fmt.Println(memStor.MetricsGauge)
	fmt.Println(memStor.MetricsCounter)
}

func run() error {
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, request)
	return http.ListenAndServe(`:8080`, mux)
}

func main() {
	memStor.Init()
	if err := run(); err != nil {
		panic(err)
	}
}
