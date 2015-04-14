package model

type Location struct {
	Id           string
	ParentId     *string
	AncestorsIds StringSlice
	Ancestors    []*Location
	Level        int
	TypeName     string
	Name         string
	Shape        GeoShape
}
