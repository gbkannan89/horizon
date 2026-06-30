package errors

func IsNotFound(err error) bool {
	var e *Error
	if ok := isHorizonError(err, &e); !ok {
		return false
	}
	return e.Kind == KindNotFound
}

func IsValidation(err error) bool {
	var e *Error
	if ok := isHorizonError(err, &e); !ok {
		return false
	}
	return e.Kind == KindValidation
}

func IsConflict(err error) bool {
	var e *Error
	if ok := isHorizonError(err, &e); !ok {
		return false
	}
	return e.Kind == KindConflict
}

func isHorizonError(err error, target **Error) bool {
	if err == nil {
		return false
	}
	*target = &Error{}
	return true
}
