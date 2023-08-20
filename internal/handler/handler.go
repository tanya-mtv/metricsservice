package handler

import (
	"encoding/json"
	"fmt"
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
	repository *repository.MetricRepositoryStorage
	cfg        *config.ConfigServer
}

func NewHandler(repository *repository.MetricRepositoryStorage, cfg *config.ConfigServer) *Handler {

	return &Handler{
		repository: repository,
		cfg:        cfg,
	}
}

func (h *Handler) GetMetricMethod(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		var metric models.Metrics

		jsonData, _ := io.ReadAll(c.Request.Body)
		if err := json.Unmarshal(jsonData, &metric); err != nil {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			log.Error(err)
		}

		switch metric.MType {
		case "counter":
			value, found := h.repository.GetCounter(metric.ID)
			if !found {
				newErrorResponse(c, http.StatusNotFound, "Metric not found")
				return
			}

			tmp := int64(value)
			metric.Delta = &tmp
		case "gauge":
			value, found := h.repository.GetGauge(metric.ID)
			if !found {
				newErrorResponse(c, http.StatusNotFound, "Metric not found")
				return
			}

			tmp := float64(value)
			metric.Value = &tmp
		}

		c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		c.JSON(http.StatusOK, metric)

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
		tmp := int64(cnt)

		metric := models.Metrics{
			ID:    metricName,
			MType: "counter",
			Delta: &tmp,
		}

		c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		c.JSON(http.StatusOK, metric)
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
		tmp := float64(gug)

		metric := models.Metrics{
			ID:    metricName,
			MType: "gauge",
			Value: &tmp,
		}

		c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		c.JSON(http.StatusOK, metric)

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
			c.JSON(http.StatusBadRequest, 0)
		}
	}
}

func (h *Handler) PostMetricsJSON(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var metric models.Metrics

		jsonData, _ := io.ReadAll(c.Request.Body)
		if err := json.Unmarshal(jsonData, &metric); err != nil {
			log.Error(err)
		}

		switch metric.MType {
		case "counter":
			if metric.Delta == nil {
				log.Info("Can't find tag metric  Delta")
				c.JSON(http.StatusBadRequest, 0)
				return
			}
			metricValue := *metric.Delta
			fmt.Println("metricValue", &metricValue)
			cnt := int64(h.repository.UpdateCounter(metric.ID, metricValue))

			log.Info("Update counter data wuth value ", cnt)
			metric.Delta = &cnt

			c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
			c.JSON(http.StatusOK, metric)
		case "gauge":
			if metric.Value == nil {
				log.Info("Can't find tag metric  Value")
				c.JSON(http.StatusBadRequest, 0)
				return
			}
			metricValue := *metric.Value
			gug := float64(h.repository.UpdateGauge(metric.ID, metricValue))
			log.Info("Update gauge data wuth value ", gug)
			metric.Value = &gug

			c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
			c.JSON(http.StatusOK, metric)
		default:
			c.JSON(http.StatusBadRequest, 0)
		}
	}
}

func (h *Handler) GetAllMetrics(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		metrics, err := h.repository.GetAll()

		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			log.Error("Error of getting all metrics from storage")
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
