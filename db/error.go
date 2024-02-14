package db

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/server"
)

// var ErrRecordNotFound = pgx.ErrNoRows

func ErrorMessage(err error) wscutils.ErrorMessage {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": //unique_violation
			return wscutils.BuildErrorMessage(server.MsgId_Invalid_Request, server.ErrCode_AlreadyExist, &pgErr.ConstraintName)
		case "23503": //foreign_key_violation
			return wscutils.BuildErrorMessage(server.MsgId_Invalid_Request, server.ErrCode_NotFound, &pgErr.ConstraintName)
		}
	} else if err.Error() == "no rows in result set" {
		feild := "slice/app/class"
		return wscutils.BuildErrorMessage(1006, server.ErrCode_InvalidRequest, &feild)
	} else {
		return wscutils.BuildErrorMessage(server.MsgId_InternalErr, server.ErrCode_DatabaseError, nil)
	}
	return wscutils.ErrorMessage{}
}

//:"ERROR: cached plan must not change result type (SQLSTATE 0A000)"}
