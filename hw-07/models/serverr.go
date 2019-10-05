package models

type ServErr struct {
	Code     int
	Err      string
	Desc     string
	Internal interface{}
}
