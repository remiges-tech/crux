package schema

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/crux/server"
	constants "github.com/remiges-tech/crux/serverBRE"
	"github.com/remiges-tech/logharbour/logharbour"
	"regexp"
)

func BRESchemaUpdate(c *gin.Context, s *service.Service) {
	l := s.LogHarbour
	l.Log("Starting execution of BRESchemaUpdate()")

	var sh updateSchema

	err := wscutils.BindJSON(c, &sh)
	if err != nil {
		l.LogActivity("Error Unmarshalling Query parameters to struct:", err.Error())
		return
	}

	// Validate request
	validationErrors := wscutils.WscValidate(sh, func(err validator.FieldError) []string { return []string{} })
	customValidationErrors := customValidationErrorsForUpdate(sh)
	validationErrors = append(validationErrors, customValidationErrors...)
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

	connectionPool, ok := s.Database.(*pgxpool.Pool)
	if !ok {
		l.Log("Error while getting query instance from service Dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}

	err = schemaUpdateWithTX(c, query, connectionPool, l, sh)
	if err != nil {
		l.LogActivity("Error while Updating schema", err.Error())
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}
	wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: nil, Messages: nil})
	l.Log("Starting execution of BRESchemaUpdate()")
}

func schemaUpdateWithTX(c context.Context, query *sqlc.Queries, connectionPool *pgxpool.Pool, l *logharbour.Logger, sh updateSchema) error {
	patternSchema, err := json.Marshal(sh.PatternSchema)
	if err != nil {
		l.LogDebug("Error while marshaling patternSchema", err)
		return err
	}

	actionSchema, err := json.Marshal(sh.ActionSchema)
	if err != nil {
		l.LogDebug("Error while marshaling actionSchema", err)
		return err
	}

	tx, err := connectionPool.Begin(c)
	if err != nil {
		return err
	}
	defer tx.Rollback(c) // error handling needed here ?
	qtx := query.WithTx(tx)
	schema, err := qtx.GetSchemaWithLock(c, sqlc.GetSchemaWithLockParams{
		Slice: sh.Slice,
		Class: sh.Class,
		App:   sh.App,
	})
	if err != nil {
		// error handling needed here ?
		tx.Rollback(c)
		return err
	}
	if patternSchema != nil && actionSchema == nil {
		err = qtx.SchemaUpdate(c, sqlc.SchemaUpdateParams{
			Slice:         sh.Slice,
			Class:         sh.Class,
			App:           sh.App,
			Brwf:          constants.BRWF,
			Patternschema: schema.Patternschema,
			Actionschema:  actionSchema,
			Editedby:      pgtype.Text{String: constants.EditedBy},
		})
		if err != nil {
			tx.Rollback(c) // error handling needed here ?
			return err
		}
	} else if actionSchema != nil && patternSchema == nil {
		err = qtx.SchemaUpdate(c, sqlc.SchemaUpdateParams{
			Slice:         sh.Slice,
			Class:         sh.Class,
			App:           sh.App,
			Brwf:          constants.BRWF,
			Patternschema: schema.Patternschema,
			Actionschema:  actionSchema,
			Editedby:      pgtype.Text{String: constants.EditedBy},
		})
		if err != nil {
			tx.Rollback(c) // error handling needed here ?
			return err
		}
	} else {
		err = qtx.SchemaUpdate(c, sqlc.SchemaUpdateParams{
			Slice:         sh.Slice,
			Class:         sh.Class,
			App:           sh.App,
			Brwf:          constants.BRWF,
			Patternschema: schema.Patternschema,
			Actionschema:  schema.Actionschema,
			Editedby:      pgtype.Text{String: constants.EditedBy},
		})
		if err != nil {
			tx.Rollback(c) // error handling needed here ?
			return err
		}

	}

	if err := tx.Commit(c); err != nil {
		return err
	}

	l.LogDataChange("Updated schema", logharbour.ChangeInfo{
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
				NewValue: actionSchema},
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
	re := regexp.MustCompile(constants.CruxIDRegExp)

	for i, attrSchema := range ps.Attr {
		i++
		if !re.MatchString(attrSchema.Name) {
			fieldName := fmt.Sprintf("attrSchema[%d].Name", i)
			vErr := wscutils.BuildErrorMessage(server.MsgId_Invalid, server.ErrCode_Invalid, &fieldName, attrSchema.Name)
			validationErrors = append(validationErrors, vErr)
		}
		if !constants.ValidTypes[attrSchema.ValType] {
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
	re := regexp.MustCompile(constants.CruxIDRegExp)

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
