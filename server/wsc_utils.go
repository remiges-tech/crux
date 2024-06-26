package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/remiges-tech/alya/router"
	"github.com/remiges-tech/alya/wscutils"

	crux "github.com/remiges-tech/crux/matching-engine"
	"github.com/remiges-tech/crux/types"
)

// CommonValidation is a generic function which setup standard validation utilizing
// validator package and Maps the errorVals based on the map parameter and
// return []errorVals
func CommonValidation(err validator.FieldError) []string {
	var vals []string
	switch err.Tag() {
	case "required":
		vals = append(vals, "not_provided")
	case "alpha":
		vals = append(vals, "only_alphabets_are_allowed")
	case "gt":
		vals = append(vals, "must_be_greater_than_zero")
	default:
		vals = append(vals, "not_valid_input")
	}
	return vals
}

func MarshalJson(data any) []byte {
	jsonData, err := json.Marshal(&data)
	if err != nil {
		log.Fatal("error marshaling")
	}
	return jsonData
}

func ReadJsonFromFile(filepath string) ([]byte, error) {
	// var err error
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("testFile path is not exist")
	}
	defer file.Close()
	jsonData, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

// to check given string is nil or not
func IsStringEmpty(s *string) bool {
	return s == nil || strings.TrimSpace(*s) == ""
}

// to check if the user has "ruleset" rights for the given app
func HasRulesetRights(app string) bool {
	userRights := GetWorkflowsByRulesetRights()
	return slices.Contains(userRights, app)
}

// to get workflows for all apps for which the user has "ruleset" rights
func GetWorkflowsByRulesetRights() []string {
	return []string{"retailBANK", "nedbank", "amazon", "myntra", "uccapp"}
}

func Authz_check(op types.OpReq, trace bool) (bool, []string) {
	caplist := op.CapNeeded
	return true, caplist
}

// To check whether requested user is exist in idshield and it is belong to valid realm
func IsValidUser(string, string) (bool, error) {
	return true, nil
}

// ExtractClaimFromJwt: this will extract the provided singleClaimName as key from the jwt token and return its value as a string
func ExtractClaimFromJwt(c *gin.Context, singleClaimName string) (string, error) {
	tokenString, err := router.ExtractToken(c.GetHeader("Authorization"))
	if err != nil {
		return "", fmt.Errorf("invalid token payload")
	}
	var name string
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return "", fmt.Errorf("invalid token payload")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		name = fmt.Sprint(claims[singleClaimName])
	}

	if name == "" {
		return "", fmt.Errorf("invalid token payload")
	}

	return name, nil
}

func ExtractRealmFromJwt(c *gin.Context) (string, error) {
	str, err := ExtractClaimFromJwt(c, "iss")
	if err != nil {
		return "", err
	}
	parts := strings.Split(str, "/realms/")
	realm := parts[1]
	return realm, nil
}

func ExtractUserNameFromJwt(c *gin.Context) (string, error) {
	return ExtractClaimFromJwt(c, "preferred_username")
}

func HandleCruxError(errs []error) []wscutils.ErrorMessage {
	var validationErrors []wscutils.ErrorMessage
	for _, err := range errs {
		var cruxErr crux.CruxError
		fmt.Println("validationErrors", err)

		if errors.As(err, &cruxErr) {
			switch cruxErr.Keyword {
			case "Empty":
				vErr := wscutils.BuildErrorMessage(MsgId_Empty, ErrCode_Empty, &cruxErr.FieldName, cruxErr.Vals)
				validationErrors = append(validationErrors, vErr)
			case "Invalid":
				vErr := wscutils.BuildErrorMessage(MsgId_Invalid, ErrCode_Invalid, &cruxErr.FieldName, cruxErr.Vals)
				validationErrors = append(validationErrors, vErr)
			case "NotAllowed":
				vErr := wscutils.BuildErrorMessage(MsgId__NotAllowed, ErrCode_NotAllowed, &cruxErr.FieldName, cruxErr.Vals)
				validationErrors = append(validationErrors, vErr)
			case "Required":
				vErr := wscutils.BuildErrorMessage(MsgId_Invalid_Request, ErrCode_RequiredOneOf, &cruxErr.FieldName, cruxErr.Vals)
				validationErrors = append(validationErrors, vErr)
			case "NotExist":
				vErr := wscutils.BuildErrorMessage(MsgId_NotFound, ErrCode_NotFound, &cruxErr.FieldName, cruxErr.Vals)
				validationErrors = append(validationErrors, vErr)
			case "NotMatch":
				vErr := wscutils.BuildErrorMessage(MsgID_NotMatched, ErrCode_Not_Match, &cruxErr.FieldName, cruxErr.Vals)
				validationErrors = append(validationErrors, vErr)
			default:
				vErr := wscutils.BuildErrorMessage(MsgId_Invalid_Request, ErrCode_InvalidRequest, &cruxErr.FieldName, cruxErr.Vals)
				validationErrors = append(validationErrors, vErr)
			}
		}

	}
	return validationErrors
}

func IsZeroOfUnderlyingType(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}
