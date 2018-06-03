# patron [![Build Status](https://travis-ci.org/mantzas/patron.svg?branch=master)](https://travis-ci.org/mantzas/patron) [![codecov](https://codecov.io/gh/mantzas/patron/branch/master/graph/badge.svg)](https://codecov.io/gh/mantzas/patron) [![Go Report Card](https://goreportcard.com/badge/github.com/mantzas/patron)](https://goreportcard.com/report/github.com/mantzas/patron) [![GoDoc](https://godoc.org/github.com/mantzas/patron?status.svg)](https://godoc.org/github.com/mantzas/patron)

Patron is a framework for creating microservices.

Patron is french for `template` or `pattern`, but it means also `boss` which we found out later (no pun intended).

Patron provides abstractions for the following functionality:

- config
- logging
- metrics and tracing (TBD)
- service
  - async message processing (RabbitMQ, Kafka)
  - sync processing (HTTP)
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

### env

The env package supports getting env vars from the system. It allows further to provide a file that contain env vars, separated by a equal sign `=`, which are then set up on the environment. In order to setup config just do the following:

```go
c,err := env.New({reader to the config file})
// error checking
config.Setup(c)
```

## Logging

The log package is designed to be a leveled logger with field support.

The log package defines two interfaces (Logger and Factory) that have to be implemented in order to set up the logging in this framework. After implementing the two interfaces you can setup logging by doing the following:

```go
  // instantiate the implemented factory f
  err := log.Setup(f)
  // handle error
```

The implementations should support following log levels:

- Debug, which should log the message with debug level
- Info, which should log the message with info level
- Warn, which should log the message with warn level
- Error, which should log the message with error level
- Panic, which should log the message with panic level and panics
- Fatal, which should log the message with fatal level and terminates the application

The first four (Debug, Info, Warn and Error) give the opportunity to differentiate the messages by severity. The last two (Panic and Fatal) do the same and do additional actions (panic and termination).

The package supports fields, which are logged along with the message, to augment the information further to ease querying in the log management system.

The following implementations are provided as sub-package:

- zerolog, which supports the excellent [zerolog](https://github.com/rs/zerolog) library

### Logger

The logger interface defines the actual logger.

```go
type Logger interface {
  Level() Level
  Fields() map[string]interface{}
  Fatal(...interface{})
  Fatalf(string, ...interface{})
  Panic(...interface{})
  Panicf(string, ...interface{})
  Error(...interface{})
  Errorf(string, ...interface{})
  Warn(...interface{})
  Warnf(string, ...interface{})
  Info(...interface{})
  Infof(string, ...interface{})
  Debug(...interface{})
  Debugf(string, ...interface{})
}
```

In order to be consistent with the design the implementation of the `Fatal(f)` have to terminate the application with an error and the `Panic(f)` need to panic.

### Factory

The factory interface defines a factory for creating a logger.

```go
type Factory interface {
  Create(map[string]interface{}) Logger
  CreateSub(Logger, map[string]interface{}) Logger
}
```

Two methods are supported:

- Create, which creates a logger with the specified fields (or nil)
- CreateSub, which creates a sub-logger that accepts a logger and fields and creates a sub-logger with the fields merged into the new one.

## Metrics and Tracing (TBD)

## Server

## Services

### Sync

### Async
