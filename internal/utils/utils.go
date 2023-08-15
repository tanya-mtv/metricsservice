package utils

import "strconv"



func CounterToBytes(number Counter) []byte {

	b := []byte(strconv.FormatInt(int64(number), 10))

	return b
}

func GaugeToBytes(number Gauge) []byte {

	b := []byte(strconv.FormatFloat(float64(number), 'f', -1, 64))

	return b
}
