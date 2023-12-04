package agent

import (
	"bytes"
	"compress/gzip"
	"github.com/JuliyaMS/service-metrics-alerting/internal/logger"
	"sync"
)

type Compressor struct {
	gz *gzip.Writer
}

func NewCompressor(compress bool) *Compressor {
	if compress {
		return &Compressor{gz: &gzip.Writer{}}
	}
	return nil
}

var gzipMutex sync.Mutex

func (c *Compressor) CompressData(data []byte) (*bytes.Buffer, error) {
	gzipMutex.Lock()
	defer gzipMutex.Unlock()

	buf := bytes.NewBuffer(nil)

	c.gz = gzip.NewWriter(buf)

	if _, err := c.gz.Write(data); err != nil {
		logger.Logger.Error("Error while compress data:", err.Error())
		return nil, err
	}
	if err := c.gz.Close(); err != nil {
		logger.Logger.Error("Error while compress data:", err.Error())
		return nil, err
	}
	return buf, nil
}
