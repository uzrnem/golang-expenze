package errors

import "fmt"

type Error struct {
	Type    Type
	Message string
}

type Errors []*Error

func New(t Type, message string) error {
	return &Error{
		Type:    t,
		Message: message,
	}
}

func Newf(t Type, format string, args ...interface{}) error {
	return &Error{
		Type:    t,
		Message: fmt.Sprintf(format, args...),
	}
}

func Encode(err error, t Type, message string) error {
	if err == nil {
		return nil
	}

	if e, ok := err.(*Error); ok {
		t = e.Type
	}
	return Newf(t, "%s | %s", message, err)
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) GetType() Type {
	return e.Type
}

func (e *Error) GetStatusCode() int {
	return e.Type.GetStatusCode()
}
