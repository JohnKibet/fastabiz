ALTER TABLE drivers
DROP COLUMN current_location;

ALTER TABLE drivers
ADD COLUMN current_location geometry(Point,4326) NOT NULL;