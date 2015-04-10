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

* **ancestors**: boolean (false), whether to return all the ancestors of the queried locations; for example, a query for California will also return USA if this flag is set.
* **shapes**: boolean (false), whether to include the polygon in GeoJSON format with the shape of each location.
* **limit**: int (0), how many records to return.
* **offset**: int (unbounded), offset of the records to be returned, use together with limit for paging.
* **scope**: string (none), all results will be limited to the locations specified in this parameter and their descendants; for example, querying for names beginning with `Ca` with `scope=USA` will return _California_ but not _Catamarca_. This parameter can be specified multiple times, and the union of all scopes will be considered.
* **set**: string (none), which locations set to query; for example, if a server contains both GADM and NaturalEarth data, the search can be restricted to the latter by specifying `set=ne`.

### /lookup

Returns all leaf locations that contain the specified point.
Requires `x` and `y` as float parameters.

### /details

Returns the details of all the locations requested by id.
Requires the parameter `id`, one or multiple times.

### /children

Returns the direct children of the specified location.
Requires the parameter `id`, once.

### /suggest

Returns all locations with a name that matches the supplied prefix.
Requires the parameter `name`, once.
