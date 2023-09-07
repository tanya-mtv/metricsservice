package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *gin.Engine, method,
	path string) (*http.Response, string) {
	req, err := http.NewRequest(method, path, nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	ts.ServeHTTP(rr, req)
	require.NoError(t, err)

	resp := rr.Result()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	fmt.Println("respBody", string(respBody))
	return resp, string(respBody)
}

func TestRouter(t *testing.T) {
	router := gin.Default()

	var testTableValue = []struct {
		url    string
		want   string
		status int
	}{

		// проверим на ошибочный запрос
		{"/value/counter1/testSetGet40", "404 page not found", http.StatusNotFound},
		{"/value/gauge1/testSetGet40", "404 page not found", http.StatusNotFound},
	}
	for _, v := range testTableValue {
		resp, get := testRequest(t, router, "GET", v.url)

		defer resp.Body.Close()
		require.Equal(t, v.status, resp.StatusCode)
		require.Equal(t, v.want, get)
	}

	var testTableUpdate = []struct {
		url    string
		want   string
		status int
	}{

		// проверим на ошибочный запрос
		{"/update/counter1/testSetGet40/1", "404 page not found", http.StatusNotFound},
		{"/update/gauge1/testSetGet40/235", "404 page not found", http.StatusNotFound},
	}
	for _, v := range testTableUpdate {
		resp, get := testRequest(t, router, "GET", v.url)

		defer resp.Body.Close()
		require.Equal(t, v.status, resp.StatusCode)
		require.Equal(t, v.want, get)
	}
}

// func TestHandler_PostMetricsList(t *testing.T) {
// 	idCounter := "PollCount"
// 	idGauge := "SomeGauge"
// 	valueCounter1, valueCounter2 := int64(rand.Int31()), int64(rand.Int31())
// 	valueGauge1, valueGauge2 := float64(rand.Float64()), float64(rand.Float64())
// 	metrics := []models.Metrics{
// 		{
// 			ID:    idCounter,
// 			MType: "counter",
// 			Delta: &valueCounter1,
// 		},
// 		{
// 			ID:    idGauge,
// 			MType: "gauge",
// 			Value: &valueGauge1,
// 		},
// 		{
// 			ID:    idCounter,
// 			MType: "counter",
// 			Delta: &valueCounter2,
// 		},
// 		{
// 			ID:    idGauge,
// 			MType: "gauge",
// 			Value: &valueGauge2,
// 		},
// 	}
// 	tests := []struct {
// 		name     string
// 		sentdata []models.Metrics
// 		status   int
// 	}{
// 		{
// 			name:     "Test post list metrics json",
// 			sentdata: metrics,
// 			status:   http.StatusOK,
// 		},
// 	}

// 	gin.SetMode(gin.TestMode)

// 	cfglog := &logger.Config{
// 		LogLevel: constants.LogLevel,
// 		DevMode:  constants.DevMode,
// 		Type:     constants.Type,
// 	}

// 	cfg := &config.ConfigServer{Port: "8080"}
// 	log := logger.NewAppLogger(cfglog)

// 	stor := repository.NewMetricStorage()

// 	h := NewHandler(stor, cfg, log)

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			w := httptest.NewRecorder()
// 			c, r := gin.CreateTestContext(w)
// 			r.POST("/", h.PostMetricsList)

// 			c.Request = httptest.NewRequest(http.MethodPost, "/updates", nil)
// 			c.Request.Body = metrics
// 			r.ServeHTTP(w, c.Request)

// 			result := w.Result()

// 			fmt.Println("111111111111", c.Request.Body)
// 			defer result.Body.Close()
// 			require.Equal(t, tt.status, result.StatusCode)

// 		})
// 	}
// }
