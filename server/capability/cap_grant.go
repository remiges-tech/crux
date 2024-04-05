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
		request    CapGrantRequest
		realmcaps  []string
		appcaps    []string
		appMap     = make(map[string][]string, 0)
		realmCapDb []string
		appCapDb   []sqlc.GetUserCapsAndAppsByRealmRow
	)
	// userID := "Admin"
	// realmName := "BSE"
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
	err = wscutils.BindJSON(c, &request)
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
	IsValidUser, err := server.IsValidUser(request.User, realmName)
	if err != nil {
		lh.Error(err).Log("error while verifying  whether user already exist and belong to realm in capgrant table")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Invalid, server.ErrCode_Invalid_User))
		return

	}

	// if user is invalid then revoke it's all existing capabilities
	if !IsValidUser {
		err = query.UpdateCapGranForUser(c, sqlc.UpdateCapGranForUserParams{
			Userid: request.User,
			Realm:  realmName,
		})
		if err != nil {
			lh.Error(err).Log("error while deleting caps for Invalid user")
			wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Invalid, server.ErrCode_Invalid_User))
			return
		}
	}

	// validating app
	if *request.App != nil {
		appNames, err := query.GetAppNames(c, realmName)
		if err != nil {
			lh.Error(err).Log("error while validating app")
			errmsg := db.HandleDatabaseError(err)
			wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
			return
		}

		for _, app := range *request.App {
			applc := strings.ToLower(app)
			if !slices.Contains(appNames, applc) {
				fieldName := "app"
				lh.Log("error while verifying whether applist contains valid apps")
				errmsg := wscutils.BuildErrorMessage(server.MsgId_Missing, server.ErrCode_App_Does_Not_Exist, &fieldName, applc)
				wscutils.SendErrorResponse(c, &wscutils.Response{Status: wscutils.ErrorStatus, Data: nil, Messages: []wscutils.ErrorMessage{errmsg}})
				return
			}

		}

	}

	// seperating capabilities
	for _, cap := range request.Cap {
		caplc := strings.ToLower(cap)
		if slices.Contains(CAPLIST_REALMLEVEL, caplc) {
			realmcaps = append(realmcaps, caplc)

		} else if slices.Contains(CAPLIST_APPLEVEL, caplc) {
			appcaps = append(appcaps, caplc)

		} else {
			fieldName := "cap"
			lh.Log("error while verifying whether caplist contains valid capability")
			errmsg := wscutils.BuildErrorMessage(server.MsgId_Missing, server.ErrCode_Capability_Does_Not_Exist, &fieldName, caplc)
			wscutils.SendErrorResponse(c, &wscutils.Response{Status: wscutils.ErrorStatus, Data: nil, Messages: []wscutils.ErrorMessage{errmsg}})
			return
		}
	}

	// adding appcaps agasinst each app in appMap
	if *request.App != nil {
		for _, app := range *request.App {
			for _, cap := range appcaps {
				appMap[app] = append(appMap[app], cap)
				lh.Debug0().LogActivity("appMap", appMap)
			}
		}
	}

	// app capabilities
	if appcaps != nil {

		// getting app and caps from db
		appCapDb, err = query.GetUserCapsAndAppsByRealm(c, sqlc.GetUserCapsAndAppsByRealmParams{
			Userid: request.User,
			Realm:  realmName,
			App:    *request.App,
		})
		if err != nil {
			errmsg := wscutils.BuildErrorMessage(server.MsgId_Missing, server.ErrCode_Capability_Does_Not_Exist, nil, err.Error())
			wscutils.SendErrorResponse(c, &wscutils.Response{Status: wscutils.ErrorStatus, Data: nil, Messages: []wscutils.ErrorMessage{errmsg}})
			return
		}

		// adding elements in appMap
		if len(appCapDb) > 0 {
			for _, app := range appCapDb {
				if slices.Contains(appMap[app.App.String], app.Cap) {
					i := slices.Index(appMap[app.App.String], app.Cap)
					appMap[app.App.String][i] = appMap[app.App.String][len(appMap[app.App.String])-1] // Copy last element to index i.
					appMap[app.App.String][len(appMap[app.App.String])-1] = ""                        // Erase last element (write zero value).
					appMap[app.App.String] = appMap[app.App.String][:len(appMap[app.App.String])-1]
				}
			}
		}
	}

	//granting App Capability

	if len(appMap) > 0 {

		for k, v := range appMap {
			for _, cap := range v {

				err = query.GrantAppCapability(c, sqlc.GrantAppCapabilityParams{
					Realm:  realmName,
					Userid: request.User,
					App:    pgtype.Text{String: k, Valid: true},
					Cap:    cap,
					From:   pgtype.Timestamp{Time: *request.From, Valid: request.From != nil},
					To:     pgtype.Timestamp{Time: *request.To, Valid: request.To != nil},
					Setby:  userID,
				})

			}
		}

	}

	// getting realmcaps from db

	if realmcaps != nil {

		realmCapDb, err = query.GetUserCapsByRealm(c, sqlc.GetUserCapsByRealmParams{
			Userid: request.User,
			Realm:  realmName,
		})
		if err != nil {
			errmsg := wscutils.BuildErrorMessage(server.MsgId_Missing, server.ErrCode_Capability_Does_Not_Exist, nil, err.Error())
			wscutils.SendErrorResponse(c, &wscutils.Response{Status: wscutils.ErrorStatus, Data: nil, Messages: []wscutils.ErrorMessage{errmsg}})
			return
		}

		// extracting realm caps which are already present in db
		if len(realmCapDb) > 0 {
			for _, v := range realmCapDb {
				for k, w := range realmcaps {
					if v == w {
						realmcaps = RemoveIndex(realmcaps, k)
						break
					}
				}
			}
		}

		// granting realmcapability
		err = query.GrantRealmCapability(c, sqlc.GrantRealmCapabilityParams{
			Realm:  realmName,
			Userid: request.User,
			Cap:    realmcaps,
			From:   pgtype.Timestamp{Time: *request.From, Valid: request.From != nil},
			To:     pgtype.Timestamp{Time: *request.To, Valid: request.To != nil},
			Setby:  userID,
		})

	}
	if err != nil {
		lh.Error(err).Log("error while granting Capabilities")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}

	lh.Log("Finished execution of capGrant")
	wscutils.SendSuccessResponse(c, wscutils.NewSuccessResponse(nil))

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

// This function removes particular value at particular index
func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}
