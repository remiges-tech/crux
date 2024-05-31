package breschema

import (
	"slices"

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

type BRESchemaStruct struct {
	Slice int32  `json:"slice" validate:"lt=15"`
	App   string `json:"app" validate:"lt=15"`
	Class string `json:"class" validate:"lt=15"`
}

var CapForList = []string{"ruleset", "schema", "root", "report"}
var schemaList []sqlc.WfSchemaListRow
var err error

// var realmName string

func BRESchemaList(c *gin.Context, s *service.Service) {
	l := s.LogHarbour
	l.Debug0().Log("starting execution of BRESchemaList()")

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

	isCapable, capList := server.Authz_check(types.OpReq{
		User:      userID,
		CapNeeded: CapForList,
	}, false)

	if !isCapable {
		l.Info().LogActivity(server.ErrCode_Unauthorized, userID)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_Unauthorized))
		return
	}

	var sh BRESchemaStruct

	err = wscutils.BindJSON(c, &sh)
	if err != nil {
		l.Debug0().Error(err).Log("error unmarshalling request payload to struct")
		return
	}

	// Validate request
	validationErrors := wscutils.WscValidate(sh, func(err validator.FieldError) []string { return []string{} })
	if len(validationErrors) > 0 {
		l.Debug0().LogDebug("standard validation errors", validationErrors)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, validationErrors))
		return
	}
	query, ok := s.Dependencies["queries"].(*sqlc.Queries)
	if !ok {
		l.Debug0().Log("error while getting query instance from service dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_Internal))
		return
	}

	if slices.Contains(capList, "root") || slices.Contains(capList, "report") {
		schemaList, err = getSchemaList(c, sh, query, realmName)
	} else if slices.Contains(capList, "ruleset") || slices.Contains(capList, "schema") {
		if sh.App != "" {
			schemaList, err = getSchemaList(c, sh, query, realmName)
		} else {
			l.Info().LogActivity(server.ErrCode_Unauthorized, userID)
			wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_Unauthorized))
			return
		}
	} else {
		l.Info().LogActivity(server.ErrCode_Unauthorized, userID)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_Unauthorized))
		return
	}
	if err != nil {
		l.Debug0().Error(err).Log("error while retrieving from db")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	} else if len(schemaList) == 0 {
		l.Debug0().Log("no record found with given details")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_NotFound, server.ErrCode_NotFound))
		return
	}
	wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: schemaList, Messages: nil})
	l.Debug0().Log("finished execution of BRESchemaList()")
}

func getSchemaList(c *gin.Context, sh BRESchemaStruct, query *sqlc.Queries, realmName string) ([]sqlc.WfSchemaListRow, error) {
	return query.WfSchemaList(c, sqlc.WfSchemaListParams{
		Relam: realmName,
		Slice: pgtype.Int4{Int32: sh.Slice, Valid: sh.Slice > 0},
		App:   pgtype.Text{String: sh.App, Valid: !server.IsStringEmpty(&sh.App)},
		Class: pgtype.Text{String: sh.Class, Valid: !server.IsStringEmpty(&sh.Class)},
		Brwf:  sqlc.BrwfEnumB,
	})
}
