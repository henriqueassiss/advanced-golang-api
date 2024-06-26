BEGIN;

CREATE TABLE IF NOT EXISTS tasks(
	id          BIGSERIAL PRIMARY KEY,
	title       TEXT NOT NULL,
	description TEXT,
	updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMIT;