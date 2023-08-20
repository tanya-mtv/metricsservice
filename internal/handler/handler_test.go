package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
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

	var testTable = []struct {
		url    string
		want   string
		status int
	}{

		// проверим на ошибочный запрос
		{"/value/counter1/testSetGet40", "404 page not found", http.StatusNotFound},
		{"/value/gauge1/testSetGet40", "404 page not found", http.StatusNotFound},
	}
	for _, v := range testTable {
		resp, get := testRequest(t, router, "GET", v.url)

		defer resp.Body.Close()
		assert.Equal(t, v.status, resp.StatusCode)
		assert.Equal(t, v.want, get)
	}

	var testTable1 = []struct {
		url    string
		want   string
		status int
	}{

		// проверим на ошибочный запрос
		{"/update/counter1/testSetGet40/1", "404 page not found", http.StatusNotFound},
		// {"/update/counter/testSetGet40/1.36", "Metric not found", http.StatusBadRequest},
		{"/update/gauge1/testSetGet40/235", "404 page not found", http.StatusNotFound},
	}
	for _, v := range testTable1 {
		resp, get := testRequest(t, router, "GET", v.url)

		defer resp.Body.Close()
		assert.Equal(t, v.status, resp.StatusCode)
		assert.Equal(t, v.want, get)
	}
}
