package handler

import (
	"compress/gzip"
	"io"

	"github.com/gin-gonic/gin"
)

type compressWriter struct {
	zw *gzip.Writer
	gin.ResponseWriter
}

// func newCompressWriter(w http.ResponseWriter) *compressWriter {
func newCompressWriter(w gin.ResponseWriter) *compressWriter {

	return &compressWriter{
		zw: gzip.NewWriter(w),
	}
}

func (c *compressWriter) Write(data []byte) (int, error) {
	return c.zw.Write(data)
}

func (c *compressWriter) Close() error {
	return c.zw.Close()
}

// compressReader реализует интерфейс io.ReadCloser и позволяет прозрачно для сервера
// декомпрессировать получаемые от клиента данные
type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}
