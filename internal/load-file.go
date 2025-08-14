package internal

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"strings"
)

// LoadFile loads environment variables from a file. The file should be formatted as a list of
// key-value pairs, one per line, separated by an equals sign. Lines that are empty or start with
// a hash (#) are ignored.
//
// # parameters
//
//	path: string   // the path to the file to loadFile
//
// # returns
//
//	error          // any error that occurrs while loading or applying variables
func LoadFile(path string) error {
	handleError := func(err error) error {
		return fmt.Errorf("env/internal.LoadFile: %w", err)
	}

	file, err := NewFileReader(path)
	if err != nil {
		return handleError(err)
	}
	defer func() { _ = file.Close() }()

	errs := []error{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || line[0] == '#' {
			continue
		}

		vname, value, ok := strings.Cut(line, "=")
		if !ok {
			errs = append(errs, fmt.Errorf("%w: %q", ErrInvalidEntry, line))
			continue
		}

		// do not replace existing variables in the env unless they are empty
		vname = strings.TrimSpace(vname)
		if value, isSet := LookupEnv(vname); isSet && len(strings.TrimSpace(value)) > 0 {
			continue
		}

		value = strings.TrimSpace(value)
		errs = append(errs, Setenv(vname, value))
	}
	return errors.Join(errs...)
}

func LoadFiles(envFileMustExist bool, files ...string) error {
	// we will be collecting any errors that occur while loading the files
	errs := []error{}

	for _, filename := range files {
		err := LoadFile(filename)
		if err == nil {
			continue
		}
		if !envFileMustExist && (filename == EnvFile || filename == "./"+EnvFile) && errors.Is(err, fs.ErrNotExist) {
			continue
		}
		errs = append(errs, fmt.Errorf("%s: %w", filename, err))
	}

	return errors.Join(errs...)
}
