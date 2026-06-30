package errors

func NewValidation(code, message string) *Error {
	return New(KindValidation, code, message)
}

func NewNotFound(code, message string) *Error {
	return New(KindNotFound, code, message)
}

func NewConflict(code, message string) *Error {
	return New(KindConflict, code, message)
}

func NewInternal(code, message string, err error) *Error {
	return Wrap(KindInternal, code, message, err)
}

func NewDependency(code, message string, err error) *Error {
	return Wrap(KindDependency, code, message, err)
}
