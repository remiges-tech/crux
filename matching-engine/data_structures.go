/*
This file contains the data structures used by the matching engine
*/

package crux

type Entity struct {
	Realm string
	App   string
	Slice string
	Class string
	Attrs map[string]string
}
type ActionSet struct {
	Tasks      []string
	Properties map[string]string
}
