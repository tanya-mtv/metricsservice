package metrics

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tanya-mtv/metricsservice/internal/config"
	"github.com/tanya-mtv/metricsservice/internal/constants"
	"github.com/tanya-mtv/metricsservice/internal/logger"
	"github.com/tanya-mtv/metricsservice/internal/models"
	"github.com/tanya-mtv/metricsservice/internal/repository"
)

func TestServiceMetrics_Post(t *testing.T) {

	cfglog := &logger.Config{
		LogLevel: constants.LogLevel,
		DevMode:  constants.DevMode,
		Type:     constants.Type,
	}
	log := logger.NewAppLogger(cfglog)

	repo := repository.NewMetricRepositoryCollector()
	sm := NewServiceMetrics(repo, &config.ConfigAgent{}, log)

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
		body         []*models.Metrics
		expectedBody string
	}{
		{"Post method gauge", []*models.Metrics{metric1}, ""},
		{"Post method counter", []*models.Metrics{metric2}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.nameTest, func(t *testing.T) {

			_, err := sm.PostJSON(tt.body, addr)

			require.NoError(t, err, "error making HTTP request")

			if tt.expectedBody != "" {
				require.NoError(t, err, "error making HTTP request")
			}
		})
	}
}

func TestServiceMetrics_Compression(t *testing.T) {
	type args struct {
		log logger.Logger
		b   []byte
	}

	cfglog := &logger.Config{
		LogLevel: constants.LogLevel,
		DevMode:  constants.DevMode,
		Type:     constants.Type,
	}

	log := logger.NewAppLogger(cfglog)
	cfg := &config.ConfigAgent{}
	repo := repository.NewMetricRepositoryCollector()
	sm := NewServiceMetrics(repo, cfg, log)

	metric := newMetric("Alloc", "gauge")
	tmp := float64(1798344)
	metric.Value = &tmp
	data, _ := json.Marshal(&metric)

	tests := []struct {
		name    string
		sm      *ServiceMetrics
		args    args
		wantErr bool
	}{
		{
			name:    "Test comprassion",
			sm:      sm,
			args:    args{log, data},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.sm.Compression(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("ServiceMetrics.Compression() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
