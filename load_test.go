package env_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	. "github.com/blugnu/test"

	"github.com/blugnu/env"
	"github.com/blugnu/env/internal"
)

func TestLoad(t *testing.T) {
	With(t)

	const cTestFilename = "test-filename"

	type testcase struct {
		fileContent map[string]any // filename: string (content) or error
		filenames   []string
		assert      func(error)
	}
	Run(Testcases(
		ForEach(func(tc testcase) {
			// arrange
			defer Restore(Original(&internal.FileExists).ReplacedBy(func(filename string) bool {
				_, ok := tc.fileContent[filename]
				return ok
			}))

			defer Restore(Original(&internal.NewFileReader).ReplacedBy(func(filename string) (io.ReadCloser, error) {
				switch v := tc.fileContent[filename].(type) {
				case error:
					return nil, v
				case string:
					return io.NopCloser(bytes.NewReader([]byte(v))), nil
				default:
					return nil, fmt.Errorf("unsupported file content type: %T", v)
				}
			}))

			// provide each test case with a clean environment including one
			// set variable (PRESET)
			defer env.State().Restore()
			os.Clearenv()
			t.Setenv("PRESET", "true")

			// act
			err := env.Load(tc.filenames...)

			// assert
			tc.assert(err)
		}),

		Case("no files specified/.env does not exist", testcase{
			assert: func(err error) {
				Expect(err).IsNil()
				Expect(os.Environ()).Should(HaveLen(1))
				Expect(os.Getenv("PRESET")).To(Equal("true"))
			},
		}),

		Case("no files specified/.env exists with valid entries", testcase{
			fileContent: map[string]any{
				internal.EnvFile: "VAR1=value1\nVAR2=value2",
			},
			assert: func(err error) {
				Expect(err).IsNil()
				Expect(os.Environ()).To(ContainItems([]string{
					"VAR1=value1",
					"VAR2=value2",
				}))
			},
		}),

		Case("no files specified/.env exists with invalid entries", testcase{
			fileContent: map[string]any{
				internal.EnvFile: "VAR1=value1\nVAR2 invalid",
			},
			assert: func(err error) {
				Expect(err).Is(env.FileError{Filename: internal.EnvFile})
				Expect(err).Is(env.ErrInvalidEntry)
				Expect(os.Environ()).Should(HaveLen(1))
				Expect(os.Getenv("PRESET")).To(Equal("true"))
			},
		}),

		Case("whitespace filename specified/.env exists with valid entries", testcase{
			fileContent: map[string]any{internal.EnvFile: "VAR1=value1"},
			filenames:   []string{"  "},
			assert: func(err error) {
				Expect(err).IsNil()
				Expect(os.Getenv("VAR1")).To(Equal("value1"))
			},
		}),

		Case("attempt to override .env value with other file by specifying .env last", testcase{
			fileContent: map[string]any{
				internal.EnvFile: "VAR1=value from .env",
				cTestFilename:    "VAR1=value from test-filename",
			},
			filenames: []string{cTestFilename, ".env"},
			assert: func(err error) {
				Expect(err).IsNil()
				Expect(os.Getenv("VAR1")).To(Equal("value from .env"))
			},
		}),

		Case("does not override existing environment variables", testcase{
			fileContent: map[string]any{
				internal.EnvFile: "PRESET=false\nVAR1=value1",
			},
			assert: func(err error) {
				Expect(err).IsNil()
				Expect(os.Getenv("PRESET")).To(Equal("true"))
				Expect(os.Getenv("VAR1")).To(Equal("value1"))
			},
		}),

		Case("loads all files or none", testcase{
			fileContent: map[string]any{
				"file1": "VAR1=value1",
				"file2": "=value2", // invalid entry: no variable name
			},
			filenames: []string{"file1", "file2"},
			assert: func(err error) {
				Expect(err).Is(env.FileError{Filename: "file2"})
				Expect(err).Is(env.ErrInvalidEntry)
				Expect(os.Environ()).Should(HaveLen(1))
				Expect(os.Getenv("PRESET")).To(Equal("true"))
			},
		}),
	))

	Run(Test("duplicate and equivalent filenames are loaded only once", func() {
		// arrange
		loaded := []string{}

		defer Restore(Original(&internal.FileExists).ReplacedBy(func(filename string) bool {
			return true
		}))
		defer Restore(Original(&internal.NewFileReader).ReplacedBy(func(filename string) (io.ReadCloser, error) {
			loaded = append(loaded, filename)
			return io.NopCloser(bytes.NewReader([]byte{})), nil
		}))
		defer env.State().Restore()
		os.Clearenv()

		// act
		err := env.Load(".env", ".env", "./.env", "./././.env", "other-file", "other-file")

		// assert
		Expect(err).IsNil()
		Expect(loaded).To(EqualSlice([]string{".env", "other-file"}), "files loaded")
	}))
}

func TestLoadFile(t *testing.T) {
	With(t)

	const cTestFilename = "test-filename"

	type testcase struct {
		reader    io.ReadCloser
		readerErr error
		assert    func(error)
	}
	Run(Testcases(
		ForEach(func(tc testcase) {
			// arrange
			defer Restore(Original(&internal.NewFileReader).ReplacedBy(func(string) (io.ReadCloser, error) {
				return tc.reader, tc.readerErr
			}))

			defer env.State().Restore()
			os.Clearenv()

			// act
			err := env.LoadFile(cTestFilename)

			// assert
			tc.assert(err)
		}),

		Case("file exists and is valid", testcase{
			reader: io.NopCloser(bytes.NewReader([]byte("VAR=is set"))),
			assert: func(err error) {
				Expect(err).Is(nil)
				Expect(os.Getenv("VAR")).To(Equal("is set"))
			},
		}),

		Case("file does not exist", testcase{
			reader:    nil,
			readerErr: os.ErrNotExist,
			assert: func(err error) {
				Expect(err).Is(env.FileError{Filename: cTestFilename})
				Expect(err).Is(os.ErrNotExist)
			},
		}),

		Case("file contains invalid entry", testcase{
			reader: io.NopCloser(bytes.NewReader([]byte("VAR1=valid\nVAR2 invalid"))),
			assert: func(err error) {
				Expect(err).Is(env.FileError{Filename: cTestFilename})
				Expect(err).Is(env.ErrInvalidEntry)

				Run(Test("valid entries are not applied", func() {
					_, ok := os.LookupEnv("VAR1")
					Expect(ok).To(BeFalse())
				}))
			},
		}),
	))
}

func TestLoadFromReader(t *testing.T) {
	With(t)

	type testcase struct {
		content string
		scanErr error
		setErr  error
		assert  func(error)
	}
	Run(Testcases(
		ForEach(func(tc testcase) {
			// arrange
			if tc.scanErr != nil {
				defer Restore(Original(&internal.ScanLines).ReplacedBy(func([]byte, bool) (int, []byte, error) {
					return 0, nil, tc.scanErr
				}))
			}
			if tc.setErr != nil {
				defer Restore(Original(&internal.Setenv).ReplacedBy(func(string, string) error {
					return tc.setErr
				}))
			}

			defer env.State().Restore()
			os.Clearenv()

			src := io.NopCloser(bytes.NewReader([]byte(tc.content)))

			// act
			err := env.LoadFromReader(src)

			// assert
			tc.assert(err)
		}),

		Case("contains empty lines", testcase{
			content: "\nVAR1=value1\n\nVAR2=value2\n",
			assert: func(err error) {
				Expect(err).IsNil()
				Expect(os.Environ()).To(ContainItems([]string{
					"VAR1=value1",
					"VAR2=value2",
				}))
			},
		}),

		Case("variable name and value contains whitespace", testcase{
			content: "VAR1 = value1\nVAR2=value2 \n",
			assert: func(err error) {
				Expect(err).IsNil()
				Expect(os.Environ()).To(ContainItems([]string{
					"VAR1= value1",
					"VAR2=value2 ",
				}))
			},
		}),

		Case("value contains equals sign", testcase{
			content: "VAR1=value=with=equals\n",
			assert: func(err error) {
				Expect(err).IsNil()
				Expect(os.Getenv("VAR1")).To(Equal("value=with=equals"))
			},
		}),

		Case("value contains hash", testcase{
			content: "VAR1=value#with#hash\n",
			assert: func(err error) {
				Expect(err).IsNil()
				Expect(os.Getenv("VAR1")).To(Equal("value#with#hash"))
			},
		}),

		Case("missing variable name", testcase{
			content: "VAR1=value1\n=value2",
			assert: func(err error) {
				Expect(err).Is(env.ErrInvalidEntry)
				Expect(os.Environ()).Should(BeEmpty())
			},
		}),

		Case("setting variable fails", testcase{
			content: "VAR1=value1\nVAR2=value2\n",
			setErr:  os.ErrInvalid,
			assert: func(err error) {
				Expect(err).Is(os.ErrInvalid)
				Expect(os.Environ()).Should(BeEmpty())
			},
		}),

		Case("ignores commented lines", testcase{
			content: "# VAR1=value1\nVAR2=value2\n",
			assert: func(err error) {
				Expect(err).IsNil()
				Expect(os.Environ()).Should(HaveLen(1))
				Expect(os.Getenv("VAR2")).To(Equal("value2"))
			},
		}),

		Case("returns any scanner error", testcase{
			content: "VAR1=value1\nVAR2=value2\n",
			scanErr: os.ErrInvalid,
			assert: func(err error) {
				Expect(err).Is(os.ErrInvalid)
				Expect(os.Environ()).Should(BeEmpty())
			},
		}),
	))
}
