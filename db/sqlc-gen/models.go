// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package sqlc

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type BrwfEnum string

const (
	BrwfEnumB BrwfEnum = "B"
	BrwfEnumW BrwfEnum = "W"
)

func (e *BrwfEnum) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = BrwfEnum(s)
	case string:
		*e = BrwfEnum(s)
	default:
		return fmt.Errorf("unsupported scan type for BrwfEnum: %T", src)
	}
	return nil
}

type NullBrwfEnum struct {
	BrwfEnum BrwfEnum `json:"brwf_enum"`
	Valid    bool     `json:"valid"` // Valid is true if BrwfEnum is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullBrwfEnum) Scan(value interface{}) error {
	if value == nil {
		ns.BrwfEnum, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.BrwfEnum.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullBrwfEnum) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.BrwfEnum), nil
}

type App struct {
	ID          int32            `json:"id"`
	Realm       string           `json:"realm"`
	Shortname   string           `json:"shortname"`
	Shortnamelc string           `json:"shortnamelc"`
	Longname    string           `json:"longname"`
	Setby       string           `json:"setby"`
	Setat       pgtype.Timestamp `json:"setat"`
}

type Capgrant struct {
	ID    int32            `json:"id"`
	Realm string           `json:"realm"`
	User  string           `json:"user"`
	App   pgtype.Text      `json:"app"`
	Cap   string           `json:"cap"`
	From  pgtype.Timestamp `json:"from"`
	To    pgtype.Timestamp `json:"to"`
	Setat pgtype.Timestamp `json:"setat"`
	Setby string           `json:"setby"`
}

type Config struct {
	Realm string           `json:"realm"`
	Slice int32            `json:"slice"`
	Name  string           `json:"name"`
	Descr string           `json:"descr"`
	Val   pgtype.Text      `json:"val"`
	Ver   pgtype.Int4      `json:"ver"`
	Setby string           `json:"setby"`
	Setat pgtype.Timestamp `json:"setat"`
}

type Deactivated struct {
	ID      int32            `json:"id"`
	Realm   string           `json:"realm"`
	User    pgtype.Text      `json:"user"`
	Deactby string           `json:"deactby"`
	Deactat pgtype.Timestamp `json:"deactat"`
}

type Realm struct {
	ID          int32            `json:"id"`
	Shortname   string           `json:"shortname"`
	Shortnamelc string           `json:"shortnamelc"`
	Longname    string           `json:"longname"`
	Setby       string           `json:"setby"`
	Setat       pgtype.Timestamp `json:"setat"`
	Payload     []byte           `json:"payload"`
}

type Realmslice struct {
	ID           int32            `json:"id"`
	Realm        string           `json:"realm"`
	Descr        string           `json:"descr"`
	Active       bool             `json:"active"`
	Activateat   pgtype.Timestamp `json:"activateat"`
	Deactivateat pgtype.Timestamp `json:"deactivateat"`
	Createdat    pgtype.Timestamp `json:"createdat"`
	Createdby    string           `json:"createdby"`
	Editedat     pgtype.Timestamp `json:"editedat"`
	Editedby     pgtype.Text      `json:"editedby"`
}

type Ruleset struct {
	ID         int32            `json:"id"`
	Realm      string           `json:"realm"`
	Slice      int32            `json:"slice"`
	App        string           `json:"app"`
	Brwf       BrwfEnum         `json:"brwf"`
	Class      string           `json:"class"`
	Setname    string           `json:"setname"`
	Schemaid   int32            `json:"schemaid"`
	IsActive   pgtype.Bool      `json:"is_active"`
	IsInternal bool             `json:"is_internal"`
	Ruleset    []byte           `json:"ruleset"`
	Createdat  pgtype.Timestamp `json:"createdat"`
	Createdby  string           `json:"createdby"`
	Editedat   pgtype.Timestamp `json:"editedat"`
	Editedby   pgtype.Text      `json:"editedby"`
}

type Schema struct {
	ID            int32            `json:"id"`
	Realm         string           `json:"realm"`
	Slice         int32            `json:"slice"`
	App           string           `json:"app"`
	Brwf          BrwfEnum         `json:"brwf"`
	Class         string           `json:"class"`
	Patternschema []byte           `json:"patternschema"`
	Actionschema  []byte           `json:"actionschema"`
	Createdat     pgtype.Timestamp `json:"createdat"`
	Createdby     string           `json:"createdby"`
	Editedat      pgtype.Timestamp `json:"editedat"`
	Editedby      pgtype.Text      `json:"editedby"`
}

type Stepworkflow struct {
	Slice    int32       `json:"slice"`
	App      pgtype.Text `json:"app"`
	Step     string      `json:"step"`
	Workflow string      `json:"workflow"`
}

type Wfinstance struct {
	ID       int32            `json:"id"`
	Entityid string           `json:"entityid"`
	Slice    int32            `json:"slice"`
	App      string           `json:"app"`
	Class    string           `json:"class"`
	Workflow string           `json:"workflow"`
	Step     string           `json:"step"`
	Loggedat pgtype.Timestamp `json:"loggedat"`
	Doneat   pgtype.Timestamp `json:"doneat"`
	Nextstep string           `json:"nextstep"`
	Parent   pgtype.Int4      `json:"parent"`
}
