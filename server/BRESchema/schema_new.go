package breschema

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	crux "github.com/remiges-tech/crux/matching-engine"
	"github.com/remiges-tech/crux/server"
	"github.com/remiges-tech/crux/types"
	"github.com/remiges-tech/logharbour/logharbour"
)

type PatternSchema struct {
	Attr      string   `json:"attr" validate:"required"`
	ShortDesc string   `json:"shortdesc"`
	LongDesc  string   `json:"longdesc"`
	ValType   string   `json:"valtype" validate:"required"`
	EnumVals  []string `json:"vals,omitempty"`
	ValMin    float64  `json:"valmin,omitempty"`
	ValMax    float64  `json:"valmax,omitempty"`
	LenMin    int      `json:"lenmin,omitempty"`
	LenMax    int      `json:"lenmax,omitempty"`
}

func (p PatternSchema) IsEmpty() bool {
	return reflect.DeepEqual(p, PatternSchema{})
}

type BRESchemaNewReq struct {
	Slice         int32               `json:"slice" validate:"required,gt=0,lt=15"`
	App           string              `json:"app" validate:"required,alpha,lt=15"`
	Class         string              `json:"class" validate:"required,lowercase,lt=15"`
	PatternSchema []PatternSchema     `json:"patternSchema" validate:"required,dive"`
	ActionSchema  crux.ActionSchema_t `json:"actionSchema"`
}

func BRESchemaNew(c *gin.Context, s *service.Service) {
	l := s.LogHarbour
	l.Debug0().Log("starting execution of BRESchemaNew()")

	// userID, err := server.ExtractUserNameFromJwt(c)
	// if err != nil {
	// 	l.Info().Log("unable to extract userID from token")
	// 	wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Missing, server.ErrCode_Token_Data_Missing))
	// 	return
	// }

	// realmName, err := server.ExtractRealmFromJwt(c)
	// if err != nil {
	// 	l.Info().Log("unable to extract realm from token")
	// 	wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Missing, server.ErrCode_Token_Data_Missing))
	// 	return
	// }

	realmName, ok := s.Dependencies["realmName"].(string)
	if !ok {
		l.Debug0().Log("error while getting realmName instance from service dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_Internal))
		return
	}

	// delete below line whie actual implementation (reason: kept for testing while writting api)
	realmName = "Ecommerce"

	caps := []string{"schema"}
	isCapable, _ := server.Authz_check(types.OpReq{
		User:      userID,
		CapNeeded: caps,
	}, false)

	if !isCapable {
		l.Info().LogActivity("unauthorized user:", userID)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_Unauthorized))
		return
	}

	var (
		req BRESchemaNewReq
	)

	err = wscutils.BindJSON(c, &req)
	if err != nil {
		l.Error(err).Log("error unmarshalling query parameters to struct:")
		return
	}
	newPatternSchema := convertPatternSchema(req.PatternSchema)
	schema := crux.Schema_t{
		Class:         req.Class,
		PatternSchema: newPatternSchema,
		ActionSchema:  req.ActionSchema,
		NChecked:      0,
	}

	// Validate request
	validationErrors := wscutils.WscValidate(req, func(err validator.FieldError) []string { return []string{} })
	customValidationErrors := customValidationErrors(schema)
	validationErrors = append(validationErrors, customValidationErrors...)
	if len(validationErrors) > 0 {
		l.Debug0().LogDebug("validation errors", validationErrors)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, validationErrors))
		return
	}

	query, ok := s.Dependencies["queries"].(*sqlc.Queries)
	if !ok {
		l.Debug0().Log("error while getting query instance from service Dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_Internal))
		return
	}
	patternSchema, err := json.Marshal(newPatternSchema)
	if err != nil {
		patternSchema := "patternSchema"
		l.Debug1().LogDebug("error while marshaling patternSchema", err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(server.MsgId_Invalid_Request, server.ErrCode_InvalidJson, &patternSchema)}))
		return
	}

	actionSchema, err := json.Marshal(req.ActionSchema)
	if err != nil {
		actionSchema := "actionSchema"
		l.Debug1().LogDebug("error while marshaling actionSchema", err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(server.MsgId_Invalid_Request, server.ErrCode_InvalidJson, &actionSchema)}))
		return
	}
	id, err := query.SchemaNew(c, sqlc.SchemaNewParams{
		RealmName:     realmName,
		Slice:         req.Slice,
		Class:         req.Class,
		App:           strings.ToLower(req.App),
		Brwf:          sqlc.BrwfEnumB,
		Patternschema: patternSchema,
		Actionschema:  actionSchema,
		Createdby:     userID,
	})
	if err != nil {
		l.Info().Error(err).Log("error while creating BRESchema")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}
	dclog := l.WithClass("BRESchema").WithInstanceId(string(id))
	dclog.LogDataChange("created BRESchema", logharbour.ChangeInfo{
		Entity: "BRESchema",
		Op:     "create",
		Changes: []logharbour.ChangeDetail{
			{
				Field:  "realm",
				OldVal: nil,
				NewVal: realmName},
			{
				Field:  "slice",
				OldVal: nil,
				NewVal: req.Slice},
			{
				Field:  "app",
				OldVal: nil,
				NewVal: req.App},
			{
				Field:  "class",
				OldVal: nil,
				NewVal: req.Class},
			{
				Field:  "brwf",
				OldVal: nil,
				NewVal: sqlc.BrwfEnumB},
			{
				Field:  "patternSchema",
				OldVal: nil,
				NewVal: newPatternSchema},
			{
				Field:  "actionSchema",
				OldVal: nil,
				NewVal: req.ActionSchema},
		},
	})
	wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: nil, Messages: nil})
	l.Debug0().Log("Finished execution of BRESchemaNew()")
}

func customValidationErrors(sh crux.Schema_t) []wscutils.ErrorMessage {
	var validationErrors []wscutils.ErrorMessage
	if len(sh.PatternSchema) > 0 {
		err := crux.VerifyPatternSchema(sh, true)
		if err != nil {
			patternSchemaError := server.HandleCruxError(err)
			validationErrors = append(validationErrors, patternSchemaError...)
		}
	}

	if server.IsZeroOfUnderlyingType(sh.ActionSchema) {
		err := crux.VerifyActionSchema(sh, true)
		if err != nil {
			actionSchemaError := server.HandleCruxError(err)
			validationErrors = append(validationErrors, actionSchemaError...)
		}
	}

	return validationErrors
}

func convertPatternSchema(oldPatternSchema []PatternSchema) []crux.PatternSchema_t {
	var newPatternSchema []crux.PatternSchema_t
	for _, patternSchema := range oldPatternSchema {
		if patternSchema.IsEmpty() {
			return newPatternSchema
		}
		newEnumVals := make(map[string]struct{})
		for _, val := range patternSchema.EnumVals {
			newEnumVals[val] = struct{}{}
		}

		patternSchema := crux.PatternSchema_t{
			Attr:      patternSchema.Attr,
			ShortDesc: patternSchema.ShortDesc,
			LongDesc:  patternSchema.LongDesc,
			ValType:   patternSchema.ValType,
			EnumVals:  newEnumVals,
			ValMin:    patternSchema.ValMin,
			ValMax:    patternSchema.ValMax,
			LenMin:    patternSchema.LenMin,
			LenMax:    patternSchema.LenMax,
		}
		newPatternSchema = append(newPatternSchema, patternSchema)
	}
	return newPatternSchema
}
