package server

const (
	MsgId_InternalErr        = 1001
	MsgId_NoSchemaFound      = 1002
	MsgId_Invalid            = 1003
	MsgId_ValTypeInvalid     = 1004
	MsgId_Empty              = 1005
	MsgId_Invalid_Request    = 1006
	MsgId_RequiredAtLeastOne = 1007
	MsgId_AlreadyExist       = 1008
	MsgId_NotFound           = 1009
	MsgId_Unauthorized       = 1010
	MsgId_StepNotFound       = 1011
)

const (
	ErrCode_NotExist       = "not_exist"
	ErrCode_Internal       = "internal_err"
	ErrCode_Internal_Retry = "internal_err_retry"
	ErrCode_Invalid        = "invalid"
	ErrCode_InvalidRequest = "invalid_request"
	ErrCode_Empty          = "empty"
	ErrCode_InvalidJson    = "invalid_json"
	ErrCode_DatabaseError  = "database_error"
	ErrCode_RequiredOne    = "required_one_field"
	ErrCode_AlreadyExist   = "already_exist"
	ErrCode_NotFound       = "not_found"
	ErrCode_Unauthorized   = "Unauthorized"
)
