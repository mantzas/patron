# patron [![Build Status](https://travis-ci.org/mantzas/patron.svg?branch=master)](https://travis-ci.org/mantzas/patron) [![codecov](https://codecov.io/gh/mantzas/patron/branch/master/graph/badge.svg)](https://codecov.io/gh/mantzas/patron) [![Go Report Card](https://goreportcard.com/badge/github.com/mantzas/patron)](https://goreportcard.com/report/github.com/mantzas/patron) [![GoDoc](https://godoc.org/github.com/mantzas/patron?status.svg)](https://godoc.org/github.com/mantzas/patron)

Patron is a framework for creating microservices.

Patron is french for `template` or `pattern`, but it means also `boss` which we found out later (no pun intended).

Patron provides abstractions for the following functionality of the framework:

- configuration
- logging
- metrics and tracing
- components and processors
  - asynchronous message processing (RabbitMQ, Kafka)
  - synchronous processing (HTTP)
- service

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

The service has some default settings tha can be changed via environment variables:

- Service HTTP port, which set's up the default HTTP components port to `50000` which can be changed via the `PATRON_HTTP_DEFAULT_PORT`
- Log level, which set's up zerolog with `INFO` as log level which can be changed via the `PATRON_LOG_LEVEL`
- Tracing, which set's up jaeger tracing with
  - agent address `0.0.0.0:6831`, which can be changed via `PATRON_JAEGER_AGENT`
  - sampler type `probabilistic`, which can be changed via `PATRON_JAEGER_SAMPLER_TYPE`
  - sampler param `0.1`, which can be changed via `PATRON_JAEGER_SAMPLER_PARAM`

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

## Metrics and Tracing

Tracing and metrics are provided by jaeger's implementation of the OpenTracing project.
Every component has been integrated with the above library and produces traces and metrics.
Metrics are provided with the default HTTP component at the `/metrics` route for Prometheus to scrape.
Tracing will be send to a jaeger agent which can be setup though environment variables mentioned in the config section.

## Processors

### Synchronous

The implementation of the processor is responsible to create a `Request` by providing everything that is needed (Headers, Fields, decoder, raw io.Reader) pass it to the implementation by invoking the `Process` method and handle the `Response` or the `error` returned by the processor.

The sync processor package contains only a interface definition of the processor along the models needed:

```go
type Processor interface {
  Process(context.Context, *Request) (*Response, error)
}
```

The `Request` model contains the following properties (which are provided when calling the "constructor" `NewRequest`)

- Headers, which may contains any headers associated with the request
- Fields, which may contain any fields associated with the request
- Raw, the raw request data (if any) in the form of a `io.Reader`
- decode, which is a function of type `encoding.Decode` that decodes the raw reader

A exported function exists for decoding the raw io.Reader in the form of

```go
Decode(v interface{}) error
```

The `Response` model contains the following properties (which are provided when calling the "constructor" `NewResponse`)

- Payload, which may hold a struct of type `interface{}`

### Asynchronous

The implementation of the async processor follows exactly the same principle as the sync processor.
The main difference is that:

- The `Request` is the `Message` and contains only data as `[]byte`
- There is no `Response`, so the processor may return a error

```go
type Processor interface {
  Process(context.Context, *Message) error
}
```

Everything else is exactly the same.

## Service

The `Service` has the role of glueing all of the above together, which are:

- setting up logging
- setting up termination by user
- starting and stopping components
- handling component errors

### Component

A `Component` is a interface that exposes the following API:

```go
type Component interface {
  Run(ctx context.Context) error
  Shutdown(ctx context.Context) error
}
```

The above API gives the `Service` the control over a component in order to start and stop it gracefully. The framework divides the components in 2 categories:

- synchronous, which are components that follow the request/response pattern and
- asynchronous, which consume messages from a source but don't respond anything back

The following component implementations are available:

- HTTP (sync)
- RabbitMQ (async)
- Kafka (async)

Adding to the above list is as easy as implementing a `Component` and a `Processor` for that component.

## Example

Setting up a new service with a HTTP `Component` is as easy as the following code:

```go
  // Set up HTTP routes
  routes := make([]sync_http.Route, 0)
  routes = append(routes, sync_http.NewRoute("/", http.MethodGet, process, true))
  
  srv, err := patron.New("test", patron.Routes(routes))
  if err != nil {
    log.Fatalf("failed to create service %v", err)
  }

  err = srv.Run()
  if err != nil {
    log.Fatalf("failed to create service %v", err)
  }
```

The above is pretty much self-explanatory.