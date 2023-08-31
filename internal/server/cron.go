package server

import (
	"context"

	"github.com/tanya-mtv/metricsservice/internal/fileservice"

	"github.com/tanya-mtv/metricsservice/internal/logger"
	"github.com/tanya-mtv/metricsservice/internal/repository"
)



func openStorage(ctx context.Context, stor *repository.MetricStorage, filename string, interval int, restore bool, log logger.Logger) *fileservice.FilesStorage {

	st := fileservice.NewFilesStorage(stor, filename, interval, log)
	if filename != "" {
		if restore {
			st.LoadLDataFromFile()
		}
		if interval != 0 {
			go st.SaveDataToFile(ctx)
		}
	}

	return st

}
