package wfinstance

import (
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

// AbortWFInstance rquest format
type WFInstanceAbortRquest struct {
	ID       *int32  `json:"id,omitempty"`
	EntityID *string `json:"entityid,omitempty"`
}

func GetWFInstanceAbort(c *gin.Context, s *service.Service) {
	lh := s.LogHarbour.WithWhatClass("wfinstance")
	lh.Log("AbortWFInstance request received")

	var (
		request  WFInstanceAbortRquest
		id       int32
		entityid string
	)

	isCapable, _ := types.Authz_check(types.OpReq{
		User: USERID,
	}, false)

	if !isCapable {
		lh.Info().LogActivity("Unauthorized user:", USERID)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_Unauthorized))
		return
	}

	// Bind request
	err := wscutils.BindJSON(c, &request)
	if err != nil {
		lh.Debug0().LogActivity("error while binding json request error:", err)
		return
	}
	// Standard validation of Incoming Request
	validationErrors := wscutils.WscValidate(request, func(err validator.FieldError) []string { return []string{} })
	if len(validationErrors) > 0 {
		lh.Debug0().LogActivity("validation error:", validationErrors)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, validationErrors))
		return
	}
	// Custom validation
	if request.ID != nil && request.EntityID != nil {
		lh.Debug0().Log("Both ID and EntityID cannot be present at the same time")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_RequiredOneOf, server.ErrCode_RequiredOne))
		return
	}
	if request.ID == nil && request.EntityID == nil {
		lh.Debug0().Log("Either ID or EntityID must be present")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_RequiredOneOf, server.ErrCode_RequiredOne))
		return
	}

	//Handle the request based on the presence of ID or EntityID
	if request.ID != nil {
		id = *request.ID
		lh.WithWhatInstanceId(string(id))
	} else {
		entityid = *request.EntityID
	}
	lh.Debug0().LogActivity("present values :", map[string]any{"ID": id, "EntityId": entityid})

	query, ok := s.Dependencies["queries"].(*sqlc.Queries)
	if !ok {
		lh.Log("Error while getting query instance from service Dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
	}

	tag, error := query.DeleteWfInstance(c, sqlc.DeleteWfInstanceParams{
		ID:       pgtype.Int4{Int32: id, Valid: id != 0},
		Entityid: pgtype.Text{String: entityid, Valid: !types.IsStringEmpty(&entityid)},
	})
	if error != nil {
		lh.LogActivity("error while deleting wfinstances  :", error.Error())
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}
	lh.Debug0().LogActivity("tag :", tag)
	if tag == -1 {
		lh.Log("no record found to delete")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_NotFound, server.ErrCode_NotFound))
		return
	}

	wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus})

}
