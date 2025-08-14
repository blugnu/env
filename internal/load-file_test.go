package internal_test

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/blugnu/test"

	"github.com/blugnu/env"
	"github.com/blugnu/env/internal"
)

const testEnv = "test.env"

func fakeFileReader(content string) internal.FileReaderCloser {
	return io.NopCloser(bytes.NewReader([]byte(content)))
}

func TestLoadFile_WhenFileDoesNotExist(t *testing.T) {
	// ARRANGE
	defer env.State().Reset()
	defer test.Using(&internal.NewFileReader, func(string) (internal.FileReaderCloser, error) {
		return nil, os.ErrNotExist
	})()

	// ACT
	err := internal.LoadFile(testEnv)

	// ASSERT
	test.Error(t, err).Is(os.ErrNotExist)
}

func TestLoadFile_WithInvalidEntry(t *testing.T) {
	// ARRANGE
	defer env.State().Reset()
	defer test.Using(&internal.NewFileReader, func(string) (internal.FileReaderCloser, error) {
		return fakeFileReader("VAR1=valid\nINVALID_LINE"), nil
	})()

	// ACT
	err := internal.LoadFile(testEnv)

	// ASSERT
	test.Error(t, err).Is(internal.ErrInvalidEntry)
	test.That(t, os.Getenv("VAR1")).Equals("valid")
}

func TestLoadFile_WithEmptyLinesAndComments(t *testing.T) {
	// ARRANGE
	defer env.State().Reset()
	defer test.Using(&internal.NewFileReader, func(string) (internal.FileReaderCloser, error) {
		return fakeFileReader("VAR1=value-1\n\n# comment\nVAR2=value-2=with-equals"), nil
	})()

	// ACT
	err := internal.LoadFile(testEnv)

	// ASSERT
	test.That(t, err).IsNil()
	test.That(t, os.Getenv("VAR1")).Equals("value-1")
	test.That(t, os.Getenv("VAR2")).Equals("value-2=with-equals")
}

func TestLoadFile_DoesNotReplaceNonEmptyVariables(t *testing.T) {
	// ARRANGE
	defer env.State().Reset()
	defer test.Using(&internal.NewFileReader, func(string) (internal.FileReaderCloser, error) {
		return fakeFileReader("VAR1=new-value"), nil
	})()

	t.Setenv("VAR1", "pre-existing-value")

	// ACT
	err := internal.LoadFile(testEnv)

	// ASSERT
	test.That(t, err).IsNil()
	test.That(t, os.Getenv("VAR1")).Equals("pre-existing-value")
}

func TestLoadFile_DoesReplaceEmptyVariables(t *testing.T) {
	// ARRANGE
	defer env.State().Reset()
	defer test.Using(&internal.NewFileReader, func(string) (internal.FileReaderCloser, error) {
		return fakeFileReader("VAR1=new-value"), nil
	})()

	t.Setenv("VAR1", "")

	// ACT
	err := internal.LoadFile(testEnv)

	// ASSERT
	test.That(t, err).IsNil()
	test.That(t, os.Getenv("VAR1")).Equals("new-value")
}

func TestLoadFiles(t *testing.T) {
	// ARRANGE
	defer env.State().Reset()
	defer test.Using(&internal.NewFileReader, func(string) (internal.FileReaderCloser, error) {
		return fakeFileReader("VAR1=value-1"), nil
	})()

	// ACT
	err := internal.LoadFiles(true, ".env")

	// ASSERT
	test.Error(t, err).IsNil()
	test.That(t, os.Getenv("VAR1")).Equals("value-1")
}

func TestLoadFiles_WhenEnvFileIsRequiredAndDoesNotExist(t *testing.T) {
	// ARRANGE
	defer env.State().Reset()
	defer test.Using(&internal.NewFileReader, func(string) (internal.FileReaderCloser, error) {
		return nil, os.ErrNotExist
	})()

	// ACT
	err := internal.LoadFiles(true, ".env")

	// ASSERT
	test.Error(t, err).Is(os.ErrNotExist)
}

func TestLoadFiles_WhenEnvFileIsOptionalAndDoesNotExist(t *testing.T) {
	// ARRANGE
	defer env.State().Reset()
	defer test.Using(&internal.NewFileReader, func(string) (internal.FileReaderCloser, error) {
		return nil, os.ErrNotExist
	})()

	// ACT
	err := internal.LoadFiles(false, ".env")

	// ASSERT
	test.That(t, err).IsNil()
}

func TestLoadFiles_WhenLoadFileFails(t *testing.T) {
	// ARRANGE
	defer env.State().Reset()
	defer test.Using(&internal.NewFileReader, func(string) (internal.FileReaderCloser, error) {
		return fakeFileReader("INVALID_ENTRY"), nil
	})()

	// ACT
	err := internal.LoadFiles(true, ".env")

	// ASSERT
	test.Error(t, err).Is(internal.ErrInvalidEntry)
}
