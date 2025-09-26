package as

import (
	"errors"
)

// notAnIntegerOrDurationError is an error type that wraps ErrNotAnIntegerOrDuration
// with the underlying errors ErrNotAnInteger and ErrNotADuration.  This allows
// callers to check for any of these errors using errors.Is while providing
// an error that concisely describes the failure.
type notAnIntegerOrDurationError struct{}

// Error implements the error interface for notAnIntegerOrDurationError by
// delegating to the wrapped ErrNotAnIntegerOrDuration.
func (ec notAnIntegerOrDurationError) Error() string {
	return ErrNotAnIntegerOrDuration.Error()
}

// Unwrap implements the Wrapper interface by returning the wrapped errors
// ErrNotAnIntegerOrDuration, ErrNotAnInteger and ErrNotADuration, allowing a
// caller to test for any of these errors using [errors.Is].
func (ec notAnIntegerOrDurationError) Unwrap() []error {
	return []error{ErrNotAnIntegerOrDuration, ErrNotAnInteger, ErrNotADuration}
}

var (
	ErrFilenameIsEmpty        = errors.New("filename is empty")
	ErrFilenameIsDirectory    = errors.New("filename is a directory")
	ErrNotAnInteger           = errors.New("not an integer")
	ErrNotADuration           = errors.New("not a duration")
	ErrNotAnIntegerOrDuration = errors.New("not an integer or duration")
	ErrNotAnAbsoluteURL       = errors.New("not an absolute URL")
)
