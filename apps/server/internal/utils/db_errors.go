package utils

import (
	"database/sql"
	"errors"

	sqlite "modernc.org/sqlite"
	sqlLib "modernc.org/sqlite/lib"
)

func IsUniqueConstraintError(err error) bool {
	var sqliteErr *sqlite.Error
	if errors.As(err, &sqliteErr) {
		return sqliteErr.Code() == sqlLib.SQLITE_CONSTRAINT_UNIQUE
	}
	return false
}

func IsNoRowsError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, sql.ErrNoRows) {
		return true
	}
	return false
}
