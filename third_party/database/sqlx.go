package database

import (
	"errors"
	"fmt"

	"github.com/henriqueassiss/advanced-golang-api/config"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func NewSqlx(cfg config.Database) (db *sqlx.DB, err error) {
	var dsn string
	switch cfg.Driver {
	case "postgres", "pgx":
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=America/Sao_Paulo",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SslMode)
	default:
		return db, errors.New("must choose a database driver")
	}

	db, err = sqlx.Open(cfg.Driver, dsn)
	if err != nil {
		return db, err
	}

	err = db.Ping()

	return db, err
}
