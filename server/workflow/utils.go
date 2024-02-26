package workflow

// this is for test cases where 'HasRootCapabilities()' = value of 'TRIGGER'
var TRIGGER bool = false

// to check if the caller has root capabilities
func HasRootCapabilities() bool {
	return TRIGGER
}
