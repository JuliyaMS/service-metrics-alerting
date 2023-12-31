package handlers

import (
	"encoding/json"
	"github.com/JuliyaMS/service-metrics-alerting/internal/config"
	"github.com/JuliyaMS/service-metrics-alerting/internal/database"
	"github.com/JuliyaMS/service-metrics-alerting/internal/file"
	"github.com/JuliyaMS/service-metrics-alerting/internal/html"
	"github.com/JuliyaMS/service-metrics-alerting/internal/logger"
	"github.com/JuliyaMS/service-metrics-alerting/internal/metrics"
	m "github.com/JuliyaMS/service-metrics-alerting/internal/middleware"
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

	if err := h.memStor.Add(metricType, metricName, metricValue); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
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

	logger.Logger.Info("Start handler:requestGetName")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")

	logger.Logger.Info("Get value from memStor")
	value := h.memStor.Get(metricType, metricName)

	if value != "-1" {
		logger.Logger.Info("Get value: ", value, " for metric: ", metricName)
		w.Write([]byte(value))
		w.WriteHeader(http.StatusOK)
		return
	}
	logger.Logger.Info("Metric with this name don`t found")
	w.WriteHeader(http.StatusNotFound)

}

func (h *Handlers) requestGetAll(w http.ResponseWriter, r *http.Request) {

	logger.Logger.Infow("Get all metrics in text/html")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	logger.Logger.Infow("Get data for HTML format")
	dataGauge, dataCounter := html.GetHTMLStructs(h.memStor.GetAll())

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	logger.Logger.Infow("Open file with html table")

	tmpl, err := template.ParseFiles("./data/html/index.html")
	if err != nil {
		logger.Logger.Error("Error while parse html file: ", err)
		return
	}

	logger.Logger.Infow("Execute for gauge metrics")

	err = tmpl.Execute(w, dataGauge)
	if err != nil {
		logger.Logger.Error("Error while execute for gauge metrics", err)
		return
	}
	logger.Logger.Infow("Execute for counter metrics")
	err = tmpl.Execute(w, dataCounter)
	if err != nil {
		logger.Logger.Error("Error while execute for counter metrics", err)
		return
	}
	logger.Logger.Infow("sending HTTP 200 response")
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
		if err := h.memStor.Add(req.MType, req.ID, strconv.FormatFloat(*req.Value, 'g', -1, 64)); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	if req.MType == "counter" {
		logger.Logger.Infow("Update counter metric", "name", req.ID, "value", *req.Delta)
		if err := h.memStor.Add(req.MType, req.ID, strconv.FormatInt(*req.Delta, 10)); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return

		}
		newDelta, err := strconv.ParseInt(h.memStor.Get("counter", req.ID), 10, 64)
		if err != nil {
			logger.Logger.Error("cannot write new Delta", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		req.Delta = &newDelta
	}

	w.WriteHeader(http.StatusOK)
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
	w.WriteHeader(http.StatusOK)

	logger.Logger.Infow("Encode data for response")
	enc := json.NewEncoder(w)
	if err := enc.Encode(req); err != nil {
		logger.Logger.Error("error encoding response", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	logger.Logger.Infow("sending HTTP 200 response")
}

func (h *Handlers) requestEmpty(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusBadRequest)
}

func (h *Handlers) PingDB(w http.ResponseWriter, r *http.Request) {
	logger.Logger.Info("start handler: PingDB")

	if err := h.memStor.CheckConnection(); err != nil {
		logger.Logger.Error("get error while check connection to Database:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	logger.Logger.Info("sending HTTP 200 response")
}

func (h *Handlers) UpdatesDB(w http.ResponseWriter, r *http.Request) {
	logger.Logger.Info("start handler: UpdatesDB")

	if r.Method != http.MethodPost {
		logger.Logger.Debug("got request with bad method", zap.String("method", r.Method))
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	logger.Logger.Infow("decoding request")
	var req []metrics.Metrics
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		logger.Logger.Error("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := h.memStor.AddAnyData(req); err != nil {
		logger.Logger.Error("Get error while execute transaction:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	logger.Logger.Infow("sending HTTP 200 response")
}

func routePost(r *chi.Mux, h *Handlers) {
	logger.Logger.Infow("Init router for function Post")
	r.Route("/update", func(r chi.Router) {
		r.Post("/", m.LoggingServer(m.CompressionGzip(h.requestUpdate)))
		r.Route("/{type}", func(r chi.Router) {
			r.Post("/", m.LoggingServer(h.requestType))
			r.Route("/{name}", func(r chi.Router) {
				r.Post("/", m.LoggingServer(h.requestName))
				r.Route("/{value}", func(r chi.Router) {
					r.Post("/", m.LoggingServer(h.requestValue))
				})
			})
		})
	})
}

func routeGet(r *chi.Mux, h *Handlers) {
	logger.Logger.Infow("Init router for function Get")
	r.Route("/value", func(r chi.Router) {
		r.Get("/", m.LoggingServer(h.requestEmpty))
		r.Route("/{type}", func(r chi.Router) {
			r.Get("/", m.LoggingServer(h.requestType))
			r.Route("/{name}", func(r chi.Router) {
				r.Get("/", m.LoggingServer(h.requestGetName))
			})
		})
	})
}

func NewRouter(DBConn *database.ConnectionDB) *chi.Mux {
	logger.Logger.Infow("init router and handlers")
	if config.Restore && config.FileStoragePath != "" {
		logger.Logger.Info("restore data from file:", config.FileStoragePath)
		err := file.ReadFromFile(config.FileStoragePath)
		if err != nil {
			logger.Logger.Errorf(err.Error(), "can't read data from file:", config.FileStoragePath)
		}
	}

	var handlers *Handlers
	if config.DatabaseDsn != "" {
		handlers = NewHandlers(DBConn)
	} else {
		handlers = NewHandlers(&storage.Storage)
	}
	handlers.memStor.Init()

	logger.Logger.Info("create new router")
	r := chi.NewRouter()
	routePost(r, handlers)
	routeGet(r, handlers)

	logger.Logger.Infow("init router another function")
	r.Post("/value/", m.LoggingServer(m.CompressionGzip(handlers.requestGetValue)))
	r.Post("/updates/", m.LoggingServer(m.CompressionGzip(handlers.UpdatesDB)))
	r.Get("/", m.LoggingServer(m.CompressionGzip(handlers.requestGetAll)))
	r.Get("/ping", m.LoggingServer(handlers.PingDB))

	return r
}
