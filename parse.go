package env

// ConversionFunc is a function that converts a string to a value of type T.
// It is the type of the conversion function used by the Parse and Override
// functions.
//
// # parameters
//
//	string string   // the string to convert
//
// # returns
//
//	T       // the converted value
//
//	error   // any error that occurs during conversion
type ConversionFunc[T any] func(string) (T, error)

// Override replaces the current value of some variable with the value obtained
// by parsing a named environment variable. If the named environment variable is
// not set, the target variable is not modified.
//
// If an error occurs during parsing or conversion, the target variable is not
// modified and the error is returned.
//
// The function returns true if the target variable was modified from its
// original value, false otherwise.
//
// # parameters
//
//	dest *T                       // a pointer to the variable to be changed
//
//	name string                   // the name of the environment variable to parse
//
//	cnv func(string) (T, error)   // a function to parse the environment variable
//
// # returns
//
//	bool    // true if the target variable was modified, false otherwise
//
//	error   // any error resulting from parsing the named environment variable
//
// The bool result should not be relied on to determine whether an error occurred;
// The function will return false when:
//
//   - an error occurs parsing or converting the named environment variable
//   - the named environment variable is not set
//   - the named environment variable is set with a value that is successfully
//     parsed and converted to yield the same value as currently held in the
//     destination variable.
//
// Only the error result can determine if an error occurred.
func Override[T comparable](dest *T, name string, cnv func(string) (T, error)) (bool, error) {
	v, err := Parse(name, cnv)
	switch {
	case err == nil && *dest != v:
		*dest = v
		return true, nil
	default:
		return false, err
	}
}

// Parse parses the environment variable with the given name and returns a value of
// type T, obtained by passing the value of the environment variable to a provided
// conversion function.
//
// # parameters
//
//	name string             // the name of the environment variable to parse
//
//	cnv ConversionFunc[T]   // a function to parse the environment variable;
//	                        // the function should return a value of type T and
//	                        // an error if the value cannot be converted
//
// # returns
//
//	T       // the value of the environment variable; if an error occurs the
//	        // zero value of T is returned
//
//	error   // any error resulting from parsing the environment variable;
//	        // if the error is ErrNotSet a nil error is returned
//
// # conversion functions
//
// The `as` package provides a number of conversion functions that can be used with
// this function.  For example, to parse an integer environment variable:
//
//	value, err := env.Parse("MY_INT_VAR", as.Int)
func Parse[T any](name string, cnv ConversionFunc[T]) (T, error) {
	if v, ok := osLookupEnv(name); ok {
		r, err := cnv(v)
		if err != nil {
			return *new(T), ParseError{VariableName: name, Err: InvalidValueError{Value: v, Err: err}}
		}
		return r, nil
	}
	return *new(T), ParseError{VariableName: name, Err: ErrNotSet}
}
