package as

import (
	"fmt"
	"math"
	"strconv"
	"strings"

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
	handleError := func(err error) (int, error) {
		return 0, fmt.Errorf("as.PortNo: %w", err)
	}

	s = strings.TrimSpace(s)
	i, err := strconv.Atoi(s)
	switch {
	case err != nil:
		return handleError(err)
	case i < 0 || i > math.MaxUint16:
		return handleError(env.RangeError[int]{Min: 0, Max: math.MaxUint16})
	}

	return i, nil
}
