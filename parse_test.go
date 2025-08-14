package env_test

import (
	"errors"
	"os"
	"strconv"
	"testing"

	"github.com/blugnu/test"

	"github.com/blugnu/env"
	"github.com/blugnu/env/internal"
)

func TestOverride(t *testing.T) {
	// ARRANGE
	var value = 42
	defer env.State().Reset()
	os.Clearenv()
	t.Setenv("VAR", "123")

	// ACT
	result, err := env.Override(&value, "VAR", strconv.Atoi)

	// ASSERT
	test.That(t, err).IsNil()
	test.IsTrue(t, result)
	test.That(t, value).Equals(123)
}

func TestOverride_WhenValueIsNotChanged(t *testing.T) {
	// ARRANGE
	var value = 123
	defer env.State().Reset()
	os.Clearenv()
	t.Setenv("VAR", "123")

	// ACT
	result, err := env.Override(&value, "VAR", strconv.Atoi)

	// ASSERT
	test.That(t, err).IsNil()
	test.IsFalse(t, result)
	test.That(t, value).Equals(123)
}

func TestOverride_WhenVariableIsNotSet(t *testing.T) {
	// ARRANGE
	var value = 42
	defer env.State().Reset()
	os.Clearenv()

	// ACT
	result, err := env.Override(&value, "VAR", strconv.Atoi)

	// ASSERT
	test.Error(t, err).Is(env.ErrNotSet)
	test.IsFalse(t, result)
	test.That(t, value).Equals(42)
}

func TestParse(t *testing.T) {
	// ARRANGE
	defer test.Using(&internal.LookupEnv, func(string) (string, bool) {
		return "123", true
	})()

	// ACT
	value, err := env.Parse("VAR", strconv.Atoi)

	// ASSERT
	test.That(t, err).IsNil()
	test.That(t, value).Equals(123)
}

func TestParse_WhenVariableNotSet(t *testing.T) {
	// ARRANGE
	defer test.Using(&internal.LookupEnv, func(string) (string, bool) {
		return "", false
	})()

	// ACT
	value, err := env.Parse("NOT_SET", func(s string) (string, error) { return s, nil })

	// ASSERT
	test.Error(t, err).Is(env.ErrNotSet)
	test.That(t, value).Equals("")
}

func TestParse_WhenConversionFails(t *testing.T) {
	// ARRANGE
	converr := errors.New("conversion error")
	defer test.Using(&internal.LookupEnv, func(string) (string, bool) {
		return "not-a-number", true
	})()

	// ACT
	value, err := env.Parse("NOT_A_NUMBER", func(s string) (int, error) { return 0, converr })

	// ASSERT
	test.Error(t, err).Is(env.InvalidValueError{Value: "not-a-number", Err: converr})
	test.That(t, value).Equals(0)
}

func TestParse_WithDefaultWhenNotSet(t *testing.T) {
	// ARRANGE
	defer test.Using(&internal.LookupEnv, func(string) (string, bool) {
		return "", false
	})()

	// ACT
	value, err := env.Parse("VAR", strconv.Atoi, 42)

	// ASSERT
	test.That(t, err).IsNil()
	test.That(t, value).Equals(42)
}

func TestParse_WithDefaultWhenSet(t *testing.T) {
	// ARRANGE
	defer test.Using(&internal.LookupEnv, func(string) (string, bool) {
		return "123", true
	})()

	// ACT
	value, err := env.Parse("VAR", strconv.Atoi, 42)

	// ASSERT
	test.That(t, err).IsNil()
	test.That(t, value).Equals(123)
}

func TestParse_WithDefaultWhenError(t *testing.T) {
	// ARRANGE
	defer test.Using(&internal.LookupEnv, func(string) (string, bool) {
		return "abc", true
	})()

	// ACT
	value, err := env.Parse("VAR", strconv.Atoi, 42)

	// ASSERT
	test.Error(t, err).Is(strconv.ErrSyntax)
	test.That(t, value).Equals(0)
}
