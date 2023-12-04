package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"github.com/JuliyaMS/service-metrics-alerting/internal/logger"
)

func GetSignature(data []byte, key string) []byte {
	logger.Logger.Info("Start encrypt data")
	h := hmac.New(sha256.New, []byte(key))
	dst := h.Sum(data)
	logger.Logger.Info("Data encrypted successfully")
	return dst
}
