package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tanya-mtv/metricsservice/internal/config"
	"github.com/tanya-mtv/metricsservice/internal/logger"
	"github.com/tanya-mtv/metricsservice/internal/models"
	"github.com/tanya-mtv/metricsservice/internal/repository"
)

type Handler struct {
	repository *repository.RepositoryStorage
	log        logger.Logger
	cfg        *config.ConfigServer
}

func NewHandler(repository *repository.RepositoryStorage, log logger.Logger, cfg *config.ConfigServer) *Handler {
	return &Handler{
		repository: repository,
		log:        log,
		cfg:        cfg,
	}
}

func (h *Handler) GetMethodCounter(repository *repository.RepositoryStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		metricName := c.Param("metricName")

		cnt, found := h.repository.GetCounter(metricName)

		if !found {
			newErrorResponse(c, http.StatusNotFound, "Metric not found")
			return
		}
		c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
		c.JSON(http.StatusOK, cnt)
	}

}

func (h *Handler) GetMethodGauge(repository *repository.RepositoryStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		metricName := c.Param("metricName")

		gug, found := h.repository.GetGauge(metricName)

		if !found {
			newErrorResponse(c, http.StatusNotFound, "Metric not found")
			return
		}
		c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
		c.JSON(http.StatusOK, gug)
	}

}

func (h *Handler) PostMetrics(repository *repository.RepositoryStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		metricType := c.Param("metricType")
		metricName := c.Param("metricName")
		switch metricType {
		case "counter":
			metricValue, err := strconv.Atoi(c.Param("metricValue"))

			if err != nil {
				newErrorResponse(c, http.StatusBadRequest, "invalid metricValue param")
				return
			}

			cnt := h.repository.UpdateCounter(metricName, int64(metricValue))

			c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
			c.JSON(http.StatusOK, cnt)
		case "gauge":
			metricValue, err := strconv.ParseFloat(c.Param("metricValue"), 64)

			if err != nil {
				newErrorResponse(c, http.StatusBadRequest, "invalid metricValue param")
				return
			}
			gug := h.repository.UpdateGauge(metricName, metricValue)

			c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
			c.JSON(http.StatusOK, gug)
		default:
			c.JSON(http.StatusBadRequest, 0)
		}
	}
}

func (h *Handler) GetAllMetrics(repository *repository.RepositoryStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		metrics, err := repository.GetAll()
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, getAllMetricResponse{
			Data: metrics,
		})
	}

}

type getAllMetricResponse struct {
	Data []models.Metrics `json:"data"`
}
