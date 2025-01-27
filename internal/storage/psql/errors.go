package psql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
)

type Error struct {
	Code    string
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

const (
	ErrNotFoundCode           = "NotFound"
	ErrAlreadyExistsCode      = "AlreadyExists"
	ErrReferenceNotExistsCode = "ReferenceNotExists"
)

var (
	ErrNotFound = &Error{
		Code:    ErrNotFoundCode,
		Message: "not found",
	}
	ErrAlreadyExists = &Error{
		Code:    ErrAlreadyExistsCode,
		Message: "already exists",
	}
	ErrReferenceNotExists = &Error{
		Code:    ErrReferenceNotExistsCode,
		Message: "reference not exists",
	}
)

func HandleDatabaseError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	var pqErr *pq.Error
	if !errors.As(err, &pqErr) {
		return fmt.Errorf("unknown error: %w", err)
	}
	switch pqErr.Code.Name() {
	case "unique_violation":
		err = ErrAlreadyExists
	case "foreign_key_violation":
		err = ErrReferenceNotExists
	default:
		err = fmt.Errorf("unknown error: %w", err)
	}

	return err
}
