package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tanya-mtv/metricsservice/internal/repository"
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

	h := Handler{
		repository: repository.NewMetricStorage(),
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
