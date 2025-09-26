package env_test

import (
	"errors"
	"strconv"
	"testing"

	. "github.com/blugnu/test"

	"github.com/blugnu/env"
	"github.com/blugnu/env/internal"
)

func TestParse(t *testing.T) {
	With(t)

	// arrange
	defer Restore(Original(&internal.LookupEnv).ReplacedBy(func(string) (string, bool) {
		return "123", true
	}))

	// act
	value, err := env.Parse("VAR", strconv.Atoi)

	// assert
	Expect(err).IsNil()
	Expect(value).To(Equal(123))
}

func TestParse_WhenVariableNotSet(t *testing.T) {
	With(t)

	// arrange
	defer Restore(Original(&internal.LookupEnv).ReplacedBy(func(string) (string, bool) {
		return "", false
	}))

	// act
	value, err := env.Parse("NOT_SET", func(s string) (string, error) { return s, nil })

	// assert
	Expect(err).Is(env.ErrNotSet)
	Expect(err).Is(env.ParseError{VariableName: "NOT_SET"})
	Expect(value).To(Equal(""))
}

func TestParse_WhenConversionFails(t *testing.T) {
	With(t)

	// arrange
	converr := errors.New("conversion error")
	defer Restore(Original(&internal.LookupEnv).ReplacedBy(func(string) (string, bool) {
		return "not-a-number", true
	}))

	// act
	value, err := env.Parse("NOT_A_NUMBER", func(s string) (int, error) { return 0, converr })

	// assert
	Expect(err).Is(converr)
	Expect(err).Is(env.ParseError{VariableName: "NOT_A_NUMBER"})
	Expect(err).Is(env.InvalidValueError{Value: "not-a-number"})
	Expect(value).To(Equal(0))
}

func TestParse_WithDefaultWhenNotSet(t *testing.T) {
	With(t)

	// arrange
	defer Restore(Original(&internal.LookupEnv).ReplacedBy(func(string) (string, bool) {
		return "", false
	}))

	// act
	value, err := env.Parse("VAR", strconv.Atoi, 42)

	// assert
	Expect(err).IsNil()
	Expect(value).To(Equal(42))
}

func TestParse_WithDefaultWhenSet(t *testing.T) {
	With(t)

	// arrange
	defer Restore(Original(&internal.LookupEnv).ReplacedBy(func(string) (string, bool) {
		return "123", true
	}))

	// act
	value, err := env.Parse("VAR", strconv.Atoi, 42)

	// assert
	Expect(err).IsNil()
	Expect(value).To(Equal(123))
}

func TestParse_WithDefaultWhenError(t *testing.T) {
	With(t)

	// arrange
	defer Restore(Original(&internal.LookupEnv).ReplacedBy(func(string) (string, bool) {
		return "abc", true
	}))

	// act
	value, err := env.Parse("VAR", strconv.Atoi, 42)

	// assert
	Expect(err).Is(strconv.ErrSyntax)
	Expect(err).Is(env.ParseError{VariableName: "VAR"})
	Expect(err).Is(env.InvalidValueError{Value: "abc"})
	Expect(value).To(Equal(0))
}

func TestParseInto(t *testing.T) {
	With(t)

	// arrange
	var value int
	defer Restore(Original(&internal.LookupEnv).ReplacedBy(func(string) (string, bool) {
		return "123", true
	}))

	// act
	err := env.ParseInto(&value, "VAR", strconv.Atoi)

	// assert
	Expect(err).IsNil()
	Expect(value).To(Equal(123))
}

func TestParseInto_WhenErrorOccurs(t *testing.T) {
	With(t)

	// arrange
	var value = 42
	defer Restore(Original(&internal.LookupEnv).ReplacedBy(func(string) (string, bool) {
		return "not a number", true
	}))

	// act
	err := env.ParseInto(&value, "VAR", strconv.Atoi)

	// assert
	Expect(err).Is(env.ParseError{VariableName: "VAR"})
	Expect(err).Is(env.InvalidValueError{Value: "not a number", Err: strconv.ErrSyntax})
	Expect(value).To(Equal(42))
}

func TestParseInto_WhenTargetIsNil(t *testing.T) {
	With(t)

	// act
	err := env.ParseInto(nil, "VAR", strconv.Atoi)

	// assert
	Expect(err).Is(env.ParseError{VariableName: "VAR", Err: env.ErrTargetIsNil})
}
