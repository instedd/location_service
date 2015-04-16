package store

import (
	"model"
)

type Store interface {
	Begin() Store
	AddLocation(*model.Location) error
	FindLocationsByPoint(x, y float64, opts model.ReqOptions) ([]*model.Location, error)
	FindLocationsByIds(ids []string, opts model.ReqOptions) ([]*model.Location, error)
	FindLocationsByParent(id string, opts model.ReqOptions) ([]*model.Location, error)
	FindLocationsByName(name string, opts model.ReqOptions) ([]*model.Location, error)
	Flush()
	Finish() error
}
