/*
This file contains the data structures used by the matching engine
*/

package main

type Entity struct {
	realm string
	app   string
	slice string
	class string
	attrs map[string]string
}
type ActionSet struct {
	tasks      []string
	properties map[string]string
}
