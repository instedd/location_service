-- +goose Up
CREATE TABLE locations (
  id SERIAL,
  parent_id INT,
  level INT,
  type_name VARCHAR(255),
  name VARCHAR(255),
  shape GEOMETRY
);


-- +goose Down
DROP TABLE locations;
