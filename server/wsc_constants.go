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
	MsgId_NotFound        = 1009 //Field <field> is not found
	MsgId_Unauthorized    = 1010
	MsgId_StepNotFound    = 1011
)

const (
	ErrCode_NotExist                                              = "not_exist"
	ErrCode_Internal                                              = "internal_err"
	ErrCode_Internal_Retry                                        = "internal_err_retry"
	ErrCode_Invalid                                               = "invalid"
	ErrCode_InvalidRequest                                        = "invalid_request"
	ErrCode_Empty                                                 = "empty"
	ErrCode_InvalidJson                                           = "invalid_json"
	ErrCode_DatabaseError                                         = "database_error"
	ErrCode_RequiredOne                                           = "required_one_field"
	ErrCode_AlreadyExist                                          = "already_exist"
	ErrCode_NotFound                                              = "not_found"
	ErrCode_Unauthorized                                          = "Unauthorized"
	ErrCode_Invalid_APP                                           = "invalid_app"
	ErrCode_Invalid_Class                                         = "invalid_class"
	ErrCode_Invalid_workflow_active_status                        = "invalid_workflow_active_status"
	ErrCode_Invalid_workflow_Internal_status                      = "invalid_workflow_internal_status"
	ErrCode_Invalid_Entity                                        = "invalid_entity"
	ErrCode_Invalid_pattern_schema                                = "invalid_pattern_schema"
	ErrCode_Invalid_action_schema                                 = "invalid_action_schema"
	ErrCode_Invalid_property_attributes                           = "invalid_property_attributes"
	ErrCode_Required                                              = "required"
	ErrCode_RequiredOneOf                                         = "required_one_of"
	ErrCode_Required_Exactly_Two_Properties                       = "required_exactly_two_properties"
	ErrCode_Does_Not_Contain_Both_Properties_Nextstep_And_Done    = "does_not_contain_both_properties_nextstep_and_done"
	ErrCode_ActionSchema_Task_Not_Same_As_PatternSchema_Step_Attr = "actionschema_tasks_are_not_same_as_'step'_in_patternschema"
	ErrCode_Invalid_NAME                                          = "invalid_name"
	Errcode_Single_Name                                           = "name_must_be_single_word"
	Errcode_Reserved_name                                         = "reserved_name"
	ErrCode_Name_Already_Exist                                    = "name_already exist"
	ERRCode_Name_Not_Exist                                        = "name_not_exist"
)
