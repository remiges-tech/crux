package crux

import "time"

/*
{
    "id": 9854,
    "entity": {
        "class": "salesvoucher",
        "voucherdate": "2024-05-20T00:00:00Z",
        "branch": "belampally",
        "voucheramt": "84213"
    },
    "step": "bankchk",
    "stepfailed": false,
    "trace": 0
}

*/

type Markdone_t struct {
	Id         int32  `json:"id"`
	Entity     Entity `json:"entity"`
	Step       string `json:"step"`
	Stepfailed bool   `json:"stepfailed"`
	Trace      int
}

type ResponseData struct {
	Id        int32              `json:"id"`
	Done      bool               `json:"done,omitempty"`
	DoneAt    time.Time          `json:"doneat,omitempty"`
	Step      string             `json:"step,omitempty"`
	Nextstep  string             `json:"nextstep,omitempty"`
	Tasks     []map[string]int32 `json:"tasks,omitempty"`
	Loggedat  time.Time          `json:"loggedat,omitempty"`
	Subflows  map[string]string  `json:"subflows,omitempty"`
	Tracedata string             `json:"tracedata,omitempty"`
}
