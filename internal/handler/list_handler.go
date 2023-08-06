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

// func (h *Handler) PostMethodCounter(c *gin.Context) {
// 	metricType := c.Param("metricType")
// 	if metricType != "counter" {
// 		c.JSON(http.StatusBadRequest, 0)
// 		return
// 	}
// 	metricName := c.Param("metricName")

// 	metricValue, err := strconv.Atoi(c.Param("metricValue"))

// 	if err != nil {
// 		newErrorResponse(c, http.StatusBadRequest, "invalid metricValue param")
// 		return
// 	}

// 	cnt := h.service.UpdateCounter(metricName, int64(metricValue))

// 	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
// 	c.JSON(http.StatusOK, cnt)
// }

// func (h *Handler) PostMethodGauge(c *gin.Context) {
// 	metricName := c.Param("metricName")
// 	metricValue, err := strconv.ParseFloat(c.Param("metricValue"), 64)

// 	if err != nil {
// 		newErrorResponse(c, http.StatusBadRequest, "invalid metricValue param")
// 		return
// 	}
// 	gug := h.service.UpdateGauge(metricName, metricValue)

// 	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
// 	c.JSON(http.StatusOK, gug)
// }

// func (h *Handler) PostMethod(c *gin.Context) {

// 	if r.Method != http.MethodPost {
// 		w.WriteHeader(http.StatusMethodNotAllowed)
// 		return
// 	}

// 	sliceURL := strings.Split(r.URL.Path, "/")

// 	if len(sliceURL) != 5 || sliceURL[1] != "update" {
// 		w.WriteHeader(http.StatusNotFound)
// 		return
// 	}

// 	metricsType := sliceURL[2]
// 	metricsName := sliceURL[3]
// 	metricsValue := sliceURL[4]

// 	if metricsType == "counter" {
// 		if value, err := strconv.ParseInt(metricsValue, 10, 64); err == nil {
// 			cnt := h.service.UpdateCounter(metricsName, value)
// 			w.Write([]byte(utils.CounterToBytes(cnt)))

// 		} else {
// 			w.WriteHeader(http.StatusBadRequest)
// 			return
// 		}
// 	} else if metricsType == "gauge" {
// 		if value, err := strconv.ParseFloat(metricsValue, 64); err == nil {
// 			gug := h.service.UpdateGauge(metricsName, value)

// 			w.Write([]byte(utils.GaugeToBytes(gug)))
// 		} else {
// 			w.WriteHeader(http.StatusBadRequest)
// 			return
// 		}
// 	} else {
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
// }

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
