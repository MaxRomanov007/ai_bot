package user

import (
	"ai-bot/internal/storage/psql"
	"errors"
	"fmt"
)

type Error struct {
	Code    string
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
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

func HandleStorageError(err error) error {
	var psqlErr *psql.Error
	if !errors.As(err, &psqlErr) {
		return fmt.Errorf("unknown error: %w", err)
	}
	switch psqlErr.Code {
	case psql.ErrNotFoundCode:
		err = ErrNotFound
	case psql.ErrAlreadyExistsCode:
		err = ErrAlreadyExists
	case psql.ErrReferenceNotExistsCode:
		err = ErrReferenceNotExists
	default:
		err = fmt.Errorf("unknown error: %w", err)
	}

	return err
}
