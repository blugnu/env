package env

import (
	"errors"
	"os"
	"testing"

	"github.com/blugnu/test"
)

func TestVars_Names(t *testing.T) {
	// ACT
	result := Vars{
		"VAR3": "value2",
		"VAR1": "value1",
		"VAR2": "value2",
	}.Names()

	// ASSERT
	test.That(t, result).Equals([]string{"VAR1", "VAR2", "VAR3"})
}

func TestVars_Set(t *testing.T) {
	// ARRANGE
	defer State().Reset()
	os.Clearenv()

	// ACT
	err := Vars{
		"VAR1": "value1",
		"VAR2": "value2",
	}.Set()

	// ASSERT
	test.That(t, err).IsNil()
	test.That(t, os.Getenv("VAR1")).Equals("value1")
	test.That(t, os.Getenv("VAR2")).Equals("value2")
}

func TestVars_Set_WhenSetenvFails(t *testing.T) {
	// ARRANGE
	defer State().Reset()

	seterr := errors.New("setenv error")
	defer test.Using(&osSetenv, func(string, string) error {
		return seterr
	})()

	// ACT
	err := Vars{"VAR1": "value1"}.Set()

	// ASSERT
	test.Error(t, err).Is(seterr)
}

func TestVars_String(t *testing.T) {
	// ACT
	result := Vars{
		"VAR3": "value3",
		"VAR1": "value1",
		"VAR2": "value2",
	}.String()

	// ASSERT
	test.That(t, result).Equals(`[VAR1="value1",VAR2="value2",VAR3="value3"]`)
}

func TestVars_String_WhenEmpty(t *testing.T) {
	// ACT
	result := Vars{}.String()

	// ASSERT
	test.That(t, result).Equals("[]")
}
