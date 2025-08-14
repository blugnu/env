package as_test

import (
	"net/url"
	"testing"

	. "github.com/blugnu/test"

	"github.com/blugnu/env/as"
)

func TestAbsoluteURL(t *testing.T) {
	With(t)

	type testcase struct {
		input  string
		assert func(*url.URL, error)
	}
	Run(Testcases(
		ForEach(func(tc testcase) {
			result, err := as.AbsoluteURL(tc.input)
			tc.assert(result, err)
		}),

		Case("valid url", testcase{
			input: "http://example.com",
			assert: func(result *url.URL, err error) {
				Expect(err).IsNil()
				Expect(result.Scheme).To(Equal("http"))
				Expect(result.Host).To(Equal("example.com"))
			},
		}),

		Case("relative url", testcase{
			input: "relative/path",
			assert: func(result *url.URL, err error) {
				Expect(err).Is(as.ErrNotAnAbsoluteURL)
				Expect(result).IsNil()
			},
		}),

		Case("invalid url", testcase{
			input: string([]byte{0x00, 0x01, 0x02}), // invalid URL
			assert: func(result *url.URL, err error) {
				Expect(any(err)).To(BeOfType[*url.Error]())
				Expect(result).IsNil()
			},
		}),
	))
}
