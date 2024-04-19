package markdone

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/crux/server"
	"github.com/remiges-tech/crux/types"
)

var (
	userID    = "1234"
	realmName = "BSE"
)

type WFInstanceMarkDoneReq struct {
	ID         int32             `json:"id" validate:"required,lt=50"`
	Entity     map[string]string `json:"entity" validate:"required,lt=50"`
	Step       string            `json:"step" validate:"required,lt=50"`
	Stepfailed bool              `json:"stepfailed,omitempty"`
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
	l.Info().Log("Authz_check completed")

	var req WFInstanceMarkDoneReq

	err := wscutils.BindJSON(c, &req)
	if err != nil {
		l.Error(err).Debug0().Log("Error Unmarshalling Query parameters to struct:")
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
		l.Debug0().Debug1().Log("Error while getting query instance from service Dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}

	connpool, ok := s.Database.(*pgxpool.Pool)
	if !ok {
		l.Debug0().Debug1().Log("Error while getting connection pool instance from service Database")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}
	tx, err := connpool.Begin(c)
	if err != nil {
		l.Error(err).Err().Log("Error while Begin tx")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}
	defer tx.Rollback(c)
	qtx := queries.WithTx(tx)

	req.Entity["step"] = req.Step
	req.Entity["stepfailed"] = strconv.FormatBool(req.Stepfailed)

	ResponseData, err := DoMarkDone(c, s, qtx, req.ID, req.Entity)
	if err != nil {
		l.Err().Error(err).Log("Error in DoMarkDone")
		wscutils.SendErrorResponse(c, &wscutils.Response{Status: "error", Data: err.Error()})
		return
	}

	if err := tx.Commit(c); err != nil {
		l.Err().Error(err).Log("Error while commits the transaction")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}
	wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: ResponseData, Messages: nil})
	l.Debug0().Log("finished execution of WFInstanceMarkDone()")
}
