package store

import (
	"model"
)

type Store interface {
	Begin() Store
	AddLocation(*model.Location) error
	Flush()
}
