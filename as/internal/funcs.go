package internal

import "net/url"

var (
	// ParseURL is an alias for [url.Parse]; it provides a test seam for mocking
	ParseURL = url.Parse
)
