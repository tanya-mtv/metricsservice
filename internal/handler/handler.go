package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"net/http"
	"strconv"

	"github.com/tanya-mtv/metricsservice/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/tanya-mtv/metricsservice/internal/config"
	"github.com/tanya-mtv/metricsservice/internal/models"
)

type Handler struct {
	storage metricStorage
	cfg     *config.ConfigServer
	log     logger.Logger
}

func NewHandler(storage metricStorage, cfg *config.ConfigServer, log logger.Logger) *Handler {
	return &Handler{
		storage: storage,
		cfg:     cfg,
		log:     log,
	}
}

func (h *Handler) GetMethodCounter(c *gin.Context) {

	metricName := c.Param("metricName")

	cnt, found := h.storage.GetCounter(metricName)

	if !found {
		newErrorResponse(c, http.StatusNotFound, "Metric not found")
		return
	}
	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.JSON(http.StatusOK, cnt)

}

func (h *Handler) GetMethodGauge(c *gin.Context) {
	metricName := c.Param("metricName")

	gug, found := h.storage.GetGauge(metricName)

	if !found {
		newErrorResponse(c, http.StatusNotFound, "Metric not found")
		return
	}
	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.JSON(http.StatusOK, gug)

}

func (h *Handler) PostMetricsList(c *gin.Context) {
	if c.ContentType() != "application/json" {
		h.log.Error("incorrect  ContentType")
		newErrorResponse(c, http.StatusBadRequest, "{}")
		return
	}
	jsonData, err := io.ReadAll(c.Request.Body)
	fmt.Printf("jsonData PostMetricsList %+v \n", string(jsonData))
	if err != nil {
		h.log.Error("PostMetricsList", err)
		newErrorResponse(c, http.StatusBadRequest, "{}")
		return
	}

	jsonDatarep := bytes.Replace(jsonData, []byte("Mtype"), []byte("type"), -1)

	var metrics []*models.Metrics
	if err := json.Unmarshal(jsonDatarep, &metrics); err != nil {
		h.log.Error(err)
		newErrorResponse(c, http.StatusBadRequest, "{}")
		return
	}

	list, err := h.storage.UpdateMetrics(metrics)
	if err != nil {
		h.log.Error("Can't update metrics", err)
		c.JSON(http.StatusBadRequest, "{}")
		return
	}

	c.JSON(http.StatusOK, list)

}

func (h *Handler) PostMetricsValueJSON(c *gin.Context) {
	if c.ContentType() != "application/json" {
		h.log.Error("PostMetricsValueJSON. Incorrect  ContentType")
		newErrorResponse(c, http.StatusBadRequest, "{}")
		return
	}
	var metric models.Metrics

	jsonData, _ := io.ReadAll(c.Request.Body)
	if err := json.Unmarshal(jsonData, &metric); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "{}")
		h.log.Error(err)
		return

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

func (h *Handler) PostMetrics(c *gin.Context) {

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

func (h *Handler) PostMetricsUpdateJSON(c *gin.Context) {

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

func (h *Handler) GetAllMetrics(c *gin.Context) {

	metrics := h.storage.GetAll()
	c.Writer.Header().Set("Content-Type", "text/html")
	c.JSON(http.StatusOK, getAllMetricResponse{
		Data: metrics,
	})

}

type getAllMetricResponse struct {
	Data []models.Metrics `json:"data"`
}
