CREATE SEQUENCE scout_id_seq;
CREATE TABLE scouts (
	id int PRIMARY KEY DEFAULT nextval('scout_id_seq'),
	uuid uuid NOT NULL,
	ip_address inet NOT NULL,
	authorised boolean DEFAULT false,
	calibration_frame bytea,
	name varchar(255) NOT NULL
);
ALTER SEQUENCE scout_id_seq OWNED BY scouts.id;