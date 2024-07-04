package as

import (
	"net/url"

	"github.com/blugnu/env"
)

// AbsoluteURL converts a string to an absolute URI. An absolute URI is a URI
// that includes a scheme, host, and path. The function returns an error if the
// string is not an absolute URI.
//
// # parameters
//
//	s string   // the string to convert
//
// # returns
//
//	*url.URL   // the converted value
//
//	error      // any error that occurs during conversion
//
// # errors
//
//   - if the string cannot be converted to a URL the function returns the
//     conversion error
//
//   - if the URL is not an absolute URI, the function returns an
//     env.InvalidValueError
func AbsoluteURL(s string) (*url.URL, error) {
	u, err := urlParse(s)
	if err != nil {
		return nil, err
	}
	if !u.IsAbs() {
		return nil, env.InvalidValueError{Value: s, Err: ErrNotAnAbsoluteURL}
	}
	return u, nil
}
