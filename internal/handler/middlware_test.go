package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tanya-mtv/metricsservice/internal/config"

	"github.com/tanya-mtv/metricsservice/internal/constants"
	"github.com/tanya-mtv/metricsservice/internal/logger"
	"github.com/tanya-mtv/metricsservice/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestHandler_GzipMiddleware(t *testing.T) {

	tests := []struct {
		name            string
		contentEncoding string
		status          int
		want            []string
	}{
		{
			name:            "Test GzipMiddleware",
			contentEncoding: "Content-Encoding",
			status:          http.StatusOK,
			want:            []string{"text/html"},
		},
		{
			name:            "Test GzipMiddleware. Empty contect type",
			contentEncoding: "",
			status:          http.StatusOK,
			want:            []string{"text/html"},
		},
	}

	gin.SetMode(gin.TestMode)

	cfglog := &logger.Config{
		LogLevel: constants.LogLevel,
		DevMode:  constants.DevMode,
		Type:     constants.Type,
	}

	cfg := &config.ConfigServer{Port: "8080"}
	log := logger.NewAppLogger(cfglog)

	repo := repository.NewMetricStorage()
	h := Handler{
		storage: repository.NewStorage(repo, cfg, log),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, r := gin.CreateTestContext(w)
			r.POST("/", h.GetAllMetrics())
			c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

			r.ServeHTTP(w, c.Request)

			result := w.Result()

			defer result.Body.Close()
			require.Equal(t, tt.status, result.StatusCode)
			require.Equal(t, tt.want, result.Header["Content-Type"])

		})
	}

}
