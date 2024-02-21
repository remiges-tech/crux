package serverBRE

var ValidTypes = map[string]bool{
	"int": true, "float": true, "str": true, "enum": true, "bool": true, "timestamps": true,
}

const (
	CreatedBy    = "admin"
	EditedBy     = "admin"
	RealmID      = 1
	CruxIDRegExp = `^[a-z][a-z0-9_]*$`
	BRWF         = "B"
)
