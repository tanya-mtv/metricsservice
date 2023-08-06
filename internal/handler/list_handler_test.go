package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method,
	// path string) (*http.Response, string) {
	path string) *http.Response {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	fmt.Println("respBody", string(respBody))
	return resp
	//, string(respBody)
}

func TestRouter(t *testing.T) {
	handl := &Handler{}
	ts := httptest.NewServer(handl.InitRoutes())
	defer ts.Close()

	var testTable = []struct {
		url    string
		want   string
		status int
	}{

		// проверим на ошибочный запрос
		{"/value/counter1/testSetGet40", "Metric not found", http.StatusNotFound},
		{"/value/gauge1/testSetGet40", "Metric not found", http.StatusNotFound},
	}
	for _, v := range testTable {
		// resp, get := testRequest(t, ts, "GET", v.url)
		resp := testRequest(t, ts, "GET", v.url)
		defer resp.Body.Close()
		assert.Equal(t, v.status, resp.StatusCode)
		// assert.Equal(t, v.want, get)
	}

	var testTable1 = []struct {
		url    string
		want   string
		status int
	}{

		// проверим на ошибочный запрос
		{"/update/counter1/testSetGet40/1", "", http.StatusBadRequest},
		// {"/update/counter/testSetGet40/1.36", "Metric not found", http.StatusBadRequest},
		// {"/update/gauge1/testSetGet40/235", "Metric not found", http.StatusBadRequest},
	}
	for _, v := range testTable1 {
		// resp, get := testRequest(t, ts, "GET", v.url)
		resp := testRequest(t, ts, "POST", v.url)
		defer resp.Body.Close()
		fmt.Println("resp.StatusCode ", resp.StatusCode)
		assert.Equal(t, v.status, resp.StatusCode)
		// assert.Equal(t, v.want, get)
	}
}
