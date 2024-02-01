package types

import "github.com/go-playground/validator/v10"

type AppConfig struct {
	DBConnURL     string `json:"db_conn_url"`
	DBHost        string `json:"db_host"`
	DBPort        int    `json:"db_port"`
	DBUser        string `json:"db_user"`
	DBPassword    string `json:"db_password"`
	DBName        string `json:"db_name"`
	DriverName    string `json:"driver_name"`
	AppServerPort string `json:"app_server_port"`
	ErrorTypeFile string `json:"error_type_file"`
}

const (
	DevEnv           Environment = "dev_env"
	ProdEnv          Environment = "prod_env"
	UATEnv           Environment = "uat_env"
	RECORD_NOT_EXIST             = "record_does_not_exist"
	OPERATION_FAILED             = "operation_failed"
)

type Environment string

func (env Environment) IsValid() bool {
	switch env {
	case DevEnv, ProdEnv, UATEnv:
		return true
	}
	return false
}

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
	default:
		vals = append(vals, "not_valid_input")
	}
	return vals
}

// func GetErrorValidationMapByAPIName(apiName string) map[string]string {
// 	var validationsMap = make(map[string]map[string]string)
// 	validationsMap["WorkflowGet"] = map[string]string{
// 		"required": "not_provided",
// 	}
// 	validationsMap["SchemaGet"] = map[string]string{
// 		"required": "not_provided",
// 	}
// 	// below is one more example ::
// 	// validationsMap["country_draft_forward"] = map[string]string{
// 	// 	"IDmin": "length must be greater than one",
// 	// }
// 	return validationsMap[apiName]
// }
