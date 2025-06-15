package error

import (
	"fmt"
)

type ObjectNotFound struct {
	Message string
}

type ValidationError struct {
	Message string
}

type EditError struct {
	Message string
}

type DBError struct {
	Message string
}

type ParserError struct {
	Message string
}

func NewObjectNotFoundError(err error) error {
	// Create a new object not found error
	return &ObjectNotFound{Message: err.Error()}
}

func NewValidationError(err error) error {
	// Create a new validation error
	return &ValidationError{Message: err.Error()}
}

func NewEditError(err error) error {
	// Create a new edit error
	return &EditError{Message: err.Error()}
}

func NewDBError(err error) error {
	// Create a new database error
	return &DBError{Message: err.Error()}
}

func NewParseError(err error) error {
	// Create a new parse error
	return &ParserError{Message: err.Error()}
}

func (e *ObjectNotFound) Error() string {
	// Return the error message
	return fmt.Sprintf("Object not found: %s", e.Message)
}

func (e *ValidationError) Error() string {
	// Return the error message
	return fmt.Sprintf("Validation error: %s", e.Message)
}

func (e *EditError) Error() string {
	// Return the error message
	return fmt.Sprintf("Edit error: %s", e.Message)
}

func (e *DBError) Error() string {
	// Return the error message
	return fmt.Sprintf("Database error: %s", e.Message)
}

func (e *ParserError) Error() string {
	// Return the error message
	return fmt.Sprintf("Parser error: %s", e.Message)
}
