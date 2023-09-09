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

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`            // Параметр кодирую строкой, принося производительность в угоду наглядности.
	Delta *int64   `json:"delta,omitempty"` //counter
	Value *float64 `json:"value,omitempty"` //gauge

}

// func TestHandler_PostMetricsList1(t *testing.T) {
// 	httpc := resty.New().SetHostURL("localhost:8080")

// 	idCounter := "CounterBatchZip1"
// 	valueCounter1, valueCounter2 := int64(1), int64(10)
// 	req := httpc.R().
// 		SetHeader("Accept-Encoding", "gzip").
// 		SetHeader("Content-Type", "application/json")

// 	metrics := []Metrics{
// 		{
// 			ID:    idCounter,
// 			MType: "counter",
// 			Delta: &valueCounter1,
// 		},
// 		{
// 			ID:    idCounter,
// 			MType: "counter",
// 			Delta: &valueCounter2,
// 		},
// 	}

// 	resp, err := req.SetBody(metrics).
// 		Post("/updates/")
// 	fmt.Printf("Data of req.Body %+v \n", req.Body)
// 	fmt.Printf("Type of req.Body %T \n", req.Body)
// 	fmt.Println(" ", resp, err)
// }
// func TestHandler_PostMetricsList(t *testing.T) {

// 	metrics := `[{
// 	    "ID":"idCounter",
// 	    "Mtype":"counter",
// 	    "Delta": 10
// 	}]`

// 	tests := []struct {
// 		name     string
// 		sentdata string
// 		status   int
// 		wantbody string
// 	}{
// 		{
// 			name:     "Test post list metrics json",
// 			sentdata: metrics,
// 			status:   http.StatusOK,
// 			wantbody: "Metrics was read",
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
// 			r.POST("/updates/", h.PostMetricsList)
// 			buf := strings.NewReader(tt.sentdata)

// 			c.Request = httptest.NewRequest(http.MethodPost, "/updates/", buf)

// 			r.ServeHTTP(w, c.Request)

// 			result := w.Result()

// 			defer result.Body.Close()
// 			require.Equal(t, tt.status, result.StatusCode)

// 		})
// 	}
// }
