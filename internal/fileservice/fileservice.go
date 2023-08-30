package fileservice

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/tanya-mtv/metricsservice/internal/repository"

	"github.com/tanya-mtv/metricsservice/internal/logger"
	"github.com/tanya-mtv/metricsservice/internal/models"
)

type FilesService struct {
	storage  *repository.Storage
	fileName string
	interval int
	log      logger.Logger
}

func NewMetricFiles(storage *repository.Storage, fileName string, interval int, log logger.Logger) *FilesService {

	return &FilesService{
		storage:  storage,
		fileName: fileName,
		interval: interval,
		log:      log,
	}
}

func (m *FilesService) LoadLDataFromFile() {

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
			m.storage.UpdateCounter(metric.ID, *metric.Delta)
		case "gauge":
			m.storage.UpdateGauge(metric.ID, *metric.Value)
		}
	}
}

func (m *FilesService) SaveDataToFile(ctx context.Context) {
	dir, _ := path.Split(m.fileName)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0666)
		if err != nil {
			m.log.Error(err)
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

func (m *FilesService) save() error {

	metrics := m.storage.GetAll()

	data, err := json.MarshalIndent(metrics, "", "   ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.fileName, data, 0666)
}
