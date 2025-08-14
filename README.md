<div align="center" style="margin-bottom:20px">
  <img src=".assets/banner.png" alt="env" />
  <!-- <hr> -->
  <div align="center">
  <h3>streamline and simplify the way you work with environment variables</h3>
  </div>
  <hr>
  <div align="center">
    <a href="https://github.com/blugnu/env/actions/workflows/release.yml">
      <img alt="build-status" src="https://github.com/blugnu/env/actions/workflows/release.yml/badge.svg"/>
    </a>
    <a href="https://goreportcard.com/report/github.com/blugnu/env" >
      <img alt="go report" src="https://goreportcard.com/badge/github.com/blugnu/env"/>
    </a>
    <a>
      <img alt="go version >= 1.14" src="https://img.shields.io/github/go-mod/go-version/blugnu/env?style=flat-square"/>
    </a>
    <a href="https://github.com/blugnu/env/blob/master/LICENSE">
      <img alt="MIT License" src="https://img.shields.io/github/license/blugnu/env?color=%234275f5&style=flat-square"/>
    </a>
    <a href="https://coveralls.io/github/blugnu/env?branch=master">
      <img alt="coverage" src="https://img.shields.io/coveralls/github/blugnu/env?style=flat-square"/>
    </a>
    <a href="https://pkg.go.dev/github.com/blugnu/env">
      <img alt="docs" src="https://pkg.go.dev/badge/github.com/blugnu/env"/>
    </a>
  </div>
</div>

## Features

- [ ] **.env File Support**: Load variables from a `.env` file and/or specified file(s)
- [ ] **Type Conversions**: Safely convert environment variable strings to Go types
- [ ] **Validation**: Use validated conversions to check common configuration
                      errors (e.g. `as.PortNo` to enforce 0 <= X <= 65535)
- [ ] **Extensible**: Implement your own type conversions
- [ ] **Testing**: Convenient testing utilities

## Installation

```bash
go get github.com/blugnu/env
```

## Examples

### Parsing Environment Values

#### 1. Parse an Environment Value with a Default Value

Demonstrates the use of the `env.Parse` function to parse an optional
value from an environment variable, with a default value:

```go
    port := 8080
    port, err := env.Parse("SERVICE_PORT", as.PortNo, port); err != nil {
        log.Fatal(err)
    }
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
```

If the environment variable is set the variable is updated with the parsed value.

If the environment variable is not set, the default value is returned.  In this case
the default value is the variable's original value so the variable is unchanged.

If the environment variable fails to parse, an error is returned.

> :warning: if the variable is set to an empty string,  whether an error is returned or
> the default value is used depends on whether an empty string can be successfully parsed.
> This is unlikely but depends on the specific `as` conversion used.

#### 2. Parse a Required Environment Variable

Demonstrates the use of the `env.Parse` function to parse a required
value (i.e. no default) from an environment variable:

```go
    authURL, err := env.Parse("AUTH_SERVICE_URL", as.AbsoluteURL)
    if err != nil {
        log.Fatal(err)
    }
```

With no default value, an error is returned if the environment variable is not set
or fails to parse.

#### 3. Parse an Environment Variable into a Variable

Demonstrates the use of the `env.ParseInto` function to parse an
environment variable into an existing variable:

```go
    var debug bool
    if err := env.ParseInto(&debug, "DEBUG", as.Bool); err != nil {
        log.Fatal(err)
    }
```

This can simplify error handling code and reduce boilerplate when parsing
multiple variables, e.g.:

```go
    var errs  []error

    errs = append(errs, env.ParseInto(&cfg.Debug, "DEBUG", as.Bool))
    errs = append(errs, env.ParseInto(&cfg.Port, "SERVICE_PORT", as.PortNo, 8080))
    errs = append(errs, env.ParseInto(&cfg.Url, "AUTH_SERVICE_URL", as.AbsoluteURL))

    return errors.Join(errs...)
```

### Get a Map of Environment Variables

Demonstrates the use of the `Vars` function to get a map of environment
variables:

```go
    // get a map containing specific environment variables (if set)
    vars := env.Vars("SERVICE_PORT", "SERVICE_HOSTNAME")
```

### Preserve Environment Variables in a Test

Although `testing.T` provides methods for setting environment variables for the
duration of the current test, this leaves other variables in the environment
unchanged.

To provide a test with a clean, known environment, use the `env.State` function
to obtain the current state of the environment, and `Restore` it at the end of
the test. This allows the environment to be cleared and set for the test using
regular `os` functions as required, without affecting other tests:

```go
    func TestSomething(t *testing.T) {
        // ARRANGE
        defer env.State().Restore()

        os.Clearenv()
        os.Setenv("SOME_VAR", "some value")

        // ACT
        SomeFuncUsingEnvVars()

        // ASSERT
        ...
    }
```

## Contributing

Contributions are welcome! Please feel free to submit a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE file](LICENSE)
for details.
