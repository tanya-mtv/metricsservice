package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/tanya-mtv/metricsservice/internal/constants"

	_ "github.com/jackc/pgx/v5/stdlib"
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
	var value int64

	query := "INSERT INTO metrics as m (id, mtype, delta, value) VALUES ($1, $2, $3, $4) ON CONFLICT (id)  DO UPDATE SET delta = (m.delta + EXCLUDED.delta) returning delta"
	row := m.db.QueryRow(query, n, "counter", v, 0)
	if err := row.Scan(&value); err != nil {
		m.log.Error("Can not scan counter value in update function ", err)
		return Counter(0)
	}
	return Counter(value)
	// retrier := NewRetrier()
	// for _, val := range retrier.retries {
	// 	needsR := m.db.Ping()

	// 	if haveToRetry(needsR) {
	// 		time.Sleep(val)
	// 	} else {

	// 		row := m.db.QueryRow(query, n, "counter", v, 0)
	// 		if err := row.Scan(&value); err != nil {
	// 			m.log.Error("Can not scan counter value in update function ", err)
	// 			return Counter(0)
	// 		}
	// 		return Counter(value)
	// 	}
	// }

	// return Counter(0)

}

func (m *DBStorage) UpdateGauge(n string, v float64) Gauge {

	var value float64

	query := "INSERT INTO metrics as m (id, mtype, delta, value) VALUES ($1, $2, $3, $4) ON CONFLICT (id)  DO UPDATE SET value =  EXCLUDED.value returning value"
	row := m.db.QueryRow(query, n, "gauge", 0, v)
	if err := row.Scan(&value); err != nil {
		m.log.Error("Can not scan counter value in update function ", err)
		return 0
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

func (m *DBStorage) GetCounter(metricName string) (Counter, bool) {
	var cnt int64
	query := "SELECT delta from metrics WHERE ID = $1"
	err := m.db.Get(&cnt, query, metricName)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			m.log.Error("Now rows in DB with metric name ", metricName)
			return 0, false
		} else {
			m.log.Error("Can't get counter from DB ", metricName)
			return 0, false
		}
	}
	return Counter(cnt), true
}

func (m *DBStorage) GetGauge(metricName string) (Gauge, bool) {
	var gug float64
	query := "SELECT value from metrics WHERE ID = $1"
	err := m.db.Get(&gug, query, metricName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			m.log.Error("Now rows in DB with metric name ", metricName)
			return 0, false
		} else {
			m.log.Error("Can't get Gauge from DB ", metricName)
			return 0, false
		}

	}
	return Gauge(gug), true
}

func (m *DBStorage) UpdateMetrics(metrics []*models.Metrics) ([]*models.Metrics, error) {

	tx, err := m.db.Begin()
	if err != nil {
		return metrics, err
	}
	// можно вызвать Rollback в defer,
	// если Commit будет раньше, то откат проигнорируется
	defer tx.Rollback()

	stmt, err := tx.Prepare(
		"INSERT INTO metrics (id, mtype, delta, value) values ($1, $2, $3, $4)" +
			"ON CONFLICT (id) DO UPDATE SET delta = $5,  value=$6")

	if err != nil {
		return metrics, err
	}
	defer stmt.Close()

	for _, v := range metrics {
		switch v.MType {
		case "counter":
			cnt := m.UpdateCounter(v.ID, *v.Delta)
			m.log.Debug("Update counter. New value is ", cnt)
			tmp := int64(cnt)
			v.Delta = &tmp

		case "gauge":
			_, err = stmt.Exec(v.ID, v.MType, 0, *v.Value, 0, *v.Value)

		}
		if err != nil {
			return metrics, err
		}
	}

	tx.Commit()

	return metrics, nil
}

type Retrier struct {
	retries []time.Duration
}

func NewRetrier() Retrier {
	return Retrier{
		retries: []time.Duration{0, constants.RetryWaitMin, constants.RetryMedium, constants.RetryWaitMax},
	}
}

// func haveToRetry(err error) bool {
// 	var pge *pgconn.PgError
// 	var nete *net.OpError

// 	if errors.As(err, &pge) {
// 		return true
// 	}
// 	if errors.Is(err, nete) {
// 		return true
// 	}
// 	return false
// }
