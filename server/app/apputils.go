package app

import (
	"regexp"
	"slices"
	"strings"

	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/server"
)

// validate app name
func validateAppName(app string) []wscutils.ErrorMessage {

	var err []wscutils.ErrorMessage
	// Check if the app name is one-word and follows identifier syntax rules
	pattern := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	if !pattern.MatchString(app) {
		//	return fmt.Errorf("%v must be a one-word and follows identifier syntax rules", app)
		err = append(err, wscutils.BuildErrorMessage(server.MsgId_Invalid, server.ErrCode_Invalid_NAME, &APP))
	}

	if strings.Contains(app, " ") {
		err = append(err, wscutils.BuildErrorMessage(server.MsgId_Invalid, server.Errcode_Single_Name, &APP))
	}

	// Check if the app name is reserved
	App := strings.ToUpper(app)

	isReservedName := slices.Contains(RESERVED_APPNAMES, App)
	if isReservedName {
		err = append(err, wscutils.BuildErrorMessage(server.MsgId_Invalid, server.Errcode_Reserved_name, &APP))
	}
	return err

}
