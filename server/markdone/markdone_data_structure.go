package markdone

import (
	"time"

	crux "github.com/remiges-tech/crux/matching-engine"
)

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
	InstanceID int32       `json:"id" validate:"required"`
	Workflow   string      `json:"workflow"`
	EntityID   string      `json:"entityID"`
	Entity     crux.Entity `json:"entity" validate:"required"`
	Trace      int         `json:"trace,omitempty"`
	Loggedat   time.Time   `json:"loggedat,omitempty"`
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
