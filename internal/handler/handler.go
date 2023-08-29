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
	repository *repository.MetricStorage
	cfg        *config.ConfigServer
	cWriter    *compressWriter
	log        logger.Logger
}

func NewHandler(repository *repository.MetricStorage, cfg *config.ConfigServer, log logger.Logger) *Handler {
	return &Handler{
		repository: repository,
		cfg:        cfg,
		cWriter:    newCompressWriter(),
		log:        log,
	}
}

func (h *Handler) GetMethodCounter() gin.HandlerFunc {
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

func (h *Handler) GetMethodGauge() gin.HandlerFunc {
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

func (h *Handler) PostMetricsValueJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		var metric models.Metrics
		jsonData, _ := io.ReadAll(c.Request.Body)
		if err := json.Unmarshal(jsonData, &metric); err != nil {
			h.log.Error(err)
		}

		switch metric.MType {
		case "counter":
			cnt, found := h.repository.GetCounter(metric.ID)
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
			gug, found := h.repository.GetGauge(metric.ID)

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

			cnt := h.repository.UpdateCounter(metricName, int64(metricValue))

			c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
			c.JSON(http.StatusOK, cnt)
		case "gauge":
			metricValue, err := strconv.ParseFloat(c.Param("metricValue"), 64)

			if err != nil {
				newErrorResponse(c, http.StatusBadRequest, "invalid metricValue param")
				return
			}
			gug := h.repository.UpdateGauge(metricName, metricValue)

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

			cnt := int64(h.repository.UpdateCounter(metric.ID, metricValue))

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
			gug := float64(h.repository.UpdateGauge(metric.ID, metricValue))
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

		metrics := h.repository.GetAll()
		c.Writer.Header().Set("Content-Type", "text/html")
		c.JSON(http.StatusOK, getAllMetricResponse{
			Data: metrics,
		})
	}

}

type getAllMetricResponse struct {
	Data []models.Metrics `json:"data"`
}
