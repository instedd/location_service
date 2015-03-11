-- +goose Up
CREATE EXTENSION postgis;

-- +goose Down
DROP EXTENSION postgis;
