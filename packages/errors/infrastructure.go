package errors

type Severity string

const (
	SeverityRetryable    Severity = "RETRYABLE"
	SeverityNonRetryable Severity = "NON_RETRYABLE"
)

type InfrastructureError struct {
	Code      string
	Message   string
	Err       error
	Retryable bool
	Severity  Severity
}

func (e *InfrastructureError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *InfrastructureError) Unwrap() error { return e.Err }

func NewRetryable(code, message string, err error) *InfrastructureError {
	return &InfrastructureError{
		Code: code, Message: message, Err: err,
		Retryable: true, Severity: SeverityRetryable,
	}
}

func NewNonRetryable(code, message string, err error) *InfrastructureError {
	return &InfrastructureError{
		Code: code, Message: message, Err: err,
		Retryable: false, Severity: SeverityNonRetryable,
	}
}
