package as

import (
	"fmt"
	"strconv"
	"strings"
)

// Bool converts a string to an Boolean.  Conversion is initially attempted
// using [strconv.ParseBool].  If that fails, the value is checked for
// "y"/"yes" or "n"/"no" (case-insensitive).
//
// If the value is not supported by [strconv.ParseBool] and is not one of
// the additional values above, an error is returned.
func Bool(s string) (bool, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	b, err := strconv.ParseBool(s)
	if err != nil {
		switch s {
		case "n", "no":
			return false, nil
		case "y", "yes":
			return true, nil
		}
		return false, fmt.Errorf("as.Bool: %w", err)
	}
	return b, nil
}
