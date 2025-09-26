package env

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidEntry is returned from [Load], [LoadFile], or [LoadFromReader] when
	// a line is encountered in a specified file or [io.Reader] that is not valid.
	ErrInvalidEntry = errors.New("invalid entry")

	// ErrNotSet is returned from [Parse] when a variable is not set
	// and no default value is specified.
	ErrNotSet = errors.New("variable not set")

	// ErrSetFailed is returned from [Load], [LoadFile], or [LoadFromReader] when
	// setting an environment variable fails.
	ErrSetFailed = errors.New("failed to set environment variable")

	// ErrTargetIsNil is returned from [ParseInto] when the target pointer is nil.
	ErrTargetIsNil = errors.New("target is nil")
)

// FileError is an error that wraps an error occurring while
// accessing a file.  It includes the name of the file being
// accessed at the time of the error.
type FileError struct {
	Filename string
	err      error
}

// NewFileError creates a new [FileError] instance associating
// the specified filename with the provided error.
func NewFileError(filename string, err error) FileError {
	return FileError{
		Filename: filename,
		err:      err,
	}
}

// Error implements the error interface.
func (err FileError) Error() string {
	switch {
	case err.err != nil && err.Filename != "":
		return fmt.Sprintf("env.FileError: %s: %v", err.Filename, err.err)
	case err.err != nil:
		return fmt.Sprintf("env.FileError: %v", err.err)
	case err.Filename != "":
		return "env.FileError: " + err.Filename
	default:
		return "env.FileError"
	}
}

// Is reports whether the target error is a match for the receiver.
// The target is considered a match if:
//
//   - it has the same filename (or no filename); and
//   - the same error (or no error)
func (err FileError) Is(target error) bool {
	isMatch := func(target *FileError) bool {
		return (target != nil) &&
			(target.Filename == "" || target.Filename == err.Filename) &&
			(target.err == nil || errors.Is(err.err, target.err))
	}

	switch t := target.(type) {
	case *FileError:
		return isMatch(t)
	case FileError:
		return isMatch(&t)
	}

	return false
}

// Unwrap returns the error that caused the FileError.
func (err FileError) Unwrap() error {
	return err.err
}

// ParseError is an error that wraps an error occurring while
// parsing an environment variable.  It includes the name of the
// variable being parsed at the time of the error.
type ParseError struct {
	VariableName string
	Err          error
}

// Error returns a string representation of the error in the form:
//
//	env.ParseError: <variable name>: <error>
//
// If the VariableName field is empty:
//
//	env.ParseError: <error>
//
// If the Err field is nil:
//
//	env.ParseError: <variable name>
//
// If both fields are empty:
//
//	env.ParseError
func (e ParseError) Error() string {
	type when struct{ hasName, hasError bool }
	fn := map[when]func() string{
		{false, false}: func() string { return "env.ParseError" },
		{false, true}:  func() string { return "env.ParseError: " + e.Err.Error() },
		{true, false}:  func() string { return fmt.Sprintf("env.ParseError: %v", e.VariableName) },
		{true, true}:   func() string { return fmt.Sprintf("env.ParseError: %v: %v", e.VariableName, e.Err) },
	}
	return fn[when{e.VariableName != "", e.Err != nil}]()
}

// Is reports whether the target error is a match for the receiver.
// To be a match, the target must:
//
//   - be a ParseError
//   - the target VariableName must match the receiver's VariableName,
//     or be empty
//   - the target Err field must satisfy errors.Is with respect to the
//     receiver Err, or be nil
func (e ParseError) Is(target error) bool {
	isMatch := func(target *ParseError) bool {
		return (target != nil) &&
			(target.VariableName == "" || e.VariableName == target.VariableName) &&
			(target.Err == nil || errors.Is(e.Err, target.Err))
	}

	switch t := target.(type) {
	case *ParseError:
		return isMatch(t)
	case ParseError:
		return isMatch(&t)
	default:
		return false
	}
}

// Unwrap returns the error that caused the env.Parse.
func (e ParseError) Unwrap() error {
	return e.Err
}

// InvalidValueError is an error type that represents an invalid value.  The Value
// field contains the invalid value, and the Err field contains the error that
// caused the value to be invalid.
type InvalidValueError struct {
	Value string
	Err   error
}

// Error returns a string representation of the error in the form:
//
//	env.InvalidValueError: <value>: <error>
//
// If the Value field is empty:
//
//	env.InvalidValueError: <error>
//
// If the Err field is nil:
//
//	env.InvalidValueError: <value>
//
// If both fields are empty:
//
//	env.InvalidValueError
func (e InvalidValueError) Error() string {
	type when struct{ hasValue, hasError bool }
	fn := map[when]func() string{
		{false, false}: func() string { return "env.InvalidValueError" },
		{false, true}:  func() string { return "env.InvalidValueError: " + e.Err.Error() },
		{true, false}:  func() string { return fmt.Sprintf("env.InvalidValueError: %v", e.Value) },
		{true, true}:   func() string { return fmt.Sprintf("env.InvalidValueError: %v: %v", e.Value, e.Err) },
	}
	return fn[when{e.Value != "", e.Err != nil}]()
}

// Is reports whether the target error is a match for the receiver.
// To be a match, the target must:
//
//   - be an InvalidValueError
//   - the target Value field must match the receiver's Value, or be empty
//   - the target Err field must satisfy errors.Is with respect to the receiver Err,
//     or be nil
func (e InvalidValueError) Is(target error) bool {
	isMatch := func(target *InvalidValueError) bool {
		return (target != nil) &&
			(target.Value == "" || e.Value == target.Value) &&
			(target.Err == nil || errors.Is(e.Err, target.Err))
	}

	switch t := target.(type) {
	case *InvalidValueError:
		return isMatch(t)
	case InvalidValueError:
		return isMatch(&t)
	default:
		return false
	}
}

// Unwrap returns the error that caused the invalid value error.
func (e InvalidValueError) Unwrap() error {
	return e.Err
}

// RangeError is an error type that represents a value that is out of range; Min and
// Max fields identify the range of valid values.
//
// If Min and Max are both the zero value of T, the error represents a general out-of-range
// error with no identified range.
type RangeError[T comparable] struct {
	Min T
	Max T
}

// Error returns a string representation of the error in the form:
//
//	out of range: <Min> <= (x) <= <Max>
//
// If Min and Max are both the zero value of T:
//
//	out of range
func (e RangeError[T]) Error() string {
	if e == (RangeError[T]{}) {
		return "env.RangeError"
	}
	return fmt.Sprintf("env.RangeError: %v <= (x) <= %v", e.Min, e.Max)
}

// Is reports whether the target error is a match for the receiver.
// To be a match, the target must:
//
//   - be a RangeError
//   - the target Min and Max fields must match the receiver's Min and Max fields,
//     or be the zero value of T
func (e RangeError[T]) Is(target error) bool {
	isMatch := func(target *RangeError[T]) bool {
		var (
			zero T
		)
		return (target != nil) &&
			((target.Min == zero) || (target.Min == e.Min)) &&
			((target.Max == zero) || (target.Max == e.Max))
	}

	switch t := target.(type) {
	case *RangeError[T]:
		return isMatch(t)
	case RangeError[T]:
		return isMatch(&t)
	default:
		return false
	}
}
