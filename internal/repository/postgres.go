package repository

import (
	"github.com/jmoiron/sqlx"
)

type DBStorage struct {
	db *sqlx.DB
}

func (d *DBStorage) Ping() error {
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
	return db, nil
}
