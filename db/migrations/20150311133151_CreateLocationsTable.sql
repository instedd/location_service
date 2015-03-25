-- +goose Up
CREATE TABLE locations (
  id VARCHAR(255),
  parent_id VARCHAR(255),
  ancestors_ids VARCHAR(255)[],
  level INT,
  type_name VARCHAR(255),
  name VARCHAR(255),
  leaf BOOLEAN,
  shape GEOMETRY
);

-- +goose Down
DROP TABLE locations;
