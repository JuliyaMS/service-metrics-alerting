package middleware

import (
	"fmt"
	"github.com/JuliyaMS/service-metrics-alerting/internal/gzip"
	"github.com/JuliyaMS/service-metrics-alerting/internal/logger"
	"net/http"
	"strings"
	"time"
)

func LoggingServer(h http.HandlerFunc) http.HandlerFunc {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		uri := r.RequestURI
		method := r.Method

		responseData := &logger.DataResponse{
			Status: 0,
			Size:   0,
		}

		lw := logger.LoggingResponse{
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

func CompressionGzip(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Logger.Infow("Start middleware")
		ow := w

		logger.Logger.Infow("Check client's Accept-Encoding")
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")

		if supportsGzip {
			logger.Logger.Info("Start compress data")
			fmt.Println(w.Header())
			cw := gzip.NewCompressWriter(w)
			ow = cw

			defer func(cw *gzip.CompressWriter) {
				err := cw.Close()
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					logger.Logger.Error("Error while write compress data", err.Error())
				}
			}(cw)
		}

		logger.Logger.Infow("Check client's Content-Encoding")
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			logger.Logger.Infow("Read compress data")
			cr, err := gzip.NewCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				logger.Logger.Error("Error while create compress reader:", err.Error())
				return
			}
			r.Body = cr

			defer func(cr *gzip.CompressReader) {
				er := cr.Close()
				if er != nil {
					w.WriteHeader(http.StatusInternalServerError)
					logger.Logger.Error("Error while read compress data", er.Error())
				}
			}(cr)
		}
		logger.Logger.Infow("End middleware")
		h.ServeHTTP(ow, r)
	}
}
