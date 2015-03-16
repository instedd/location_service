-- +goose Up
CREATE INDEX idx_locations_on_shape ON locations USING GIST(shape) WHERE leaf = true;

-- +goose Down
DROP INDEX idx_locations_on_shape;
