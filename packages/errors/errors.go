package errors

import "fmt"

type Kind string

const (
	KindValidation     Kind = "VALIDATION_ERROR"
	KindAuthentication Kind = "AUTHENTICATION_ERROR"
	KindAuthorization  Kind = "AUTHORIZATION_ERROR"
	KindNotFound       Kind = "NOT_FOUND"
	KindConflict       Kind = "CONFLICT"
	KindInternal       Kind = "INTERNAL_ERROR"
	KindDependency     Kind = "DEPENDENCY_ERROR"
)

type Error struct {
	Kind    Kind
	Code    string
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Kind, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Kind, e.Message)
}

func (e *Error) Unwrap() error { return e.Err }

func New(kind Kind, code, message string) *Error {
	return &Error{Kind: kind, Code: code, Message: message}
}

func Wrap(kind Kind, code, message string, err error) *Error {
	return &Error{Kind: kind, Code: code, Message: message, Err: err}
}
