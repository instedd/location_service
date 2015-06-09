-- +goose Up
CREATE EXTENSION unaccent;
CREATE TEXT SEARCH CONFIGURATION location_name (COPY = simple);
ALTER TEXT SEARCH CONFIGURATION location_name
  ALTER MAPPING FOR asciihword, asciiword, hword_asciipart, hword, hword_part, word WITH unaccent, simple;
CREATE INDEX idx_locations_on_name ON locations USING gin(to_tsvector('location_name', name));

-- +goose Down
DROP INDEX idx_locations_on_name;
DROP TEXT SEARCH CONFIGURATION location_name;
DROP EXTENSION unaccent;


