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
	"github.com/remiges-tech/crux/server"
	"github.com/remiges-tech/crux/types"
)

type wfschemagetRow struct {
	Slice         int32               `json:"slice"`
	App           string              `json:"app"`
	Class         string              `json:"class"`
	Longname      string              `json:"longname"`
	Patternschema types.PatternSchema `json:"patternschema"`
	Actionschema  types.ActionSchema  `json:"actionschema"`
	Createdat     pgtype.Timestamp    `json:"createdat"`
	Createdby     string              `json:"createdby"`
	Editedat      pgtype.Timestamp    `json:"editedat"`
	Editedby      pgtype.Text         `json:"editedby"`
}

// SchemaGet will be responsible for processing the /wfschemaget request that comes through as a POST
func SchemaGet(c *gin.Context, s *service.Service) {
	lh := s.LogHarbour
	lh.Log("SchemaGet request received")

	// implement the user realm here
	var userRealm int32 = 1

	var (
		request  SchemaGetReq
		response wfschemagetRow
	)

	isCapable, _ := types.Authz_check(types.OpReq{
		User:      userID,
		CapNeeded: CapForList,
	}, false)

	if !isCapable {
		lh.Info().LogActivity("Unauthorized user:", userID)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_Unauthorized))
		return
	}

	// step 1: json request binding with a struct
	err := wscutils.BindJSON(c, &request)
	if err != nil {
		lh.Debug0().LogActivity("error while binding json request error:", err.Error())
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
		lh.Log("Error while getting query instance from service Dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}
	dbResponse, err := query.Wfschemaget(c, sqlc.WfschemagetParams{
		Slice: request.Slice,
		App:   request.App,
		Class: request.Class,
		Realm: userRealm,
	})
	if err != nil {
		lh.Debug0().LogActivity("failed to get data from DB:", err.Error())
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}

	errors := response.bindSchemaGetResp(s, dbResponse)
	if len(errors) > 0 {
		lh.Debug0().LogActivity("error while converting byte patternschema or action schema to struct:", err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, errors))
		return
	}

	lh.Debug0().Log("Record found finished execution of SchemaGet()")
	wscutils.SendSuccessResponse(c, wscutils.NewSuccessResponse(response))
}

func (response *wfschemagetRow) bindSchemaGetResp(s *service.Service, dbResponse sqlc.WfschemagetRow) []wscutils.ErrorMessage {
	lh := s.LogHarbour
	lh.Log("bindSchemaGetResp request received")
	var (
		pattrn *types.PatternSchema
		action *types.ActionSchema
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
	response.Patternschema = *pattrn
	response.Actionschema = *action
	return errors
}

// To convert byte data to patternschema struct
func byteToStruct(byteData []byte, v any) error {
	err := json.Unmarshal(byteData, &v)
	if err != nil {
		return fmt.Errorf("error decoding JSON: %v", err)
	}
	return nil
}
