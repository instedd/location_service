-- +goose Up
ALTER TABLE locations ADD COLUMN simple_shape GEOMETRY;

-- +goose Down
ALTER TABLE locations DROP COLUMN simple_shape;
