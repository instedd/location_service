package store

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"model"
	"strings"
)

var debug bool

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

func SetDebug(val bool) {
	debug = val
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
		location.Id, location.ParentId, location.Level, location.TypeName, location.Name, &location.Shape, &ancestors)
	return err
}

func (self sqlStore) FindLocationsByPoint(x, y float64, opts model.ReqOptions) ([]*model.Location, error) {
	return self.doQuery("l.leaf = TRUE AND ST_Within(ST_SetSRID(ST_Point($1, $2), 4326), l.shape)", opts, x, y)
}

func (self sqlStore) FindLocationsByIds(ids []string, opts model.ReqOptions) ([]*model.Location, error) {
	if len(ids) == 0 {
		return make([]*model.Location, 0), nil
	}

	placeholders := placeholdersFor(ids, 0)
	args := argsFor(ids)
	return self.doQuery(fmt.Sprintf("l.id IN (%s)", placeholders), opts, args...)
}

func (self sqlStore) FindLocationsByParent(parentId string, opts model.ReqOptions) ([]*model.Location, error) {
	opts.Ancestors = false
	if len(parentId) == 0 {
		return self.doQuery("l.parent_id IS NULL", opts)
	} else {
		return self.doQuery("l.parent_id = $1", opts, parentId)
	}

}

func (self sqlStore) FindLocationsByName(name string, opts model.ReqOptions) ([]*model.Location, error) {
	return self.doQuery(`l.name LIKE ($1 || '%')`, opts, name)
}

func (self sqlStore) doQuery(predicate string, opts model.ReqOptions, queryArgs ...interface{}) ([]*model.Location, error) {
	var query string
	setPredicate, setArgs := setFor(opts, len(queryArgs))
	scope, scopeArgs := scopeFor(opts, len(queryArgs)+len(setArgs))

	query = fmt.Sprintf(`
		SELECT %s
		FROM locations AS l
		WHERE %s%s%s
		ORDER BY l.id%s`, fieldsFor(opts, "l"), predicate, setPredicate, scope, pagingFor(opts))

	args := append(queryArgs, setArgs...)
	args = append(args, scopeArgs...)

	if debug {
		fmt.Printf("\nExecuting query:%s\nWith params: %s\n", strings.Replace(query, "				", " ", -1), args)
	}

	rows, err := self.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	locations, err := readLocations(rows, opts)
	if err != nil {
		return nil, err
	}

	if opts.Ancestors {
		return self.addAncestors(locations, opts)
	} else {
		return locations, nil
	}

}

func (self sqlStore) addAncestors(locations []*model.Location, opts model.ReqOptions) ([]*model.Location, error) {
	ancestors := make(map[string](*model.Location))
	for _, location := range locations {
		for _, ancestorId := range location.AncestorsIds {
			ancestors[ancestorId] = nil
		}
	}

	ancestorIds := make([]string, 0, len(ancestors))
	for ancestorId := range ancestors {
		ancestorIds = append(ancestorIds, ancestorId)
	}

	placeholders := placeholdersFor(ancestorIds, 0)
	args := argsFor(ancestorIds)

	ancestorOpts := opts
	ancestorOpts.Ancestors = false
	query := fmt.Sprintf(`
		SELECT %s
		FROM locations AS l
		WHERE l.id IN (%s)
		ORDER BY l.id`, fieldsFor(ancestorOpts, "l"), placeholders)

	rows, err := self.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	ancestorsList, err := readLocations(rows, ancestorOpts)
	if err != nil {
		return nil, err
	}

	for _, ancestor := range ancestorsList {
		ancestors[(*ancestor).Id] = ancestor
	}

	for _, location := range locations {
		location.Ancestors = make([]*model.Location, 0, len(location.AncestorsIds))
		for _, ancestorId := range (*location).AncestorsIds {
			location.Ancestors = append(location.Ancestors, ancestors[ancestorId])
		}
	}

	return locations, nil
}

func setFor(opts model.ReqOptions, argsBase int) (string, []interface{}) {
	if len(opts.Set) > 0 {
		return fmt.Sprintf(" AND l.id LIKE ($%d || ':%%')", argsBase+1), []interface{}{opts.Set}
	} else {
		return " ", make([]interface{}, 0)
	}
}

func fieldsFor(opts model.ReqOptions, alias string) string {
	fields := fmt.Sprintf(`%s.id, %s.parent_id, %s.name, %s.type_name, %s.level, ST_AsBinary(%s.center::geometry) as bincenter, %s.ancestors_ids`,
		alias, alias, alias, alias, alias, alias, alias)
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

func readLocations(rows *sql.Rows, opts model.ReqOptions) ([]*model.Location, error) {
	locations := make([]*model.Location, 0, 20)

	for rows.Next() {
		var location model.Location
		var err error

		if opts.Shapes {
			err = rows.Scan(&location.Id, &location.ParentId, &location.Name, &location.TypeName, &location.Level, &location.Center, &location.AncestorsIds, &location.Shape)
		} else {
			err = rows.Scan(&location.Id, &location.ParentId, &location.Name, &location.TypeName, &location.Level, &location.Center, &location.AncestorsIds)
		}

		if err != nil {
			return nil, err
		}

		locations = append(locations, &location)
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

func (self sqlStore) Finish() error {
	var err error
	_, err = self.db.Exec("UPDATE locations SET leaf = NOT EXISTS (SELECT 1 FROM locations l2 WHERE l2.parent_id = locations.id)")
	if err != nil {
		return err
	}

	_, err = self.db.Exec("UPDATE locations SET center = ST_PointOnSurface(shape)::point WHERE center IS NULL AND shape IS NOT NULL")
	return err
}
