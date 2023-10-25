package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/JuliyaMS/service-metrics-alerting/internal/html"
	"github.com/JuliyaMS/service-metrics-alerting/internal/logger"
	"github.com/JuliyaMS/service-metrics-alerting/internal/metrics"
	"github.com/JuliyaMS/service-metrics-alerting/internal/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"html/template"
	"net/http"
	"strconv"
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

func (h *Handlers) requestUpdate(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		logger.Logger.Infow("got request with bad method", zap.String("method", r.Method))
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	logger.Logger.Infow("decoding request")

	var req metrics.Metrics
	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(&req); err != nil {
		logger.Logger.Error("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if req.MType == "gauge" {
		logger.Logger.Infow("Update gauge metric", "name", req.ID, "value", *req.Value)
		w.WriteHeader(h.memStor.Add(req.MType, req.ID, strconv.FormatFloat(*req.Value, 'g', -1, 64)))
	}

	if req.MType == "counter" {
		logger.Logger.Infow("Update counter metric", "name", req.ID, "value", *req.Delta)
		w.WriteHeader(h.memStor.Add(req.MType, req.ID, strconv.FormatInt(*req.Delta, 10)))
		newDelta, err := strconv.ParseInt(h.memStor.Get("counter", req.ID), 10, 64)
		if err != nil {
			logger.Logger.Error("cannot write new Delta", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		req.Delta = &newDelta
	}

	w.Header().Set("Content-Type", "application/json")

	logger.Logger.Infow("Encode data for response")
	enc := json.NewEncoder(w)
	if err := enc.Encode(req); err != nil {
		logger.Logger.Error("error encoding response", zap.Error(err))
		return
	}
	logger.Logger.Infow("sending HTTP 200 response")
}

func (h *Handlers) requestGetValue(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		logger.Logger.Debug("got request with bad method", zap.String("method", r.Method))
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	logger.Logger.Infow("decoding request")
	var req metrics.Metrics
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		logger.Logger.Error("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if req.MType == "gauge" {
		logger.Logger.Infow("Get gauge metric value", "name", req.ID)
		Value, err := strconv.ParseFloat(h.memStor.Get(req.MType, req.ID), 64)
		if err != nil {
			logger.Logger.Error("cannot write Value", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if Value == -1 {
			Value = 0
		}
		req.Value = &Value
	}

	if req.MType == "counter" {
		logger.Logger.Infow("Get counter metric value", "name", req.ID)
		Delta, err := strconv.ParseInt(h.memStor.Get(req.MType, req.ID), 10, 64)
		if err != nil {
			logger.Logger.Error("cannot write Delta", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if Delta == -1 {
			Delta = 0
		}
		req.Delta = &Delta
	}
	w.Header().Set("Content-Type", "application/json")

	logger.Logger.Infow("Encode data for response")

	enc := json.NewEncoder(w)
	if err := enc.Encode(req); err != nil {
		logger.Logger.Error("error encoding response", zap.Error(err))
		return
	}
	logger.Logger.Infow("sending HTTP 200 response")
}

func (h *Handlers) requestEmpty(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusBadRequest)
}

func routePost(r *chi.Mux, h *Handlers) {
	logger.Logger.Infow("Init router for function Post")
	r.Route("/update", func(r chi.Router) {
		r.Post("/", logger.LoggingServer(h.requestUpdate))
		r.Route("/{type}", func(r chi.Router) {
			r.Post("/", logger.LoggingServer(h.requestType))
			r.Route("/{name}", func(r chi.Router) {
				r.Post("/", logger.LoggingServer(h.requestName))
				r.Route("/{value}", func(r chi.Router) {
					r.Post("/", logger.LoggingServer(h.requestValue))
				})
			})
		})
	})
}

func routeGet(r *chi.Mux, h *Handlers) {
	logger.Logger.Infow("Init router for function Get")
	r.Route("/value", func(r chi.Router) {
		r.Get("/", logger.LoggingServer(h.requestEmpty))
		r.Route("/{type}", func(r chi.Router) {
			r.Get("/", logger.LoggingServer(h.requestType))
			r.Route("/{name}", func(r chi.Router) {
				r.Get("/", logger.LoggingServer(h.requestGetName))
			})
		})
	})
}

func NewRouter() *chi.Mux {
	logger.Logger.Infow("Init router and handlers")
	handlers := NewHandlers(&storage.MemStorage{})
	handlers.memStor.Init()

	r := chi.NewRouter()
	routePost(r, handlers)
	routeGet(r, handlers)
	logger.Logger.Infow("Init router another function")
	r.Post("/value/", logger.LoggingServer(handlers.requestGetValue))
	r.Get("/", logger.LoggingServer(handlers.requestGetAll))
	return r
}
