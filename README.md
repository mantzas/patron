# patron [![Build Status](https://travis-ci.org/mantzas/patron.svg?branch=master)](https://travis-ci.org/mantzas/patron) [![codecov](https://codecov.io/gh/mantzas/patron/branch/master/graph/badge.svg)](https://codecov.io/gh/mantzas/patron) [![Go Report Card](https://goreportcard.com/badge/github.com/mantzas/patron)](https://goreportcard.com/report/github.com/mantzas/patron) [![GoDoc](https://godoc.org/github.com/mantzas/patron?status.svg)](https://godoc.org/github.com/mantzas/patron)

Patron is a framework for creating microservices.

Patron is french for `template` or `pattern`, but it means also `boss` which we found out later (no pun intended).

Patron provides abstractions for the following functionality:

- config
- logging
- metric
- tracing (TBD)
- service
  - async message processing
  - sync processing (http)
- server

## Config

The config package defines a interface that has to be implemented in order to be used inside the application.

```go
type Config interface {
  Set(key string, value interface{}) error
  Get(key string) (interface{}, error)
  GetBool(key string) (bool, error)
  GetInt64(key string) (int64, error)
  GetString(key string) (string, error)
  GetFloat64(key string) (float64, error)
}
```

After implementing the interface a instance has to be provided to the `Setup` method of the package in order to be used directly from the package eg `config.GetBool()`.

The following implementations are provided as sub-packages:

- env, support for env files and env vars
- viper, support for the excellent [viper](https://github.com/spf13/viper) library

### env

The env package supports getting env vars from the system. It allows further to provide a file that contain env vars, separated by a `=`, which are set up on the environment. In order to setup config just do the following:

```go
c,err := env.New({reader to the config file})
// error checking
config.Setup(c)
```

### viper

[Viper](https://github.com/spf13/viper) is a much more complete configuration library. Check the documentation to see what's available. In order to setup config just do the following:

```go
// setup viper by using the library
c,err := viper.New({reader to the config file})
// error checking
config.Setup(c)
```

## Logging