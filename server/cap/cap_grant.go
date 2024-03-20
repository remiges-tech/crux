package cap

// import (
// 	"github.com/gin-gonic/gin"
// 	"github.com/go-playground/validator10"
// 	"github.com/jackc/pgtype"
// 	"github.com/remiges-tech/alya/service"
// 	"github.com/remiges-tech/alya/wscutils"
// 	"github.com/remiges-tech/crux/db/sqlc-gen"
// 	"github.com/remiges-tech/crux/server"
// 	"github.com/remiges-tech/crux/types"
// )

// type CapGrantRequest struct {
// 	User string           `json:"user" validate:"required"`
// 	App  []string         `json:"app,omitempty"`
// 	Cap  []string         `json:"cap validate:required"`
// 	From pgtype.Timestamp `json:"from,omitempty"`
// 	To   pgtype.Timestamp `json:"to,omitempty"`
// }

// // This call grants a capability or multiple capabilities to a user.
// // CapGrant will be responsible for processing the /capgrant request that comes through as a POST
// func CapGrant(c *gin.Context, s *service.Service) {
// 	lh := s.LogHarbour
// 	lh.Log("Capgrant request received")

// 	var (
// 		request CapGrantRequest
// 	)

// 	isCapable, _ := server.Authz_check(types.OpReq{
// 		User:      userID,
// 		CapNeeded: CapForList,
// 	}, false)

// 	if !isCapable {
// 		lh.Info().LogActivity("unauthorized user:", userID)
// 		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_Unauthorized))
// 		return
// 	}

// 	// step 1: json request binding with a struct
// 	err := wscutils.BindJSON(c, &request)
// 	if err != nil {
// 		lh.Debug0().Error(err).Log("error while binding json request error:")
// 		return
// 	}

// 	// step 2: standard validation
// 	valError := wscutils.WscValidate(request, func(err validator.FieldError) []string { return []string{} })
// 	if len(valError) > 0 {
// 		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, valError))
// 		lh.Debug0().LogActivity("validation error:", valError)
// 		return
// 	}
// 	query, ok := s.Dependencies["queries"].(*sqlc.Queries)
// 	if !ok {
// 		lh.Log("error while getting query instance from service Dependencies")
// 		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
// 		return
// 	}
// }
