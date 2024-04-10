package breschema

const (
	CruxIDRegExp = `^[a-z][a-z0-9_]*$`
	BRWF         = "B"
)

var (
	ValidTypes = map[string]bool{
		"int": true, "float": true, "str": true, "enum": true, "bool": true, "timestamps": true,
	}
)
