package schema

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	crux "github.com/remiges-tech/crux/matching-engine"
	"github.com/remiges-tech/crux/server"
	"github.com/remiges-tech/crux/types"
)

type patternSchema_t struct {
	Attr      string   `json:"attr" validate:"required"`
	ShortDesc string   `json:"shortdesc" validate:"required"`
	LongDesc  string   `json:"longdesc" validate:"required"`
	ValType   string   `json:"valtype" validate:"required"`
	EnumVals  []string `json:"vals,omitempty"`
	ValMin    float64  `json:"valmin,omitempty"`
	ValMax    float64  `json:"valmax,omitempty"`
	LenMin    int      `json:"lenmin,omitempty"`
	LenMax    int      `json:"lenmax,omitempty"`
}
type wfschemagetRow struct {
	Slice         int32               `json:"slice"`
	App           string              `json:"app"`
	Class         string              `json:"class"`
	Longname      string              `json:"longname"`
	Patternschema []patternSchema_t   `json:"patternschema"`
	Actionschema  crux.ActionSchema_t `json:"actionschema"`
	Createdat     pgtype.Timestamp    `json:"createdat"`
	Createdby     string              `json:"createdby"`
	Editedat      pgtype.Timestamp    `json:"editedat"`
	Editedby      pgtype.Text         `json:"editedby"`
}

// SchemaGet will be responsible for processing the /wfschemaget request that comes through as a POST
func SchemaGet(c *gin.Context, s *service.Service) {
	lh := s.LogHarbour
	lh.Log("SchemaGet request received")

	// userID, err := server.ExtractUserNameFromJwt(c)
	// if err != nil {
	// 	lh.Info().Log("unable to extract userID from token")
	// 	wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Missing, server.ErrCode_Token_Data_Missing))
	// 	return
	// }

	// realmName, err := server.ExtractRealmFromJwt(c)
	// if err != nil {
	// 	lh.Info().Log("unable to extract realm from token")
	// 	wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Missing, server.ErrCode_Token_Data_Missing))
	// 	return
	// }

	var (
		request  SchemaGetReq
		response wfschemagetRow
	)

	isCapable, _ := server.Authz_check(types.OpReq{
		User:      userID,
		CapNeeded: CapForList,
	}, false)

	// isCapable := true
	if !isCapable {
		lh.Info().LogActivity("unauthorized user:", userID)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_Unauthorized))
		return
	}

	// step 1: json request binding with a struct
	err = wscutils.BindJSON(c, &request)
	if err != nil {
		lh.Debug0().Error(err).Log("error while binding json request error:")
		return
	}

	// step 2: standard validation
	valError := wscutils.WscValidate(request, func(err validator.FieldError) []string { return []string{} })
	if len(valError) > 0 {
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, valError))
		lh.Debug0().LogActivity("validation error:", valError)
		return
	}
	query, ok := s.Dependencies["queries"].(*sqlc.Queries)
	if !ok {
		lh.Log("error while getting query instance from service Dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}
	dbResponse, err := query.Wfschemaget(c, sqlc.WfschemagetParams{
		Slice: request.Slice,
		App:   request.App,
		Class: request.Class,
		Realm: realmName,
		Brwf:  sqlc.BrwfEnumW,
	})
	if err != nil {
		lh.Debug0().Error(err).Log("failed to get data from db")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}

	errors := response.bindSchemaGetResp(s, dbResponse)
	if len(errors) > 0 {
		lh.Debug0().LogActivity("error while converting byte patternschema or action schema to struct:", errors)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, errors))
		return
	}

	lh.Debug0().Log("record found finished execution of SchemaGet()")
	wscutils.SendSuccessResponse(c, wscutils.NewSuccessResponse(response))
}

func (response *wfschemagetRow) bindSchemaGetResp(s *service.Service, dbResponse sqlc.WfschemagetRow) []wscutils.ErrorMessage {
	lh := s.LogHarbour
	lh.Log("bindSchemaGetResp request received")
	var (
		pattrn []crux.PatternSchema_t
		action crux.ActionSchema_t
		errors []wscutils.ErrorMessage
	)
	response.Slice = dbResponse.Slice
	response.App = dbResponse.App
	response.Class = dbResponse.Class
	response.Longname = dbResponse.Longname
	response.Createdat = dbResponse.Createdat
	response.Createdby = dbResponse.Createdby
	response.Editedat = dbResponse.Editedat
	response.Editedby = dbResponse.Editedby

	err := byteToStruct(dbResponse.Patternschema, &pattrn)
	if err != nil {
		errors = append(errors, wscutils.BuildErrorMessage(server.MsgId_NoSchemaFound, server.ErrCode_Invalid_pattern_schema, nil))
	}
	err = byteToStruct(dbResponse.Actionschema, &action)
	if err != nil {
		errors = append(errors, wscutils.BuildErrorMessage(server.MsgId_NoSchemaFound, server.ErrCode_Invalid_action_schema, nil))
	}

	for _, v := range pattrn {
		var t patternSchema_t
		t.bindPattrnSchemaResp(v)
		for k, _ := range v.EnumVals {
			t.EnumVals = append(t.EnumVals, k)
		}
		response.Patternschema = append(response.Patternschema, t)
	}

	// response.Patternschema = pattrn
	response.Actionschema = action
	return errors
}

func (t *patternSchema_t) bindPattrnSchemaResp(v crux.PatternSchema_t) {
	t.Attr = v.Attr
	t.ShortDesc = v.ShortDesc
	t.LongDesc = v.LongDesc
	t.ValType = v.ValType
	t.ValMin = v.ValMin
	t.ValMax = v.ValMax
	t.LenMin = v.LenMin
	t.LenMax = v.LenMax
}

// To convert byte data to patternschema struct
func byteToStruct(byteData []byte, v any) error {
	err := json.Unmarshal(byteData, &v)
	if err != nil {
		return fmt.Errorf("error decoding JSON: %v", err)
	}
	return nil
}
