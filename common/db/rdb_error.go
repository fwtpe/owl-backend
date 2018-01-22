package db

import (
	"github.com/fwtpe/owl-backend/common/utils"
)

// Defines the type of database error
type DbError struct {
	*utils.StackError
}

// Panic with database error if the error is vialbe
func PanicIfError(err error) {
	if !utils.IsViable(err) {
		return
	}

	// Skip this frame of callers
	panic(NewDatabaseErrorWithDepth(err, 1))
}

// Constructs an error of database
func NewDatabaseError(err error) *DbError {
	stackError, ok := err.(*utils.StackError)
	if ok {
		return &DbError{stackError}
	}

	// Skip this frame of callers
	return NewDatabaseErrorWithDepth(err, 1)
}

func NewDatabaseErrorWithDepth(err error, depth int) *DbError {
	stackError, ok := err.(*utils.StackError)
	if ok {
		return &DbError{stackError}
	}

	// Skip this frame of callers
	return &DbError{utils.BuildErrorWithCallerDepth(err, depth + 1)}
}
