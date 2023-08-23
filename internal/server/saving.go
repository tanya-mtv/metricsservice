package server

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/tanya-mtv/metricsservice/internal/logger"

	"github.com/tanya-mtv/metricsservice/internal/config"

	"github.com/tanya-mtv/metricsservice/internal/models"
	"github.com/tanya-mtv/metricsservice/internal/repository"
)

func loadLDataFromFile(repo *repository.MetricRepositoryStorage, log logger.Logger, filePath string) {
	file, err := os.ReadFile(filePath)
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
			repo.UpdateCounter(metric.ID, *metric.Delta)
		case "gauge":
			repo.UpdateGauge(metric.ID, *metric.Value)
		}
	}
}

func saveDataToFile(repo *repository.MetricRepositoryStorage, log logger.Logger, cfg *config.ConfigServer) {
	dir, _ := path.Split(cfg.FileName)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0666)
		if err != nil {
			log.Error(err)
		}
	}
	pollTicker := time.NewTicker(time.Duration(cfg.Interval) * time.Second)
	defer pollTicker.Stop()
	for range pollTicker.C {
		save(repo, cfg.FileName)
	}
}

func save(repo *repository.MetricRepositoryStorage, filePath string) error {

	metrics := repo.GetAll()

	data, err := json.MarshalIndent(metrics, "", "   ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0666)
}
