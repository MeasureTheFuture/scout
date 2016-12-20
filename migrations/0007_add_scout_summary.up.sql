CREATE TABLE scout_summaries (
	scout_id int references scouts(id),
	visitor_count int NOT NULL DEFAULT 0,
	visit_time_buckets real[20][20] NOT NULL
);