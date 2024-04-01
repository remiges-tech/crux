package capability

import (
	"fmt"
	"slices"
	"strings"
	"time"

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

type CapGrantRequest struct {
	User string     `json:"user" validate:"required"`
	App  *[]string  `json:"app" validate:"omitempty,gt=0"`
	Cap  []string   `json:"cap" validate:"required,gt=0"`
	From *time.Time `json:"from,omitempty"`
	To   *time.Time `json:"to,omitempty"`
}

// This call grants a capability or multiple capabilities to a user.
// CapGrant will be responsible for processing the /capgrant request that comes through as a POST
func CapGrant(c *gin.Context, s *service.Service) {
	lh := s.LogHarbour
	lh.Log("Capgrant request received")

	var (
		request   CapGrantRequest
		realmcaps []string
		appcaps   []string
	)
	userID := "Admin"
	realmName := "BSE"
	// userID, err := server.ExtractUserNameFromJwt(c)
	// if err != nil {
	// 	lh.Info().Log("unable to extract userID from token")
	// 	wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Missing, server.ErrCode_Token_Data_Missing))
	// 	return
	// }

	// realmName, err := server.ExtractRealmFromJwt(c)
	// if err != nil {
	// 	lh.Info().Log("unable to extract realm from token")
	// 	wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Missing, server.ErrCode_Token_Data_Missing))
	// 	return
	// }

	capNeeded := []string{"auth"}
	isCapable, capList := server.Authz_check(types.OpReq{
		User:      userID,
		CapNeeded: capNeeded,
	}, false)

	authRights := slices.Contains(capList, "auth")
	if !isCapable && !authRights {
		lh.Info().LogActivity("unauthorized user:", userID)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_Unauthorized))
		return
	}

	// json request binding with a struct
	err := wscutils.BindJSON(c, &request)
	if err != nil {
		lh.Debug0().Error(err).Log("error while binding json request error:")
		return
	}

	// standard validation
	valError := wscutils.WscValidate(request, func(err validator.FieldError) []string { return []string{} })
	if len(valError) > 0 {
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, valError))
		lh.Debug0().LogActivity("validation error:", valError)
		return
	}

	// verify timestamp
	err = validateTimestamp(request.From, request.To)
	if err != nil {
		lh.Log("error while validating timestamp")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Invalid, server.ErrCode_Invalid_Timestamp))
		return
	}

	query, ok := s.Dependencies["queries"].(*sqlc.Queries)
	if !ok {
		lh.Log("error while getting query instance from service Dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}

	// validate whether user contains valid realm
	realmArr, err := query.GetUserRealm(c, request.User)
	if err != nil {
		lh.Error(err).Log("error while geting count to verify whether user already exist and belong to realm in capgrant table")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Invalid, server.ErrCode_Invalid_User))
		return

	}

	if !slices.Contains(realmArr, realmName) {
		lh.Log("error while validating  whether userID belongs to valid realm")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Invalid, server.ErrCode_Invalid_User))
		return
		// err := query.UpdateCapGranForUser(c, userID)
		// if err != nil {
		// 	lh.Error(err).Log("error while deleting caps for Invalid user")
		// 	wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Invalid, server.ErrCode_Invalid_User))
		// 	return
		// }
	}

	// caplist
	for _, cap := range request.Cap {
		caplc := strings.ToLower(cap)
		if slices.Contains(CAPLIST_REALMLEVEL, caplc) {
			realmcaps = append(realmcaps, caplc)
		}
		if slices.Contains(CAPLIST_APPLEVEL, caplc) {
			appcaps = append(appcaps, caplc)
		}
	}

	// granting Realm Capability
	if realmcaps != nil {
		err = query.GrantRealmCapability(c, sqlc.GrantRealmCapabilityParams{
			Realm:  realmName,
			Userid: request.User,
			Cap:    realmcaps,
			From:   pgtype.Timestamp{Time: *request.From, Valid: request.From != nil},
			To:     pgtype.Timestamp{Time: *request.To, Valid: request.To != nil},
			Setby:  userID,
		})
		if err != nil {
			lh.Error(err).Log("error while granting realmCapability")
			errmsg := db.HandleDatabaseError(err)
			wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
			return
		}
	}
	if appcaps != nil {
		//granting App Capability
		err = query.GrantAppCapability(c, sqlc.GrantAppCapabilityParams{
			Realm:  realmName,
			Userid: request.User,
			App:    *request.App,
			Cap:    appcaps,
			From:   pgtype.Timestamp{Time: *request.From, Valid: request.From != nil},
			To:     pgtype.Timestamp{Time: *request.To, Valid: request.To != nil},
			Setby:  userID,
		})
		if err != nil {
			lh.Error(err).Log("error while granting appCapability")
			errmsg := db.HandleDatabaseError(err)
			wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
			return
		}
	}

}

// validating timestamp
func validateTimestamp(fromTS, toTS *time.Time) error {
	// Get current time in UTC
	currentTime := time.Now().UTC()

	// Check if both timestamps are provided
	if fromTS != nil && toTS != nil {

		// Check if toTS is after fromTS
		if !toTS.After(*fromTS) {
			return fmt.Errorf("toTS must be after fromTS")
		}
		// Check if fromTS is  in the future
		if fromTS.Before(currentTime) {
			return fmt.Errorf("fromTS must be in the future")
		}
		// Check if toTS is in the future
		if toTS.Before(currentTime) {
			return fmt.Errorf("toTS must be in the future")
		}
	}

	// check if only fromTs provided
	if fromTS != nil && toTS == nil {
		if fromTS.Before(currentTime) {
			return fmt.Errorf("fromTS must be in the future")
		}
	}

	// check if only toTs provided
	if fromTS == nil && toTS != nil {
		if toTS.Before(currentTime) {
			return fmt.Errorf("fromTS must be in the future")
		}
	}
	// if both timestmps are not provided
	if fromTS == nil && toTS == nil {
		return nil
	}

	return nil
}
