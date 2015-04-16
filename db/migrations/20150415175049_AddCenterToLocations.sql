-- +goose Up
ALTER TABLE locations ADD COLUMN center POINT;
CREATE INDEX idx_locations_on_center ON locations USING GIST(center);

-- +goose Down
DROP INDEX idx_locations_on_center;
ALTER TABLE locations DROP COLUMN center;
