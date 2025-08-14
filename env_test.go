package env_test

import (
	"errors"
	"os"
	"testing"

	. "github.com/blugnu/test"

	"github.com/blugnu/env"
	"github.com/blugnu/env/internal"
)

func TestClear(t *testing.T) {
	With(t)

	// arrange
	defer env.State().Restore()
	t.Setenv("VAR", "value")

	// act
	env.Clear()

	// assert
	Expect(env.IsSet("VAR")).To(BeFalse())
}

func TestGet(t *testing.T) {
	With(t)

	// arrange
	defer env.State().Restore()
	os.Clearenv()
	t.Setenv("VAR", "value")

	// act
	result := env.Get("VAR")

	// assert
	Expect(result).To(Equal("value"), "variable present")

	// act
	result = env.Get("NOTSET")

	// assert
	Expect(result).To(Equal(""), "variable not present")
}

func TestIsSet(t *testing.T) {
	With(t)

	defer env.State().Restore()
	os.Clearenv()

	t.Setenv("VAR1", "value")

	set1 := env.IsSet("VAR1")
	set2 := env.IsSet("VAR2")

	Expect(set1).To(BeTrue(), "var 1")
	Expect(set2).To(BeFalse(), "var 2")
}

func TestLookup(t *testing.T) {
	With(t)

	// arrange
	defer env.State().Restore()
	os.Clearenv()
	t.Setenv("VAR", "value")
	t.Setenv("EMPTY", "")

	Run(Test("variable that is set", func() {
		// act
		result, ok := env.Lookup("VAR")

		// assert
		Expect(result).To(Equal("value"), "result")
		Expect(ok).To(BeTrue(), "ok")
	}))

	Run(Test("variable that is not set", func() {
		// act
		result, ok := env.Lookup("NOTSET")

		// assert
		Expect(result).To(Equal(""), "result")
		Expect(ok).To(BeFalse(), "ok")
	}))

	Run(Test("variable that is empty", func() {
		// act
		result, ok := env.Lookup("EMPTY")

		// assert
		Expect(result).To(Equal(""), "result")
		Expect(ok).To(BeTrue(), "ok")
	}))
}

func TestSet(t *testing.T) {
	With(t)

	// arrange
	defer env.State().Restore()
	os.Clearenv()

	// act
	err := env.Set("VAR1", "value1")

	// assert
	Expect(err).Should(BeNil())
	Expect(os.Getenv("VAR1")).To(Equal("value1"))
}

func TestUnset(t *testing.T) {
	With(t)

	Run(Testcases(
		ForEach(func(testcase func()) {
			defer env.State().Restore()
			os.Clearenv()

			T().Setenv("VAR1", "value1")
			T().Setenv("VAR2", "value2")

			testcase()
		}),

		Case("no names specified", func() {
			// act
			err := env.Unset()

			// assert
			Expect(err).IsNil()
			Expect(os.Getenv("VAR1")).To(Equal("value1"))
			Expect(os.Getenv("VAR2")).To(Equal("value2"))
		}),

		Case("names specified", func() {
			// act
			err := env.Unset("VAR1")

			// assert
			Expect(err).IsNil()
			Expect(os.Getenv("VAR1")).To(Equal(""))
			Expect(os.Getenv("VAR2")).To(Equal("value2"))
		}),

		Case("error when unsetting", func() {
			// arrange
			unsetErr := errors.New("unset error")
			defer Restore(Original(&internal.Unsetenv).ReplacedBy(func(string) error { return unsetErr }))

			// act
			err := env.Unset("VAR1", "VAR2")

			// assert
			Expect(err).Is(unsetErr)
			Expect(os.Getenv("VAR1")).To(Equal("value1"))
			Expect(os.Getenv("VAR2")).To(Equal("value2"))
		}),
	))
}
