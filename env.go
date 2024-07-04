package env

import (
	"os"
	"strings"
)

// Clear removes all environment variables.
func Clear() {
	os.Clearenv()
}

// Get returns the value of the environment variable with the given name.  If the
// variable is not set an empty string is returned.
//
// To differentiate between a variable that is not set and a variable that is set to
// an empty string, use the `Lookup` function.
//
// # parameters
//
//	name string   // the name of the environment variable
//
// # returns
//
//	string   // the value of the environment variable
func Get(name string) string {
	return os.Getenv(name)
}

// GetVars returns a map of environment variables.  If no variable names are provided
// all environment variables are returned.  If variable names are provided, only
// those variables are returned (if set).
//
// # parameters
//
//	names ...string   // (optional) names of environment variables to return;
//	                 // if no names are provided the returned map contains all
//	                 // environment variables.
//
// If a name is provided that is not set in the environment it is not included in
// the returned map.
//
// # returns
//
//	Vars   // a map of environment variables
//
// The returned map is a `map[string]string` where the key is the name of the
// environment variable and the value is the value of the environment variable.
//
// If no environment variables are set or all specified variables names are not set,
// the returned map is empty.
func GetVars(names ...string) Vars {
	var result Vars

	switch len(names) {
	case 0: // all environment variables
		env := os.Environ()
		result = make(Vars, len(env))
		for _, s := range env {
			k, v, _ := strings.Cut(s, "=")
			result[k] = v
		}
	default: // only the named variables (if set)
		result = make(Vars, len(names))
		for _, k := range names {
			if v, ok := os.LookupEnv(k); ok {
				result[k] = v
			}
		}
	}

	return result
}

// Lookup returns the value of the environment variable with the given name and a
// boolean indicating whether the variable is set.  If the variable is not set the
// returned value is an empty string and the boolean is `false`.
//
// If you do not need to differentiate between a variable that is not set and a
// variable that is set to an empty string, use the `Get` function.
//
// # parameters
//
//	name string   // the name of the environment variable
//
// # returns
//
//	string   // the value of the environment variable
//
//	bool     // true if the variable is set, false otherwise
func Lookup(name string) (string, bool) {
	return os.LookupEnv(name)
}

// Set sets the value of the environment variable with the given name.  If the variable
// does not exist it is created.
//
// # parameters
//
//	name  string   // the name of the environment variable
//
//	value string   // the value to set
//
// # returns
//
//	error   // any error that occurs while setting the environment variable
func Set(name, value string) error {
	return os.Setenv(name, value)
}

// Unset removes the environment variables with the given names.  If a variable does
// not exist (already not set) it is ignored.
//
// # parameters
//
//	names ...string   // the names of the environment variables to remove
//
// # returns
//
//	error   // any error that occurs while unsetting the environment variables
//
// On Unix systems (including Linux and macOS) the error is always `nil`, but
// on Windows systems the error may be non-nil.
func Unset(name ...string) error {
	for _, k := range name {
		if err := osUnsetenv(k); err != nil {
			return err
		}
	}
	return nil
}
