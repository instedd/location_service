package store

import (
	"model"
)

type Store interface {
	Begin() Store
	AddLocation(*model.Location) error
	FindLocationsByPoint(x, y float64, includeShape bool) ([]model.Location, error)
	Flush()
	Finish()
}
