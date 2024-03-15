package schema

import "regexp"

const (
	typeBool  = "bool"
	typeInt   = "int"
	typeFloat = "float"
	typeStr   = "str"
	typeEnum  = "enum"
	typeTS    = "ts"

	step         = "step"
	stepFailed   = "stepfailed"
	start        = "START"
	nextStep     = "nextstep"
	done         = "done"
	cruxIDRegExp = `^[a-z][a-z0-9_]*$`
)

var (
	// userID     = "1234"
	capForNew  = []string{"schema"}
	// realmName  = int32(1)
	validTypes = map[string]bool{
		"int": true, "float": true, "str": true, "enum": true, "bool": true, "timestamps": true,
	}
	capForUpdate = []string{"schema"}
	re           = regexp.MustCompile(cruxIDRegExp)
)
