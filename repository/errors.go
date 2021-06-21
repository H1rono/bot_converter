package repository

import (
	"errors"
	"gorm.io/gorm"
)

var (
	// ErrNilID is a general error that indicates the given UUID was uuid.Nil.
	ErrNilID = errors.New("nil id")
	// ErrNotFound is a general error that indicates specified record was not found.
	ErrNotFound = errors.New("not found")
	// ErrAlreadyExists is a general error that indicates the given record already exists.
	ErrAlreadyExists = errors.New("already exists")
	// ErrForbidden is a general error that indicates the operation is forbidden.
	ErrForbidden = errors.New("forbidden")
)

// convertError converts gorm error to repository error.
func convertError(err error) error {
	switch err {
	case gorm.ErrRecordNotFound:
		return ErrNotFound
	default:
		return err
	}
}
