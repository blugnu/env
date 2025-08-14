package env_test

import (
	"errors"
	"testing"

	. "github.com/blugnu/test"

	"github.com/blugnu/env"
)

func TestFileError_Error(t *testing.T) {
	With(t)

	type testcase struct {
		sut    env.FileError
		assert func(string)
	}
	Run(Testcases(
		ForEach(func(tc testcase) {
			result := tc.sut.Error()
			tc.assert(result)
		}),

		Case("zero value", testcase{
			sut: env.FileError{},
			assert: func(result string) {
				Expect(result).To(Equal("env.FileError"))
			},
		}),

		Case("with Filename", testcase{
			sut: env.NewFileError("file.txt", nil),
			assert: func(result string) {
				Expect(result).To(Equal("env.FileError: file.txt"))
			},
		}),

		Case("with Err", testcase{
			sut: env.NewFileError("", errors.New("some error")),
			assert: func(result string) {
				Expect(result).To(Equal("env.FileError: some error"))
			},
		}),

		Case("with Filename and Err", testcase{
			sut: env.NewFileError("file.txt", errors.New("some error")),
			assert: func(result string) {
				Expect(result).To(Equal("env.FileError: file.txt: some error"))
			},
		}),
	))
}

func TestFileError_Is(t *testing.T) {
	With(t)

	sut := env.NewFileError("file.txt", errors.New("some error"))

	type testcase struct {
		scenario string
		target   error
		expected bool
	}
	Run(Testcases(
		ForEach(func(tc testcase) {
			result := sut.Is(tc.target)
			Expect(result).To(Equal(tc.expected))
		}),

		Cases([]testcase{
			{scenario: "FileError/no Filename/nil error",
				target:   env.FileError{},
				expected: true,
			},
			{scenario: "FileError/no Filename/matching error",
				target:   env.NewFileError("", sut.Unwrap()),
				expected: true,
			},
			{scenario: "FileError/no Filename/non-matching error",
				target:   env.NewFileError("", errors.New("other error")),
				expected: false,
			},
			{scenario: "FileError/same Filename/nil error",
				target:   env.NewFileError("file.txt", nil),
				expected: true,
			},
			{scenario: "FileError/same Filename/matching error",
				target:   env.NewFileError("file.txt", sut.Unwrap()),
				expected: true,
			},
			{scenario: "FileError/same Filename/non-matching error",
				target:   env.NewFileError("file.txt", errors.New("other error")),
				expected: false,
			},
			{scenario: "FileError/different Filename/nil error",
				target:   env.NewFileError("other.txt", nil),
				expected: false,
			},
			{scenario: "FileError/different Filename/matching error",
				target:   env.NewFileError("other.txt", sut.Unwrap()),
				expected: false,
			},
			{scenario: "*FileError",
				target:   &env.FileError{},
				expected: true,
			},
			{scenario: "other error",
				target:   errors.New("some error"),
				expected: false,
			},
		}),
	))
}

func TestFileError_Unwrap(t *testing.T) {
	With(t)

	// arrange
	err := errors.New("some error")
	sut := env.NewFileError("file.txt", err)

	// act
	result := sut.Unwrap()

	// assert
	Expect(result).To(Equal(err))
}

func TestParseError_Error(t *testing.T) {
	With(t)

	type testcase struct {
		sut    env.ParseError
		assert func(string)
	}
	Run(Testcases(
		ForEach(func(tc testcase) {
			result := tc.sut.Error()
			tc.assert(result)
		}),

		Case("zero value", testcase{
			sut: env.ParseError{},
			assert: func(result string) {
				Expect(result).To(Equal("env.ParseError"))
			},
		}),

		Case("with VariableName", testcase{
			sut: env.ParseError{
				VariableName: "VAR",
			},
			assert: func(result string) {
				Expect(result).To(Equal("env.ParseError: VAR"))
			},
		}),

		Case("with Err", testcase{
			sut: env.ParseError{
				Err: errors.New("some error"),
			},
			assert: func(result string) {
				Expect(result).To(Equal("env.ParseError: some error"))
			},
		}),

		Case("with VariableName and Err", testcase{
			sut: env.ParseError{
				VariableName: "VAR",
				Err:          errors.New("some error"),
			},
			assert: func(result string) {
				Expect(result).To(Equal("env.ParseError: VAR: some error"))
			},
		}),
	))
}

func TestParseError_Is(t *testing.T) {
	With(t)

	// arrange
	sut := env.ParseError{
		VariableName: "VAR",
		Err:          errors.New("some error"),
	}

	type testcase struct {
		target error
		assert func(result bool)
	}
	Run(Testcases(
		ForEach(func(tc testcase) {
			result := sut.Is(tc.target)
			tc.assert(result)
		}),

		Case("target: not a env.ParseError", testcase{
			target: errors.New("some error"),
			assert: func(result bool) {
				Expect(result).To(BeFalse())
			},
		}),

		// value targets
		Case("target: zero-value env.ParseError", testcase{
			target: env.ParseError{},
			assert: func(result bool) {
				Expect(result).To(BeTrue())
			},
		}),

		Case("target: env.ParseError with different VariableName", testcase{
			target: env.ParseError{VariableName: "OTHER", Err: errors.New("some error")},
			assert: func(result bool) {
				Expect(result).To(BeFalse())
			},
		}),

		Case("target: env.ParseError with different Err", testcase{
			target: env.ParseError{VariableName: "VAR", Err: errors.New("other error")},
			assert: func(result bool) {
				Expect(result).To(BeFalse())
			},
		}),

		Case("target: env.ParseError with same VariableName and Err", testcase{
			target: env.ParseError{VariableName: "VAR", Err: sut.Err},
			assert: func(result bool) {
				Expect(result).To(BeTrue())
			},
		}),

		// pointer targets
		Case("target: zero-value *env.ParseError", testcase{
			target: &env.ParseError{},
			assert: func(result bool) {
				Expect(result).To(BeTrue())
			},
		}),

		Case("target: *env.ParseError with different VariableName", testcase{
			target: &env.ParseError{VariableName: "OTHER", Err: errors.New("some error")},
			assert: func(result bool) {
				Expect(result).To(BeFalse())
			},
		}),

		Case("target: *env.ParseError with different Err", testcase{
			target: &env.ParseError{VariableName: "VAR", Err: errors.New("other error")},
			assert: func(result bool) {
				Expect(result).To(BeFalse())
			},
		}),

		Case("target: *env.ParseError with same VariableName and Err", testcase{
			target: &env.ParseError{VariableName: "VAR", Err: sut.Err},
			assert: func(result bool) {
				Expect(result).To(BeTrue())
			},
		}),
	))
}

func TestParseError_Unwrap(t *testing.T) {
	With(t)

	// arrange
	sut := env.ParseError{
		VariableName: "VAR",
		Err:          errors.New("some error"),
	}

	// act
	result := sut.Unwrap()

	// assert
	Expect(result).To(Equal(sut.Err))
}

//nolint:dupl // false negative; not duplicate lines
func TestInvalidValueError_Error(t *testing.T) {
	With(t)

	type testcase struct {
		sut    env.InvalidValueError
		assert func(string)
	}
	Run(Testcases(
		ForEach(func(tc testcase) {
			result := tc.sut.Error()
			tc.assert(result)
		}),

		Case("with Value and Err", testcase{
			sut: env.InvalidValueError{
				Value: "VAL",
				Err:   errors.New("some error"),
			},
			assert: func(result string) {
				Expect(result).To(Equal("env.InvalidValueError: VAL: some error"))
			},
		}),

		Case("with Value and nil Err", testcase{
			sut: env.InvalidValueError{
				Value: "VAL",
			},
			assert: func(result string) {
				Expect(result).To(Equal("env.InvalidValueError: VAL"))
			},
		}),

		Case("with empty Value and Err", testcase{
			sut: env.InvalidValueError{
				Err: errors.New("some error"),
			},
			assert: func(result string) {
				Expect(result).To(Equal("env.InvalidValueError: some error"))
			},
		}),

		Case("with empty Value and nil Err", testcase{
			sut: env.InvalidValueError{},
			assert: func(result string) {
				Expect(result).To(Equal("env.InvalidValueError"))
			},
		}),
	))
}

func TestInvalidValueError_Is(t *testing.T) {
	With(t)

	sut := env.InvalidValueError{
		Value: "VAL",
		Err:   errors.New("some error"),
	}
	type testcase struct {
		target error
		assert func(result bool)
	}
	Run(Testcases(
		ForEach(func(tc testcase) {
			result := sut.Is(tc.target)
			tc.assert(result)
		}),

		Case("target: not a env.InvalidValueError", testcase{
			target: errors.New("some error"),
			assert: func(result bool) {
				Expect(result).To(BeFalse())
			},
		}),

		// value targets
		Case("target: zero-value env.InvalidValueError", testcase{
			target: env.InvalidValueError{},
			assert: func(result bool) {
				Expect(result).To(BeTrue())
			},
		}),

		Case("target: env.InvalidValueError with different Value", testcase{
			target: env.InvalidValueError{Value: "OTHER", Err: errors.New("some error")},
			assert: func(result bool) {
				Expect(result).To(BeFalse())
			},
		}),

		Case("target: env.InvalidValueError with different Err", testcase{
			target: env.InvalidValueError{Value: "VAL", Err: errors.New("other error")},
			assert: func(result bool) {
				Expect(result).To(BeFalse())
			},
		}),

		Case("target: env.InvalidValueError with same Value and Err", testcase{
			target: env.InvalidValueError{Value: "VAL", Err: sut.Err},
			assert: func(result bool) {
				Expect(result).To(BeTrue())
			},
		}),

		// pointer targets
		Case("target: zero-value *env.InvalidValueError", testcase{
			target: &env.InvalidValueError{},
			assert: func(result bool) {
				Expect(result).To(BeTrue())
			},
		}),

		Case("target: *env.InvalidValueError with different Value", testcase{
			target: &env.InvalidValueError{Value: "OTHER", Err: errors.New("some error")},
			assert: func(result bool) {
				Expect(result).To(BeFalse())
			},
		}),

		Case("target: *env.InvalidValueError with different Err", testcase{
			target: &env.InvalidValueError{Value: "VAL", Err: errors.New("other error")},
			assert: func(result bool) {
				Expect(result).To(BeFalse())
			},
		}),

		Case("target: *env.InvalidValueError with same Value and Err", testcase{
			target: &env.InvalidValueError{Value: "VAL", Err: sut.Err},
			assert: func(result bool) {
				Expect(result).To(BeTrue())
			},
		}),
	))
}

func TestInvalidValueError_Unwrap(t *testing.T) {
	With(t)

	// arrange
	sut := env.InvalidValueError{
		Value: "VAL",
		Err:   errors.New("some error"),
	}

	// act
	result := sut.Unwrap()

	// assert
	Expect(result).To(Equal(sut.Err))
}

func TestRangeError_Error(t *testing.T) {
	With(t)

	type testcase struct {
		sut    env.RangeError[int]
		assert func(result string)
	}
	Run(Testcases(
		ForEach(func(tc testcase) {
			result := tc.sut.Error()
			tc.assert(result)
		}),

		Case("with Min and Max", testcase{
			sut: env.RangeError[int]{Min: 1, Max: 10},
			assert: func(result string) {
				Expect(result).To(Equal("env.RangeError: 1 <= (x) <= 10"))
			},
		}),

		Case("with zero Min and Max", testcase{
			sut: env.RangeError[int]{Min: 0, Max: 0},
			assert: func(result string) {
				Expect(result).To(Equal("env.RangeError"))
			},
		}),
	))
}

func TestRangeError_Is(t *testing.T) {
	With(t)

	sut := env.RangeError[int]{Min: 1, Max: 10}
	type testcase struct {
		target error
		assert func(result bool)
	}
	Run(Testcases(
		ForEach(func(tc testcase) {
			result := sut.Is(tc.target)
			tc.assert(result)
		}),

		Case("target: not a env.RangeError", testcase{
			target: errors.New("some error"),
			assert: func(result bool) {
				Expect(result).To(BeFalse())
			},
		}),

		// value targets
		Case("target: zero-value env.RangeError", testcase{
			target: env.RangeError[int]{},
			assert: func(result bool) {
				Expect(result).To(BeTrue())
			},
		}),

		Case("target: env.RangeError with different Min", testcase{
			target: env.RangeError[int]{Min: 2, Max: 10},
			assert: func(result bool) {
				Expect(result).To(BeFalse())
			},
		}),

		Case("target: env.RangeError with different Max", testcase{
			target: env.RangeError[int]{Min: 1, Max: 20},
			assert: func(result bool) {
				Expect(result).To(BeFalse())
			},
		}),

		Case("target: env.RangeError with same Min and Max", testcase{
			target: env.RangeError[int]{Min: 1, Max: 10},
			assert: func(result bool) {
				Expect(result).To(BeTrue())
			},
		}),

		// pointer targets
		Case("target: zero-value *env.RangeError", testcase{
			target: &env.RangeError[int]{},
			assert: func(result bool) {
				Expect(result).To(BeTrue())
			},
		}),

		Case("target: *env.RangeError with different Min", testcase{
			target: &env.RangeError[int]{Min: 2, Max: 10},
			assert: func(result bool) {
				Expect(result).To(BeFalse())
			},
		}),

		Case("target: *env.RangeError with different Max", testcase{
			target: &env.RangeError[int]{Min: 1, Max: 20},
			assert: func(result bool) {
				Expect(result).To(BeFalse())
			},
		}),

		Case("target: *env.RangeError with same Min and Max", testcase{
			target: &env.RangeError[int]{Min: 1, Max: 10},
			assert: func(result bool) {
				Expect(result).To(BeTrue())
			},
		}),
	))
}
