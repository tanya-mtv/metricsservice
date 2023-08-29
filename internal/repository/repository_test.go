package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateCounter(t *testing.T) {
	mem := NewMetricStorage()

	testStruct := []struct {
		testname   string
		metricname string
		value      int64
		result     int64
	}{
		{testname: "TestCounter1", metricname: "CounterMenric1", value: 1, result: 1},
		{testname: "TestCounter2", metricname: "CounterMenric1", value: 1, result: 2},
		{testname: "TestCounter3", metricname: "CounterMenric1", value: 100, result: 102},
		{testname: "TestCounter4", metricname: "CounterMenric2", value: 1, result: 1},
		{testname: "TestCounter5", metricname: "CounterMenric2", value: 0, result: 1},
	}

	for _, test := range testStruct {
		t.Run(test.testname, func(t *testing.T) {
			mem.UpdateCounter(test.metricname, test.value)
			assert.Equal(t, Counter(test.result), mem.counterData[test.metricname])
		})

	}
}

func TestUpdateGauge(t *testing.T) {
	mem := NewMetricStorage()

	testStruct := []struct {
		testname   string
		metricname string
		value      float64
		result     float64
	}{
		{testname: "TestGauge1", metricname: "GaugeMenric1", value: 1.25, result: 1.25},
		{testname: "TestGauge2", metricname: "GaugeMenric1", value: 1.26, result: 1.26},
		{testname: "TestGauge3", metricname: "GaugeMenric2", value: 100.00, result: 100.00},
		{testname: "TestGauge4", metricname: "GaugeMenric2", value: 0, result: 0.00},
	}

	for _, test := range testStruct {
		t.Run(test.testname, func(t *testing.T) {
			mem.UpdateGauge(test.metricname, test.value)
			assert.Equal(t, Gauge(test.result), Gauge(mem.gaugeData[test.metricname]))
		})

	}

}
