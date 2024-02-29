package wfinstance

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/crux/server"
	"github.com/remiges-tech/crux/types"
	"github.com/remiges-tech/logharbour/logharbour"
)

// AbortWFInstance rquest format
type WFInstanceAbortRquest struct {
	ID       *int32  `json:"id" validate:"omitempty,gt=0"`
	EntityID *string `json:"entityid" validate:"omitempty,gt=0,lt=40"`
}

func GetWFInstanceAbort(c *gin.Context, s *service.Service) {
	lh := s.LogHarbour.WithClass("wfinstance")
	lh.Log("GetWFInstanceAbort request received")

	var (
		request  WFInstanceAbortRquest
		id       int32
		entityid string
	)

	isCapable, _ := server.Authz_check(types.OpReq{
		User: USERID,
	}, false)

	if !isCapable {
		lh.Info().LogActivity("GetWFInstanceAbort||unauthorized user:", USERID)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_Unauthorized))
		return
	}

	// Bind request
	err := wscutils.BindJSON(c, &request)
	if err != nil {
		lh.Error(err).Log("GetWFInstanceAbort||error while binding json request error:")
		return
	}
	// Standard validation of Incoming Request
	validationErrors := wscutils.WscValidate(request, func(err validator.FieldError) []string { return []string{} })
	if len(validationErrors) > 0 {
		lh.Debug0().LogActivity("GetWFInstanceAbort||validation error:", validationErrors)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, validationErrors))
		return
	}
	// Custom validation
	if request.ID != nil && request.EntityID != nil {
		lh.Debug0().Log("GetWFInstanceAbort||both id and entityid cannot be present at the same time")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_RequiredOneOf, server.ErrCode_RequiredOne))
		return
	}
	if request.ID == nil && request.EntityID == nil {
		lh.Debug0().Log("GetWFInstanceAbort||either id or entityid must be present")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_RequiredOneOf, server.ErrCode_RequiredOne))
		return
	}

	//Handle the request based on the presence of ID or EntityID
	if request.ID != nil {
		id = *request.ID
		lh.WithInstanceId(string(id))
	} else {
		entityid = *request.EntityID
	}
	lh.Debug0().LogActivity("GetWFInstanceAbort||present values :", map[string]any{"ID": id, "EntityId": entityid})

	query, ok := s.Dependencies["queries"].(*sqlc.Queries)
	if !ok {
		lh.Log("GetWFInstanceAbort||error while getting query instance from service Dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
	}
	// delete wfinstance by ID or entityID
	deletedWfinstaceListByID, error := query.DeleteWfinstanceByID(c, sqlc.DeleteWfinstanceByIDParams{
		ID:       pgtype.Int4{Int32: id, Valid: id != 0},
		Entityid: pgtype.Text{String: entityid, Valid: !server.IsStringEmpty(&entityid)},
	})
	if error != nil {
		lh.Error(error).Log("GetWFInstanceAbort||error while deleting wfinstances by id or entityid")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}

	if len(deletedWfinstaceListByID) == 0 {
		lh.Log("GetWFInstanceAbort||no record found to delete")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_NotFound, server.ErrCode_NotFound))
		return
	}

	// To Get parentList from WFInstanceList data
	parentList := getParentList(deletedWfinstaceListByID)
	lh.Debug0().LogActivity("GetWFInstanceAbort||parentlist form WFInstanceListByID data:", parentList)

	for parentList != nil {
		lh.Debug0().Log("GetWFInstanceAbort||inside for loop : if parentList is Not Nil ")
		// To get GetWFInstanceList by parentList
		deletedWfinstanceListByParents, err := query.DeleteWFInstanceListByParents(c, sqlc.DeleteWFInstanceListByParentsParams{
			ID:     parentList,
			Parent: parentList,
		})
		if err != nil {
			lh.Error(err).Log("GetWFInstanceAbort||error while getting wfinstance List by parentList")
			errmsg := db.HandleDatabaseError(err)
			wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
			return
		}

		// Append wfinstanceListByParents data
		deletedWfinstaceListByID = append(deletedWfinstaceListByID, deletedWfinstanceListByParents...)

		// Update parentList using getParentList function
		parentList = getParentList(deletedWfinstanceListByParents)
		lh.Debug0().LogActivity("GetWFInstanceAbort||updated ParentList :", parentList)
	}

	// data change log
	for _, val := range deletedWfinstaceListByID {
		dclog := lh.WithClass("wfinstance").WithInstanceId(strconv.Itoa(int(val.ID)))
		dclog.LogDataChange("deleted wfinstance ", logharbour.ChangeInfo{
			Entity: "wfinstance",
			Op:     "delete",
			Changes: []logharbour.ChangeDetail{
				{
					Field:  "row",
					OldVal: val,
					NewVal: nil},
			},
		})
	}

	wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus})

}
