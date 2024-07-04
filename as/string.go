package as

// String converts a string to a string. Effectively a NO-OP conversion,
// this function is provided for completeness and to allow for consistent
// use of the `as` package with `env.Parse` and `env.Override` functions.
//
// # parameters
//
//	s string   // the string to 'convert'
//
// # returns
//
//	string     // the 'converted' value (i.e. the input string)
func String(s string) (string, error) {
	return s, nil
}
