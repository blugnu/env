package as_test

import (
	"math"
	"strconv"
	"testing"

	. "github.com/blugnu/test"

	"github.com/blugnu/env/as"
)

func TestInt(t *testing.T) {
	With(t)

	Run(Testcases(
		For(func(input string, assert func(int, error)) {
			result, err := as.Int(input)
			assert(result, err)
		}),

		Case("-2147483648", func(result int, err error) {
			Expect(err).IsNil()
			Expect(result).To(Equal(math.MinInt32))
		}),

		Case("0000", func(result int, err error) {
			Expect(err).IsNil()
			Expect(result).To(Equal(0))
		}),

		Case("2147483647", func(result int, err error) {
			Expect(err).IsNil()
			Expect(result).To(Equal(math.MaxInt32))
		}),

		Case("not-a-number", func(result int, err error) {
			Expect(err).Is(strconv.ErrSyntax)
			Expect(err).Is(as.ErrNotAnInteger)
			Expect(result).To(Equal(0))
		}),
	))
}
