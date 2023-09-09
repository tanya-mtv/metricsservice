package server

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/tanya-mtv/metricsservice/internal/handler"
)

func (s *server) NewRouter(db *sqlx.DB) *gin.Engine {
	hp := handler.NewHandlerPing(db)
	h := handler.NewHandler(s.stor, s.cfg, s.log)

	router := gin.New()

	router.Use(h.GzipMiddleware)
	router.Use(h.WithLogging)

	router.GET("/", h.GetAllMetrics)
	router.GET("/ping", hp.Ping)

	router.POST("/updates/", h.PostMetricsList)
	router.POST("/update", h.PostMetricsUpdateJSON)
	router.POST("/update/:metricType/:metricName/:metricValue", h.PostMetrics)

	value := router.Group("/value")
	{
		value.POST("/", h.PostMetricsValueJSON)
		value.GET("/counter/:metricName", h.GetMethodCounter)
		value.GET("/gauge/:metricName", h.GetMethodGauge)
	}

	return router
}
