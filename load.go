package env

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"
)

// Load loads environment variables from one or more files.  Files should be formatted as a list
// of key-value pairs, one per line, separated by an equals sign. Lines that are empty or start
// with a hash (#) are ignored.
//
// # example file format
//
//	# this is a comment
//	NAME1=value1
//	NAME2=value2
//
//	# this is another comment
//	NAME3=value3
//
// # parameters
//
//	files: ...string    // 0..n file path(s)
//
// # returns
//
//	error      // an error that wraps all errors that occurred while loading variables;
//	           // if no errors occurred the result is nil
//
// The joined errors will be in the order that the files were specified and will be wrapped
// with the file path that caused the error:
//
//	"path/to/file: error"
//
// # .env file
//
// The function will always attempt to load variables from a ".env" file.
//
// If ".env" (or "./.env") is included in the files parameter it will be loaded
// in the order specified relative to other files; if the ".env" file does not exist
// an error will be included in the returned error.
//
// If ".env" is not explicitly specified it will be loaded before any other files, if it
// exists; if it does not exist it is ignored without error.
//
// If no files are specified the function will attempt to load variables from ".env"
// and will return an error if the file does not exist.
//
// # example: loading .env:
//
//	if err := Load(); err != nil {
//		log.Fatal(err) // possibly because .env does not exist
//	}
//
// # example: loading .env and a specified file:
//
//	if err := Load("test.env"); err != nil {
//		log.Fatal(err) // will not be because .env did not exist; could be because test.env does not exist
//	}
func Load(files ...string) error {
	// determine if ".env" has been explicitly specified and if it is required
	filenames := map[string]bool{}
	for _, f := range files {
		filenames[f] = true
	}
	dotenvRequired := len(filenames) == 0 || (filenames[".env"] || filenames["./.env"])

	// if ".env" has not been explicitly specified we will load it before loading
	// any other files
	if !filenames["./.env"] && !filenames[".env"] {
		files = append([]string{".env"}, files...)
	}

	// we will be collecting any errors that occur while loading the files
	errs := []error{}

	for _, filename := range files {
		err := loadFile(filename)
		if err == nil {
			continue
		}
		if !dotenvRequired && (filename == ".env" || filename == "./.env") && errors.Is(err, fs.ErrNotExist) {
			continue
		}
		errs = append(errs, fmt.Errorf("%s: %w", filename, err))
	}

	return errors.Join(errs...)
}

// loadFile loads environment variables from a file. The file should be formatted as a list of
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
func loadFile(path string) error {
	file, err := newFileReader(path)
	if err != nil {
		return err
	}
	defer file.Close()

	errs := []error{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || line[0] == '#' {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		vname := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		errs = append(errs, os.Setenv(vname, value))
	}
	return errors.Join(errs...)
}
