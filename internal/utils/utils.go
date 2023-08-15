package utils

import (
	"strconv"

	"github.com/tanya-mtv/metricsservice/internal/repository"
)

func CounterToBytes(number repository.Counter) []byte {

	b := []byte(strconv.FormatInt(int64(number), 10))

	return b
}

func GaugeToBytes(number repository.Gauge) []byte {

	b := []byte(strconv.FormatFloat(float64(number), 'f', -1, 64))

	return b
}
