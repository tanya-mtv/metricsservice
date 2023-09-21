package server

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/tanya-mtv/metricsservice/internal/repository"

	"github.com/tanya-mtv/metricsservice/internal/models"
)

func (s *server) openStorage(ctx context.Context, db *sqlx.DB) {

	if db == nil {
		if s.cfg.FileName != "" {

			s.stor = repository.NewMetricFiles()
			if s.cfg.Restore {
				s.LoadLDataFromFile(s.cfg.FileName)
			}
			if s.cfg.Interval != 0 {
				go s.SaveDataToFile(ctx, s.cfg.FileName, s.cfg.Interval)
			}
		} else {
			s.stor = repository.NewMetricStorage()

		}
	} else {

		s.stor = repository.NewDBStorage(db, s.log)

	}

}

func (s *server) LoadLDataFromFile(fileName string) {

	var file []byte
	var err error

	for _, val := range s.ret.Retries {
		file, err = os.ReadFile(fileName)
		if s.ret.Next(err, val) {
			s.log.Error("Can not read file ", err)
		} else {
			break
		}

	}

	data := make([]models.Metrics, 0, 29)

	if err := json.Unmarshal(file, &data); err != nil {
		fmt.Println(err)
	}

	for _, metric := range data {

		switch metric.MType {
		case "counter":
			s.stor.UpdateCounter(metric.ID, *metric.Delta)
		case "gauge":
			s.stor.UpdateGauge(metric.ID, *metric.Value)
		}
	}
}

func (s *server) SaveDataToFile(ctx context.Context, fileName string, interval int) {
	dir, _ := path.Split(fileName)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0666)
		if err != nil {
			s.log.Error(err)
		}
	}
	pollTicker := time.NewTicker(time.Duration(interval) * time.Second)
	defer pollTicker.Stop()
	for {
		select {
		case <-pollTicker.C:
			_ = s.save(fileName)
		case <-ctx.Done():
			return
		}
	}

}

func (s *server) save(fileName string) error {

	metrics := s.stor.GetAll()

	data, err := json.MarshalIndent(metrics, "", "   ")
	if err != nil {
		return err
	}

	return os.WriteFile(fileName, data, 0666)
}
