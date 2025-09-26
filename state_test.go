package env_test

import (
	"os"
	"strings"
	"testing"

	. "github.com/blugnu/test"

	"github.com/blugnu/env"
)

func TestState(t *testing.T) {
	With(t)

	// arrange
	environ := os.Environ()

	// act
	result := env.State()

	// assert
	Expect([]string(result)).To(DeepEqual(environ))
}

func TestState_Restore(t *testing.T) {
	With(t)

	// os-level implementation to capture and restore current env
	// and restore on test completion
	orig := os.Environ()
	defer func(environ []string) {
		os.Clearenv()
		for _, kv := range environ {
			if i := strings.IndexByte(kv, '='); i > 0 {
				_ = os.Setenv(kv[:i], kv[i+1:]) //nolint:usetesting // must use os.Setenv for test cleanup
			}
		}
	}(orig)

	// arrange
	t.Setenv("VAR1", "value1")
	t.Setenv("VAR2", "value2")

	state := env.State()

	os.Clearenv()

	// act
	t.Setenv("NEWVAR", "newvalue")
	state.Restore()

	// assert: restores variables to their original state
	Expect(os.Getenv("VAR1")).To(Equal("value1"))
	Expect(os.Getenv("VAR2")).To(Equal("value2"))
	// assert: removes variables not present in the original state
	Expect(env.IsSet("NEWVAR")).To(BeFalse())
}
