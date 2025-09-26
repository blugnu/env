package env

import (
	"fmt"
	"os"

	"github.com/blugnu/env/internal"
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

// IsSet returns true if the environment variable with the given name is set.
// The variable only needs to be present in the environment, not necessarily
// set to a non-empty value.
func IsSet(name string) bool {
	_, ok := internal.LookupEnv(name)
	return ok
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
//
//nolint:wrapcheck // thin wrapper over internal.Setenv
func Set(name, value string) error {
	return internal.Setenv(name, value)
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
		if err := internal.Unsetenv(k); err != nil {
			return fmt.Errorf("env.Unset: %s: %w", k, err)
		}
	}
	return nil
}
