package db

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/server"
)

// The function HandleDatabaseError first checks if the error is a PostgreSQL-specific error
//
//	by attempting to cast it to a *pgconn.PgError. If successful, it examines the PostgreSQL
//	error code to determine the nature of the error.If none of the above conditions are met,
//	it constructs a generic error message indicating an internal server error related to
//
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
			return wscutils.BuildErrorMessage(server.MsgId_Invalid_Request, server.ErrCode_Invalid, &pgErr.ColumnName)
		case "0A000": //ERROR: cached plan must not change result type (SQLSTATE 0A000)

			return wscutils.BuildErrorMessage(server.MsgId_InternalErr, server.ErrCode_Internal_Retry, nil)
		case "XX000": //ERROR: cache lookup failed for type 67119 (SQLSTATE XX000)
			return wscutils.BuildErrorMessage(server.MsgId_InternalErr, server.ErrCode_Internal_Retry, nil)
		default:
			return wscutils.BuildErrorMessage(server.MsgId_InternalErr, server.ErrCode_DatabaseError, nil)
		}
	} else if err != nil {
		if err.Error() == "no rows in result set" {
			// field := "slice/app/class"
			return wscutils.BuildErrorMessage(server.MsgId_Invalid_Request, server.ErrCode_No_record_Found, nil)
		} else {
			return wscutils.BuildErrorMessage(server.MsgId_InternalErr, server.ErrCode_DatabaseError, nil)
		}
	} else {
		return wscutils.BuildErrorMessage(server.MsgId_InternalErr, server.ErrCode_DatabaseError, nil)
	}
}
