package markdone

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	crux "github.com/remiges-tech/crux/matching-engine"
	"github.com/remiges-tech/crux/server"
	"github.com/remiges-tech/crux/types"
)

var (
	userID    = "1234"
	realmName = "BSE"
)

type WFInstanceMarkDoneReq struct {
	ID         int32             `json:"id" validate:"required"`
	Entity     map[string]string `json:"entity" validate:"required"`
	Step       string            `json:"step" validate:"required,alpha"`
	Stepfailed bool              `json:"stepfailed"`
	Trace      int               `json:"trace,omitempty"`
}

func WFInstanceMarkDone(c *gin.Context, s *service.Service) {
	l := s.LogHarbour
	l.Debug0().Log("starting execution of WFInstanceMarkDone()")

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
	reqCaps := []string{"root"}
	isCapable, _ := server.Authz_check(types.OpReq{
		User:      userID,
		CapNeeded: reqCaps,
	}, false)

	if !isCapable {
		l.Info().LogActivity("Unauthorized user:", userID)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_Unauthorized))
		return
	}

	var req WFInstanceMarkDoneReq

	err := wscutils.BindJSON(c, &req)
	if err != nil {
		l.Error(err).Log("Error Unmarshalling Query parameters to struct:")
		return
	}
	// Validate request
	validationErrors := wscutils.WscValidate(req, func(err validator.FieldError) []string { return []string{} })
	if len(validationErrors) > 0 {
		l.Debug0().LogDebug("validation errors", validationErrors)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, validationErrors))
		return
	}

	queries, ok := s.Dependencies["queries"].(*sqlc.Queries)
	if !ok {
		l.Debug0().Log("Error while getting query instance from service Dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}

	connpool, ok := s.Database.(*pgxpool.Pool)
	if !ok {
		l.Debug0().Log("Error while getting connection pool instance from service Database")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}
	tx, err := connpool.Begin(c)
	if err != nil {
		l.LogActivity("Error while Begin tx", err.Error())
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}
	defer tx.Rollback(c)
	qtx := queries.WithTx(tx)

	// get instance record
	wfinstance, err := queries.GetWFInstanceFromId(c, req.ID)
	if err != nil {
		l.Error(err).Log("Error while GetWFInstanceFromId() in WFInstanceMarkDone")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}

	req.Entity["step"] = req.Step
	req.Entity["stepfailed"] = strconv.FormatBool(req.Stepfailed)

	var DoMarkDoneParam = Markdone_t{
		InstanceID: wfinstance.ID,
		EntityID:   wfinstance.Entityid,
		Workflow:   wfinstance.Workflow,
		Loggedat:   wfinstance.Loggedat.Time,
		Entity: crux.Entity{
			Realm: realmName,
			App:   wfinstance.App,
			Slice: wfinstance.Slice,
			Class: wfinstance.Class,
			Attrs: req.Entity,
		},
	}

	ResponseData, err := DoMarkDone(c, s, qtx, DoMarkDoneParam)
	if err != nil {
		l.Debug1().LogDebug("Error while marshaling patternSchema", err)
		wscutils.SendErrorResponse(c, &wscutils.Response{Status: "error", Data: err.Error()})
		return
	}

	if err := tx.Commit(c); err != nil {
		l.LogActivity("Error while commits the transaction", err.Error())
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}
	wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: ResponseData, Messages: nil})

}
