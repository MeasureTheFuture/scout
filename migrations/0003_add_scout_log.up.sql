CREATE TABLE scout_logs (
	scout_id int REFERENCES scouts(id),
	log bytea NOT NULL,
	created_at timestamp NOT NULL
);
CREATE INDEX scout_logs_idx ON scout_logs (created_at);