package handler

import (
	"encoding/json"
	"io"

	"net/http"
	"strconv"

	"github.com/tanya-mtv/metricsservice/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/tanya-mtv/metricsservice/internal/config"
	"github.com/tanya-mtv/metricsservice/internal/models"
	"github.com/tanya-mtv/metricsservice/internal/repository"
)

type Handler struct {
	storage *repository.Storage
	cfg     *config.ConfigServer
	log     logger.Logger
}

func NewHandler(storage *repository.Storage, cfg *config.ConfigServer, log logger.Logger) *Handler {
	return &Handler{
		storage: storage,
		cfg:     cfg,
		log:     log,
	}
}

func (h *Handler) GetMethodCounter() gin.HandlerFunc {
	return func(c *gin.Context) {
		metricName := c.Param("metricName")

		cnt, found := h.storage.GetCounter(metricName)

		if !found {
			newErrorResponse(c, http.StatusNotFound, "Metric not found")
			return
		}
		c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
		c.JSON(http.StatusOK, cnt)
	}

}

func (h *Handler) GetMethodGauge() gin.HandlerFunc {
	return func(c *gin.Context) {
		metricName := c.Param("metricName")

		gug, found := h.storage.GetGauge(metricName)

		if !found {
			newErrorResponse(c, http.StatusNotFound, "Metric not found")
			return
		}
		c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
		c.JSON(http.StatusOK, gug)
	}

}

func (h *Handler) PostMetricsValueJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		var metric models.Metrics
		jsonData, _ := io.ReadAll(c.Request.Body)
		if err := json.Unmarshal(jsonData, &metric); err != nil {
			h.log.Error(err)
		}

		switch metric.MType {
		case "counter":
			cnt, found := h.storage.GetCounter(metric.ID)
			if !found {
				newErrorResponse(c, http.StatusNotFound, "Metric not found")
				return
			}

			tmp := int64(cnt)

			metric := models.Metrics{
				ID:    metric.ID,
				MType: "counter",
				Delta: &tmp,
			}

			c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
			c.JSON(http.StatusOK, metric)
		case "gauge":
			gug, found := h.storage.GetGauge(metric.ID)

			if !found {
				newErrorResponse(c, http.StatusNotFound, "Metric not found")
				return
			}
			tmp := float64(gug)

			metric := models.Metrics{
				ID:    metric.ID,
				MType: "gauge",
				Value: &tmp,
			}

			c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
			c.JSON(http.StatusOK, metric)
		default:
			c.JSON(http.StatusBadRequest, 0)
		}

	}

}

func (h *Handler) PostMetrics() gin.HandlerFunc {
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

			cnt := h.storage.UpdateCounter(metricName, int64(metricValue))

			c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
			c.JSON(http.StatusOK, cnt)
		case "gauge":
			metricValue, err := strconv.ParseFloat(c.Param("metricValue"), 64)

			if err != nil {
				newErrorResponse(c, http.StatusBadRequest, "invalid metricValue param")
				return
			}
			gug := h.storage.UpdateGauge(metricName, metricValue)

			c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
			c.JSON(http.StatusOK, gug)
		default:
			c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
			c.JSON(http.StatusBadRequest, 0)
		}
	}
}

func (h *Handler) PostMetricsUpdateJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		var metric models.Metrics

		jsonData, _ := io.ReadAll(c.Request.Body)
		if err := json.Unmarshal(jsonData, &metric); err != nil {
			h.log.Error(err)
		}

		switch metric.MType {
		case "counter":
			if metric.Delta == nil {
				h.log.Info("Can't find  metric tag Delta")
				c.JSON(http.StatusBadRequest, 0)
				return
			}
			metricValue := *metric.Delta

			cnt := int64(h.storage.UpdateCounter(metric.ID, metricValue))

			metric.Delta = &cnt

			c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
			c.JSON(http.StatusOK, metric)
		case "gauge":
			if metric.Value == nil {
				h.log.Info("Can't find tag metric  Value")
				c.JSON(http.StatusBadRequest, 0)
				return
			}
			metricValue := *metric.Value
			gug := float64(h.storage.UpdateGauge(metric.ID, metricValue))
			h.log.Info("Update gauge data with value ", gug)
			metric.Value = &gug

			c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
			c.JSON(http.StatusOK, metric)
		default:
			c.JSON(http.StatusBadRequest, 0)
		}
	}
}

func (h *Handler) GetAllMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		metrics := h.storage.GetAll()
		c.Writer.Header().Set("Content-Type", "text/html")
		c.JSON(http.StatusOK, getAllMetricResponse{
			Data: metrics,
		})
	}

}

type getAllMetricResponse struct {
	Data []models.Metrics `json:"data"`
}
