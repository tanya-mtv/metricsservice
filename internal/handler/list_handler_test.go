package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tanya-mtv/metricsservice/internal/repository"
)

func testRequest(t *testing.T, ts *httptest.Server, method,
	path string) (*http.Response, string) {

	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestRouter(t *testing.T) {
	repos := &repository.Repository{
		MetricStorage: repository.NewMetricRepository(),
	}

	handl := &Handler{
		repository: repos,
		router:     gin.New(),
	}
	ts := httptest.NewServer(handl.InitRoutes())
	defer ts.Close()

	var testTableStatus = []struct {
		url    string
		want   string
		status int
	}{
		{"/value/counter1/testSetGet40", "404 page not found", http.StatusNotFound},
		{"/value/gauge1/testSetGet40", "404 page not found", http.StatusNotFound},
		{"/update/counter1/testSetGet40/1", "404 page not found", http.StatusNotFound},
		{"/update/gauge1/testSetGet40/235", "404 page not found", http.StatusNotFound},
	}
	for _, v := range testTableStatus {
		resp, get := testRequest(t, ts, "GET", v.url)
		defer resp.Body.Close()

		assert.Equal(t, v.status, resp.StatusCode)
		assert.Equal(t, v.want, get)
	}

	var testTablePost = []struct {
		url    string
		want   string
		status int
	}{

		{"/update/counter/testSetGet40/1", "1", http.StatusOK},
		{"/update/counter/testSetGet40/1", "2", http.StatusOK},
		{"/update/counter/testSetGet40/1.36", "{\"message\":\"invalid metricValue param\"}", http.StatusBadRequest},

		{"/update/gauge/testSetGet40/236.66", "236.66", http.StatusOK},
		{"/update/gauge/testSetGet40/105.66", "105.66", http.StatusOK},
	}
	for _, v := range testTablePost {
		resp, get := testRequest(t, ts, "POST", v.url)
		defer resp.Body.Close()

		assert.Equal(t, v.status, resp.StatusCode)
		assert.Equal(t, v.want, get)
	}

	var testTableGET = []struct {
		url    string
		want   string
		status int
	}{
		{"/value/counter/testSetGet40", "2", http.StatusOK},
		{"/value/gauge/testSetGet40", "105.66", http.StatusOK},
	}
	for _, v := range testTableGET {
		resp, get := testRequest(t, ts, "GET", v.url)
		defer resp.Body.Close()

		assert.Equal(t, v.status, resp.StatusCode)
		assert.Equal(t, v.want, get)
	}
}
