package handler

import (
	"time"

	"github.com/tanya-mtv/metricsservice/internal/logger"

	"github.com/gin-gonic/gin"
)

func (h *Handler) WithLogging(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		req := c.Request
		res := c.Writer

		c.Next()

		duration := time.Since(start)

		log.Infoln(
			"uri:", req.RequestURI,
			"method:", req.Method,
			"duration:", duration,
			"status:", res.Status(),
			"size:", res.Size(),
		)

	}

}
