package middleware

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"github.com/JuliyaMS/service-metrics-alerting/internal/config"
	"github.com/JuliyaMS/service-metrics-alerting/internal/hash"
	"github.com/JuliyaMS/service-metrics-alerting/internal/logger"
	"io"
	"net/http"
)

type signatureResponseWriter struct {
	http.ResponseWriter
}

func newSignatureResponseWriter(w http.ResponseWriter) *signatureResponseWriter {
	return &signatureResponseWriter{
		ResponseWriter: w,
	}
}

func (c *signatureResponseWriter) Write(data []byte) (int, error) {
	hash := hash.GetSignature(data, config.HashKeyServer)

	c.ResponseWriter.Header().Set("HashSHA256", hex.EncodeToString(hash))

	return c.ResponseWriter.Write(data)
}

func checkSignature(r *http.Request) error {
	logger.Logger.Info("Check signature from agent")
	data, err := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	sign, err := base64.StdEncoding.DecodeString(r.Header.Get("HashSHA256"))
	if err != nil {
		return err
	}

	if !bytes.Equal(sign, hash.GetSignature(data, config.HashKeyServer)) {
		return err
	}
	return nil
}

func SignatureData(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ow := w
		if config.HashKeyServer != "" {
			err := checkSignature(r)
			if err != nil {
				logger.Logger.Error("Get error while check signature: ", err.Error())
				w.WriteHeader(http.StatusBadRequest)
			}
			ow = newSignatureResponseWriter(w)
		}
		h.ServeHTTP(ow, r)
	}
}
