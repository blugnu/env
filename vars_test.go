package env_test

import (
	"os"
	"testing"

	. "github.com/blugnu/test"

	"github.com/blugnu/env"
)

func TestVars(t *testing.T) {
	With(t)

	Run(Testcases(
		ForEach(func(testcase func()) {
			defer env.State().Restore()
			os.Clearenv()

			T().Setenv("VAR1", "value1")
			T().Setenv("VAR2", " value2 ")

			testcase()
		}),

		Case("loads all variables when none are specified", func() {
			// act
			result := env.Vars()

			// assert
			Expect(result).To(EqualMap(map[string]string{
				"VAR1": "value1",
				"VAR2": " value2 ",
			}))
		}),

		Case("specified variables (including ones not set)", func() {
			// act
			result := env.Vars("VAR1", "VAR3")

			// assert
			Expect(result).To(EqualMap(map[string]string{
				"VAR1": "value1",
			}))
		}),

		Case("trims whitespace from names", func() {
			// act
			result := env.Vars(" VAR1\t")

			// assert
			Expect(result).To(EqualMap(map[string]string{
				"VAR1": "value1",
			}))
		}),

		Case("does not trim whitespace from values", func() {
			// act
			result := env.Vars("VAR2")

			// assert
			Expect(result).To(EqualMap(map[string]string{
				"VAR2": " value2 ",
			}))
		}),
	))
}
