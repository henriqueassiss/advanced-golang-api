package mock

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/henriqueassiss/advanced-golang-api/internal/domain/task"
	"github.com/henriqueassiss/advanced-golang-api/internal/domain/task/repository"
	"github.com/henriqueassiss/advanced-golang-api/internal/utils/errorMsg"
	"github.com/jmoiron/sqlx"
)

func verifyMockAuthenticity(db *sqlx.DB) error {
	var count uint64
	rows, err := db.Queryx(repository.Count)
	if err != nil {
		return err
	}

	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return err
		}
	}

	if count > 1 {
		return errorMsg.ErrTableIsPopulated
	}

	return nil
}

func Mock(db *sqlx.DB) error {
	err := verifyMockAuthenticity(db)
	if err != nil {
		return err
	}

	var ts []task.Schema
	for i := 0; i < 10; i++ {
		var t task.Schema

		t.Title = gofakeit.Sentence(3)
		t.Description = gofakeit.SentenceSimple()

		ts = append(ts, t)
	}

	_, err = db.NamedExec(`INSERT INTO tasks (title,
	description)
	VALUES (:title,
	:description)`, ts)

	return err
}
