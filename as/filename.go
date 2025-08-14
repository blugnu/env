package as

import (
	"fmt"
	"os"
)

// Filename verifies that the provided value is a valid, existing filename.
func Filename(v string) (string, error) {
	handleError := func(err error) (string, error) {
		return "", fmt.Errorf("as.Filename: %w", err)
	}

	if v == "" {
		return handleError(ErrFilenameIsEmpty)
	}

	fi, err := os.Stat(v)
	switch {
	case err != nil:
		return handleError(err)

	case fi.IsDir():
		return handleError(ErrFilenameIsDirectory)
	}

	return v, nil
}
