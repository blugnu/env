package env

import (
	"errors"
	"fmt"
)

var (
	ErrNotSet            = errors.New("not set")
	ErrSetVariableFailed = errors.New("set variable failed")
)

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
	if target, ok := target.(ParseError); ok {
		return (target.VariableName == "" || e.VariableName == target.VariableName) &&
			(target.Err == nil || errors.Is(e.Err, target.Err))
	}
	return false
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
	if target, ok := target.(InvalidValueError); ok {
		return (target.Value == "" || e.Value == target.Value) &&
			(target.Err == nil || errors.Is(e.Err, target.Err))
	}
	return false
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
	if target, ok := target.(RangeError[T]); ok {
		return e == target || (target == RangeError[T]{})
	}
	return false
}
