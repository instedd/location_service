package model

import (
	"github.com/foobaz/geom"
)

type Location struct {
	Id       int
	ParentId int
	Level    int
	TypeName string
	Name     string
	Shape    geom.T
}
