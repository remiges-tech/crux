package wfinstanceserv

const (
	SLICE_REQUIRED                   = "slice is required"
	ENTITY_REQUIRED                  = "Entity is required"
	ENTITYID_REQUIRED                = "Entity ID is required"
	APP_NAME_REQUIRED                = "App Name is required"
	WORKFLOW_REQUIRED                = "Workflow is required"
	CLASS_REQUIRED                   = "class is required"
	RECORD_NOT_EXIST                 = "record_does_not_exist"
	INVALID_ENTITY                   = "invalid_entity"
	INVALID_APP                      = "invalid_app"
	INVALID_CLASS                    = "invalid_class"
	INVALID_WORKFLOW_ACTIVE_STATUS   = "invalid_wf_active_status"
	INVALID_WORKFLOW_INTERNAL_STATUS = "invalid_wf_internal_status"
	INSTANCE_ALREADY_EXIST           = "instance_already_exist"
	SCHEMA_PATTERN_NOT_FOUND         = "schema_pattern_not_found"
	INVALID_PATTERN                  = "entity_does_not_match_with_pattern"
	INVALID_PROPERTY_ATTRIBUTES      = "invalid_property_attributes"
	INSERT_OPERATION_FAILED          = "insert_operation_failed"
	INVALID_DATABASE_DEPENDENCY      = "invalid_database_dependency"

	typeBool   = "bool"
	typeInt    = "int"
	typeFloat  = "float"
	typeStr    = "str"
	typeEnum   = "enum"
	typeTS     = "ts"
	timeLayout = "2006-01-02T15:04:05Z"
)

var APP, CLASS, SLICE, ENTITY, WORKFLOW, ENTITYID string = "app", "class", "slice", "entity", "workflow", "entityid"
var ACTIONSET_PROPERTIES, TASK = "actionset_properties", "task"
var DONE, NEXTSTEP = "done", "nextstep"
