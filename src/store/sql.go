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
	return self.doQuery("l.leaf = TRUE AND ST_Within(ST_SetSRID(ST_Point($1, $2), 4326), l.shape)", opts, x, y)
}

func (self sqlStore) FindLocationsByIds(ids []string, opts model.ReqOptions) ([]model.Location, error) {
	if len(ids) == 0 {
		return make([]model.Location, 0), nil
	}

	placeholders := placeholdersFor(ids, 0)
	args := argsFor(ids)
	return self.doQuery(fmt.Sprintf("l.id IN (%s)", placeholders), opts, args...)
}

func (self sqlStore) FindLocationsByParent(parentId string, opts model.ReqOptions) ([]model.Location, error) {
	opts.Ancestors = false
	return self.doQuery("l.parent_id = $1", opts, parentId)
}

func (self sqlStore) FindLocationsByName(name string, opts model.ReqOptions) ([]model.Location, error) {
	return self.doQuery(`l.name LIKE $1 || '%'`, opts, name)
}

func (self sqlStore) doQuery(predicate string, opts model.ReqOptions, queryArgs ...interface{}) ([]model.Location, error) {
	var query string
	scope, scopeArgs := scopeFor(opts, len(queryArgs))

	if opts.Ancestors {
		query = fmt.Sprintf(`
				SELECT DISTINCT %s
				FROM locations AS l
					INNER JOIN locations as t
					ON t.id = ANY(l.ancestors_ids) OR t.id = l.id
				WHERE %s%s
				ORDER BY t.id%s`, fieldsFor(opts, "t"), predicate, scope, pagingFor(opts))
	} else {
		query = fmt.Sprintf(`
				SELECT %s
				FROM locations AS l
				WHERE %s%s
				ORDER BY l.id%s`, fieldsFor(opts, "l"), predicate, scope, pagingFor(opts))
	}

	args := append(queryArgs, scopeArgs...)

	fmt.Println("Query")
	fmt.Println(query)
	fmt.Println("Args")
	fmt.Println(args)

	rows, err := self.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return readLocations(rows, opts)
}

func fieldsFor(opts model.ReqOptions, alias string) string {
	fields := fmt.Sprintf(`%s.id, %s.parent_id, %s.name, %s.ancestors_ids`, alias, alias, alias, alias)
	if opts.Shapes {
		fields = fmt.Sprintf(`%s, ST_AsBinary(%s.shape) as binshape`, fields, alias)
	}
	return fields
}

func pagingFor(opts model.ReqOptions) string {
	if opts.Limit > 0 {
		return fmt.Sprintf(" LIMIT %d OFFSET %d", opts.Limit, opts.Offset)
	} else {
		return " "
	}
}

func scopeFor(opts model.ReqOptions, argsBase int) (string, []interface{}) {
	if len(opts.Scope) > 0 {
		placeholders := placeholdersFor(opts.Scope, argsBase)
		query := fmt.Sprintf(" AND (l.id IN (%s) OR (l.ancestors_ids && ARRAY[%s::varchar]))", placeholders, placeholders)
		return query, argsFor(opts.Scope)
	} else {
		return " ", make([]interface{}, 0)
	}
}

func readLocations(rows *sql.Rows, opts model.ReqOptions) ([]model.Location, error) {
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

func placeholdersFor(arr []string, base int) string {
	placeholders := ""
	for i, _ := range arr {
		if i == 0 {
			placeholders = fmt.Sprintf("$%d", base+1)
		} else {
			placeholders = fmt.Sprintf("%s,$%d", placeholders, i+base+1)
		}
	}
	return placeholders
}

func argsFor(strArgs []string) []interface{} {
	args := make([]interface{}, len(strArgs))
	for i, str := range strArgs {
		args[i] = str
	}
	return args
}

func (self sqlStore) Begin() Store {
	return self
}

func (self sqlStore) Flush() {
}

func (self sqlStore) Finish() {
	self.db.Exec("UPDATE locations SET leaf = NOT EXISTS (SELECT 1 FROM locations l2 WHERE l2.parent_id = locations.id)")
}
