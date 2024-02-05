package schema

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/logharbour/logharbour"
)

type attribute struct {
	Name      string   `json:"name,omitempty"`
	ShortDesc string   `json:"shortdesc,omitempty"`
	LongDesc  string   `json:"longdesc,omitempty"`
	ValType   string   `json:"valtype,omitempty"`
	Vals      []string `json:"vals,omitempty"`
	Enumdesc  []string `json:"enumdesc,omitempty"`
	ValMax    int32    `json:"valmax,omitempty"`
	ValMin    int32    `json:"valmin,omitempty"`
	LenMax    int32    `json:"lenmax,omitempty"`
	LenMin    int32    `json:"lenmin,omitempty"`
}
type patternSchema struct {
	Class *string      `json:"class,omitempty"`
	Attr  []*attribute `json:"attr,omitempty"`
}
type actionSchema struct {
	Class      string   `json:"class,omitempty"`
	Tasks      []string `json:"tasks,omitempty"`
	Properties []string `json:"properties,omitempty"`
}
type updateSchema struct {
	Slice         int32          `json:"slice" validate:"required,gt=0"`
	App           string         `json:"App" validate:"required,alpha"`
	Class         string         `json:"class" validate:"required,lowercase"`
	PatternSchema *patternSchema `json:"patternSchema,omitempty"`
	ActionSchema  *actionSchema  `json:"actionSchema,omitempty"`
}

const (
	cruxIDRegExp = `^[a-z][a-z0-9_]*$`
)

var (
	editedBy   = "admin"
	validTypes = map[string]bool{
		"int": true, "float": true, "str": true, "enum": true, "bool": true, "timestamps": true,
	}
)

func SchemaUpdate(c *gin.Context, s *service.Service) {
	l := s.LogHarbour
	l.Log("Starting execution of SchemaUpdate()")

	var sh updateSchema
	// check the capgrant table to see if the calling user has the capability to perform the
	// operation
	// isCapable, _ := utils.Authz_check(types.OpReq{
	// 	User:      username,
	// 	CapNeeded: []string{"schema"},
	// }, false)

	// if !isCapable {
	// 	l.Log("Unauthorized user:")
	// 	wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(utils.ErrUnauthorized))
	// 	return
	// }

	// The system will check whether there are any rulesets in the ruleset table whose
	// (slice,app,class) match this record. If this is true, then the call will fail.
	// In other words, updating a schema is not allowed once rulesets referring to it are defined.

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
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, validationErrors))
		return
	}

	query, ok := s.Dependencies["queries"].(*sqlc.Queries)
	if !ok {
		l.Log("Error while getting query instance from service Dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(wscutils.ErrcodeDatabaseError))
		return
	}

	connpool, ok := s.Database.(*pgxpool.Pool)
	if !ok {
		l.Log("Error while getting query instance from service Dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(wscutils.ErrcodeDatabaseError))
		return
	}

	err = schemaUpdateWithTX(c, query, connpool, l, sh)
	if err != nil {
		l.LogActivity("Error while Updating schema", err.Error())
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(wscutils.ErrcodeDatabaseError))
		return
	}
	wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: "updated successfully", Messages: nil})
	l.Log("Starting execution of SchemaUpdate()")
}

func schemaUpdateWithTX(c context.Context, query *sqlc.Queries, connpool *pgxpool.Pool, l *logharbour.Logger, sh updateSchema) error {
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

	tx, err := connpool.Begin(c)
	if err != nil {
		return err
	}
	defer tx.Rollback(c)
	qtx := query.WithTx(tx)
	schema, err := qtx.UpdateSchemaWithLock(c, sqlc.UpdateSchemaWithLockParams{
		Slice: sh.Slice,
		Class: sh.Class,
		App:   sh.App,
	})
	if err != nil {
		tx.Rollback(c)
		return err
	}
	if patternSchema != nil {
		_, err = qtx.SchemaUpdate(c, sqlc.SchemaUpdateParams{
			Slice:         sh.Slice,
			Class:         sh.Class,
			App:           sh.App,
			Brwf:          "W",
			Patternschema: schema.Patternschema,
			Actionschema:  actionSchema,
			Editedby:      editedBy,
		})
		if err != nil {
			tx.Rollback(c)
			return err
		}
	} else if actionSchema != nil {
		_, err = qtx.SchemaUpdate(c, sqlc.SchemaUpdateParams{
			Slice:         sh.Slice,
			Class:         sh.Class,
			App:           sh.App,
			Brwf:          "W",
			Patternschema: schema.Patternschema,
			Actionschema:  actionSchema,
			Editedby:      editedBy,
		})
		if err != nil {
			tx.Rollback(c)
			return err
		}
	} else {
		_, err = qtx.SchemaUpdate(c, sqlc.SchemaUpdateParams{
			Slice:         sh.Slice,
			Class:         sh.Class,
			App:           sh.App,
			Brwf:          "W",
			Patternschema: schema.Patternschema,
			Actionschema:  schema.Actionschema,
			Editedby:      editedBy,
		})
		if err != nil {
			tx.Rollback(c)
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
	if sh.PatternSchema != nil {
		patternSchemaError := verifyPatternSchema(sh.PatternSchema)
		validationErrors = append(validationErrors, patternSchemaError...)
	}
	if sh.ActionSchema != nil {
		actionSchemaError := verifyActionSchema(sh.ActionSchema)
		validationErrors = append(validationErrors, actionSchemaError...)
	}
	if sh.PatternSchema == nil && sh.ActionSchema == nil {
		fieldName := fmt.Sprintln("PatternSchema/ActionSchema")
		vErr := wscutils.BuildErrorMessage("at_least_one_must_be_supplied", &fieldName)
		validationErrors = append(validationErrors, vErr)
	}

	return validationErrors
}
func verifyPatternSchema(ps *patternSchema) []wscutils.ErrorMessage {
	var validationErrors []wscutils.ErrorMessage
	re := regexp.MustCompile(cruxIDRegExp)

	for i, attrSchema := range ps.Attr {
		i++
		if !re.MatchString(attrSchema.Name) {
			fieldName := fmt.Sprintf("attrSchema[%d].Name", i)
			vErr := wscutils.BuildErrorMessage("not_valid", &fieldName, attrSchema.Name)
			validationErrors = append(validationErrors, vErr)
		}
		if !validTypes[attrSchema.ValType] {
			fieldName := fmt.Sprintf("attrSchema[%d].ValType", i)
			vErr := wscutils.BuildErrorMessage("not_valid", &fieldName, attrSchema.ValType)
			validationErrors = append(validationErrors, vErr)
		}
		if attrSchema.ValType == "enum" && len(attrSchema.Vals) == 0 {
			fieldName := fmt.Sprintf("attrSchema[%d].Vals", i)
			vErr := wscutils.BuildErrorMessage("empty", &fieldName)
			validationErrors = append(validationErrors, vErr)
		}
	}
	return validationErrors
}

func verifyActionSchema(as *actionSchema) []wscutils.ErrorMessage {
	var validationErrors []wscutils.ErrorMessage
	re := regexp.MustCompile(cruxIDRegExp)

	for i, task := range as.Tasks {
		if !re.MatchString(task) {
			fieldName := fmt.Sprintf("actionSchema.Tasks[%d]", i)
			vErr := wscutils.BuildErrorMessage("not_valid", &fieldName, task)
			validationErrors = append(validationErrors, vErr)
		}
	}
	for i, propName := range as.Properties {
		if !re.MatchString(propName) {
			fieldName := fmt.Sprintf("actionSchema.Properties[%d]", i)
			vErr := wscutils.BuildErrorMessage("not_valid", &fieldName, propName)
			validationErrors = append(validationErrors, vErr)
		}
	}
	return validationErrors
}
