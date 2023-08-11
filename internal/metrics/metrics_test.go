package metrics

import (
	"net/http"
	"testing"

	"github.com/tanya-mtv/metricsservice/internal/config"
	"github.com/tanya-mtv/metricsservice/internal/repository"
	"github.com/tanya-mtv/metricsservice/internal/servise"
)

func TestServiceMetrics_Post(t *testing.T) {
	repos := &repository.Repository{
		MetricStorageAgent: repository.NewMetricStorageAgent(),
	}
	serv := servise.NewServise(repos)
	sm := NewServiceMetrics(&config.ConfigAgent{}, serv)

	// }
	addr := "http://localhost:8080/update"
	var tests = []struct {
		nameTest    string
		metricType  string
		metricName  string
		metricValue string
		// want        string
		status int
	}{
		{"Post method gauge", "gauge", "Mallocs", "1277", http.StatusOK},
		{"Post method counter", "counter", "PollCount", "15", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.nameTest, func(t *testing.T) {
			sm.Post(tt.metricType, tt.metricName, tt.metricValue, addr)
		})
	}
}
