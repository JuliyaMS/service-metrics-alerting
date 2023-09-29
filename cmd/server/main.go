package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
)

type MemStorage struct {
	metrics map[string]float64
}

var storage = MemStorage{}

func paths(url string) []string {

	reg, _ := regexp.Compile("\\w+")
	paths := reg.FindAllString(url, -1)
	return paths
}

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
	case 4:
		if !checkType(p[1]) || !checkValue(p[3]) {
			return http.StatusBadRequest

		} else {
			value, err := strconv.ParseFloat(p[3], 64)
			if err == nil {
				if p[1] == "counter" {
					if _, ok := storage.metrics[p[2]]; ok {
						storage.metrics[p[2]] += value
					} else {
						storage.metrics[p[2]] = value
					}
				} else {
					storage.metrics[p[2]] = value
				}
				return http.StatusOK
			}
		}
	case 3:
		if _, err := strconv.Atoi(p[2]); err == nil || checkType(p[1]) {
			return http.StatusNotFound
		}
	case 2:
		if checkType(p[1]) {
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
