/*
This file contains the functions that represent Cache tests for Load()/Purge()/Reload(). These functions are called
inside TestCache()) in do_matchest.go.

Some of the definitions of rulesets below deliberately use a lot of whitespace to keep the code consistent
and to make it easier to understand, add to, and edit these tests
*/

package main

import "log"

func testCache(tests *[]doMatchTest) {

	testLoad(tests)
	//testPurge(tests)
	//testReload(tests)
}

func testLoad(tests *[]doMatchTest) {

	err := Load()
	if err != nil {
		log.Println("ERROR LOAD", err)
	}
}
