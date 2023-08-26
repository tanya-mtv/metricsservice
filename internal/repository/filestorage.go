package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/tanya-mtv/metricsservice/internal/config"
	"github.com/tanya-mtv/metricsservice/internal/logger"
	"github.com/tanya-mtv/metricsservice/internal/models"
)

type MetricRepositoryFiles struct {
	repository *MetricRepositoryStorage
	cfg        *config.ConfigServer
}

func NewMetricMetricRepositoryFiles(repository *MetricRepositoryStorage, cfg *config.ConfigServer) *MetricRepositoryFiles {

	return &MetricRepositoryFiles{
		repository: repository,
		cfg:        cfg,
	}
}
func (m *MetricRepositoryFiles) LoadLDataFromFile(log logger.Logger) {
	file, err := os.ReadFile(m.cfg.FileName)
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

func (m *MetricRepositoryFiles) SaveDataToFile(log logger.Logger) {
	dir, _ := path.Split(m.cfg.FileName)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0666)
		if err != nil {
			log.Error(err)
		}
	}
	pollTicker := time.NewTicker(time.Duration(m.cfg.Interval) * time.Second)
	defer pollTicker.Stop()
	for range pollTicker.C {
		m.save()
	}
}

func (m *MetricRepositoryFiles) save() error {

	metrics := m.repository.GetAll()

	data, err := json.MarshalIndent(metrics, "", "   ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.cfg.FileName, data, 0666)
}
