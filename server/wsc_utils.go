package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/remiges-tech/crux/types"
)


// CommonValidation is a generic function which setup standard validation utilizing
// validator package and Maps the errorVals based on the map parameter and
// return []errorVals
func CommonValidation(err validator.FieldError) []string {
	var vals []string
	switch err.Tag() {
	case "required":
		vals = append(vals, "not_provided")
	case "alpha":
		vals = append(vals, "only_alphabets_are_allowed")
	case "gt":
		vals = append(vals, "must_be_greater_than_zero")
	default:
		vals = append(vals, "not_valid_input")
	}
	return vals
}

func MarshalJson(data any) []byte {
	jsonData, err := json.Marshal(&data)
	if err != nil {
		log.Fatal("error marshaling")
	}
	return jsonData
}

func ReadJsonFromFile(filepath string) ([]byte, error) {
	// var err error
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("testFile path is not exist")
	}
	defer file.Close()
	jsonData, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

// to check given string is nil or not
func IsStringEmpty(s *string) bool {
	return s == nil || strings.TrimSpace(*s) == ""
}

// to check if the user has "ruleset" rights for the given app
func HasRulesetRights(app string) bool {
	userRights := GetWorkflowsByRulesetRights()
	return slices.Contains(userRights, app)
}

// to get workflows for all apps for which the user has "ruleset" rights
func GetWorkflowsByRulesetRights() []string {
	return []string{"retailBANK", "nedbank"}
}

func Authz_check(op types.OpReq, trace bool) (bool, []string) {
	caplist := op.CapNeeded
	return true, caplist
}
