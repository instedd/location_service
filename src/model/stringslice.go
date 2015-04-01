package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"strings"
)

type StringSlice []string

func (s *StringSlice) Scan(src interface{}) error {
	asBytes, ok := src.([]byte)
	if !ok {
		return error(errors.New("Scan source was not []bytes"))
	}

	asString := strings.TrimRight(strings.TrimLeft(string(asBytes), "{"), "}")
	if len(asString) == 0 {
		(*s) = StringSlice(make([]string, 0))
		return nil
	}

	values := strings.Split(asString, ",")
	unwrapped := make([]string, len(values))
	for i, str := range values {
		unwrapped[i] = strings.Trim(str, "\"")
	}

	(*s) = StringSlice(unwrapped)
	return nil
}

func (p *StringSlice) Value() (driver.Value, error) {
	array := *p
	wrapped := make([]string, len(array))
	for i, str := range array {
		wrapped[i] = fmt.Sprintf("\"%s\"", strings.Replace(str, "\"", "\\\"", -1))
	}
	text := fmt.Sprintf("{%s}", strings.Join(wrapped, ", "))
	return text, nil
}
