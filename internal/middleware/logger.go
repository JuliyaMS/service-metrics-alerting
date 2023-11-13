package middleware

import (
	"github.com/JuliyaMS/service-metrics-alerting/internal/logger"
	"net/http"
	"time"
)

type DataResponse struct {
	Status int
	Size   int
}

type LoggingResponse struct {
	http.ResponseWriter
	ResponseData *DataResponse
}

func (r *LoggingResponse) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.ResponseData.Size += size
	return size, err
}

func (r *LoggingResponse) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.ResponseData.Status = statusCode
}

func LoggingServer(h http.HandlerFunc) http.HandlerFunc {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		uri := r.RequestURI
		method := r.Method

		responseData := &DataResponse{
			Status: 0,
			Size:   0,
		}

		lw := LoggingResponse{
			ResponseWriter: w,
			ResponseData:   responseData,
		}

		h.ServeHTTP(&lw, r)

		duration := time.Since(start)

		logger.Logger.Infoln(
			"uri", uri,
			"method", method,
			"duration", duration,
			"size", responseData.Size,
			"status", responseData.Status,
		)
	}
	return logFn
}
