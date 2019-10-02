package models

// ServErr - Server error response
type ServErr struct {
	Code     int         `json:"code"`
	Err      string      `json:"err"`
	Desc     string      `json:"desc"`
	Internal interface{} `json:"internal"`
}
