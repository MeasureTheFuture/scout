CREATE TABLE scout_healths (
	scout_id int REFERENCES scouts(id),
	cpu real NOT NULL,
	memory real NOT NULL,
	total_memory real NOT NULL,
	storage real NOT NULL,
	created_at timestamp NOT NULL
);
CREATE INDEX scout_healths_idx ON scout_healths (created_at);
