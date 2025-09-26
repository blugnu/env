package as_test

import (
	"errors"
	"os"
	"testing"

	. "github.com/blugnu/test"

	"github.com/blugnu/env/as"
)

func TestFilename(t *testing.T) {
	With(t)

	type testcase struct {
		input  string
		assert func(string, error)
	}
	Run(Testcases(
		ForEach(func(tc testcase) {
			// act
			result, err := as.Filename(tc.input)

			// assert
			tc.assert(result, err)
		}),

		Case("valid filename that exists", testcase{
			input: "./filename_test.go",
			assert: func(result string, err error) {
				Expect(err).IsNil()
				Expect(result).To(Equal("./filename_test.go"))
			},
		}),

		Case("valid filename that does not exist", testcase{
			input: "./file/does/not/exist",
			assert: func(result string, err error) {
				Expect(err).Is(os.ErrNotExist)
				Expect(result).To(Equal(""))
			},
		}),

		Case("invalid filename", testcase{
			input: string([]byte{0x00, 0x01, 0x02}), // invalid filename
			assert: func(result string, err error) {
				Expect(err).IsNotNil()
				Expect(errors.Is(err, os.ErrNotExist)).To(BeFalse()) // may be a different error on different OSes
				Expect(result).To(Equal(""))
			},
		}),

		Case("empty filename", testcase{
			input: "",
			assert: func(result string, err error) {
				Expect(err).Is(as.ErrFilenameIsEmpty)
				Expect(result).To(Equal(""))
			},
		}),

		Case("filename that is a directory", testcase{
			input: ".",
			assert: func(result string, err error) {
				Expect(err).Is(as.ErrFilenameIsDirectory)
				Expect(result).To(Equal(""))
			},
		}),
	))
}
