package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tanya-mtv/metricsservice/internal/models"
)

func (h *Handler) GetMethodCounter(c *gin.Context) {

	metricName := c.Param("metricName")

	cnt, found := h.service.GetCounter(metricName)

	if !found {
		newErrorResponse(c, http.StatusNotFound, "Metric not found")
		return
	}
	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.JSON(http.StatusOK, cnt)
}

func (h *Handler) GetMethodGauge(c *gin.Context) {
	metricName := c.Param("metricName")

	gug, found := h.service.GetGauge(metricName)

	if !found {
		newErrorResponse(c, http.StatusNotFound, "Metric not found")
		return
	}
	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.JSON(http.StatusOK, gug)
}

func (h *Handler) PostMethod(c *gin.Context) {
	metricType := c.Param("metricType")
	metricName := c.Param("metricName")
	switch metricType {
	case "counter":
		metricValue, err := strconv.Atoi(c.Param("metricValue"))

		if err != nil {
			newErrorResponse(c, http.StatusBadRequest, "invalid metricValue param")
			return
		}

		cnt := h.service.UpdateCounter(metricName, int64(metricValue))

		c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
		c.JSON(http.StatusOK, cnt)
	case "gauge":
		metricValue, err := strconv.ParseFloat(c.Param("metricValue"), 64)

		if err != nil {
			newErrorResponse(c, http.StatusBadRequest, "invalid metricValue param")
			return
		}
		gug := h.service.UpdateGauge(metricName, metricValue)

		c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
		c.JSON(http.StatusOK, gug)
	default:
		c.JSON(http.StatusBadRequest, 0)
	}
}

func (h *Handler) getAllMetrics(c *gin.Context) {

	metrics, err := h.service.MetricStorage.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, getAllMetricResponse{
		Data: metrics,
	})
}

type getAllMetricResponse struct {
	Data []models.Metrics `json:"data"`
}
