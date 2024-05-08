package breschema

type patternSchema struct {
	Class string      `json:"class,omitempty"`
	Attr  []attribute `json:"attr,omitempty"`
}
// used in test cases no need for validation
type BRESchemaListStruct struct {
	Slice *int32  `form:"slice"`
	App   *string `form:"app"`
	Class *string `form:"class"`
}

type attribute struct {
	Name      string   `json:"name" validate:"required"`
	ShortDesc string   `json:"shortdesc" validate:"required"`
	LongDesc  string   `json:"longdesc" validate:"required"`
	ValType   string   `json:"valtype" validate:"required"`
	Vals      []string `json:"vals,omitempty"`
	Enumdesc  []string `json:"enumdesc,omitempty"`
	ValMax    *int32   `json:"valmax,omitempty"`
	ValMin    *int32   `json:"valmin,omitempty"`
	LenMax    *int32   `json:"lenmax,omitempty"`
	LenMin    *int32   `json:"lenmin,omitempty"`
}


