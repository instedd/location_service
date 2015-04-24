# InSTEDD Location Service

InSTEDD Location Service hosts locations information from different sources, and provides a simple HTTP API for querying them based on multiple criteria.

## Development

1. Install [go](https://golang.org/)
2. Run `make` to install dependencies and build

## Database

The service uses a PostgreSQL backend that is configured in `db/dbconf.yml`. Default database name is `location_service`. Migrations are managed via [goose](https://bitbucket.org/liamstask/goose), and are executed by running `bin/goose up`.

## Import

Data can be imported from multiple data sources via the executable `bin/importer`. Run the executable for more info on the accepted options.

The following data sources are currently supported:

* GADM shapefiles
* Natural Earth shapefiles

## API

API server is started via `bin/api`, optionally providing a `-port XXXX` flag to start on a specific port (8080 by default). All endpoints are queried via HTTP GET and provide responses in JSON format.

All endpoints can be invoked with the following optional parameters:

* `ancestors` boolean (false), whether to return all the ancestors of the queried locations; if set, an additional field `ancestors` will be included for each result
* `shapes` boolean (false), whether to include the polygon in GeoJSON format with the shape of each location; if set, an additional field `shape` will be included for each result
* `limit` int (0), how many records to return
* `offset` int (unbounded), offset of the records to be returned, use together with limit for paging.
* `scope` string (none), all results will be limited to the locations specified in this parameter and their descendants; for example, querying for names beginning with `Ca` with `scope=gadm:USA,gadm:MEX` will return _California_ but not _Catamarca_
* `set` string (none), which locations set to query; for example, if a server contains both GADM and NaturalEarth data, the search can be restricted to the latter by specifying `set=ne`.
* `object` string (none), if set, will return a map of locations, indexed by their ids, nested within the key with name `object`, useful for supplying locations with the format expected by [notifiable-diseases](github.com/instedd/notifiable-diseases); if not set, will return an array
* `simplify` float (none), if set, shapes will be dynamically simplified using this value as tolerance (higher values imply a higher simplification); note that this is an expensive operation to perform, and overly complex shapes are already cached in a simplified version, so this parameter is often unneeded

### /lookup

Returns all leaf locations that contain the specified point.

- `x` float
- `y` float

### /details

Returns the details of all the locations requested by id.

- `ids` string, comma separated list of ids

### /children

Returns the direct children of the specified location, or all roots if no parent location is specified.

- `id` string

### /suggest

Returns all locations with a name that matches the supplied prefix.

- `name` string

### /list

Returns all locations in the service

