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

- [ ] **.env File Support**: Load variables from a `.env` file and/or any other file(s)
- [ ] **Type Conversions**: Easily convert environment variables to Go types
- [ ] **Validation**: Check common configuration errors (e.g. `as.PortNo` to enforce 0 <= X <= 65535)
- [ ] **Testing**: Convenient testing utilities

## Installation

```bash
go get github.com/blugnu/env
```

## Example Usage

### Override Default Configuration

Demonstrates the use of the `env.Override` function to replace a default
configuration value with a value parsed from an environment variable:

```go
    port := 8080
    if _, err := env.Override(&port, "SERVICE_PORT", as.PortNo); err != nil {
        log.Fatal(err)
    }
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
```

### Parse a Required Configuration Value

Demonstrates the use of the `env.Parse` function to parse a required
configuration value from an environment variable:

```go
    authURL, err := env.Parse("AUTH_SERVICE_URL", as.AbsoluteURL)
    if err != nil {
        log.Fatal(err)
    }
```

### Load Configuration from a File

Demonstrates the use of the `env.Load` function to load configuration:

> by default, with no filename(s) specified, the `Load()` function
> loads configuration from a `.env` file.

```go
    if err := env.Load(); err != nil {
        log.Fatal(err)
    }
```

### Preserve Environment Variables in a Test

Demonstrates the use of `defer env.State().Reset()` to preserve environment
variables during a test:

```go
    func TestSomething(t *testing.T) {
        // ARRANGE
        defer env.State().Reset()
        env.Vars{
            "SOME_VAR": "some value",
            "ANOTHER_VAR": "another value",
        }.Set()

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
