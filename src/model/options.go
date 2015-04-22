package model

type ReqOptions struct {
	Ancestors bool
	Shapes    bool
	Simplify  float32
	Offset    int
	Limit     int
	Set       string
	Object    string
	Scope     []string
}
