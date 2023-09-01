package server

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/tanya-mtv/metricsservice/internal/config"
	"github.com/tanya-mtv/metricsservice/internal/handler"
	"github.com/tanya-mtv/metricsservice/internal/logger"
	"github.com/tanya-mtv/metricsservice/internal/repository"
)

func NewRouter(stor *repository.MetricStorage, db *sqlx.DB, cfg *config.ConfigServer, log logger.Logger) *gin.Engine {

	h := handler.NewHandler(stor, db, cfg, log)
	router := gin.New()

	router.Use(h.GzipMiddleware)
	router.Use(h.WithLogging)

	router.GET("/", h.GetAllMetrics)
	router.GET("/ping", h.Ping)

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
