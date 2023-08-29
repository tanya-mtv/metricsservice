package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/tanya-mtv/metricsservice/internal/logger"
	"github.com/tanya-mtv/metricsservice/internal/models"
)

type MetricRepositoryFiles struct {
	repository *MetricStorage
	fileName   string
	interval   int
}

func NewMetricMetricRepositoryFiles(repository *MetricStorage, fileName string, interval int) *MetricRepositoryFiles {

	return &MetricRepositoryFiles{
		repository: repository,
		fileName:   fileName,
		interval:   interval,
	}
}
func (m *MetricRepositoryFiles) LoadLDataFromFile() {
	file, err := os.ReadFile(m.fileName)
	if err != nil {
		fmt.Println(err)
	}

	data := make([]models.Metrics, 0, 29)

	if err := json.Unmarshal(file, &data); err != nil {
		fmt.Println(err)
	}

	for _, metric := range data {
		switch metric.MType {
		case "counter":
			m.repository.UpdateCounter(metric.ID, *metric.Delta)
		case "gauge":
			m.repository.UpdateGauge(metric.ID, *metric.Value)
		}
	}
}

func (m *MetricRepositoryFiles) SaveDataToFile(log logger.Logger, ctx context.Context) {
	dir, _ := path.Split(m.fileName)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0666)
		if err != nil {
			log.Error(err)
		}
	}
	pollTicker := time.NewTicker(time.Duration(m.interval) * time.Second)
	defer pollTicker.Stop()
	for {
		select {
		case <-pollTicker.C:
			m.save()
		case <-ctx.Done():
			return
		}
	}

}

func (m *MetricRepositoryFiles) save() error {

	metrics := m.repository.GetAll()

	data, err := json.MarshalIndent(metrics, "", "   ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.fileName, data, 0666)
}
