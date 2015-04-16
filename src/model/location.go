package model

type Location struct {
	Id           string
	ParentId     *string
	AncestorsIds StringSlice
	Ancestors    []*Location
	Level        int
	TypeName     string
	Name         string
	Center       GeoShape
	Shape        GeoShape
}
