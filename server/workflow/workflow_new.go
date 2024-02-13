package workflow

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/crux/server"
	"github.com/remiges-tech/crux/server/schema"
)

type WorkflowNew struct {
	Slice      int32   `json:"slice" validate:"required,gt=0"`
	App        string  `json:"app" validate:"required,alpha"`
	Class      string  `json:"class" validate:"required,lowercase"`
	Name       string  `json:"name" validate:"required,lowercase"`
	IsInternal bool    `json:"is_internal" validate:"required"`
	Flowrules  []Rules `json:"flowrules" validate:"required,dive"`
}

func WorkFlowNew(c *gin.Context, s *service.Service) {
	l := s.LogHarbour
	l.Log("Starting execution of WorkFlowNew()")

	var wf WorkflowNew
	var ruleSchema schema.Schema

	err := wscutils.BindJSON(c, &wf)
	if err != nil {
		l.LogActivity("Error Unmarshalling Query parameters to struct:", err.Error())
		return
	}

	validationErrors := wscutils.WscValidate(wf, func(err validator.FieldError) []string { return []string{} })
	if len(validationErrors) > 0 {
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, validationErrors))
		return
	}

	query, ok := s.Dependencies["queries"].(*sqlc.Queries)
	if !ok {
		l.Log("Error while getting query instance from service Dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}

	connpool, ok := s.Database.(*pgxpool.Pool)
	if !ok {
		l.Log("Error while getting connection pool instance from service Dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}
	tx, err := connpool.Begin(c)
	if err != nil {
		l.LogActivity("Error while Begin tx", err.Error())
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}
	defer tx.Rollback(c)
	qtx := query.WithTx(tx)

	schema, err := qtx.GetSchemaWithLock(c, sqlc.GetSchemaWithLockParams{
		Slice: wf.Slice,
		App:   wf.App,
		Class: wf.Class,
	})
	if err != nil {
		l.LogActivity("failed to get schema from DB:", err.Error())
		errmsg := db.ErrorMessage(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}
	ruleSchema.Slice = wf.Slice
	ruleSchema.App = wf.App
	ruleSchema.Class = wf.Class
	err = json.Unmarshal([]byte(schema.Patternschema), &ruleSchema.PatternSchema)
	if err != nil {
		l.LogActivity("Error while Unmarshalling PatternSchema", err)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}
	err = json.Unmarshal(schema.Actionschema, &ruleSchema.ActionSchema)
	if err != nil {
		l.LogActivity("Error while Unmarshaling ActionSchema", err)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}

	// custom Validation
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
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(server.MsgId_Invalid_Request, server.ErrCode_InvalidRequest, &patternSchema)}))
		return
	}

	err = qtx.WorkFlowNew(c, sqlc.WorkFlowNewParams{
		Realm:      realmID,
		Slice:      wf.Slice,
		App:        wf.App,
		Brwf:       brwf,
		Class:      wf.Class,
		Setname:    setBy,
		Schemaid:   schema.ID,
		IsActive:   pgtype.Bool{Bool: false, Valid: true},
		IsInternal: wf.IsInternal,
		Ruleset:    ruleset,
		Createdby:  setBy,
	})
	if err != nil {
		l.LogActivity("Error while Inserting data in ruleset", err.Error())
		errmsg := db.ErrorMessage(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}
	if err := tx.Commit(c); err != nil {
		l.LogActivity("Error while commits the transaction", err.Error())
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}
	wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: nil, Messages: nil})
	l.Log("Finished execution of WorkFlowNew()")
}

func customValidationErrors(wf WorkflowNew, ruleSchema schema.Schema) []wscutils.ErrorMessage {
	var validationErrors []wscutils.ErrorMessage
	if len(wf.Flowrules) == 0 {
		fieldName := "flowrules"
		vErr := wscutils.BuildErrorMessage(server.MsgId_Empty, server.ErrCode_Empty, &fieldName)
		validationErrors = append(validationErrors, vErr)
	}
	rulePatternError := verifyRulePatterns(wf, ruleSchema)
	validationErrors = append(validationErrors, rulePatternError...)

	ruleActionError := verifyRuleActions(wf, ruleSchema)
	validationErrors = append(validationErrors, ruleActionError...)
	return validationErrors
}

func verifyRulePatterns(ruleSet WorkflowNew, ruleSchema schema.Schema) []wscutils.ErrorMessage {
	var validationErrors []wscutils.ErrorMessage

	for i, rule := range ruleSet.Flowrules {
		i++
		for j, term := range rule.RulePattern {
			j++
			valType := getType(ruleSchema, term.Attr)
			if valType == "" {
				// If the attribute name is not in the pattern-schema, we check if it's a task "tag"
				// by checking for its presence in the action-schema
				if !isStringInArray(term.Attr, ruleSchema.ActionSchema.Tasks) {
					fieldName := fmt.Sprintf("flowrules[%d].rulepattern[%d].attr", i, j)
					vErr := wscutils.BuildErrorMessage(server.MsgId_Invalid, server.ErrCode_Invalid, &fieldName, term.Attr)
					validationErrors = append(validationErrors, vErr)
				} else {
					// If it is a tag, the value type is set to bool
					term.Val = typeBool
				}
			}
			if len(valType) != 0 {
				if !verifyType(term.Val, valType) {
					fieldName := fmt.Sprintf("flowrules[%d].rulepattern[%d].val", i, j)
					// term.Val
					vErr := wscutils.BuildErrorMessage(server.MsgId_Invalid, server.ErrCode_Invalid, &fieldName)
					validationErrors = append(validationErrors, vErr)
				}
			}
			if !validOps[term.Op] {
				fieldName := fmt.Sprintf("flowrules[%d].rulepattern[%d].term.Op", i, j)
				vErr := wscutils.BuildErrorMessage(server.MsgId_Invalid, server.ErrCode_Invalid, &fieldName, term.Op)
				validationErrors = append(validationErrors, vErr)
			}
		}

		stepFound := false
		for _, term := range rule.RulePattern {
			if term.Attr == step {
				stepFound = true
				break
			}
		}
		if !stepFound {
			fieldName := fmt.Sprintf("flowrules[%d].rulepattern", i)
			vErr := wscutils.BuildErrorMessage(server.MsgId_StepNotFound, server.ErrCode_Invalid, &fieldName)
			validationErrors = append(validationErrors, vErr)
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

func verifyRuleActions(ruleSet WorkflowNew, ruleSchema schema.Schema) []wscutils.ErrorMessage {
	var validationErrors []wscutils.ErrorMessage
	for i, rule := range ruleSet.Flowrules {
		i++
		for j, t := range rule.RuleActions.Tasks {
			j++
			if !isStringInArray(t, ruleSchema.ActionSchema.Tasks) {
				fieldName := fmt.Sprintf("flowrules[%d].tasks[%d]", i, j)
				vErr := wscutils.BuildErrorMessage(server.MsgId_NotFound, server.ErrCode_NotFound, &fieldName, t)
				validationErrors = append(validationErrors, vErr)
			}
		}

		for j, p := range rule.RuleActions.Properties {
			j++
			if !isStringInArray(p.Name, ruleSchema.ActionSchema.Properties) {
				fieldName := fmt.Sprintf("flowrules[%d].properties[%d]", i, j)
				vErr := wscutils.BuildErrorMessage(server.MsgId_NotFound, server.ErrCode_NotFound, &fieldName, p.Name)
				validationErrors = append(validationErrors, vErr)
			}
		}

		if rule.RuleActions.WillReturn && rule.RuleActions.WillExit {
			fieldName := fmt.Sprintf("flowrules[%d].ruleactions(WillReturn/WillExit)", i)
			vErr := wscutils.BuildErrorMessage(server.MsgId_RequiredOneOf, server.ErrCode_RequiredOne, &fieldName)
			validationErrors = append(validationErrors, vErr)
		}

		nsFound, doneFound := areNextStepAndDoneInProps(rule.RuleActions.Properties)
		if !nsFound && !doneFound {
			fieldName := "properties(nextstep/done)"
			vErr := wscutils.BuildErrorMessage(server.MsgId_RequiredOneOf, server.ErrCode_RequiredOne, &fieldName)
			validationErrors = append(validationErrors, vErr)
		}

		if doneFound && !(len(rule.RuleActions.Tasks) == 0) {
			fieldName := fmt.Sprintf("flowrules[%d].properties{done}", i)
			vErr := wscutils.BuildErrorMessage(server.MsgId_Empty, server.ErrCode_Empty, &fieldName)
			validationErrors = append(validationErrors, vErr)
		}
		currNS := getNextStep(rule.RuleActions.Properties)
		if len(currNS) > 0 && !isStringInArray(currNS, rule.RuleActions.Tasks) {
			fieldName := fmt.Sprintf("flowrules[%d].properties{nextstep}", i)
			vErr := wscutils.BuildErrorMessage(server.MsgId_NotFound, server.ErrCode_NotFound, &fieldName)
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
