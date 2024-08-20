/*
This file contains the data structures used by the matching engine
*/

package crux

type Entity struct {
	Realm string            `json:"realm"`
	App   string            `json:"app" validator:"required,alpha,lt=15"`
	Slice int32             `json:"slice" validator:"required"`
	Class string            `json:"class" validator:"required,alpha,lt=15"`
	Attrs map[string]string `json:"attrs" validator:"required,alpha,lt=30"`
}
type ActionSet struct {
	Tasks      []string
	Properties map[string]string
}
