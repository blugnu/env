package env

import (
	"io"
	"os"
)

// function variables to facilitate testing
var (
	osLookupEnv = os.LookupEnv
	osOpen      = os.Open
	osSetenv    = os.Setenv
	osUnsetenv  = os.Unsetenv
)

type fileReader = interface {
	io.Reader
	Close() error
}

var newFileReader = func(path string) (fileReader, error) { return osOpen(path) }
