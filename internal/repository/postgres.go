package repository

import (
	"database/sql"
	"errors"
	"net"
	"time"

	"github.com/tanya-mtv/metricsservice/internal/constants"

	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/tanya-mtv/metricsservice/internal/logger"
	"github.com/tanya-mtv/metricsservice/internal/models"
)

func NewPostgresDB(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS metrics (id serial PRIMARY KEY, name varchar(30), mtype char(30),  delta bigint, value double precision, UNIQUE (name));")

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

	query := "INSERT INTO metrics as m (name, mtype, delta, value) VALUES ($1, $2, $3, $4) ON CONFLICT (name)  DO UPDATE SET delta = (m.delta + EXCLUDED.delta) returning delta"
	retrier := NewRetrier()
	for _, val := range retrier.retries {
		err := m.db.Ping()
		//check type of err
		if haveToRetry(err) {
			time.Sleep(val)
		} else {
			row := m.db.QueryRow(query, n, "counter", v, 0)
			if err := row.Scan(&value); err != nil {
				m.log.Error("Can not scan counter value in update function ", err)
				return Counter(v)
			}
		}
	}

	return Counter(value)

}

func (m *DBStorage) UpdateGauge(n string, v float64) Gauge {

	var value float64

	query := "INSERT INTO metrics as m (name, mtype, delta, value) VALUES ($1, $2, $3, $4) ON CONFLICT (name)  DO UPDATE SET value =  EXCLUDED.value returning value"

	retrier := NewRetrier()
	for _, val := range retrier.retries {
		err := m.db.Ping()
		//check type of err
		if haveToRetry(err) {
			time.Sleep(val)
		} else {

			row := m.db.QueryRow(query, n, "gauge", 0, v)
			if err := row.Scan(&value); err != nil {
				m.log.Error("Can not scan counter value in update function ", err)
				return Gauge(v)
			}
			return Gauge(value)
		}
	}

	return Gauge(v)
}

func (m *DBStorage) GetAll() []models.Metrics {
	metricsSlice := make([]models.Metrics, 0, 29)
	query := "SELECT  name id, mtype, delta, value from metrics"

	retrier := NewRetrier()
	for _, val := range retrier.retries {
		err := m.db.Ping()
		//check type of err
		if haveToRetry(err) {
			time.Sleep(val)
		} else {

			err := m.db.Select(&metricsSlice, query)
			if err != nil {
				m.log.Error("Can't get all metric ")
				return metricsSlice
			}

			return metricsSlice
		}
	}

	return metricsSlice
}

func (m *DBStorage) GetCounter(metricName string) (Counter, bool) {
	var cnt int64
	query := "SELECT delta from metrics WHERE name = $1"
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
	query := "SELECT value from metrics WHERE name = $1"
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
		"INSERT INTO metrics (name, mtype, delta, value) values ($1, $2, $3, $4)" +
			"ON CONFLICT (name) DO UPDATE SET delta = $5,  value=$6")

	if err != nil {
		return metrics, err
	}

	stmtc, err := tx.Prepare(
		"INSERT INTO metrics as m (name, mtype, delta, value) VALUES ($1, $2, $3, $4)" +
			" ON CONFLICT (name)  DO UPDATE SET delta = (m.delta + EXCLUDED.delta)")

	if err != nil {
		return metrics, err
	}
	defer stmt.Close()

	for _, v := range metrics {
		switch v.MType {
		case "counter":
			_, err = stmtc.Exec(v.ID, v.MType, *v.Delta, 0)

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

func haveToRetry(err error) bool {
	var pge *pgconn.PgError
	var nete *net.OpError

	if errors.As(err, &pge) {
		return true
	}
	if errors.Is(err, nete) {
		return true
	}
	return false
}
