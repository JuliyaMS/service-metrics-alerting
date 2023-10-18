package handlers

import (
	"fmt"
	"github.com/JuliyaMS/service-metrics-alerting/internal/html"
	"github.com/JuliyaMS/service-metrics-alerting/internal/metrics"
	"github.com/JuliyaMS/service-metrics-alerting/internal/storage"
	"github.com/go-chi/chi/v5"
	"html/template"
	"net/http"
)

type Handlers struct {
	memStor storage.Repositories
}

func NewHandlers(stor storage.Repositories) *Handlers {
	return &Handlers{
		memStor: stor,
	}
}

func (h *Handlers) requestValue(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")
	metricValue := chi.URLParam(r, "value")

	w.WriteHeader(h.memStor.Add(metricType, metricName, metricValue))

}

func (h *Handlers) requestName(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")

	if metrics.CheckDigit(metricType, metricName) || metrics.CheckType(metricType) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (h *Handlers) requestGetName(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")

	value := h.memStor.Get(metricType, metricName)
	fmt.Println(value)
	if value != "-1" {
		w.Write([]byte(value))
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusNotFound)

}

func (h *Handlers) requestGetAll(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	dataGauge, dataCounter := html.GetHTMLStructs(h.memStor.GetAll())
	tmpl, _ := template.ParseFiles("../../data/html/index.html")
	tmpl.Execute(w, dataGauge)
	tmpl.Execute(w, dataCounter)
	w.WriteHeader(http.StatusOK)

}

func (h *Handlers) requestType(w http.ResponseWriter, r *http.Request) {

	metricType := chi.URLParam(r, "type")
	if metrics.CheckType(metricType) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (h *Handlers) requestEmpty(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusBadRequest)
}

func routePost(r *chi.Mux, h *Handlers) {

	r.Route("/update", func(r chi.Router) {
		r.Post("/", h.requestEmpty)
		r.Route("/{type}", func(r chi.Router) {
			r.Post("/", h.requestType)
			r.Route("/{name}", func(r chi.Router) {
				r.Post("/", h.requestName)
				r.Route("/{value}", func(r chi.Router) {
					r.Post("/", h.requestValue)
				})
			})
		})
	})
}

func routeGet(r *chi.Mux, h *Handlers) {
	r.Route("/value", func(r chi.Router) {
		r.Get("/", h.requestEmpty)
		r.Route("/{type}", func(r chi.Router) {
			r.Get("/", h.requestType)
			r.Route("/{name}", func(r chi.Router) {
				r.Get("/", h.requestGetName)
			})
		})
	})
}

func NewRouter() *chi.Mux {
	handlers := NewHandlers(&storage.MemStorage{})
	handlers.memStor.Init()

	r := chi.NewRouter()
	routePost(r, handlers)
	routeGet(r, handlers)
	r.Get("/", handlers.requestGetAll)
	return r
}
