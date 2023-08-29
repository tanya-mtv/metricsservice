package handler

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"io"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

type compressWriter struct {
	zw *gzip.Writer
	gin.ResponseWriter
	buf bytes.Buffer
}

func newCompressWriter() *compressWriter {
	var bf bytes.Buffer
	return &compressWriter{
		zw: gzip.NewWriter(&bf),
	}
}

func (c *compressWriter) Write(data []byte) (int, error) {
	return c.zw.Write(data)
}

func (c *compressWriter) Close() error {
	return c.zw.Close()
}

func (c *compressWriter) CloseNotify() <-chan bool {
	return c.ResponseWriter.CloseNotify()
}

func (c *compressWriter) Flush() {
	c.ResponseWriter.Flush()

}

func (c *compressWriter) Header() http.Header {
	return c.ResponseWriter.Header()
}

func (c *compressWriter) Size() int {
	return c.ResponseWriter.Size()
}

func (c *compressWriter) Status() int {
	return c.ResponseWriter.Status()
}

func (c *compressWriter) WriteHeader(i int) {
	c.ResponseWriter.WriteHeader(i)
}

func (c *compressWriter) Written() bool {
	return c.ResponseWriter.Written()
}

func (c *compressWriter) WriteHeaderNow() {
	c.ResponseWriter.WriteHeaderNow()
}

func (c *compressWriter) WriteString(s string) (n int, err error) {
	return c.ResponseWriter.WriteString(s)
}

func (c *compressWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return c.ResponseWriter.(http.Hijacker).Hijack()
}

func (c *compressWriter) Pusher() (pusher http.Pusher) {
	if pusher, ok := c.ResponseWriter.(http.Pusher); ok {
		return pusher
	}
	return nil
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
