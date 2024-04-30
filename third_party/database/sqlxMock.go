package database

import (
	"testing"

	"github.com/jmoiron/sqlx"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
)

func NewSqlxMock(t *testing.T) (*sqlx.DB, sqlxmock.Sqlmock) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatal(err)
	}

	return db, mock
}
