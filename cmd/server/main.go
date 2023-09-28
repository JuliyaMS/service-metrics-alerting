package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type MemStorage struct {
	metrics map[string]float64
}

var storage = MemStorage{}

func checkType(value string) bool {
	types := []string{"gauge", "counter"}
	for _, tp := range types {
		if tp == value {
			return true
		}
	}
	return false
}

func checkValue(value string) bool {
	if _, err := strconv.ParseFloat(value, 64); err == nil {
		return true
	}
	return false
}

func control(p []string) int {
	switch count := len(p); count {
	case 5:
		if !checkType(p[2]) || !checkValue(p[4]) {
			return http.StatusBadRequest

		} else {
			value, err := strconv.ParseFloat(p[4], 64)
			if err == nil {
				if p[2] == "counter" {
					if _, ok := storage.metrics[p[3]]; ok {
						storage.metrics[p[3]] += value
					} else {
						storage.metrics[p[3]] = value
					}
				} else {
					storage.metrics[p[3]] = value
				}
				return http.StatusOK
			}
		}
	case 4:
		if _, err := strconv.Atoi(p[3]); err == nil {
			return http.StatusNotFound
		}
	case 3:
		if checkType(p[2]) {
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
	p := strings.Split(r.URL.Path, "/")
	w.WriteHeader(control(p))
	fmt.Println(storage.metrics)
}

func run() error {
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, request)
	return http.ListenAndServe(`:8080`, mux)
}

func main() {
	storage.metrics = make(map[string]float64)
	if err := run(); err != nil {
		panic(err)
	}
}
