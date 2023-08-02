package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/tanya-mtv/metricsservice/internal/utils"

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

func (h *Handler) handleMethod(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	sliceURL := strings.Split(r.URL.Path, "/")

	if len(sliceURL) != 5 || sliceURL[1] != "update" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	metricsType := sliceURL[2]
	metricsName := sliceURL[3]
	metricsValue := sliceURL[4]

	if metricsType == "counter" {
		if value, err := strconv.ParseInt(metricsValue, 10, 64); err == nil {
			cnt := h.service.UpdateCounter(metricsName, value)
			w.Write([]byte(utils.CounterToBytes(cnt)))

		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else if metricsType == "gauge" {
		if value, err := strconv.ParseFloat(metricsValue, 64); err == nil {
			gug := h.service.UpdateGauge(metricsName, value)

			w.Write([]byte(utils.GaugeToBytes(gug)))
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
}

func (h *Handler) InitRoutes() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/", h.handleMethod)

	return router
}
