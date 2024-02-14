package db

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/server"
)

// The function HandleDatabaseError first checks if the error is a PostgreSQL-specific error
//	by attempting to cast it to a *pgconn.PgError. If successful, it examines the PostgreSQL
//	error code to determine the nature of the error.If none of the above conditions are met,
//	it constructs a generic error message indicating an internal server error related to

// the database.
func HandleDatabaseError(err error) wscutils.ErrorMessage {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": //unique_violation
			return wscutils.BuildErrorMessage(server.MsgId_Invalid_Request, server.ErrCode_AlreadyExist, &pgErr.ConstraintName)
		case "23503": //foreign_key_violation
			return wscutils.BuildErrorMessage(server.MsgId_Invalid_Request, server.ErrCode_NotFound, &pgErr.ConstraintName)
		case "23502": //not_null_violation
			return wscutils.BuildErrorMessage(server.MsgId_Invalid_Request, server.ErrCode_Empty, &pgErr.ConstraintName)
		case "40001": //serialization_failure
			return wscutils.BuildErrorMessage(server.MsgId_InternalErr, server.ErrCode_Internal_Retry, &pgErr.ConstraintName)
		case "40P01": //	deadlock_detected
			return wscutils.BuildErrorMessage(server.MsgId_InternalErr, server.ErrCode_Internal_Retry, &pgErr.ConstraintName)
		case "08006": //	connection_failure
			return wscutils.BuildErrorMessage(server.MsgId_InternalErr, server.ErrCode_Internal_Retry, &pgErr.ConstraintName)
		case "57P01": //	admin_shutdown
			return wscutils.BuildErrorMessage(server.MsgId_InternalErr, server.ErrCode_Internal_Retry, &pgErr.ConstraintName)
		case "55P03": //	lock_not_available  Occurs when a transaction waits too long for a lock to be released, exceeding the configured timeout period.
			return wscutils.BuildErrorMessage(server.MsgId_InternalErr, server.ErrCode_Internal_Retry, &pgErr.ConstraintName)
		}
	} else if err.Error() == "no rows in result set" {
		feild := "slice/app/class"
		return wscutils.BuildErrorMessage(1006, server.ErrCode_InvalidRequest, &feild)
	} else {
		return wscutils.BuildErrorMessage(server.MsgId_InternalErr, server.ErrCode_DatabaseError, nil)
	}
	return wscutils.ErrorMessage{}
}