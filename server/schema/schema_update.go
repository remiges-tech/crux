package schema

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/crux/server"
	"github.com/remiges-tech/crux/types"
	"github.com/remiges-tech/logharbour/logharbour"
)

type updateSchema struct {
	Slice         int32          `json:"slice" validate:"required,gt=0"`
	App           string         `json:"App" validate:"required,alpha"`
	Class         string         `json:"class" validate:"required,lowercase"`
	PatternSchema *patternSchema `json:"patternSchema,omitempty"`
	ActionSchema  *actionSchema  `json:"actionSchema,omitempty"`
}

func SchemaUpdate(c *gin.Context, s *service.Service) {
	l := s.LogHarbour
	l.Debug0().Log("Starting execution of SchemaUpdate()")

	isCapable, _ := types.Authz_check(types.OpReq{
		User:      userID,
		CapNeeded: capForUpdate,
	}, false)

	if !isCapable {
		l.Info().LogActivity("Unauthorized user:", userID)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_Unauthorized))
		return
	}

	var sh updateSchema

	err := wscutils.BindJSON(c, &sh)
	if err != nil {
		l.LogActivity("Error Unmarshalling Query paramaeters to struct:", err.Error())
		return
	}

	// Validate request
	validationErrors := wscutils.WscValidate(sh, func(err validator.FieldError) []string { return []string{} })
	customValidationErrors := customValidationErrorsForUpdate(sh)
	validationErrors = append(validationErrors, customValidationErrors...)
	if len(validationErrors) > 0 {
		l.Debug0().LogDebug("standard validation errors", validationErrors)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, validationErrors))
		return
	}

	query, ok := s.Dependencies["queries"].(*sqlc.Queries)
	if !ok {
		l.Debug0().Log("Error while getting query instance from service Dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}

	connpool, ok := s.Database.(*pgxpool.Pool)
	if !ok {
		l.Debug0().Log("Error while getting query instance from service Database")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}

	err = schemaUpdateWithTX(c, query, connpool, l, sh)
	if err != nil {
		l.LogActivity("Error while Updating schema", err.Error())
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}
	wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: nil, Messages: nil})
	l.Debug0().Log("Starting execution of SchemaUpdate()")
}

func schemaUpdateWithTX(c context.Context, query *sqlc.Queries, connpool *pgxpool.Pool, l *logharbour.Logger, sh updateSchema) error {
	patternSchema, err := json.Marshal(sh.PatternSchema)
	if err != nil {
		l.Debug1().LogDebug("Error while marshaling patternSchema", err)
		return err
	}

	actionSchema, err := json.Marshal(sh.ActionSchema)
	if err != nil {
		l.Debug1().LogDebug("Error while marshaling actionSchema", err)
		return err
	}

	tx, err := connpool.Begin(c)
	if err != nil {
		return err
	}
	defer tx.Rollback(c)
	qtx := query.WithTx(tx)
	schema, err := qtx.GetSchemaWithLock(c, sqlc.GetSchemaWithLockParams{
		Slice: sh.Slice,
		Class: sh.Class,
		App:   sh.App,
	})
	if err != nil {
		tx.Rollback(c)
		return err
	}
	if patternSchema != nil && actionSchema == nil {
		err = qtx.SchemaUpdate(c, sqlc.SchemaUpdateParams{
			Slice:         sh.Slice,
			Class:         sh.Class,
			App:           sh.App,
			Brwf:          sqlc.BrwfEnumW,
			Patternschema: schema.Patternschema,
			Actionschema:  actionSchema,
			Editedby:      pgtype.Text{String: userID},
		})
		if err != nil {
			tx.Rollback(c)
			return err
		}
	} else if actionSchema != nil && patternSchema == nil {
		err = qtx.SchemaUpdate(c, sqlc.SchemaUpdateParams{
			Slice:         sh.Slice,
			Class:         sh.Class,
			App:           sh.App,
			Brwf:          sqlc.BrwfEnumW,
			Patternschema: schema.Patternschema,
			Actionschema:  actionSchema,
			Editedby:      pgtype.Text{String: userID},
		})
		if err != nil {
			tx.Rollback(c)
			return err
		}
	} else {
		err = qtx.SchemaUpdate(c, sqlc.SchemaUpdateParams{
			Slice:         sh.Slice,
			Class:         sh.Class,
			App:           sh.App,
			Brwf:          sqlc.BrwfEnumW,
			Patternschema: schema.Patternschema,
			Actionschema:  schema.Actionschema,
			Editedby:      pgtype.Text{String: userID},
		})
		if err != nil {
			tx.Rollback(c)
			return err
		}

	}

	if err := tx.Commit(c); err != nil {
		return err
	}
	dclog := l.WithWhatClass("schema").WithWhatInstanceId(string(schema.ID))
	dclog.LogDataChange("Updated schema", logharbour.ChangeInfo{
		Entity:    "schema",
		Operation: "Update",
		Changes: []logharbour.ChangeDetail{
			{
				Field:    "patternSchema",
				OldValue: string(schema.Patternschema),
				NewValue: patternSchema},
			{
				Field:    "actionSchema",
				OldValue: string(schema.Actionschema),
				NewValue: sh.ActionSchema},
		},
	})

	return nil
}

func customValidationErrorsForUpdate(sh updateSchema) []wscutils.ErrorMessage {
	var validationErrors []wscutils.ErrorMessage
	if sh.PatternSchema != nil && sh.ActionSchema == nil {
		patternSchemaError := verifyPatternSchemaUpdate(sh.PatternSchema)
		validationErrors = append(validationErrors, patternSchemaError...)
	} else if sh.ActionSchema != nil && sh.PatternSchema == nil {
		actionSchemaError := verifyActionSchemaUpdate(sh.ActionSchema)
		validationErrors = append(validationErrors, actionSchemaError...)
	} else if sh.PatternSchema == nil && sh.ActionSchema == nil {
		fieldName := fmt.Sprintln("PatternSchema/ActionSchema")
		vErr := wscutils.BuildErrorMessage(server.MsgId_RequiredOneOf, server.ErrCode_RequiredOne, &fieldName)
		validationErrors = append(validationErrors, vErr)
	} else {
		patternSchemaError := verifyPatternSchemaUpdate(sh.PatternSchema)
		validationErrors = append(validationErrors, patternSchemaError...)
		actionSchemaError := verifyActionSchemaUpdate(sh.ActionSchema)
		validationErrors = append(validationErrors, actionSchemaError...)
	}

	return validationErrors
}
func verifyPatternSchemaUpdate(ps *patternSchema) []wscutils.ErrorMessage {
	var validationErrors []wscutils.ErrorMessage

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
	}
	return validationErrors
}

func verifyActionSchemaUpdate(as *actionSchema) []wscutils.ErrorMessage {
	var validationErrors []wscutils.ErrorMessage
	re := regexp.MustCompile(cruxIDRegExp)

	for i, task := range as.Tasks {
		if !re.MatchString(task) {
			fieldName := fmt.Sprintf("actionSchema.Tasks[%d]", i)
			vErr := wscutils.BuildErrorMessage(server.MsgId_Invalid, server.ErrCode_Invalid, &fieldName, task)
			validationErrors = append(validationErrors, vErr)
		}
	}
	for i, propName := range as.Properties {
		if !re.MatchString(propName) {
			fieldName := fmt.Sprintf("actionSchema.Properties[%d]", i)
			vErr := wscutils.BuildErrorMessage(server.MsgId_Invalid, server.ErrCode_Invalid, &fieldName, propName)
			validationErrors = append(validationErrors, vErr)
		}
	}
	return validationErrors
}
