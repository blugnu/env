package env

import (
	"errors"
	"os"
	"testing"

	"github.com/blugnu/test"
)

func TestClear(t *testing.T) {
	// ARRANGE
	defer State().Reset()
	os.Setenv("VAR", "value")

	// ACT
	Clear()

	// ASSERT
	_, isSet := os.LookupEnv("VAR")
	test.IsFalse(t, isSet)
}

func TestGet(t *testing.T) {
	// ARRANGE
	defer State().Reset()
	os.Clearenv()
	os.Setenv("VAR", "value")

	// ACT
	result := Get("VAR")

	// ASSERT
	test.That(t, result, "variable present").Equals("value")

	// ACT
	result = Get("NOTSET")

	// ASSERT
	test.That(t, result, "variable not present").Equals("")
}

func TestGetVars(t *testing.T) {
	// ARRANGE
	testcases := []struct {
		scenario string
		exec     func(t *testing.T)
	}{
		{scenario: "all variables",
			exec: func(t *testing.T) {
				// ARRANGE
				defer State().Reset()
				os.Clearenv()
				os.Setenv("VAR1", "value1")
				os.Setenv("VAR2", "value2")

				// ACT
				result := GetVars()

				// ASSERT
				test.Map(t, result).Equals(Vars{"VAR1": "value1", "VAR2": "value2"})
			},
		},
		{scenario: "specified variables (including ones not set)",
			exec: func(t *testing.T) {
				// ARRANGE
				defer State().Reset()
				os.Clearenv()
				os.Setenv("VAR1", "value1")
				os.Setenv("VAR2", "value2")

				// ACT
				result := GetVars("VAR1", "VAR3")

				// ASSERT
				test.Map(t, result).Equals(Vars{"VAR1": "value1"})
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {
			tc.exec(t)
		})
	}
}

func TestLookup(t *testing.T) {
	// ARRANGE
	defer State().Reset()
	os.Clearenv()
	os.Setenv("VAR", "value")

	// ACT
	result, ok := Lookup("VAR")

	// ASSERT
	test.That(t, result, "variable present").Equals("value")
	test.IsTrue(t, ok, "variable present")

	// ACT
	result, ok = Lookup("NOTSET")

	// ASSERT
	test.That(t, result, "variable not present").Equals("")
	test.IsFalse(t, ok, "variable not present")
}

func TestSet(t *testing.T) {
	// ARRANGE
	defer State().Reset()
	os.Clearenv()

	// ACT
	err := Set("VAR1", "value1")

	// ASSERT
	test.That(t, err).IsNil()
	test.That(t, os.Getenv("VAR1")).Equals("value1")
}

func TestUnset(t *testing.T) {
	// ARRANGE
	testEnv := map[string]string{
		"VAR1": "value1",
		"VAR2": "value2",
	}
	testcases := []struct {
		scenario string
		exec     func(t *testing.T)
	}{
		{scenario: "no names specified",
			exec: func(t *testing.T) {
				// ACT
				err := Unset()

				// ASSERT
				test.That(t, err).IsNil()
				test.That(t, os.Getenv("VAR1")).Equals("value1")
				test.That(t, os.Getenv("VAR2")).Equals("value2")
			},
		},
		{scenario: "name specified",
			exec: func(t *testing.T) {
				// ACT
				err := Unset("VAR1")

				// ASSERT
				test.That(t, err).IsNil()
				_, isSet := os.LookupEnv("VAR1")
				test.IsFalse(t, isSet)
				test.That(t, os.Getenv("VAR2")).Equals("value2")
			},
		},
		{
			scenario: "error when unsetting",
			exec: func(t *testing.T) {
				// ARRANGE
				unseterr := errors.New("unset error")
				defer test.Using(&osUnsetenv, func(string) error { return unseterr })()

				// ACT
				err := Unset("VAR1", "VAR2")

				// ASSERT
				test.Error(t, err).Is(unseterr)
				test.That(t, os.Getenv("VAR1")).Equals("value1")
				test.That(t, os.Getenv("VAR2")).Equals("value2")
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {
			// ARRANGE
			defer State().Reset()
			os.Clearenv()
			for k, v := range testEnv {
				os.Setenv(k, v)
			}

			// ACT & ASSERT
			tc.exec(t)
		})
	}
}
