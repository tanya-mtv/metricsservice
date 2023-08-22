package handler

import (
	"bufio"
	"compress/gzip"
	"io"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

type compressWriter struct {
	zw *gzip.Writer
	w  gin.ResponseWriter
}

// func newCompressWriter(w http.ResponseWriter) *compressWriter {
func newCompressWriter(w gin.ResponseWriter) *compressWriter {

	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (c *compressWriter) Write(data []byte) (int, error) {
	return c.zw.Write(data)
}

func (c *compressWriter) Close() error {
	return c.zw.Close()
}

func (c *compressWriter) CloseNotify() <-chan bool {
	return c.w.CloseNotify()
}

func (c *compressWriter) Flush() {
	c.w.Flush()

}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) Size() int {
	return c.w.Size()
}

func (c *compressWriter) Status() int {
	return c.w.Status()
}

func (c *compressWriter) WriteHeader(i int) {
	c.w.WriteHeader(i)
}

func (c *compressWriter) Written() bool {
	return c.w.Written()
}

func (c *compressWriter) WriteHeaderNow() {
	c.w.WriteHeaderNow()
}

func (c *compressWriter) WriteString(s string) (n int, err error) {
	return c.w.WriteString(s)
}

func (c *compressWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return c.w.(http.Hijacker).Hijack()
}

func (c *compressWriter) Pusher() (pusher http.Pusher) {
	if pusher, ok := c.w.(http.Pusher); ok {
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
