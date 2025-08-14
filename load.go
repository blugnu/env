package env

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/blugnu/env/internal"
)

// Load loads environment variables from 0 or more files. Files are loaded in the order
// specified and will be loaded only once even if specified multiple times. Files that do not
// exist are ignored without error.
//
// The function will not replace or overwrite a variable already present in the environment
// (even if its value is empty).
//
// All errors encountered while loading files are returned as a single error value using
// [errors.Join].  If any error occurs while loading files, the environment is restored to
// its original state.
//
// The contents of any loaded file must be formatted as a list of key-value pairs, one
// per line, separated by an equals sign. Lines that are empty or start with a hash (#)
// are ignored.
//
// # .env File
//
// The .env file is a special file that is always loaded first (if it exists) whether it is
// included in the list of filenames or not.
func Load(filenames ...string) error {
	initialState := State()

	cleaned := make([]string, 0, len(filenames))
	seen := make(map[string]struct{}, len(filenames))

	// clean up filenames: remove empty strings, duplicates or files that
	// don't exist, ensuring that .env is loaded first
	for _, filename := range append([]string{internal.EnvFile}, filenames...) {
		if filename = strings.TrimSpace(filename); filename == "" {
			continue
		}

		filename = filepath.Clean(filename)
		if _, ok := seen[filename]; ok {
			continue
		}
		seen[filename] = struct{}{}

		if !internal.FileExists(filename) {
			continue
		}

		cleaned = append(cleaned, filename)
	}

	if len(cleaned) == 0 {
		return nil
	}

	errs := make([]error, len(cleaned))
	for i, file := range cleaned {
		errs[i] = LoadFile(file)
	}

	if err := errors.Join(errs...); err != nil {
		initialState.Restore()
		return fmt.Errorf("env.Load: %w", err)
	}

	return nil
}

// LoadFile loads environment variables from a file. Actual loading is performed by
// the [LoadFromReader] function.
//
// # parameters
//
//	filename: string   // the filename to load
//
// # returns
//
//	error
func LoadFile(filename string) error {
	handleError := func(err error) error {
		return NewFileError(filename, err)
	}

	reader, err := internal.NewFileReader(filename)
	if err != nil {
		return handleError(err)
	}
	defer func() { _ = reader.Close() }()

	if err := LoadFromReader(reader); err != nil {
		return handleError(err)
	}

	return nil
}

// LoadFromReader loads environment variables from an [io.Reader]. The reader
// should provide a list of key-value pairs, one per line, separated by an
// equals sign. Lines that are empty, all whitespace or start with a '#'
// character are ignored.
//
// If any error occurs during loading, the environment is restored to the
// state it was in before the function was called.  i.e. all values from
// the reader are successfully loaded or none at all.
func LoadFromReader(reader io.Reader) error {
	initialState := State()

	handleError := func(err error) error {
		initialState.Restore()
		return fmt.Errorf("env.LoadFromReader: %w", err)
	}

	errs := []error{}
	scanner := bufio.NewScanner(reader)
	scanner.Split(internal.ScanLines)

	// allow longer lines than the default (~64KiB)
	//
	// this is primarily to cater for the use of environment
	// variable files using base64 encoded certificates and keys,
	//
	// NOTE: use of environment variables for such things is generally
	// discouraged in favor of more robust secrets management solutions,
	// but may be used in non-production or constrained environments
	const initialCap = 64 * 1024   // 64Kb
	const maxCap = 4 * 1024 * 1024 // 4Mb
	scanner.Buffer(make([]byte, 0, initialCap), maxCap)

	for scanner.Scan() {
		line := scanner.Text()
		if trimmed := strings.TrimSpace(line); trimmed == "" || trimmed[0] == '#' {
			continue
		}

		vname, value, ok := strings.Cut(line, "=")
		if !ok {
			errs = append(errs, fmt.Errorf("%w: %q", ErrInvalidEntry, line))
			continue
		}

		// variable names are trimmed (values are not)
		if vname = strings.TrimSpace(vname); vname == "" {
			errs = append(errs, fmt.Errorf("%w: no variable name: %q", ErrInvalidEntry, line))
			continue
		}

		// do not replace existing variables (even if empty)
		if _, isSet := internal.LookupEnv(vname); isSet {
			continue
		}

		if err := internal.Setenv(vname, value); err != nil {
			errs = append(errs, fmt.Errorf("%w: %w", ErrSetFailed, err))
		}
	}

	errs = append(errs, scanner.Err())

	if err := errors.Join(errs...); err != nil {
		return handleError(err)
	}

	return nil
}
