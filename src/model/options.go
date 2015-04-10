package model

type ReqOptions struct {
	Ancestors bool
	Shapes    bool
	Offset    int
	Limit     int
	Set       string
	Scope     []string
}
