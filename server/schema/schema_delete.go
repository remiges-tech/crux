package schema

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/crux/server"
	"github.com/remiges-tech/crux/types"
	"github.com/remiges-tech/logharbour/logharbour"
)

// SchemaDelete will be responsible for processing the /WFschemaDelete request that comes through as a POST
func SchemaDelete(c *gin.Context, s *service.Service) {
	lh := s.LogHarbour
	lh.Log("SchemaDelete request received")
	userID, err := server.ExtractUserNameFromJwt(c)
	if err != nil {
		lh.Info().Log("unable to extract userID from token")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Missing, server.ErrCode_Token_Data_Missing))
		return
	}

	realmName, err := server.ExtractRealmFromJwt(c)
	if err != nil {
		lh.Info().Log("unable to extract realm from token")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Missing, server.ErrCode_Token_Data_Missing))
		return
	}

	isCapable, _ := server.Authz_check(types.OpReq{
		User:      userID,
		CapNeeded: CapForList,
	}, false)

	if !isCapable {
		lh.Info().LogActivity("Unauthorized user:", userID)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_Unauthorized))
		return
	}

	var request SchemaGetReq
	err = wscutils.BindJSON(c, &request)
	if err != nil {
		lh.Debug0().LogActivity("error while binding json request error:", err.Error)
		return
	}

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
	deletedSchemaData, err := query.Wfschemadelete(c, sqlc.WfschemadeleteParams{
		Slice: request.Slice,
		App:   request.App,
		Class: request.Class,
		Realm: realmName,
		Brwf:  sqlc.BrwfEnumW,
	})
	if err != nil {
		lh.Debug0().LogActivity("failed while deleting record:", err.Error())
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}

	// data change log
	for _, val := range deletedSchemaData {
		dclog := lh.WithClass("WFSchema").WithInstanceId(strconv.Itoa(int(val.ID)))
		dclog.LogDataChange("deleted WFSchema ", logharbour.ChangeInfo{
			Entity: "WFSchema",
			Op:     "delete",
			Changes: []logharbour.ChangeDetail{
				{
					Field:  "row",
					OldVal: val,
					NewVal: nil,
				},
			},
		})
	}

	lh.Debug0().Log("Record deleted finished execution of SchemaDelete()")
	wscutils.SendSuccessResponse(c, wscutils.NewSuccessResponse(err))
}
