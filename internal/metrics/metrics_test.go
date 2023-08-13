package metrics

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tanya-mtv/metricsservice/internal/config"
	"github.com/tanya-mtv/metricsservice/internal/repository"
)

func TestServiceMetrics_Post(t *testing.T) {
	repos := &repository.Repository{
		MetricRepositoryAgent: repository.NewMetricRepositoryAgent(),
	}

	sm := NewServiceMetrics(&config.ConfigAgent{}, repos)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	}))
	addr := server.URL + "/update/"
	fmt.Println("addr ", addr)
	var tests = []struct {
		nameTest     string
		metricType   string
		metricName   string
		metricValue  string
		expectedBody string
	}{
		{"Post method gauge", "gauge", "Mallocs", "1277", "1277"},
		{"Post method counter", "counter", "PollCount", "15", "15"},
	}

	for _, tt := range tests {
		t.Run(tt.nameTest, func(t *testing.T) {
			_, err := sm.Post(tt.metricType, tt.metricName, tt.metricValue, addr)

			assert.NoError(t, err, "error making HTTP request")

			if tt.expectedBody != "" {
				assert.NoError(t, err, "error making HTTP request")
			}
		})
	}
}
