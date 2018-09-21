# patron [![CircleCI](https://circleci.com/gh/mantzas/patron.svg?style=svg)](https://circleci.com/gh/mantzas/patron) [![codecov](https://codecov.io/gh/mantzas/patron/branch/master/graph/badge.svg)](https://codecov.io/gh/mantzas/patron) [![Go Report Card](https://goreportcard.com/badge/github.com/mantzas/patron)](https://goreportcard.com/report/github.com/mantzas/patron) [![GoDoc](https://godoc.org/github.com/mantzas/patron?status.svg)](https://godoc.org/github.com/mantzas/patron)

Patron is a framework for creating microservices.

`Patron` is french for `template` or `pattern`, but it means also `boss` which we found out later (no pun intended).

The entry point of the framework is the `Service`. The `Service` uses `Components` to handle the processing of sync and async requests. The `Service` starts by default a `HTTP Component` which hosts the debug, health and metric endpoints. Any other endpoints will be added to the default `HTTP Component` as `Routes`. The service set's up by default logging with `zerolog`, tracing and metrics with `jaeger` and `prometheus`.

`Patron` provides abstractions for the following functionality of the framework:

- service, which orchestrates everything
- components and processors, which provide a abstraction of adding processing functionality to the service
  - asynchronous message processing (RabbitMQ, Kafka)
  - synchronous processing (HTTP)
- metrics and tracing
- logging

`Patron` provides same defaults for making the usage as simple as possible.

## How to Contribute

Please see [CONTRIBUTING.md](CONTRIBUTING.md).

## patron-cli

The framework supplies a cli in order to simplify repository generation with the following features:

- git repository creation
- cmd folder and main.go creation with build version support (`go build -ldflags '-X main.version=1.0.0' main.go`)
- go module support and vendoring
- Dockerfile with version support (`docker build --build-arg version=1.0.0`)

The latest version can be installed with

```go
go get github.com/mantzas/patron/cmd/patron
```

## Service

The `Service` has the role of glueing all of the above together, which are:

- setting up logging
- setting up default HTTP component with the following endpoints configured:
  - profiling via pprof
  - health check
  - info endpoint for returning information about the service
- setting up termination by os signal
- starting and stopping components
- handling component errors
- setting up metrics and tracing

The service has some default settings which can be changed via environment variables:

- Service HTTP port, for setting the default HTTP components port to `50000` with `PATRON_HTTP_DEFAULT_PORT`
- Log level, for setting zerolog with `INFO` log level with `PATRON_LOG_LEVEL`
- Tracing, for setting up jaeger tracing with
  - agent address `0.0.0.0:6831` with `PATRON_JAEGER_AGENT`
  - sampler type `probabilistic`with `PATRON_JAEGER_SAMPLER_TYPE`
  - sampler param `0.1` with `PATRON_JAEGER_SAMPLER_PARAM`

### Component

A `Component` is a interface that exposes the following API:

```go
type Component interface {
  Run(ctx context.Context) error  
}
```

The above API gives the `Service` the ability to start and gracefully shutdown a `component` via context cancellation. The framework divides the components in 2 categories:

- synchronous, which are components that follow the request/response pattern and
- asynchronous, which consume messages from a source but don't respond anything back

The following component implementations are available:

- HTTP (sync)
- RabbitMQ consumer (async)
- Kafka consumer (async)

Adding to the above list is as easy as implementing a `Component` and a `Processor` for that component.

## Example

Setting up a new service with a HTTP `Component` is as easy as the following code:

```go
  // Set up HTTP routes
  routes := make([]sync_http.Route, 0)
  routes = append(routes, sync_http.NewRoute("/", http.MethodGet, processor, true))
  
  srv, err := patron.New("test", patron.Routes(routes))
  if err != nil {
    log.Fatalf("failed to create service %v", err)
  }

  err = srv.Run()
  if err != nil {
    log.Fatalf("failed to create service %v", err)
  }
```

The above is pretty much self-explanatory. The processor follows the sync pattern.

## Processors

### Synchronous

The implementation of the processor is responsible to create a `Request` by providing everything that is needed (Headers, Fields, decoder, raw io.Reader) pass it to the implementation by invoking the `Process` method and handle the `Response` or the `error` returned by the processor.

The sync package contains only a function definition along with the models needed:

```go
type ProcessorFunc func(context.Context, *Request) (*Response, error)
```

The `Request` model contains the following properties (which are provided when calling the "constructor" `NewRequest`)

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
type ProcessorFunc func(context.Context, *Message) error
```

Everything else is exactly the same.

## Metrics and Tracing

Tracing and metrics are provided by jaeger's implementation of the OpenTracing project.
Every component has been integrated with the above library and produces traces and metrics.
Metrics are provided with the default HTTP component at the `/metrics` route for Prometheus to scrape.
Tracing will be send to a jaeger agent which can be setup though environment variables mentioned in the config section. Sane defaults are applied for making the use easy.
We have included some clients inside the trace package which are instrumented and allow propagation of tracing to
downstream systems. The tracing information is added to each implementations header. These clients are:

- HTTP
- AMQP
- Kafka

## Logging

The log package is designed to be a leveled logger with field support.

The log package defines the logger interface and a factory function type that needs to be implemented in order to set up the logging in this framework.

```go
  // instantiate the implemented factory func type and fields (map[string]interface{})
  err := log.Setup(factory, fields)
  // handle error
```

`If the setup is omitted the package will not setup any logging!`

From there logging is as simple as

```go
  log.Info("Hello world!")
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

The following implementations are provided as sub-package and are by default wired up in the framework:

- zerolog, which supports the excellent [zerolog](https://github.com/rs/zerolog) library and is set up by default

### Logger

The logger interface defines the actual logger.

```go
type Logger interface {
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

The factory function type defines a factory for creating a logger.

```go
type FactoryFunc func(map[string]interface{}) Logger
```