// Package env provides functions for working with environment variables.
//
// The package provides functions to parse environment variables using
// strongly-typed, generic conversion functions and for loading environment
// variables from files or an [io.Reader].  Several common conversion
// functions are provided in the [as] sub-package.
//
// In addition, simple functions are provided to wrap the standard library
// functions for getting, setting, and unsetting environment variables.
//
// For testing purposes, the package also provides functions to capture the
// current state of the environment and conveniently restore it at the
// completion of a test. This functionality differs from the testing.T
// methods [testing.T.Setenv] and [testing.T.Cleanup] in that it captures
// the entire environment at a point in time to be restored later.  This
// allows tests to be structured in a way that ensures a clean environment
// for each sub-test.  By contrast, the testing.T methods only provide
// for changes made to the environment to be automatically reverted.
//
// SPDX-License-Identifier: MIT
// This file is intentionally left blank aside from the package comment.
package env
