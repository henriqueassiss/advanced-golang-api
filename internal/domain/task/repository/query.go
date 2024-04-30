package repository

var (
	Count = `SELECT count(*) FROM tasks`

	Select = `SELECT ? FROM tasks t`

	InsertInto = `INSERT INTO tasks (?) VALUES (?) RETURNING id`

	Update = `UPDATE tasks SET ?, updated_at = CURRENT_TIMESTAMP WHERE id = $1`

	Delete = `DELETE FROM tasks WHERE id = $1`
)
