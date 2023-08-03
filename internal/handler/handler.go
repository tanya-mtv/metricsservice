package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tanya-mtv/metricsservice/internal/config"
	"github.com/tanya-mtv/metricsservice/internal/logger"
	"github.com/tanya-mtv/metricsservice/internal/servise"
)

type Handler struct {
	service *servise.Service
	log     logger.Logger
	cfg     *config.Config
}

func NewHandler(service *servise.Service, log logger.Logger, cfg *config.Config) *Handler {
	return &Handler{
		service: service,
		log:     log,
		cfg:     cfg,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	// router := http.NewServeMux()
	// router.HandleFunc("/", h.handleMethod)
	router := gin.New()

	router.GET("", h.getAllMetrics)

	update := router.Group("/update")
	{
		update.POST("/", func(c *gin.Context) {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Неверный путь URL",
			})
		})

		update.POST("/:metricType/:metricName/:metricValue", h.PostMethod)
		// update.POST("/:metricType/:metricName/:metricValue", h.PostMethodCounter)
		// update.POST("/:metricType/:metricName/:metricValue", h.PostMethodGauge)

		update.GET("/counter/:metricName", h.GetMethodCounter)
		update.GET("/gauge/:metricName", h.GetMethodGauge)
	}

	return router
}
