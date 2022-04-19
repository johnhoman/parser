package object

import "fmt"

type ErrorType string

const (
	ErrorTypeException = "Exception"
	ErrorTypeTypeError = "TypeError"
)

type Error struct {
	Message   string
	ErrorType ErrorType
}

func (e *Error) Type() Type { return TypeError }
func (e *Error) Inspect() string { return e.Message }

var _ Object = &Error{}

func NewTypeError(message string, a ...interface{}) *Error {
	if len(a) > 0 {
		message = fmt.Sprintf(message, a...)
	}
	return &Error{Message: message, ErrorType: ErrorTypeTypeError}
}
