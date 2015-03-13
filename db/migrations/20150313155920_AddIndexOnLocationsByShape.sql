-- +goose Up
CREATE INDEX idx_locations_on_shape ON locations USING GIST(shape);

-- +goose Down
DROP INDEX idx_locations_on_shape;
