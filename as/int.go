package as

import (
	"fmt"
	"strconv"
)

// Int converts a string to an integer.
//
// # parameters
//
//	s string   // the string to convert
//
// # returns
//
//	int        // the converted value
//
//	error      // any error that occurs during conversion
func Int(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("as.Int: %w: %w", ErrNotAnInteger, err)
	}
	return i, nil
}
