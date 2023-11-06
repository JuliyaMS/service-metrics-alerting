package logger

import (
	"go.uber.org/zap"
	"net/http"
)

var Logger = NewLogger()

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

func NewLogger() *zap.SugaredLogger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	return logger.Sugar()
}
