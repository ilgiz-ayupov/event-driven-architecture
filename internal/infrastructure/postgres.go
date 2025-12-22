package infrastructure

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgres(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// проверить соединение
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
