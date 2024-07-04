package env

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"testing"

	"github.com/blugnu/test"
)

func fakeFile(content string) fileReader {
	return io.NopCloser(bytes.NewReader([]byte(content)))
}

func TestLoad(t *testing.T) {
	// ARRANGE
	testcases := []struct {
		scenario string
		args     []any
		exec     func(t *testing.T)
	}{
		{scenario: "no arguments/.env exists",
			exec: func(t *testing.T) {
				// ARRANGE
				readsDotEnv := false
				defer State().Reset()
				defer test.Using(&newFileReader, func(path string) (fileReader, error) {
					readsDotEnv = readsDotEnv || path == ".env"
					return fakeFile("VAR1=loaded-value-1\nVAR2=loaded-value-2"), nil
				})()
				os.Clearenv()
				os.Setenv("VAR1", "env-value")

				// ACT
				err := Load()

				// ASSERT
				test.That(t, err).IsNil()
				test.IsTrue(t, readsDotEnv)
				test.That(t, os.Getenv("VAR1")).Equals("loaded-value-1")
				test.That(t, os.Getenv("VAR2")).Equals("loaded-value-2")
			},
		},
		{scenario: "no arguments/.env does not exist",
			exec: func(t *testing.T) {
				// ARRANGE
				defer State().Reset()
				defer test.Using(&newFileReader, func(string) (fileReader, error) {
					return nil, os.ErrNotExist
				})()
				os.Clearenv()
				os.Setenv("VAR1", "env-value")

				// ACT
				err := Load()

				// ASSERT
				test.Error(t, err).Is(os.ErrNotExist)
				test.That(t, os.Getenv("VAR1")).Equals("env-value")
			},
		},
		{scenario: ".env argument/.env does not exist",
			exec: func(t *testing.T) {
				// ARRANGE
				defer State().Reset()
				defer test.Using(&newFileReader, func(string) (fileReader, error) {
					return nil, os.ErrNotExist
				})()
				os.Clearenv()
				os.Setenv("VAR1", "env-value")

				// ACT
				err := Load(".env")

				// ASSERT
				test.Error(t, err).Is(os.ErrNotExist)
				test.That(t, os.Getenv("VAR1")).Equals("env-value")
			},
		},
		{scenario: "file path argument/valid file/.env does not exist",
			exec: func(t *testing.T) {
				// ARRANGE
				filesLoaded := []string{}
				defer State().Reset()
				defer test.Using(&newFileReader, func(path string) (fileReader, error) {
					filesLoaded = append(filesLoaded, path)
					switch path {
					case ".env":
						return nil, fs.ErrNotExist
					case "test.env":
						return fakeFile("VAR1=loaded-value-1\nVAR2=loaded-value-2"), nil
					default:
						panic("unexpected file path: " + path)
					}
				})()
				os.Clearenv()
				os.Setenv("VAR1", "env-value")

				// ACT
				err := Load("test.env")

				// ASSERT
				test.That(t, err).IsNil()
				test.Slice(t, filesLoaded).Equals([]string{".env", "test.env"})
				test.That(t, os.Getenv("VAR1")).Equals("loaded-value-1")
				test.That(t, os.Getenv("VAR2")).Equals("loaded-value-2")
			},
		},
		{scenario: "file path argument/valid file/.env exists",
			exec: func(t *testing.T) {
				// ARRANGE
				filesLoaded := []string{}
				defer State().Reset()
				defer test.Using(&newFileReader, func(path string) (fileReader, error) {
					filesLoaded = append(filesLoaded, path)
					switch path {
					case ".env":
						return fakeFile("VAR3=dotenv-value-3"), nil
					case "test.env":
						return fakeFile("VAR1=loaded-value-1\nVAR2=loaded-value-2"), nil
					default:
						panic("unexpected file path: " + path)
					}
				})()
				os.Clearenv()
				os.Setenv("VAR1", "env-value")

				// ACT
				err := Load("test.env")

				// ASSERT
				test.That(t, err).IsNil()
				test.That(t, filesLoaded).Equals([]string{".env", "test.env"})
				test.That(t, os.Getenv("VAR1")).Equals("loaded-value-1")
				test.That(t, os.Getenv("VAR2")).Equals("loaded-value-2")
				test.That(t, os.Getenv("VAR3")).Equals("dotenv-value-3")
			},
		},
		{scenario: "file path argument/valid file/.env error",
			exec: func(t *testing.T) {
				// ARRANGE
				dotenverr := errors.New("error reading .env")
				filesLoaded := []string{}
				defer State().Reset()
				defer test.Using(&newFileReader, func(path string) (fileReader, error) {
					filesLoaded = append(filesLoaded, path)
					switch path {
					case ".env":
						return nil, dotenverr
					case "test.env":
						return fakeFile("VAR1=loaded-value-1\nVAR2=loaded-value-2"), nil
					default:
						panic("unexpected file path: " + path)
					}
				})()
				os.Clearenv()
				os.Setenv("VAR1", "env-value")

				// ACT
				err := Load("test.env")

				// ASSERT
				test.Error(t, err).Is(dotenverr)
				test.That(t, filesLoaded).Equals([]string{".env", "test.env"})
				test.That(t, os.Getenv("VAR1")).Equals("loaded-value-1")
				test.That(t, os.Getenv("VAR2")).Equals("loaded-value-2")
			},
		},
		{scenario: "explicit .env/before other files",
			exec: func(t *testing.T) {
				// ARRANGE
				filesLoaded := []string{}
				defer State().Reset()
				defer test.Using(&newFileReader, func(path string) (fileReader, error) {
					filesLoaded = append(filesLoaded, path)
					switch path {
					case ".env":
						return fakeFile("VAR=dotenv-value"), nil
					case "test.env":
						return fakeFile("VAR=test-value"), nil
					default:
						panic("unexpected file path: " + path)
					}
				})()
				os.Clearenv()
				os.Setenv("VAR", "env-value")

				// ACT
				err := Load(".env", "test.env")

				// ASSERT
				test.That(t, err).IsNil()
				test.That(t, filesLoaded).Equals([]string{".env", "test.env"})
				test.That(t, os.Getenv("VAR")).Equals("test-value")
			},
		},
		{scenario: "explicit .env/after other files",
			exec: func(t *testing.T) {
				// ARRANGE
				filesLoaded := []string{}
				defer State().Reset()
				defer test.Using(&newFileReader, func(path string) (fileReader, error) {
					filesLoaded = append(filesLoaded, path)
					switch path {
					case ".env":
						return fakeFile("VAR=dotenv-value"), nil
					case "test.env":
						return fakeFile("VAR=test-value"), nil
					default:
						panic("unexpected file path: " + path)
					}
				})()
				os.Clearenv()
				os.Setenv("VAR", "env-value")

				// ACT
				err := Load("test.env", ".env")

				// ASSERT
				test.That(t, err).IsNil()
				test.That(t, filesLoaded).Equals([]string{"test.env", ".env"})
				test.That(t, os.Getenv("VAR")).Equals("dotenv-value")
			},
		},
		{scenario: "multiple errors",
			exec: func(t *testing.T) {
				// ARRANGE
				dotenverr := errors.New("error reading .env")
				testenverr := errors.New("error reading test.env")
				defer test.Using(&newFileReader, func(path string) (fileReader, error) {
					switch path {
					case ".env":
						return nil, dotenverr
					case "test.env":
						return nil, testenverr
					default:
						panic("unexpected file path: " + path)
					}
				})()

				// ACT
				err := Load("test.env")

				// ASSERT
				test.Error(t, err).Is(dotenverr)
				test.Error(t, err).Is(testenverr)
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {
			tc.exec(t)
		})
	}
}

func TestLoadFile_WithEmptyLinesAndComments(t *testing.T) {
	// ARRANGE
	defer State().Reset()
	defer test.Using(&newFileReader, func(string) (fileReader, error) {
		return fakeFile("VAR1=value-1\n\n# comment\nVAR2=value-2=with-equals"), nil
	})()

	// ACT
	err := loadFile("test.env")

	// ASSERT
	test.That(t, err).IsNil()
	test.That(t, os.Getenv("VAR1")).Equals("value-1")
	test.That(t, os.Getenv("VAR2")).Equals("value-2=with-equals")
}
