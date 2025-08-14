package as_test

import (
	"strconv"
	"testing"
	"time"

	. "github.com/blugnu/test"

	"github.com/blugnu/env/as"
)

func TestDuration(t *testing.T) {
	With(t)

	Run(Testcases(
		For(func(input string, assert func(time.Duration, error)) {
			result, err := as.Duration(input)
			assert(result, err)
		}),

		Case("1h10m", func(result time.Duration, err error) {
			Expect(err).IsNil()
			Expect(result).To(Equal(1*time.Hour + 10*time.Minute))
		}),

		Case("1.5h", func(result time.Duration, err error) {
			Expect(err).IsNil()
			Expect(result).To(Equal(1*time.Hour + 30*time.Minute))
		}),

		Case("1", func(result time.Duration, err error) {
			Expect(err).Is(as.ErrNotADuration)
			Expect(result).To(Equal(time.Duration(0)))
		}),

		Case("not-a-duration", func(result time.Duration, err error) {
			Expect(err).Is(as.ErrNotADuration)
			Expect(result).To(Equal(time.Duration(0)))
		}),
	))

	Run(Test("valid duration with whitespace", func() {
		result, err := as.Duration(" 1h ")
		Expect(err).IsNil()
		Expect(result).To(Equal(1 * time.Hour))
	}))

	Run(Test("empty duration", func() {
		result, err := as.Duration("  ")
		Expect(err).Is(as.ErrNotADuration)
		Expect(result).To(Equal(time.Duration(0)))
	}))
}

func TestDurationIn(t *testing.T) {
	With(t)

	type testcase struct {
		input  string
		unit   time.Duration
		assert func(time.Duration, error)
	}
	Run(Testcases(
		ForEach(func(tc testcase) {
			// act
			result, err := as.DurationIn(tc.unit)(tc.input)

			// assert
			tc.assert(result, err)
		}),

		Case("applies specified unit to plain integer value", testcase{
			input: "1",
			unit:  time.Hour,
			assert: func(result time.Duration, err error) {
				Expect(err).IsNil()
				Expect(result).To(Equal(time.Hour))
			},
		}),

		Case("ignores unit for valid duration expression", testcase{
			input: "1h",
			unit:  time.Second,
			assert: func(result time.Duration, err error) {
				Expect(err).IsNil()
				Expect(result).To(Equal(time.Hour))
			},
		}),

		Case("with invalid value", testcase{
			input: "abc",
			assert: func(result time.Duration, err error) {
				Expect(err).Is(strconv.ErrSyntax)
				Expect(result).To(Equal(time.Duration(0)))
			},
		}),

		// MARK: edge cases

		Case("with zero value", testcase{
			input: "0",
			unit:  time.Minute,
			assert: func(result time.Duration, err error) {
				Expect(err).IsNil()
				Expect(result).To(Equal(time.Duration(0)))
			},
		}),

		Case("with negative value", testcase{
			input: "-2",
			unit:  time.Minute,
			assert: func(result time.Duration, err error) {
				Expect(err).IsNil()
				Expect(result).To(Equal(-2 * time.Minute))
			},
		}),

		Case("with negative duration expression", testcase{
			input: "-2h30m",
			unit:  time.Minute,
			assert: func(result time.Duration, err error) {
				Expect(err).IsNil()
				Expect(result).To(Equal(-2*time.Hour - 30*time.Minute))
			},
		}),

		Case("with decimal unit", testcase{
			input: "1.5",
			unit:  time.Second,
			assert: func(result time.Duration, err error) {
				Expect(err).Is(as.ErrNotADuration)
				Expect(err).Is(as.ErrNotAnInteger)
				Expect(err).Is(as.ErrNotAnIntegerOrDuration)
				Expect(result).To(Equal(time.Duration(0)))
			},
		}),
	))
}
