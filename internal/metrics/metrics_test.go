package metrics

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tanya-mtv/metricsservice/internal/models"

	"github.com/tanya-mtv/metricsservice/internal/constants"
	"github.com/tanya-mtv/metricsservice/internal/logger"

	"github.com/stretchr/testify/assert"
	"github.com/tanya-mtv/metricsservice/internal/config"
	"github.com/tanya-mtv/metricsservice/internal/repository"
)

func TestServiceMetrics_Post(t *testing.T) {
	repos := &repository.MetricRepositoryCollector{}

	sm := NewServiceMetrics(&config.ConfigAgent{}, repos)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	}))
	addr := server.URL + "/update/"
	fmt.Println("addr ", addr)

	tmp1 := float64(1.222233)
	metric1 := &models.Metrics{
		ID:    "Alloc",
		MType: "gauge",
		Value: &tmp1,
	}

	tmp2 := int64(1)
	metric2 := &models.Metrics{
		ID:    "pollCount",
		MType: "counter",
		Delta: &tmp2,
	}

	var tests = []struct {
		nameTest     string
		body         *models.Metrics
		expectedBody string
	}{
		{"Post method gauge", metric1, ""},
		{"Post method counter", metric2, ""},
	}

	cfglog := &logger.Config{
		LogLevel: constants.LogLevel,
		DevMode:  constants.DevMode,
		Type:     constants.Type,
	}
	log := logger.NewAppLogger(cfglog)

	for _, tt := range tests {
		t.Run(tt.nameTest, func(t *testing.T) {

			_, err := sm.Post(tt.body, addr, log)

			assert.NoError(t, err, "error making HTTP request")

			if tt.expectedBody != "" {
				assert.NoError(t, err, "error making HTTP request")
			}
		})
	}
}
