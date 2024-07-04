package env

import (
	"os"
	"strings"
)

// state represents the state of all environment variables captured at
// a point in time.  It is used to reset the environment to a known
// state, typically after a test has modified the environment.
type state []string

// State captures the current state of the environment variables. It
// is typically used to capture the state before a test modifies the
// environment so that it can be reset to the original state after
// the test has run (by deferring a call to the Reset method on the
// returned state value).
//
// # returns
//
//	state   // the current state of the environment variables
//
// The state is captured using os.Environ().
//
// # example
//
//	func TestSomething(t *testing.T) {
//		// ARRANGE
//		defer env.State().Reset()
//		env.Vars{
//			"VAR1", "value1",
//			"VAR2", "value2",
//		}.Set()
//
//		// ACT
//		// ...
//	}
func State() state {
	return os.Environ()
}

// Reset resets the environment variables to the state captured when
// the State function was called.
//
// This function is typically used in a defer statement to ensure the
// environment is reset after a test has modified the environment.
//
// # example
//
//	func TestSomething(t *testing.T) {
//		// ARRANGE
//		defer env.State().Reset()
//		env.Vars{
//			"VAR1", "value1",
//			"VAR2", "value2",
//		}.Set()
//
//		// ACT
//		// ...
//	}
func (s state) Reset() {
	os.Clearenv()
	for _, e := range s {
		k, v, _ := strings.Cut(e, "=")
		os.Setenv(k, v)
	}
}
