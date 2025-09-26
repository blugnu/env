package env

import (
	"github.com/blugnu/env/internal"
)

// ConversionFunc is the type of a function that converts a string to a
// value of type T.  The function should return an error if the conversion
// cannot be performed.
type ConversionFunc[T any] func(string) (T, error)

// Parse parses the environment variable with the given name and returns a value of
// type T, obtained by passing the value of the environment variable to a provided
// conversion function.
//
// # Default Value
//
// A default value may be provided (optional) which will be returned if the specified
// environment variable is not set.
//
// If the environment variable is set but cannot be parsed or converted, an error
// will be returned and any default value (if provided) is ignored.
//
// The default value is not used if the variable is set to an empty string.
//
// # parameters
//
//	name string             // the name of the environment variable to parse
//
//	cnv ConversionFunc[T]   // a function to parse the environment variable;
//	                        // the function should return a value of type T and
//	                        // an error if the value cannot be converted
//
//	def ...T                // optional default value to use if the environment
//	                        // variable is not set. Only the first default value
//	                        // is used. The argument is variadic to allow for no
//	                        // default value.
//
// # returns
//
//	T       // the value of the environment variable; if an error occurs the
//	        // zero value of T is returned
//
//	error   // any [ParseError] resulting from parsing the environment variable;
//	        // if the variable is not set and no default is provided, the error
//	        // wraps [ErrNotSet]; if variable is set but cannot be converted,
//			// the error wraps an [InvalidValueError].
//
// # conversion functions
//
// The `as` package provides a number of conversion functions that can be used.
// For example, to parse an integer environment variable:
//
//	value, err := env.Parse("MY_INT_VAR", as.Int)
func Parse[T any](name string, cnv ConversionFunc[T], def ...T) (T, error) {
	handleError := func(err error) (T, error) {
		return *new(T), ParseError{VariableName: name, Err: err}
	}

	if v, ok := internal.LookupEnv(name); ok {
		r, err := cnv(v)
		if err != nil {
			return handleError(InvalidValueError{Value: v, Err: err})
		}
		return r, nil
	}

	if len(def) > 0 {
		return def[0], nil
	}

	return handleError(ErrNotSet)
}

// ParseInto is a wrapper around [Parse] that stores the result into a variable
// identified by a provided target pointer and returns only an error (or nil).
// If the target is nil a [ParseError] is returned with an underlying
// [ErrTargetIsNil] error.
//
// ParseInto is a convenience function to reduce boilerplate code when parsing
// multiple environment variables.
//
// # parameters
//
//	target *T               // pointer to the variable to store the result
//
//	name string             // the name of the environment variable to parse
//
//	cnv ConversionFunc[T]   // a function to parse the environment variable;
//	                        // the function should return a value of type T and
//	                        // an error if the value cannot be converted
//
//	def ...T                // optional default value to use if the environment
//	                        // variable is not set. Only the first default value
//	                        // is used. The argument is variadic to allow for no
//	                        // default value.
func ParseInto[T any](target *T, name string, cnv ConversionFunc[T], def ...T) error {
	if target == nil {
		return ParseError{VariableName: name, Err: ErrTargetIsNil}
	}

	v, err := Parse(name, cnv, def...)
	if err != nil {
		return err
	}

	*target = v
	return nil
}
