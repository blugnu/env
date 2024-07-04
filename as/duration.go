package as

import (
	"time"
)

// Duration parses a string into a time.Duration. If a unit is provided, the string is first
// converted to an integer and then multiplied by the unit.
//
// Attempting to convert a duration string using a specified unit will fail due to the
// duration string contain non-numeric characters.
//
// # parameters
//
//	s string             // the string to convert
//
//	u ...time.Duration   // the unit to multiply the duration by; if no unit is provided
//	                     // the string is parsed as a duration.  If multiple units are
//	                     // provided, only the first is used
//
// # returns
//
//	time.Duration   // the converted value
//
//	error           // any error that occurs during conversion
//
// # example: parse a duration string
//
//	d, err := as.Duration("1h30m")
//
// # example: parse a duration string with a unit
//
//	d, err := as.Duration("1", time.Hour)
//
// # example: parse a duration string with a unit
//
//	// this will fail, returning an integer conversion error
//	d, err := as.Duration("1h", time.Hour)
func Duration(s string, u ...time.Duration) (time.Duration, error) {
	if len(u) == 0 {
		return time.ParseDuration(s)
	}
	i, err := Int(s)
	if err != nil {
		return 0, err
	}
	return time.Duration(i) * u[0], nil
}
