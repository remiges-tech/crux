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

type markdone_t struct {
	Id         int32  `json:"id"`
	Entity     Entity `json:"entity"`
	Step       string `json:"step"`
	Stepfailed bool   `json:"stepfailed"`
	Trace      int
}

type response_t struct {
	Task     map[string]int32 `json:"tasks"`
	Loggedat time.Time        `json:"loggedat"`
}
