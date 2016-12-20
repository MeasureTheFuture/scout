CREATE SEQUENCE scout_interaction_id_seq;
CREATE TABLE scout_interactions (
	id int PRIMARY KEY DEFAULT nextval('scout_interaction_id_seq'),
	scout_id int REFERENCES scouts(id),
	duration real NOT NULL,
	waypoints path NOT NULL,
	waypoint_widths path NOT NULL,
	waypoint_times real[] NOT NULL,
	processed boolean NOT NULL DEFAULT false,
	entered_at timestamp NOT NULL,
	created_at timestamp NOT NULL
);
ALTER SEQUENCE scout_interaction_id_seq OWNED BY scout_interactions.id;
CREATE INDEX scout_interactions_idx ON scout_interactions (created_at);