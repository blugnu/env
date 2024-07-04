package as

import (
	"strconv"
)

// Int converts a string to an integer.
//
// # parameters
//
// 	s string   // the string to convert
//
// # returns
//
// 	int        // the converted value
//
// 	error      // any error that occurs during conversion
func Int(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return i, nil
}
