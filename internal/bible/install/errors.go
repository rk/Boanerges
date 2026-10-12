package install

import "errors"

var (
	ErrNotFound       = errors.New("not found")
	ErrConflict       = errors.New("conflict")
	ErrUnprocessable  = errors.New("unprocessable")
)

type Error struct {
	Kind    error
	Message string
}

func (e *Error) Error() string {
	if e.Message != "" {
		return e.Message
	}

	return e.Kind.Error()
}

func (e *Error) Is(target error) bool {
	return errors.Is(e.Kind, target)
}

func notFound(msg string) *Error {
	return &Error{Kind: ErrNotFound, Message: msg}
}

func conflict(msg string) *Error {
	return &Error{Kind: ErrConflict, Message: msg}
}

func unprocessable(msg string) *Error {
	return &Error{Kind: ErrUnprocessable, Message: msg}
}
