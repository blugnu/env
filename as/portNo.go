package as

import (
	"strconv"

	"github.com/blugnu/env"
)

// PortNo converts a string to a port number. A port number is an integer in the
// range 0 to 65535.
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
//
// # errors
//
//   - if the string cannot be converted to an integer the function returns the
//     conversion error
//
//   - if the integer is outside the valid range, the function returns an
//     env.RangeError
func PortNo(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	if i < 0 || i > 65535 {
		return 0, env.RangeError[int]{Min: 0, Max: 65535}
	}
	return i, nil
}
