package as_test

import (
	"errors"
	"strconv"
	"testing"

	. "github.com/blugnu/test"

	"github.com/blugnu/env"
	"github.com/blugnu/env/as"
)

func TestPortNo(t *testing.T) {
	With(t)

	Run(Testcases(
		For(func(input string, assert func(int, error)) {
			result, err := as.PortNo(input)
			assert(result, err)
		}),

		Case("-1", func(result int, err error) {
			rangeErr := env.RangeError[int]{}
			Expect(errors.As(err, &rangeErr)).To(BeTrue())
			Expect(rangeErr.Min).To(Equal(0))
			Expect(rangeErr.Max).To(Equal(65535))
			Expect(result).To(Equal(0))
		}),

		Case("0", func(result int, err error) {
			Expect(err).IsNil()
			Expect(result).To(Equal(0))
		}),

		Case("8080", func(result int, err error) {
			Expect(err).IsNil()
			Expect(result).To(Equal(8080))
		}),

		Case("65535", func(result int, err error) {
			Expect(err).IsNil()
			Expect(result).To(Equal(65535))
		}),

		Case("65536", func(result int, err error) {
			Expect(err).Is(env.RangeError[int]{})
			Expect(result).To(Equal(0))
		}),

		Case("not-a-number", func(result int, err error) {
			Expect(err).Is(strconv.ErrSyntax)
			Expect(result).To(Equal(0))
		}),
	))
}
