-- +goose Up
CREATE UNIQUE INDEX idx_locations_on_id     ON locations (id);
CREATE INDEX idx_locations_on_parent_id     ON locations (parent_id);
CREATE INDEX idx_locations_on_ancestors_ids ON locations (ancestors_ids);
CREATE INDEX idx_locations_on_level         ON locations (level);

-- +goose Down
DROP INDEX idx_locations_on_id;
DROP INDEX idx_locations_on_parent_id;
DROP INDEX idx_locations_on_ancestors_ids;
DROP INDEX idx_locations_on_level;
