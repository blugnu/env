package env

import (
	"os"
	"strings"

	"github.com/blugnu/env/internal"
)

// Vars returns a map containing environment variables.
//
// If no variable names are provided, the map is initialized with all
// variables in the current environment.
//
// Otherwise, the returned map contains entries only for those variables
// that are specified and which are set.
//
// # parameters
//
//	names ...string   // (optional) names of environment variables to return;
//	                  // if no names are provided the returned map contains all
//	                  // current environment variables.
//
// # returns
//
//	map[string]string   // a map of environment variables where the key is the
//	                    // name of the environment variable; keys are trimmed of
//	                    // leading and trailing whitespace; values are not trimmed
func Vars(names ...string) map[string]string {
	var (
		src []string
		fn  func(string) (string, string, bool)
	)

	switch len(names) {
	case 0: // all environment variables
		src = os.Environ()
		fn = func(s string) (string, string, bool) {
			k, v, _ := strings.Cut(s, "=")
			return strings.TrimSpace(k), v, true
		}

	default: // only the named variables (if set)
		src = names
		fn = func(s string) (string, string, bool) {
			k := strings.TrimSpace(s)
			v, isSet := internal.LookupEnv(k)
			return k, v, isSet
		}
	}

	result := make(map[string]string, len(src))
	for _, s := range src {
		if k, v, isSet := fn(s); isSet {
			result[k] = v
		}
	}

	return result
}
