package model

type Location struct {
	Id           string
	ParentId     *string
	AncestorsIds StringSlice
	Level        int
	TypeName     string
	Name         string
	Shape        GeoShape
}
