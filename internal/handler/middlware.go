package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) WithLogging(c *gin.Context) {

	start := time.Now()

	req := c.Request
	res := c.Writer

	c.Next()

	duration := time.Since(start)

	h.log.Infoln(
		"uri:", req.RequestURI,
		"method:", req.Method,
		"duration:", duration,
		"status:", res.Status(),
		"size:", res.Size(),
	)

}

func (h *Handler) GzipMiddleware(c *gin.Context) {

	// проверяем, что клиент отправил серверу сжатые данные в формате gzip
	contentEncoding := c.GetHeader("Content-Encoding")
	sendsGzip := strings.Contains(contentEncoding, "gzip")
	if sendsGzip {
		// оборачиваем тело запроса в io.Reader с поддержкой декомпрессии
		cr, err := newCompressReader(c.Request.Body)

		if err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		// меняем тело запроса на новое
		c.Request.Body = cr
		defer cr.Close()
	}

	// проверяем, что клиент умеет получать от сервера сжатые данные в формате gzip
	acceptEncoding := c.GetHeader("Accept-Encoding")
	supportsGzip := strings.Contains(acceptEncoding, "gzip")
	contentType := c.GetHeader("Content-Type")

	if supportsGzip && (strings.Contains(contentType, "application/json") || strings.Contains(contentType, "text/html") || len(contentType) == 0) {
		cw := newCompressWriter(c.Writer)
		cw.Header().Add("Content-Encoding", "gzip")

		defer cw.Close()

		c.Writer = cw

	}

	c.Next()

}
