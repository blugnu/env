package as

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Duration parses a string into a time.Duration. The string must contain
// a valid duration expression (e.g., "1h30m"). If the string does not
// contain a valid duration expression, an error is returned with a zero
// duration.
func Duration(s string) (time.Duration, error) {
	s = strings.TrimSpace(s)
	dur, err := time.ParseDuration(s)
	if err != nil {
		return 0, fmt.Errorf("as.Duration: %w: %w", ErrNotADuration, err)
	}
	return dur, nil
}

// DurationIn returns a function that parses a string into a [time.Duration] applying
// specified [time.Duration] units.
//
// If the string fails to parse as a plain integer value, an attempt is made to parse
// it as a duration expression. If this is successful, the parsed duration is returned
// and the specified units are ignored.
//
// If the string does not hold either a simple integer or a valid duration expression,
// an error is returned with a zero duration.
func DurationIn(units time.Duration) func(string) (time.Duration, error) {
	return func(s string) (time.Duration, error) {
		s = strings.TrimSpace(s)
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			// attempting to parse as a duration expression is a fallback,
			// but if that fails then the issue is that the value is not
			// an integer, so any duration parsing error is discarded in
			// favor of the original integer parsing error
			if dur, err := time.ParseDuration(s); err == nil {
				return dur, nil
			}

			return 0, fmt.Errorf("as.DurationIn(%s): %w: %w", units, notAnIntegerOrDurationError{}, err)
		}

		return time.Duration(i) * units, nil
	}
}
