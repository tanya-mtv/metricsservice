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
	cfg     *config.ConfigServer
	router  *gin.Engine
}

func NewHandler(service *servise.Service, log logger.Logger, cfg *config.ConfigServer, router *gin.Engine) *Handler {
	return &Handler{
		service: service,
		log:     log,
		cfg:     cfg,
		router:  router,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {

	h.router.GET("", h.getAllMetrics)

	update := h.router.Group("/update")
	{
		update.POST("/", func(c *gin.Context) {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Неверный путь URL",
			})
		})

		update.POST("/:metricType/:metricName/:metricValue", h.PostMethod)

	}

	value := h.router.Group("/value")
	{
		value.GET("/counter/:metricName", h.GetMethodCounter)
		value.GET("/gauge/:metricName", h.GetMethodGauge)
	}

	return h.router
}
