package wfinstance

const (
	// error messages
	INVALID_DATABASE_DEPENDENCY = "invalid_database_dependency"

	// fields
	STEP      = "step"
	STEPFALED = "stepfailed"
	DONE      = "done"
	NEXTSTEP  = "nextstep"
	START     = "start"
	FALSE     = "false"
	TRUE      = "true"
	realm     = "Nova"
	userID    = "Admin"
)

//var ENTITYREALM = "Nova"

// feilds for error messages
var APP, CLASS, SLICE, ENTITY, WORKFLOW, ENTITYID string = "app", "class", "slice", "entity", "workflow", "entityid"
var ACTIONSET_PROPERTIES, TASK = "actionset_properties", "task"
