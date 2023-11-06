package gzip

import (
	"fmt"
	"github.com/JuliyaMS/service-metrics-alerting/internal/logger"
	"net/http"
	"strings"
)

func GzipCompression(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Logger.Infow("Start midlware")
		ow := w

		logger.Logger.Infow("Check client's Accept-Encoding")
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")

		if supportsGzip {
			logger.Logger.Info("Start compress data")
			fmt.Println(w.Header())
			cw := newCompressWriter(w)
			ow = cw

			defer func(cw *compressWriter) {
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
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				logger.Logger.Error("Error while create compress reader:", err.Error())
				return
			}
			r.Body = cr

			defer func(cr *compressReader) {
				er := cr.Close()
				if er != nil {
					w.WriteHeader(http.StatusInternalServerError)
					logger.Logger.Error("Error while read compress data", er.Error())
				}
			}(cr)
		}
		logger.Logger.Infow("End midlware")
		h.ServeHTTP(ow, r)
	}
}
