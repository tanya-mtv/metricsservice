package repository

import (
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/tanya-mtv/metricsservice/internal/logger"
	"github.com/tanya-mtv/metricsservice/internal/models"
)

type DBPing struct {
	db *sqlx.DB
}

func (d *DBPing) Ping() error {
	err := d.db.Ping()
	if err != nil {

		return err
	}
	return nil
}

func NewPostgresDB(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS metrics (ID char(30) UNIQUE, mtype char(30),  delta integer, value double precision);")

	if err != nil {
		return db, err
	}

	return db, nil
}

type DBStorage struct {
	db  *sqlx.DB
	log logger.Logger
}

func NewDBStorage(db *sqlx.DB, log logger.Logger) *DBStorage {

	return &DBStorage{
		db:  db,
		log: log,
	}
}

func (m *DBStorage) UpdateCounter(n string, v int64) Counter {

	var id string
	var delta int64
	query := "SELECT id, delta from metrics WHERE ID = $1"
	row := m.db.QueryRow(query, n)

	row.Scan(&id, &delta)

	if id == "" {
		query = "INSERT INTO metrics (id, mtype, delta, value) values ($1, $2, $3, $4) returning delta"
		row = m.db.QueryRow(query, n, "counter", v, 0)
		if err := row.Scan(&delta); err != nil {
			return 0
		} else {
			return Counter(v)
		}
	}

	newValue := v + delta

	query = "UPDATE metrics set delta = $1 WHERE id = $2"
	_, err := m.db.Exec(query, newValue, n)

	if err != nil {
		m.log.Error("Can't update data from DB with metric ", n)
		return Counter(v)
	}

	return Counter(newValue)
}

func (m *DBStorage) UpdateGauge(n string, v float64) Gauge {
	var id string
	var value float64
	query := "SELECT id, value from metrics WHERE ID = $1"
	row := m.db.QueryRow(query, n)

	row.Scan(&id, &value)

	if id == "" {
		query = "INSERT INTO metrics (id, mtype, delta, value) values ($1, $2, $3, $4) returning value"
		row = m.db.QueryRow(query, n, "gauge", 0, v)
		if err := row.Scan(&value); err != nil {
			return 0
		} else {
			return Gauge(v)
		}
	}

	query = "UPDATE metrics set value = $1 WHERE id = $2"
	_, err := m.db.Exec(query, v, n)

	if err != nil {
		m.log.Error("Can't update data from DB with metric ", n)
		return Gauge(v)
	}

	return Gauge(v)
}

func (m *DBStorage) GetAll() []models.Metrics {
	metricsSlice := make([]models.Metrics, 0, 29)
	query := "SELECT id, mtype, delta, value from metrics"

	err := m.db.Select(&metricsSlice, query)
	if err != nil {
		m.log.Error("Can't get all metric ")
		return metricsSlice
	}

	return metricsSlice
}

func (m *DBStorage) GetAllCounter() map[string]int64 {
	metricsSlice := make([]models.Metrics, 0, 29)
	elementMap := make(map[string]int64, 0)
	query := "SELECT id, mtype, delta, value from metrics where mtype='counter'"

	err := m.db.Select(&metricsSlice, query)
	if err != nil {
		m.log.Error("Can't get all metric ")
		return elementMap
	}

	for _, value := range metricsSlice {
		name := strings.TrimSpace(value.ID)
		elementMap[name] = *value.Delta
	}
	return elementMap
}

func (m *DBStorage) GetCounter(metricName string) (Counter, bool) {
	var cnt int64
	query := "SELECT delta from metrics WHERE ID = $1"
	err := m.db.Get(&cnt, query, metricName)
	if err != nil {
		m.log.Error("Can't get counter from DB ")
		return 0, false
	}
	return Counter(cnt), true
}

func (m *DBStorage) GetGauge(metricName string) (Gauge, bool) {
	var gug float64
	query := "SELECT value from metrics WHERE ID = $1"
	err := m.db.Get(&gug, query, metricName)
	if err != nil {
		m.log.Error("Can't get counter from DB ")
		return 0, false
	}
	return Gauge(gug), true
}

func (m *DBStorage) UpdateMetrics(metrics []models.Metrics) error {
	cntMap := m.GetAllCounter()

	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	// можно вызвать Rollback в defer,
	// если Commit будет раньше, то откат проигнорируется
	defer tx.Rollback()

	stmt, err := tx.Prepare(
		"INSERT INTO metrics (id, mtype, delta, value) values ($1, $2, $3, $4)" +
			"ON CONFLICT (id) DO UPDATE SET delta = $5,  value=$6")

	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, v := range metrics {
		switch v.MType {
		case "counter":
			newDelta := cntMap[v.ID] + *v.Delta
			_, err = stmt.Exec(v.ID, v.MType, newDelta, 0, newDelta, 0)
		case "gauge":
			_, err = stmt.Exec(v.ID, v.MType, 0, *v.Value, 0, *v.Value)
		}
		if err != nil {
			return err
		}
	}
	tx.Commit()
	return nil
}
