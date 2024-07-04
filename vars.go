package env

import (
	"fmt"
	"slices"
	"strings"
)

// Vars is a map of environment variables.
type Vars map[string]string

// Names returns the names of all variables in the map as a sorted slice.
func (v Vars) Names() []string {
	names := make([]string, 0, len(v))
	for k := range v {
		names = append(names, k)
	}
	slices.Sort(names)
	return names
}

// Set applies the variables to the environment.  If an error occurs while setting
// a variable, the error is returned without any further variables being set.
//
// # returns
//
//	error   // any error that occurs while setting the variables
func (v Vars) Set() error {
	for k, v := range v {
		if err := osSetenv(k, v); err != nil {
			return fmt.Errorf("failed to set environment variable %s: %w", k, err)
		}
	}
	return nil
}

// String returns a string representation of the variables. The result is a string of
// comma delimited NAME="VALUE" entries, sorted by NAME and enclosed in [-]'s.
//
// # result
//
//	string   // a string representation of the variables in the form:
//
//		[NAME1="VALUE1",NAME2="VALUE2",NAME3="VALUE3"]
//
// If the map is empty the result is "[]".
func (v Vars) String() string {
	n := v.Names()
	ls := make([]string, 0, len(n))
	for _, k := range n {
		ls = append(ls, k+`="`+v[k]+`"`)
	}
	return "[" + strings.Join(ls, ",") + "]"
}
