package env

import (
	"os"
	"strings"
	"testing"

	"github.com/blugnu/test"
)

func TestState(t *testing.T) {
	// ARRANGE
	env := os.Environ()

	// ACT
	result := State()

	// ASSERT
	test.Slice(t, result).Equals(env)
}

func TestState_Reset(t *testing.T) {
	// ARRANGE
	og := os.Environ()
	defer func() {
		os.Clearenv()
		for _, s := range og {
			k, v, _ := strings.Cut(s, "=")
			os.Setenv(k, v)
		}
	}()

	os.Clearenv()

	// ACT
	state{
		"VAR1=value1",
		"VAR2=value2",
	}.Reset()

	// ASSERT
	test.That(t, os.Getenv("VAR1")).Equals("value1")
	test.That(t, os.Getenv("VAR2")).Equals("value2")
}
