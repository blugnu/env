package internal

import (
	"bufio"
	"errors"
	"io"
	"os"
)

//nolint:gochecknoglobals // test seams
var (
	// Open is an alias for [os.Open]; it provides a test seam for mocking
	Open = os.Open

	// LookupEnv is an alias for [os.LookupEnv]; it provides a test seam for mocking
	LookupEnv = os.LookupEnv

	// ScanLines is an alias for [bufio.ScanLines]; it provides a test seam for mocking
	ScanLines = bufio.ScanLines

	// Setenv is an alias for [os.Setenv]; it provides a test seam for mocking
	Setenv = os.Setenv

	// Stat is an alias for [os.Stat]; it provides a test seam for mocking
	Stat = os.Stat

	// Unsetenv is an alias for [os.Unsetenv]; it provides a test seam for mocking
	Unsetenv = os.Unsetenv
)

// FileExists returns true if a file with the given filename exists.
// Note that this function returns false for directories and may return
// true for a file that exists but is not accessible.  Access errors should
// be handled by the caller when they occur upon attempting to open the file.
//
// NOTE: This function variable is a wrapper around os.Stat that provides a
// test seam for mocking in tests
var FileExists = func(filename string) bool {
	fi, err := Stat(filename)
	switch {
	case errors.Is(err, os.ErrNotExist):
		return false
	case fi != nil:
		return !fi.IsDir()
	default:
		return false
	}
}

// NewFileReader returns a new file reader for the given path.
//
// This function is a wrapper around os.Open that returns an io.ReadCloser,
// making it easier to substitute fake readers in tests instead of having to
// use test files or create temporary files:
//
// e.g.
//
//	defer Restore(Original(&internal.NewFileReader).ReplacedBy(func(string) (io.ReadCloser, error) {
//	    return io.NopCloser(bytes.NewReader([]byte(content))), nil
//	}))
var NewFileReader = func(path string) (io.ReadCloser, error) {
	return Open(path)
}
