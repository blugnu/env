package as_test

import (
	"strconv"
	"strings"
	"testing"

	. "github.com/blugnu/test"

	"github.com/blugnu/env/as"
)

func TestBool(t *testing.T) {
	With(t)

	type testcase struct {
		values string
		assert func(bool, error)
	}
	Run(Testcases(
		ForEach(func(tc testcase) {
			for _, s := range strings.Split(tc.values, ", ") {
				result, err := as.Bool(s)
				tc.assert(result, err)
			}
		}),

		Case("true values", testcase{
			values: "t, T, true, TRUE, True, 1, yes, YES, Yes, y, Y, YEs, TrUE",
			assert: func(result bool, err error) {
				Expect(err).IsNil()
				Expect(result).To(Equal(true))
			},
		}),

		Case("true values with whitespace", testcase{
			values: "  t , T  ,  true  ,  TRUE ,  YES   ",
			assert: func(result bool, err error) {
				Expect(err).IsNil()
				Expect(result).To(Equal(true))
			},
		}),

		Case("false values", testcase{
			values: "f, F, false, FALSE, False, 0, no, NO, No, n, N, FalSe, nO",
			assert: func(result bool, err error) {
				Expect(err).IsNil()
				Expect(result).To(Equal(false))
			},
		}),

		Case("false values with whitespace", testcase{
			values: "  f , F  ,  false  ,  FALSE ,  NO   ",
			assert: func(result bool, err error) {
				Expect(err).IsNil()
				Expect(result).To(Equal(false))
			},
		}),

		Case("invalid, empty and whitespace values", testcase{
			values: "not a valid boolean, 01, 00, 1.0, 0.0, ,     ",
			assert: func(result bool, err error) {
				Expect(err).Is(strconv.ErrSyntax)
				Expect(result).To(Equal(false))
			},
		}),
	))
}
