CREATE TYPE scout_state AS ENUM ('idle', 'calibrating', 'calibrated', 'measuring');
ALTER TABLE scouts ADD COLUMN state scout_state NOT NULL DEFAULT 'idle';