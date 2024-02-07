package workflow

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/crux/server/schema"
)

func WorkFlowNew(c *gin.Context, s *service.Service) {
	l := s.LogHarbour
	l.Log("Starting execution of WorkFlowNew()")

	var wf workflowNew
	var ruleSchema schema.Schema

	err := wscutils.BindJSON(c, &wf)
	if err != nil {
		l.LogActivity("Error Unmarshalling Query parameters to struct:", err.Error())
		return
	}

	query, ok := s.Dependencies["queries"].(*sqlc.Queries)
	if !ok {
		l.Log("Error while getting query instance from service Dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(wscutils.ErrcodeDatabaseError))
		return
	}

	schema, err := query.WfSchemaGet(c, sqlc.WfSchemaGetParams{
		Slice: wf.Slice,
		App:   wf.App,
		Class: wf.Class,
	})
	if err != nil {
		l.LogActivity("failed to get data from DB:", err.Error())
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(wscutils.ErrcodeDatabaseError))
		return
	}
	ruleSchema.Slice = schema.Slice
	ruleSchema.App = schema.App
	ruleSchema.Class = schema.Class
	err = json.Unmarshal([]byte(schema.Patternschema), &ruleSchema.PatternSchema)
	if err != nil {
		l.LogActivity("Error while Unmarshalling PatternSchema", err)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(wscutils.ErrcodeDatabaseError))
		return
	}
	err = json.Unmarshal(schema.Actionschema, &ruleSchema.ActionSchema)
	if err != nil {
		l.LogActivity("Error while Unmarshaling ActionSchema", err)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(wscutils.ErrcodeDatabaseError))
		return
	}

	// Validate request
	validationErrors := wscutils.WscValidate(wf, func(err validator.FieldError) []string { return []string{} })
	customValidationErrors := customValidationErrors(wf, ruleSchema)
	validationErrors = append(validationErrors, customValidationErrors...)
	if len(validationErrors) > 0 {
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, validationErrors))
		return
	}

	ruleset, err := json.Marshal(wf.Flowrules)
	if err != nil {
		patternSchema := "flowrules"
		l.LogDebug("Error while marshaling Flowrules", err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(wscutils.ErrcodeInvalidJson, &patternSchema)}))
		return
	}

	_, err = query.WorkFlowNew(c, sqlc.WorkFlowNewParams{
		Realm:      realmID,
		Slice:      wf.Slice,
		App:        wf.App,
		Brwf:       brwf,
		Class:      wf.Class,
		Setname:    setBy,
		Schemaid:   schema.ID,
		IsActive:   pgtype.Bool{Bool: isActive},
		IsInternal: wf.IsInternal,
		Ruleset:    ruleset,
		Createdby:  setBy,
	})
	if err != nil {
		l.LogActivity("Error while creating schema", err.Error())
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(wscutils.ErrcodeDatabaseError))
		return
	}
	wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: "Created successfully", Messages: nil})
	l.Log("Finished execution of SchemaNew()")
}

func customValidationErrors(wf workflowNew, ruleSchema schema.Schema) []wscutils.ErrorMessage {
	var validationErrors []wscutils.ErrorMessage
	if len(wf.Flowrules) == 0 {
		fieldName := "flowrules"
		vErr := wscutils.BuildErrorMessage("empty", &fieldName)
		validationErrors = append(validationErrors, vErr)
	}
	rulePatternError := verifyRulePatterns(wf, ruleSchema)
	validationErrors = append(validationErrors, rulePatternError...)

	ruleActionError := verifyRuleActions(wf, ruleSchema)
	validationErrors = append(validationErrors, ruleActionError...)
	return validationErrors
}

func verifyRulePatterns(ruleSet workflowNew, ruleSchema schema.Schema) []wscutils.ErrorMessage {
	var validationErrors []wscutils.ErrorMessage

	for _, rule := range ruleSet.Flowrules {
		for i, term := range rule.RulePattern {
			i++
			if !(rule.RulePattern[0].AttrName == step) {
				fieldName := fmt.Sprintf("flowrules.rulepattern[%d].AttrName", i)
				vErr := wscutils.BuildErrorMessage("step_not_exist", &fieldName)
				validationErrors = append(validationErrors, vErr)
			}
			valType := getType(ruleSchema, term.AttrName)
			if valType == "" {
				// If the attribute name is not in the pattern-schema, we check if it's a task "tag"
				// by checking for its presence in the action-schema
				if !isStringInArray(term.AttrName, ruleSchema.ActionSchema.Tasks) {
					fieldName := fmt.Sprintf("RulePattern[%d].AttrName", i)
					vErr := wscutils.BuildErrorMessage("not_exist", &fieldName, term.AttrName)
					validationErrors = append(validationErrors, vErr)
				}
				// If it is a tag, the value type is set to bool
				valType = typeBool
			}
			if !verifyType(term.AttrVal, valType) {
				fieldName := fmt.Sprintf("RulePattern[%d].AttrVal", i)
				// term.AttrVal
				vErr := wscutils.BuildErrorMessage("not_support", &fieldName)
				validationErrors = append(validationErrors, vErr)
			}
			if !validOps[term.Op] {
				fieldName := fmt.Sprintf("RulePattern[%d].term.Op", i)
				vErr := wscutils.BuildErrorMessage("not_support", &fieldName, term.Op)
				validationErrors = append(validationErrors, vErr)
			}
		}
	}
	return validationErrors
}

func getType(ruleSchema schema.Schema, name string) string {
	for _, as := range ruleSchema.PatternSchema.Attr {
		if as.Name == name {
			return as.ValType
		}
	}
	return ""
}

func isStringInArray(s string, arr []string) bool {
	for _, a := range arr {
		if a == s {
			return true
		}
	}
	return false
}

// Returns whether or not the type of "val" is the same as "valType"
func verifyType(val any, valType string) bool {
	var ok bool
	switch valType {
	case typeBool:
		_, ok = val.(bool)
	case typeInt:
		_, ok = val.(int)
	case typeFloat:
		_, ok = val.(float64)
	case typeStr, typeEnum:
		_, ok = val.(string)
	case typeTS:
		s, _ := val.(string)
		_, err := time.Parse(timeLayout, s)
		ok = (err == nil)
	}
	return ok
}

func verifyRuleActions(ruleSet workflowNew, ruleSchema schema.Schema) []wscutils.ErrorMessage {
	var validationErrors []wscutils.ErrorMessage
	for _, rule := range ruleSet.Flowrules {

		for i, t := range rule.RuleActions.Tasks {
			i++
			if !isStringInArray(t, ruleSchema.ActionSchema.Tasks) {
				fieldName := fmt.Sprintf("Tasks[%d]", i)
				vErr := wscutils.BuildErrorMessage("not_exist", &fieldName, t)
				validationErrors = append(validationErrors, vErr)
			}
		}

		for i, p := range rule.RuleActions.Properties {
			i++
			if !isStringInArray(p.Name, ruleSchema.ActionSchema.Properties) {
				fieldName := fmt.Sprintf("Properties[%d]", i)
				vErr := wscutils.BuildErrorMessage("not_exist", &fieldName, p.Name)
				validationErrors = append(validationErrors, vErr)
			}
		}
		// if rule.RuleActions.WillReturn && rule.RuleActions.WillExit {
		// 	// return false, fmt.Errorf("there is a rule with both the RETURN and EXIT instructions in ruleset %v", ruleSet.setName)
		// }

		nsFound, doneFound := areNextStepAndDoneInProps(rule.RuleActions.Properties)
		if !nsFound {
			fieldName := "Properties"
			vErr := wscutils.BuildErrorMessage("nextstep_not_exist", &fieldName)
			validationErrors = append(validationErrors, vErr)
		}
		if !doneFound {
			fieldName := "Properties"
			vErr := wscutils.BuildErrorMessage("done_not_exist", &fieldName)
			validationErrors = append(validationErrors, vErr)
		}
		if !doneFound && len(rule.RuleActions.Tasks) == 0 {
			// return false, fmt.Errorf("no tasks and no 'done=true' in a rule in ruleset %v", ruleSet.setName)
			fieldName := "Properties"
			vErr := wscutils.BuildErrorMessage("empty_task", &fieldName)
			validationErrors = append(validationErrors, vErr)
		}
		currNS := getNextStep(rule.RuleActions.Properties)
		if len(currNS) > 0 && !isStringInArray(currNS, rule.RuleActions.Tasks) {
			// return false, fmt.Errorf("`nextstep` value not found in `tasks` in a rule in ruleset %v", ruleSet.setName)
			fieldName := "Properties"
			vErr := wscutils.BuildErrorMessage("task_not_exist", &fieldName)
			validationErrors = append(validationErrors, vErr)
		}

	}
	return validationErrors
}

func areNextStepAndDoneInProps(props []Property) (bool, bool) {
	nsFound, doneFound := false, false
	for _, p := range props {
		if p.Name == nextStep {
			nsFound = true
		}
		if p.Name == done && p.Val == trueStr {
			doneFound = true
		}
	}
	return nsFound, doneFound
}

func getNextStep(props []Property) string {
	for _, p := range props {
		if p.Name == nextStep {
			return p.Val
		}
	}
	return ""
}
