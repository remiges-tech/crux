package realmslice

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/crux/server"
	"github.com/remiges-tech/crux/types"
)

type RealmSliceNewRequest struct {
	CopyOf int64    `json:"copyof" validate:"omitempty,gt=0"`
	Descr  string   `json:"descr" validate:"omitempty"`
	App    []string `json:"app" validate:"omitempty,lt=15"`
}

func RealmSliceNew(c *gin.Context, s *service.Service) {
	l := s.LogHarbour
	l.Debug0().Log("Starting execution of RealmSliceNew()")

	userID, err := server.ExtractUserNameFromJwt(c)
	if err != nil {
		l.Info().Log("unable to extract userID from token")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Missing, server.ERRCode_Token_Data_Missing))
		return
	}

	realmName, err := server.ExtractRealmFromJwt(c)
	if err != nil {
		l.Info().Log("unable to extract realm from token")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Missing, server.ERRCode_Token_Data_Missing))
		return
	}

	isCapable, _ := server.Authz_check(types.OpReq{
		User:      userID,
		CapNeeded: capForNew,
	}, false)

	if !isCapable {
		l.Info().LogActivity("Unauthorized user:", userID)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_Unauthorized))
		return
	}

	var req RealmSliceNewRequest

	err = wscutils.BindJSON(c, &req)
	if err != nil {
		l.Error(err).Log("Error Unmarshalling Query parameters to struct:")
		return
	}

	// Validate request
	validationErrors := wscutils.WscValidate(req, func(err validator.FieldError) []string { return []string{} })
	if len(validationErrors) > 0 {
		l.Debug0().LogDebug("standard validation errors", validationErrors)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, validationErrors))
		return
	}

	query, ok := s.Dependencies["queries"].(*sqlc.Queries)
	if !ok {
		l.Info().Log("Error while getting query instance from service Dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_Internal))
		return
	}

	connpool, ok := s.Database.(*pgxpool.Pool)
	if !ok {
		l.Info().Log("Error while getting query instance from service Database")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}

	if req.CopyOf == 0 && req.App == nil {
		newSliceID, err := query.InsertNewRecordInRealmSlice(c, sqlc.InsertNewRecordInRealmSliceParams{
			Realm:     realmName,
			Descr:     req.Descr,
			Createdby: userID,
		})
		if err != nil {
			l.Info().Error(err).Log("Error while creating new realmslice")
			errmsg := db.HandleDatabaseError(err)
			wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
			return
		}
		wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: newSliceID, Messages: nil})
		l.Debug0().Log("Finished execution of RealmSliceNew()")
		return
	}

	tx, err := connpool.Begin(c)
	if err != nil {
		l.Info().Log("Error while Begin transaction")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}
	defer tx.Rollback(c)
	qtx := query.WithTx(tx)

	newSliceID, err := qtx.CloneRecordInRealmSliceBySliceID(c, sqlc.CloneRecordInRealmSliceBySliceIDParams{
		ID:        int32(req.CopyOf),
		Realm:     realmName,
		Createdby: userID,
		Descr:     pgtype.Text{String: req.Descr, Valid: true},
	})
	if err != nil {
		tx.Rollback(c)
		l.Info().Error(err).Log("Error while creating clone of realmslice")
		if err.Error() == "no rows in result set" {
			feild := "CopyOf"
			wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{wscutils.BuildErrorMessage(1006, server.ErrCode_NotExist, &feild)}))
			return
		}
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}

	tag, err := qtx.CloneRecordInConfigBySliceID(c, sqlc.CloneRecordInConfigBySliceIDParams{
		Slice:   int32(req.CopyOf),
		Slice_2: newSliceID,
		Setby:   userID,
	})
	if err != nil {
		tx.Rollback(c)
		l.Info().Error(err).Log("Error while creating clone of config")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}
	if strings.Contains(tag.String(), "0") {
		l.Log("no record found in config to clone in config")
	}

	tag, err = qtx.CloneRecordInSchemaBySliceID(c, sqlc.CloneRecordInSchemaBySliceIDParams{
		Slice:     int32(req.CopyOf),
		Slice_2:   newSliceID,
		Createdby: userID,
		App:       req.App,
	})
	if err != nil {
		tx.Rollback(c)
		l.Info().Error(err).Log("Error while creating clone of schema")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}
	if strings.Contains(tag.String(), "0") {
		l.Log("no record found to clone in schema")

	}

	tag, err = qtx.CloneRecordInRulesetBySliceID(c, sqlc.CloneRecordInRulesetBySliceIDParams{
		Slice:     int32(req.CopyOf),
		Slice_2:   newSliceID,
		Createdby: userID,
		App:       req.App,
	})

	if err != nil {
		tx.Rollback(c)
		l.Info().Error(err).Log("Error while creating clone of ruleset")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}
	if strings.Contains(tag.String(), "0") {
		l.Log("no record found to clone in ruleset")
	}

	if err := tx.Commit(c); err != nil {
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}

	wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: newSliceID, Messages: nil})
	l.Debug0().Log("Finished execution of RealmSliceNew()")
}
