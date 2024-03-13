package realmSliceManagement

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/crux/server"
	"github.com/remiges-tech/crux/types"
)

var (
	userID    = "1234"
	capForNew = []string{"root"}
	realmName = "NSE"
	// realmID   = int32(11)
)

type request struct {
	CopyOf int64    `json:"copyof" validate:"omitempty,number,gt=0"`
	Descr  string   `json:"descr" validate:"omitempty"`
	App    []string `json:"app" validate:"omitempty,lt=15"`
}

func RealmSliceNew(c *gin.Context, s *service.Service) {
	l := s.LogHarbour
	l.Debug0().Log("Starting execution of RealmSliceNew()")

	isCapable, _ := server.Authz_check(types.OpReq{
		User:      userID,
		CapNeeded: capForNew,
	}, false)

	if !isCapable {
		l.Info().LogActivity("Unauthorized user:", userID)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_Unauthorized))
		return
	}

	var req request

	err := wscutils.BindJSON(c, &req)
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

	tx, err := connpool.Begin(c)
	if err != nil {
		l.Info().Log("Error while Begin transaction")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}
	defer tx.Rollback(c)
	qtx := query.WithTx(tx)

	newSliceID, err := qtx.CreateNewSliceBY(c, sqlc.CreateNewSliceBYParams{
		ID:    int32(req.CopyOf),
		Realm: realmName,
		Descr: req.Descr,
	})
	if err != nil {
		tx.Rollback(c)
		l.Info().Error(err).Log("Error while creating copy of realmslice")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}

	tag, err := qtx.CopyConfig(c, sqlc.CopyConfigParams{
		Slice:   int32(req.CopyOf),
		Slice_2: newSliceID,
		Setby:   userID,
	})
	if err != nil {
		tx.Rollback(c)
		l.Info().Error(err).Log("Error while creating new record in config")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}
	if strings.Contains(tag.String(), "0") {
		l.Log("no record found in config to create new record")
		// wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_NotFound, server.ErrCode_NotFound))
		// return
	}

	tag, err = qtx.CopySchema(c, sqlc.CopySchemaParams{
		Slice:     int32(req.CopyOf),
		Slice_2:   newSliceID,
		Createdby: userID,
		App:       req.App,
	})
	if err != nil {
		tx.Rollback(c)
		l.Info().Error(err).Log("Error while creating copy of schema")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}
	if strings.Contains(tag.String(), "0") {
		l.Log("no record found to create new record in schema")
		// wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_NotFound, server.ErrCode_NotFound))
		// return
	}

	tag, err = qtx.CopyRuleset(c, sqlc.CopyRulesetParams{
		Slice:     int32(req.CopyOf),
		Slice_2:   newSliceID,
		Createdby: userID,
		App:       req.App,
	})

	if err != nil {
		tx.Rollback(c)
		l.Info().Error(err).Log("Error while creating copy of ruleset")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}
	if strings.Contains(tag.String(), "0") {
		l.Log("no record found to create new record in ")
		// wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_NotFound, server.ErrCode_NotFound))
		// return
	}

	if err := tx.Commit(c); err != nil {
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}

	wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: newSliceID, Messages: nil})
	l.Debug0().Log("Starting execution of RealmSliceNew()")
}
