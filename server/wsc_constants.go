package server

const (
	MsgId_InternalErr     = 1001
	MsgId_NoSchemaFound   = 1002
	MsgId_Invalid         = 1003 // Field <field> is invalid
	MsgId_ValTypeInvalid  = 1004
	MsgId_Empty           = 1005 // Field <field> is empty
	MsgId_Invalid_Request = 1006
	MsgId_RequiredOneOf   = 1007 // Field <field> must have either <val1> or <val2>
	MsgId_AlreadyExist    = 1008
	MsgId_NotFound        = 1009
	MsgId_Unauthorized    = 1010
	MsgId_StepNotFound    = 1011
)

const (
	ErrCode_NotExist                         = "not_exist"
	ErrCode_Internal                         = "internal_err"
	ErrCode_Internal_Retry                   = "internal_err_retry"
	ErrCode_Invalid                          = "invalid"
	ErrCode_InvalidRequest                   = "invalid_request"
	ErrCode_Empty                            = "empty"
	ErrCode_InvalidJson                      = "invalid_json"
	ErrCode_DatabaseError                    = "database_error"
	ErrCode_RequiredOne                      = "required_one_field"
	ErrCode_AlreadyExist                     = "already_exist"
	ErrCode_NotFound                         = "not_found"
	ErrCode_Unauthorized                     = "Unauthorized"
	ErrCode_Invalid_APP                      = "invalid_app"
	ErrCode_Invalid_Class                    = "invalid_class"
	ErrCode_Invalid_workflow_active_status   = "invalid_workflow_active_status"
	ErrCode_Invalid_workflow_Internal_status = "invalid_workflow_intenal_status"
	ErrCode_Invalid_Entity                   = "invalid_entity"
	ErrCode_Invalid_pattern_schema           = "invalid_pattern_schema"
	ErrCode_Invalid_property_attributes      = "invalid_property_attributes"
)
