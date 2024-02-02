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

	typeBool   = "bool"
	typeInt    = "int"
	typeFloat  = "float"
	typeStr    = "str"
	typeEnum   = "enum"
	typeTS     = "ts"
	timeLayout = "2006-01-02T15:04:05Z"
)

var APP, CLASS, SLICE, ENTITY, WORKFLOW, ENTITYID string = "app", "class", "slice", "entity", "workflow", "entityid"
