package env

import (
	"errors"
	"testing"

	"github.com/blugnu/test"
)

func TestParseError_Error(t *testing.T) {
	// ARRANGE
	testcases := []struct {
		scenario string
		varname  string
		err      error
		assert   func(t *testing.T, result string)
	}{
		{scenario: "with VariableName and Err",
			varname: "VAR",
			err:     errors.New("some error"),
			assert: func(t *testing.T, result string) {
				test.That(t, result).Equals("env.ParseError: VAR: some error")
			},
		},
		{scenario: "with VariableName and nil Err",
			varname: "VAR",
			assert: func(t *testing.T, result string) {
				test.That(t, result).Equals("env.ParseError: VAR")
			},
		},
		{scenario: "with empty VariableName and Err",
			err: errors.New("some error"),
			assert: func(t *testing.T, result string) {
				test.That(t, result).Equals("env.ParseError: some error")
			},
		},
		{scenario: "with empty VariableName and nil Err",
			assert: func(t *testing.T, result string) {
				test.That(t, result).Equals("env.ParseError")
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {
			// ARRANGE
			sut := ParseError{
				VariableName: tc.varname,
				Err:          tc.err,
			}

			// ACT
			result := sut.Error()

			// ASSERT
			tc.assert(t, result)
		})
	}
}

func TestParseError_Is(t *testing.T) {
	// ARRANGE
	sut := ParseError{
		VariableName: "VAR",
		Err:          errors.New("some error"),
	}
	testcases := []struct {
		scenario string
		target   error
		assert   func(t *testing.T, result bool)
	}{
		{scenario: "target: non-ParseError",
			target: errors.New("some error"),
			assert: func(t *testing.T, result bool) {
				test.Bool(t, result).IsFalse()
			},
		},
		{scenario: "target: zero-value ParseError",
			target: ParseError{},
			assert: func(t *testing.T, result bool) {
				test.Bool(t, result).IsTrue()
			},
		},
		{scenario: "target: ParseError with different VariableName",
			target: ParseError{VariableName: "OTHER", Err: errors.New("some error")},
			assert: func(t *testing.T, result bool) {
				test.Bool(t, result).IsFalse()
			},
		},
		{scenario: "target: ParseError with different Err",
			target: ParseError{VariableName: "VAR", Err: errors.New("other error")},
			assert: func(t *testing.T, result bool) {
				test.Bool(t, result).IsFalse()
			},
		},
		{scenario: "target: ParseError with same VariableName and Err",
			target: ParseError{VariableName: "VAR", Err: sut.Err},
			assert: func(t *testing.T, result bool) {
				test.Bool(t, result).IsTrue()
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {
			// ACT
			result := sut.Is(tc.target)

			// ASSERT
			tc.assert(t, result)
		})
	}
}

func TestParseError_Unwrap(t *testing.T) {
	// ARRANGE
	sut := ParseError{
		VariableName: "VAR",
		Err:          errors.New("some error"),
	}

	// ACT
	result := sut.Unwrap()

	// ASSERT
	test.Value(t, result).Equals(sut.Err)
}

func TestInvalidValueError_Error(t *testing.T) {
	// ARRANGE
	testcases := []struct {
		scenario string
		value    string
		err      error
		assert   func(t *testing.T, result string)
	}{
		{scenario: "with Value and Err",
			value: "VAL",
			err:   errors.New("some error"),
			assert: func(t *testing.T, result string) {
				test.That(t, result).Equals("env.InvalidValueError: VAL: some error")
			},
		},
		{scenario: "with Value and nil Err",
			value: "VAL",
			assert: func(t *testing.T, result string) {
				test.That(t, result).Equals("env.InvalidValueError: VAL")
			},
		},
		{scenario: "with empty Value and Err",
			err: errors.New("some error"),
			assert: func(t *testing.T, result string) {
				test.That(t, result).Equals("env.InvalidValueError: some error")
			},
		},
		{scenario: "with empty Value and nil Err",
			assert: func(t *testing.T, result string) {
				test.That(t, result).Equals("env.InvalidValueError")
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {
			// ARRANGE
			sut := InvalidValueError{
				Value: tc.value,
				Err:   tc.err,
			}

			// ACT
			result := sut.Error()

			// ASSERT
			tc.assert(t, result)
		})
	}
}

func TestInvalidValueError_Is(t *testing.T) {
	// ARRANGE
	sut := InvalidValueError{
		Value: "VAL",
		Err:   errors.New("some error"),
	}
	testcases := []struct {
		scenario string
		target   error
		assert   func(t *testing.T, result bool)
	}{
		{scenario: "target: non-InvalidValueError",
			target: errors.New("some error"),
			assert: func(t *testing.T, result bool) {
				test.Bool(t, result).IsFalse()
			},
		},
		{scenario: "target: zero-value InvalidValueError",
			target: InvalidValueError{},
			assert: func(t *testing.T, result bool) {
				test.Bool(t, result).IsTrue()
			},
		},
		{scenario: "target: InvalidValueError with different Value",
			target: InvalidValueError{Value: "OTHER", Err: errors.New("some error")},
			assert: func(t *testing.T, result bool) {
				test.Bool(t, result).IsFalse()
			},
		},
		{scenario: "target: InvalidValueError with different Err",
			target: InvalidValueError{Value: "VAL", Err: errors.New("other error")},
			assert: func(t *testing.T, result bool) {
				test.Bool(t, result).IsFalse()
			},
		},
		{scenario: "target: InvalidValueError with same Value and Err",
			target: InvalidValueError{Value: "VAL", Err: sut.Err},
			assert: func(t *testing.T, result bool) {
				test.Bool(t, result).IsTrue()
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {
			// ACT
			result := sut.Is(tc.target)

			// ASSERT
			tc.assert(t, result)
		})
	}
}

func TestInvalidValueError_Unwrap(t *testing.T) {
	// ARRANGE
	sut := InvalidValueError{
		Value: "VAL",
		Err:   errors.New("some error"),
	}

	// ACT
	result := sut.Unwrap()

	// ASSERT
	test.Value(t, result).Equals(sut.Err)
}

func TestRangeError_Error(t *testing.T) {
	// ARRANGE
	testcases := []struct {
		scenario string
		min      int
		max      int
		assert   func(t *testing.T, result string)
	}{
		{scenario: "with Min and Max",
			min: 1,
			max: 10,
			assert: func(t *testing.T, result string) {
				test.That(t, result).Equals("env.RangeError: 1 <= (x) <= 10")
			},
		},
		{scenario: "with zero Min and Max",
			assert: func(t *testing.T, result string) {
				test.That(t, result).Equals("env.RangeError")
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {
			// ARRANGE
			sut := RangeError[int]{
				Min: tc.min,
				Max: tc.max,
			}

			// ACT
			result := sut.Error()

			// ASSERT
			tc.assert(t, result)
		})
	}
}

func TestRangeError_Is(t *testing.T) {
	// ARRANGE
	sut := RangeError[int]{Min: 1, Max: 10}

	testcases := []struct {
		scenario string
		target   error
		assert   func(t *testing.T, result bool)
	}{
		{scenario: "target: non-RangeError",
			target: errors.New("some error"),
			assert: func(t *testing.T, result bool) {
				test.Bool(t, result).IsFalse()
			},
		},
		{scenario: "target: zero-value RangeError",
			target: RangeError[int]{},
			assert: func(t *testing.T, result bool) {
				test.Bool(t, result).IsTrue()
			},
		},
		{scenario: "target: RangeError with different Min",
			target: RangeError[int]{Min: 2, Max: 10},
			assert: func(t *testing.T, result bool) {
				test.Bool(t, result).IsFalse()
			},
		},
		{scenario: "target: RangeError with different Max",
			target: RangeError[int]{Min: 1, Max: 20},
			assert: func(t *testing.T, result bool) {
				test.Bool(t, result).IsFalse()
			},
		},
		{scenario: "target: RangeError with same Min and Max",
			target: RangeError[int]{Min: 1, Max: 10},
			assert: func(t *testing.T, result bool) {
				test.Bool(t, result).IsTrue()
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {
			// ACT
			result := sut.Is(tc.target)

			// ASSERT
			tc.assert(t, result)
		})
	}
}
