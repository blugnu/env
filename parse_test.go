package env

import (
	"errors"
	"os"
	"strconv"
	"testing"

	"github.com/blugnu/test"
)

func TestParse(t *testing.T) {
	// ARRANGE
	defer test.Using(&osLookupEnv, func(string) (string, bool) {
		return "123", true
	})()

	// ACT
	value, err := Parse("VAR", strconv.Atoi)

	// ASSERT
	test.That(t, value).Equals(123)
	test.That(t, err).IsNil()
}

func TestParse_WhenVariableNotSet(t *testing.T) {
	// ARRANGE
	defer test.Using(&osLookupEnv, func(string) (string, bool) {
		return "", false
	})()

	// ACT
	value, err := Parse("NOT_SET", func(s string) (string, error) { return s, nil })

	// ASSERT
	test.That(t, value).Equals("")
	test.Error(t, err).Is(ErrNotSet)
}

func TestParse_WhenConversionFails(t *testing.T) {
	// ARRANGE
	converr := errors.New("conversion error")
	defer test.Using(&osLookupEnv, func(string) (string, bool) {
		return "not-a-number", true
	})()

	// ACT
	value, err := Parse("NOT_A_NUMBER", func(s string) (int, error) { return 0, converr })

	// ASSERT
	test.That(t, value).Equals(0)
	test.Error(t, err).Is(InvalidValueError{Value: "not-a-number", Err: converr})
}

func TestOverride(t *testing.T) {
	// ARRANGE
	var value = 42
	defer State().Reset()
	os.Clearenv()
	os.Setenv("VAR", "123")

	// ACT
	result, err := Override(&value, "VAR", strconv.Atoi)

	// ASSERT
	test.That(t, err).IsNil()
	test.IsTrue(t, result)
	test.That(t, value).Equals(123)
}

func TestOverride_WhenValueIsNotChanged(t *testing.T) {
	// ARRANGE
	var value = 123
	defer State().Reset()
	os.Clearenv()
	os.Setenv("VAR", "123")

	// ACT
	result, err := Override(&value, "VAR", strconv.Atoi)

	// ASSERT
	test.That(t, err).IsNil()
	test.IsFalse(t, result)
	test.That(t, value).Equals(123)
}

func TestOverride_WhenVariableIsNotSet(t *testing.T) {
	// ARRANGE
	var value = 42
	defer State().Reset()
	os.Clearenv()

	// ACT
	result, err := Override(&value, "VAR", strconv.Atoi)

	// ASSERT
	test.Error(t, err).Is(ErrNotSet)
	test.IsFalse(t, result)
	test.That(t, value).Equals(42)
}
