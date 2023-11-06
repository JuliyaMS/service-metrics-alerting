package gzip

import (
	"compress/gzip"
	"net/http"
)

type compressWriter struct {
	res http.ResponseWriter
	zw  *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		res: w,
		zw:  gzip.NewWriter(w),
	}
}

func (c *compressWriter) Header() http.Header {
	return c.res.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.res.Header().Set("Content-Encoding", "gzip")
	}
	c.res.WriteHeader(statusCode)
}

func (c *compressWriter) Close() error {
	return c.zw.Close()
}
