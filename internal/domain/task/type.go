package task

import "database/sql"

type Schema struct {
	ID          uint64       `db:"id"`
	Title       string       `db:"title"`
	Description string       `db:"description"`
	UpdatedAt   sql.NullTime `db:"updated_at"`
}
