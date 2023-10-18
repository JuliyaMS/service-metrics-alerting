package headers

import (
	"fmt"
	"github.com/JuliyaMS/service-metrics-alerting/internal/checks"
	"github.com/JuliyaMS/service-metrics-alerting/internal/storage"
	"github.com/go-chi/chi/v5"
	"html/template"
	"net/http"
)

var memStor = storage.MemStorage{}

func requestValue(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")
	metricValue := chi.URLParam(r, "value")

	w.WriteHeader(memStor.Add(metricType, metricName, metricValue))
	memStor.Print()

}

func requestName(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")

	if checks.CheckDigit(metricType, metricName) || checks.CheckType(metricType) {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

}

func requestGetName(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")

	value := memStor.Get(metricType, metricName)
	fmt.Println(value)
	if value != "-1" {
		w.Write([]byte(value))
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusNotFound)

}

func requestGetAll(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	dataCounter, dataGauge := memStor.GetHTMLStructs()
	tmpl, _ := template.ParseFiles("../../data/html/index.html")
	tmpl.Execute(w, dataGauge)
	tmpl.Execute(w, dataCounter)
	w.WriteHeader(http.StatusOK)

}

func requestType(w http.ResponseWriter, r *http.Request) {

	metricType := chi.URLParam(r, "type")
	if checks.CheckType(metricType) {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

}

func requestEmpty(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusBadRequest)
}

func routePost(r *chi.Mux) {

	r.Route("/update", func(r chi.Router) {
		r.Post("/", requestEmpty)
		r.Route("/{type}", func(r chi.Router) {
			r.Post("/", requestType)
			r.Route("/{name}", func(r chi.Router) {
				r.Post("/", requestName)
				r.Route("/{value}", func(r chi.Router) {
					r.Post("/", requestValue)
				})
			})
		})
	})
}

func routeGet(r *chi.Mux) {
	r.Route("/value", func(r chi.Router) {
		r.Get("/", requestEmpty)
		r.Route("/{type}", func(r chi.Router) {
			r.Get("/", requestType)
			r.Route("/{name}", func(r chi.Router) {
				r.Get("/", requestGetName)
			})
		})
	})
}

func Router() *chi.Mux {
	memStor.Init()

	r := chi.NewRouter()
	routePost(r)
	routeGet(r)
	r.Get("/", requestGetAll)
	return r
}
