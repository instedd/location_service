package store

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"model"
)

type sqlStore struct {
	db *sql.DB
}

func NewSqlStore() (Store, error) {
	db, err := sql.Open("postgres", "dbname=location_service sslmode=disable")
	if err != nil {
		return nil, err
	}

	return sqlStore{db: db}, nil
}

func (self sqlStore) AddLocation(location *model.Location) error {
	var ancestors model.StringSlice
	var err error

	if location.ParentId != nil {
		err = self.db.QueryRow("SELECT ancestors_ids FROM locations WHERE id = $1", *location.ParentId).Scan(&ancestors)
		if err != nil {
			return err
		}

		ancestors = append(ancestors, *location.ParentId)
	} else {
		ancestors = make([]string, 0)
	}

	_, err = self.db.Exec(`INSERT INTO locations(id, parent_id, level, type_name, name, shape, ancestors_ids) VALUES ($1, $2, $3, $4, $5, ST_GeomFromWKB($6, 4326), $7)`,
		location.Id, location.ParentId, location.Level, location.TypeName, location.Name, location.Shape, &ancestors)
	return err
}

func (self sqlStore) FindLocationsByPoint(x, y float64, opts model.ReqOptions) ([]model.Location, error) {
	query := QueryFor("l.leaf = TRUE AND ST_Within(ST_SetSRID(ST_Point($1, $2), 4326), l.shape)", opts)
	rows, err := self.db.Query(query, x, y)
	if err != nil {
		return nil, err
	}

	return ReadLocations(rows, opts)
}

func (self sqlStore) FindLocationsByIds(ids []string, opts model.ReqOptions) ([]model.Location, error) {
	if len(ids) == 0 {
		return make([]model.Location, 0), nil
	}

	placeholders := ""
	for i, _ := range ids {
		if i == 0 {
			placeholders = "$1"
		} else {
			placeholders = fmt.Sprintf("%s,$%d", placeholders, i+1)
		}
	}

	query := QueryFor(fmt.Sprintf("l.id IN (%s)", placeholders), opts)

	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	rows, err := self.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return ReadLocations(rows, opts)
}

func (self sqlStore) FindLocationsByParent(parentId string, opts model.ReqOptions) ([]model.Location, error) {
	fields := FieldsFor(opts, "l")

	var query string
	query = fmt.Sprintf(`
		SELECT %s
		FROM locations AS l
		WHERE l.parent_id = $1
		ORDER BY id %s`, fields, PagingFor(opts))

	rows, err := self.db.Query(query, parentId)
	if err != nil {
		return nil, err
	}

	return ReadLocations(rows, opts)
}

func (self sqlStore) FindLocationsByName(name string, opts model.ReqOptions) ([]model.Location, error) {
	query := QueryFor(`l.name LIKE $1 || '%%'`, opts)
	rows, err := self.db.Query(query, name)
	if err != nil {
		return nil, err
	}

	return ReadLocations(rows, opts)
}

func FieldsFor(opts model.ReqOptions, alias string) string {
	fields := fmt.Sprintf(`%s.id, %s.parent_id, %s.name, %s.ancestors_ids`, alias, alias, alias, alias)
	if opts.Shapes {
		fields = fmt.Sprintf(`%s, ST_AsBinary(%s.shape) as binshape`, fields, alias)
	}
	return fields
}

func PagingFor(opts model.ReqOptions) string {
	if opts.Limit > 0 {
		return fmt.Sprintf(" LIMIT %d OFFSET %d", opts.Limit, opts.Offset)
	} else {
		return ""
	}
}

func QueryFor(predicate string, opts model.ReqOptions) string {
	if opts.Ancestors {
		return fmt.Sprintf(`
				SELECT DISTINCT %s
				FROM locations AS l
					INNER JOIN locations as t
					ON t.id = ANY(l.ancestors_ids) OR t.id = l.id
				WHERE %s
				ORDER BY t.id
				%s`, FieldsFor(opts, "t"), predicate, PagingFor(opts))
	} else {
		return fmt.Sprintf(`
				SELECT %s
				FROM locations AS l
				WHERE %s
				ORDER BY l.id
				%s`, FieldsFor(opts, "l"), predicate, PagingFor(opts))
	}
}

func ReadLocations(rows *sql.Rows, opts model.ReqOptions) ([]model.Location, error) {
	locations := make([]model.Location, 0, 20)

	for rows.Next() {
		var location model.Location

		var err error
		if opts.Shapes {
			err = rows.Scan(&location.Id, &location.ParentId, &location.Name, &location.AncestorsIds, &location.Shape)
		} else {
			err = rows.Scan(&location.Id, &location.ParentId, &location.Name, &location.AncestorsIds)
		}

		if err != nil {
			return nil, err
		}

		locations = append(locations, location)
	}

	return locations, nil
}

func (self sqlStore) Begin() Store {
	return self
}

func (self sqlStore) Flush() {
}

func (self sqlStore) Finish() {
	self.db.Exec("UPDATE locations SET leaf = NOT EXISTS (SELECT 1 FROM locations l2 WHERE l2.parent_id = locations.id)")
}
