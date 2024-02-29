package schema

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/crux/server"
	"github.com/remiges-tech/crux/types"
	"github.com/remiges-tech/logharbour/logharbour"
)

func SchemaNew(c *gin.Context, s *service.Service) {
	l := s.LogHarbour
	l.Debug0().Log("Starting execution of SchemaNew()")

	isCapable, _ := server.Authz_check(types.OpReq{
		User:      userID,
		CapNeeded: capForNew,
	}, false)

	if !isCapable {
		l.Info().LogActivity("Unauthorized user:", userID)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_Unauthorized))
		return
	}

	var sh Schema

	err := wscutils.BindJSON(c, &sh)
	if err != nil {
		l.Error(err).Log("Error Unmarshalling Query parameters to struct:")
		return
	}

	// Validate request
	validationErrors := wscutils.WscValidate(sh, func(err validator.FieldError) []string { return []string{} })
	customValidationErrors := customValidationErrors(sh)
	validationErrors = append(validationErrors, customValidationErrors...)
	if len(validationErrors) > 0 {
		l.Debug0().LogDebug("validation errors", validationErrors)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, validationErrors))
		return
	}

	query, ok := s.Dependencies["queries"].(*sqlc.Queries)
	if !ok {
		l.Debug0().Log("Error while getting query instance from service Dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_Internal))
		return
	}
	patternSchema, err := json.Marshal(sh.PatternSchema)
	if err != nil {
		patternSchema := "patternSchema"
		l.Debug1().LogDebug("Error while marshaling patternSchema", err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(server.MsgId_Invalid_Request, server.ErrCode_InvalidJson, &patternSchema)}))
		return
	}

	actionSchema, err := json.Marshal(sh.ActionSchema)
	if err != nil {
		actionSchema := "actionSchema"
		l.Debug1().LogDebug("Error while marshaling actionSchema", err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(server.MsgId_Invalid_Request, server.ErrCode_InvalidJson, &actionSchema)}))
		return
	}
	id, err := query.SchemaNew(c, sqlc.SchemaNewParams{
		Realm:         realmID,
		Slice:         sh.Slice,
		Class:         sh.Class,
		App:           sh.App,
		Brwf:          sqlc.BrwfEnumW,
		Patternschema: patternSchema,
		Actionschema:  actionSchema,
		Createdby:     userID,
	})
	if err != nil {
		l.Info().Error(err).Log("Error while creating schema")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}
	dclog := l.WithClass("schema").WithInstanceId(string(id))
	dclog.LogDataChange("created schema", logharbour.ChangeInfo{
		Entity: "schema",
		Op:     "create",
		Changes: []logharbour.ChangeDetail{
			{
				Field:  "realm",
				OldVal: nil,
				NewVal: realmID},
			{
				Field:  "slice",
				OldVal: nil,
				NewVal: sh.Slice},
			{
				Field:  "app",
				OldVal: nil,
				NewVal: sh.App},
			{
				Field:  "class",
				OldVal: nil,
				NewVal: sh.Class},
			{
				Field:  "brwf",
				OldVal: nil,
				NewVal: sqlc.BrwfEnumW},
			{
				Field:  "patternSchema",
				OldVal: nil,
				NewVal: patternSchema},
			{
				Field:  "actionSchema",
				OldVal: nil,
				NewVal: sh.ActionSchema},
		},
	})
	wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: nil, Messages: nil})
	l.Debug0().Log("Finished execution of SchemaNew()")
}

func customValidationErrors(sh Schema) []wscutils.ErrorMessage {
	var validationErrors []wscutils.ErrorMessage
	patternSchemaError := verifyPatternSchema(sh.PatternSchema)
	validationErrors = append(validationErrors, patternSchemaError...)

	actionSchemaError := verifyActionSchema(sh)
	validationErrors = append(validationErrors, actionSchemaError...)
	return validationErrors
}

func verifyPatternSchema(ps types.PatternSchema) []wscutils.ErrorMessage {
	var validationErrors []wscutils.ErrorMessage
	stepFound, stepFailedFound := false, false
	for i, attrSchema := range ps.Attr {
		i++
		if !re.MatchString(attrSchema.Name) {
			fieldName := fmt.Sprintf("attrSchema[%d].Name", i)
			vErr := wscutils.BuildErrorMessage(server.MsgId_Invalid, server.ErrCode_Invalid, &fieldName, attrSchema.Name)
			validationErrors = append(validationErrors, vErr)
		}
		if !validTypes[attrSchema.ValType] {
			fieldName := fmt.Sprintf("attrSchema[%d].ValType", i)
			vErr := wscutils.BuildErrorMessage(server.MsgId_Invalid, server.ErrCode_Invalid, &fieldName, attrSchema.ValType)
			validationErrors = append(validationErrors, vErr)
		}
		if attrSchema.ValType == "enum" && len(attrSchema.Vals) == 0 {
			fieldName := fmt.Sprintf("attrSchema[%d].Vals", i)
			vErr := wscutils.BuildErrorMessage(server.MsgId_Empty, server.ErrCode_Empty, &fieldName)
			validationErrors = append(validationErrors, vErr)
		}
		if attrSchema.Name == step && attrSchema.ValType == typeEnum {
			stepFound = true
		}
		val := sliceToMap(attrSchema.Vals)
		if attrSchema.Name == step && !val[start] {
			fieldName := fmt.Sprintf("attrSchema[%d].Vals", i)
			vErr := wscutils.BuildErrorMessage(server.MsgId_Invalid_Request, server.ErrCode_Required, &fieldName)
			validationErrors = append(validationErrors, vErr)
		}
		if attrSchema.Name == stepFailed && attrSchema.ValType == typeBool {
			stepFailedFound = true
		}
	}
	if !stepFound || !stepFailedFound {
		fieldName := "attr.Name"
		vErr := wscutils.BuildErrorMessage(server.MsgId_Invalid_Request, server.ErrCode_RequiredOneOf, &fieldName)
		validationErrors = append(validationErrors, vErr)
	}
	return validationErrors
}

func verifyActionSchema(sh Schema) []wscutils.ErrorMessage {
	var validationErrors []wscutils.ErrorMessage
	nextStepFound, doneFound := false, false
	for i, task := range sh.ActionSchema.Tasks {
		if !re.MatchString(task) {
			fieldName := fmt.Sprintf("actionSchema.Tasks[%d]", i)
			vErr := wscutils.BuildErrorMessage(server.MsgId_Invalid, server.ErrCode_Invalid, &fieldName, task)
			validationErrors = append(validationErrors, vErr)
		}
	}
	if len(sh.ActionSchema.Properties) != 2 {
		fieldName := "Properties"
		vErr := wscutils.BuildErrorMessage(server.MsgId_Invalid_Request, server.ErrCode_Required_Exactly_Two_Properties, &fieldName)
		validationErrors = append(validationErrors, vErr)
	}
	for i, propName := range sh.ActionSchema.Properties {
		if !re.MatchString(propName) {
			fieldName := fmt.Sprintf("actionSchema.Properties[%d]", i)
			vErr := wscutils.BuildErrorMessage(server.MsgId_Invalid, server.ErrCode_Invalid, &fieldName, propName)
			validationErrors = append(validationErrors, vErr)
		} else if propName == nextStep {
			nextStepFound = true
		} else if propName == done {
			doneFound = true
		}
	}

	if !nextStepFound || !doneFound {
		fieldName := "actionSchema.Properties"
		vErr := wscutils.BuildErrorMessage(server.MsgId_Invalid, server.ErrCode_Does_Not_Contain_Both_Properties_Nextstep_And_Done, &fieldName)
		validationErrors = append(validationErrors, vErr)
	}
	if !reflect.DeepEqual(getTasksMapForWF(sh.ActionSchema.Tasks), getStepAttrVals(sh)) {
		fieldName := "actionSchema.Properties"
		vErr := wscutils.BuildErrorMessage(server.MsgId_Invalid, server.ErrCode_ActionSchema_Task_Not_Same_As_PatternSchema_Step_Attr, &fieldName)
		validationErrors = append(validationErrors, vErr)
	}
	return validationErrors
}

func getTasksMapForWF(tasks []string) map[string]bool {
	tm := map[string]bool{}
	for _, t := range tasks {
		tm[t] = true
	}
	// To allow comparison with the set of valid values for the 'step' attribute, which includes "START"
	tm[start] = true
	return tm
}

func getStepAttrVals(sh Schema) map[string]bool {
	for _, ps := range sh.PatternSchema.Attr {
		if ps.Name == step {
			val := sliceToMap(ps.Vals)
			return val
		}
	}
	return nil
}

func sliceToMap(slice []string) map[string]bool {
	stringMap := make(map[string]bool)
	for _, s := range slice {
		stringMap[s] = true
	}
	return stringMap
}
