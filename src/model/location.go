package model

import (
	"github.com/foobaz/geom"
)

type Location struct {
	Id           string
	ParentId     *string
	AncestorsIds []string
	Level        int
	TypeName     string
	Name         string
	Shape        geom.T
}
