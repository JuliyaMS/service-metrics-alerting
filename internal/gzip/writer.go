package gzip

import (
	"compress/gzip"
	"net/http"
)

type CompressWriter struct {
	res http.ResponseWriter
	zw  *gzip.Writer
}

func NewCompressWriter(w http.ResponseWriter) *CompressWriter {
	return &CompressWriter{
		res: w,
		zw:  gzip.NewWriter(w),
	}
}

func (c *CompressWriter) Header() http.Header {
	return c.res.Header()
}

func (c *CompressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

func (c *CompressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.res.Header().Set("Content-Encoding", "gzip")
	}
	c.res.WriteHeader(statusCode)
}

func (c *CompressWriter) Close() error {
	return c.zw.Close()
}
