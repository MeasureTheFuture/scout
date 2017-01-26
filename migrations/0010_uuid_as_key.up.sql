ALTER TABLE scout_healths DROP CONSTRAINT scout_healths_scout_id_fkey;
ALTER TABLE scout_healths ADD COLUMN scout_uuid uuid;
UPDATE scout_healths T SET scout_uuid = (SELECT uuid FROM scouts WHERE id = T.scout_id);
ALTER TABLE scout_healths ALTER COLUMN scout_uuid SET NOT NULL;
ALTER TABLE scout_healths DROP COLUMN scout_id;

ALTER TABLE scout_interactions DROP CONSTRAINT scout_interactions_scout_id_fkey;
ALTER TABLE scout_interactions ADD COLUMN scout_uuid uuid;
UPDATE scout_interactions T SET scout_uuid = (SELECT uuid FROM scouts WHERE id = T.scout_id);
ALTER TABLE scout_interactions ALTER COLUMN scout_uuid SET NOT NULL;
ALTER TABLE scout_interactions DROP COLUMN scout_id;

ALTER TABLE scout_logs DROP CONSTRAINT scout_logs_scout_id_fkey;
ALTER TABLE scout_logs ADD COLUMN scout_uuid uuid;
UPDATE scout_logs T SET scout_uuid = (SELECT uuid FROM scouts WHERE id = T.scout_id);
ALTER TABLE scout_logs ALTER COLUMN scout_uuid SET NOT NULL;
ALTER TABLE scout_logs DROP COLUMN scout_id;

ALTER TABLE scout_summaries DROP CONSTRAINT scout_summaries_scout_id_fkey;
ALTER TABLE scout_summaries ADD COLUMN scout_uuid uuid;
UPDATE scout_summaries T SET scout_uuid = (SELECT uuid FROM scouts WHERE id = T.scout_id);
ALTER TABLE scout_summaries ALTER COLUMN scout_uuid SET NOT NULL;
ALTER TABLE scout_summaries DROP COLUMN scout_id;

ALTER TABLE scouts DROP CONSTRAINT scouts_pkey;
ALTER TABLE scouts ADD PRIMARY KEY (uuid);
ALTER TABLE scouts DROP COLUMN id;