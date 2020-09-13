# patron ![CI](https://github.com/beatlabs/patron/workflows/CI/badge.svg) [![codecov](https://codecov.io/gh/beatlabs/patron/branch/master/graph/badge.svg)](https://codecov.io/gh/beatlabs/patron) [![Go Report Card](https://goreportcard.com/badge/github.com/beatlabs/patron)](https://goreportcard.com/report/github.com/beatlabs/patron) [![GoDoc](https://godoc.org/github.com/beatlabs/patron?status.svg)](https://godoc.org/github.com/beatlabs/patron) ![GitHub release](https://img.shields.io/github/release/beatlabs/patron.svg)

Patron is a framework for creating microservices, originally created by Sotiris Mantzaris (https://github.com/mantzas). This fork is maintained by Beat Engineering (https://thebeat.co)

`Patron` is french for `template` or `pattern`, but it means also `boss` which we found out later (no pun intended).

The entry point of the framework is the `Service`. The `Service` uses `Components` to handle the processing of sync and async requests. The `Service` starts by default an `HTTP Component` which hosts the debug, alive, ready and metric endpoints. Any other endpoints will be added to the default `HTTP Component` as `Routes`. Alongside `Routes` one can specify middleware functions to be applied ordered to all routes as `MiddlewareFunc`. The service set's up by default logging with `zerolog`, tracing and metrics with `jaeger` and `prometheus`.

`Patron` provides abstractions for the following functionality of the framework:

- service, which orchestrates everything
- components and processors, which provide an abstraction of adding processing functionality to the service
  - asynchronous message processing (RabbitMQ, Kafka, AWS SQS)
  - synchronous processing (HTTP)
  - gRPC support
- metrics and tracing
- logging

`Patron` provides the same defaults for making the usage as simple as possible.
`Patron` needs Go 1.13 as a minimum.

## How to Contribute

1. **Contributor**: An issue has to be created with a problem description and possible solutions using the github template. The better the problem and solution are described, the easier for the **Curators and Others** to understand it and the faster the process. In case of a bug, steps to reproduce will help a lot.
2. **Curators and Others**: The curators will engage in a discussion about the problem and the possible solution. Others can join the discussion to bring other solutions and insights at any point.
3. **Curators**: After the discussion mentioned above, it will be determined if the proposed solution will be implemented or not. Appropriate tags will be applied to the issue.
4. **Contributor**: The contributor will work as follows:
    - assign the issue to himself
    - create a fork
    - git clone the fork on your machine
    - enable [signing your work](SIGNYOURWORK.md)
    - create a PR. use WIP to mark unfinished work e.g. __WIP: Fixing a bug (fixes #1)__
    - after development has finished remove the WIP if applied
5. **Curators**: The curators will conduct a full code review
6. **Curators**: After at least 2 curators have approved the PR, it will be merged to master
7. **Curators**: A release will follow after that at some point

PR's should have the following requirements:

- Tests are required (where applicable, terms may vary)
  - Unit
  - Component
  - Integration
- High code coverage
- Coding style (go fmt)
- Linting we use [golangci-lint](https://github.com/golangci/golangci-lint)

## Code of conduct

Please note that this project is released with a [Contributor Code of Conduct](https://www.contributor-covenant.org/adopters). By participating in this project and its community you agree to abide by those terms.

## patron-cli

The framework supplies a cli in order to simplify repository generation with the following features:

- git repository creation
- cmd folder and main.go creation with build version support (`go build -ldflags '-X main.version=1.0.0' main.go`)
- go module support and vendoring
- Dockerfile with version support (`docker build --build-arg version=1.0.0`)

The latest version can be installed with

    go get github.com/beatlabs/patron/cmd/patron

Below is an example of a service created with the cli that has a module name `github.com/beatlabs/test` and will be created in the test folder in the current directory.

    patron -m "github.com/beatlabs/test" -p "test"


## Service

The `Service` has the role of glueing all of the above together, which are:

- setting up logging
- setting up default HTTP component with the following endpoints configured:
  - profiling via pprof
  - liveness check
  - readiness check
- setting up termination by os signal
- setting up SIGHUP custom hook if provided by an option
- starting and stopping components
- handling component errors
- setting up metrics and tracing

The service has some default settings which can be changed via environment variables:

- Service HTTP port, for setting the default HTTP components port to `50000` with `PATRON_HTTP_DEFAULT_PORT`
- Service HTTP read and write timeout, for changing the default HTTP read and write timeout to any value, use `PATRON_HTTP_READ_TIMEOUT`, `PATRON_HTTP_WRITE_TIMEOUT` respectively. For acceptable values check [here](https://golang.org/pkg/time/#ParseDuration).
- Log level, for setting zerolog with `INFO` log level with `PATRON_LOG_LEVEL`
- Tracing, for setting up jaeger tracing with
  - agent host `0.0.0.0` with `PATRON_JAEGER_AGENT_HOST`
  - agent port `6831` with `PATRON_JAEGER_AGENT_PORT`
  - sampler type `probabilistic`with `PATRON_JAEGER_SAMPLER_TYPE`
  - sampler param `0.0` with `PATRON_JAEGER_SAMPLER_PARAM`, which means that traces are not initiated here.

### Component

A `Component` is an interface that exposes the following API:

    type Component interface {
      Run(ctx context.Context) error  
    }

The above API gives the `Service` the ability to start and gracefully shutdown a `component` via context cancellation. Furthermore, the component describes itself by implementing the `Info` method and thus giving the service the ability to report the information of all components. The framework divides the components in 2 categories:

- synchronous, which are components that follow the request/response pattern and
- asynchronous, which consume messages from a source but don't respond anything back

The following component implementations are available:

- HTTP (sync)
- gRPC
- RabbitMQ consumer (async)
- Kafka consumer (async)
- AWS SQS (async)

Adding to the above list is as easy as implementing a `Component` and a `Processor` for that component.

### HTTP Middlewares

A `MiddlewareFunc` preserves the default net/http middleware pattern.
You can create new middleware functions and pass them to Service to be chained on all routes in the default Http Component.

    type MiddlewareFunc func(next http.Handler) http.Handler
    
    // Setup a simple middleware for CORS
    newMiddleware := func(h http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Add("Access-Control-Allow-Origin", "*")
            // Next
            h.ServeHTTP(w, r)
        })
    }

### gRPC

On the server side, the gRPC component injects a `UnaryInterceptor` which handles tracing and log propagation.
On the client side, we inject a `UnaryInterceptor` which handles tracing and log propagation.

## Examples

Detailed examples can be found in the [examples](/examples) folder with the following components involved:

- [HTTP Component, HTTP Tracing, HTTP middleware](/examples/first/main.go)
- [Kafka Component, HTTP Component, HTTP Authentication, Kafka Tracing](/examples/second/main.go)
- [Kafka Component, AMQP Tracing](/examples/third/main.go)
- [AMQP Component, AWS SNS](/examples/fourth/main.go)
- [AWS SQS](/examples/fifth/main.go)
- [gRPC](/examples/sixth/main.go)

## Processors

### Synchronous

The implementation of the processor is responsible to create a `Request` by providing everything that is needed (Headers, Fields, decoder, raw io.Reader) pass it to the implementation by invoking the `Process` method and handle the `Response` or the `error` returned by the processor.

The sync package contains only a function definition along with the models needed:

    type ProcessorFunc func(context.Context, *Request) (*Response, error)

The `Request` model contains the following properties (which are provided when calling the "constructor" `NewRequest`)

- Fields, which may contain any fields associated with the request
- Raw, the raw request data (if any) in the form of a `io.Reader`
- Headers, the request headers in the form of `map[string]string`
- decode, which is a function of type `encoding.Decode` that decodes the raw reader

An exported function exists for decoding the raw io.Reader in the form of

    Decode(v interface{}) error

The `Response` model contains the following properties (which are provided when calling the "constructor" `NewResponse`)

- Payload, which may hold a struct of type `interface{}`

### Middlewares per Route

Middlewares can also run per routes using the processor as Handler.
So using the `Route` helpers:

    // A route with ...MiddlewareFunc that will run for this route only + tracing
    route := NewRoute("/index", "GET" ProcessorFunc, true, ...MiddlewareFunc)
    // A route with ...MiddlewareFunc that will run for this route only + auth + tracing
    routeWithAuth := NewAuthRoute("/index", "GET" ProcessorFunc, true, Authendicator, ...MiddlewareFunc)

### HTTP Caching

The caching layer for HTTP routes is specified per Route.

    // RouteCache is the builder needed to build a cache for the corresponding route
    type RouteCache struct {
        // cache is the ttl cache implementation to be used
        cache cache.TTLCache
        // age specifies the minimum and maximum amount for max-age and min-fresh header values respectively
        // regarding the client cache-control requests in seconds
        age age
    }
    
    func NewRouteCache(ttlCache cache.TTLCache, age Age) *RouteCache

#### server cache

- The **cache key** is based on the route path and the url request parameters.
- The server caches only **GET requests**.
- The server implementation must specify an **Age** parameters upon construction.
- Age with **Min=0** and **Max=0** effectively disables caching
- The route should return always the most fresh object instance.
- An **ETag header** must be always in responses that are part of the cache, representing the hash of the response.
- Requests within the time-to-live threshold, will be served from the cache. 
Otherwise the request will be handled as usual by the route processor function. 
The resulting response will be cached for future requests.
- Requests where the client control header requirements cannot be met i.e. **very low max-age** or **very high min-fresh** parameters,
will be returned to the client with a `Warning` header present in the response. 

```
Note : When a cache is used, the handler execution might be skipped.
That implies that all generic handler functionalities MUST be delegated to a custom middleware.
i.e. counting number of server client requests etc ... 
```

### Usage

- provide the cache in the route builder

    NewRouteBuilder("/", handler).
        WithRouteCache(cache, http.Age{
            Min: 30 * time.Minute,
            Max: 1 * time.Hour,
        }).
        MethodGet()

- use the cache as a middleware

    NewRouteBuilder("/", handler).
        WithMiddlewares(NewCachingMiddleware(NewRouteCache(cc, Age{Max: 10 * time.Second}))).
        MethodGet()

#### client cache-control
The client can control the cache with the appropriate Headers
- `max-age=?` 

returns the cached instance only if the age of the instance is lower than the max-age parameter.
This parameter is bounded from below by the server option `minAge`.
This is to avoid chatty clients with no cache control policy (or very aggressive max-age policy) to effectively disable the cache
- `min-fresh=?` 
 
returns the cached instance if the time left for expiration is lower than the provided parameter.
This parameter is bounded from above by the server option `maxFresh`.
This is to avoid chatty clients with no cache control policy (or very aggressive min-fresh policy) to effectively disable the cache

- `no-cache` / `no-store`

returns a new response to the client by executing the route processing function.
NOTE : Except for cases where a `minAge` or `maxFresh` parameter has been specified in the server.
This is again a safety mechanism to avoid 'aggressive' clients put unexpected load on the server.
The server is responsible to cap the refresh time, BUT must respond with a `Warning` header in such a case.
- `only-if-cached`

expects any response that is found in the cache, otherwise returns an empty response

#### metrics

The http cache exposes several metrics, used to 
- assess the state of the cache
- help trim the optimal time-to-live policy
- identify client control interference

By default we are using prometheus as the the pre-defined metrics framework.

- `additions = misses + evictions`

Always , the cache addition operations (objects added to the cache), 
must be equal to the misses (requests that were not cached) plus the evictions (expired objects).
Otherwise we would expect to notice also an increased amount of errors or having the cache misbehaving in a different manner.

- `additions ~ misses`

If the additions and misses are comparable e.g. misses are almost as many as the additions, 
it would point to some cleanup of the cache itself. In that case the cache seems to not be able to support
the request patterns and control headers.

- `hits ~ additions`

The cache hit count represents how well the cache performs for the access patterns of client requests. 
If this number is rather low e.g. comparable to the additions, 
this would signify that probably a cache is not a good option for the access patterns at hand.

- `eviction age`

The age at which the objects are evicted from the cache is a very useful indicator. 
If the vast amount of evictions are close to the time to live setting, it would indicate a nicely working cache.
If we find that many evictions happen before the time to live threshold, clients would be making use cache-control headers.
 

#### cache design reference
- https://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html
- https://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.9

#### improvement considerations
- we can could the storing of the cached objects and their age counter. That way we would avoid loading the whole object in memory,
if the object is already expired. This approach might provide considerable performance (in terms of memory utilisation) 
improvement for big response objects. 
- we could extend the metrics to use the key of the object as a label as well for more fine-grained tuning.
But this has been left out for now, due to the potentially huge number of metric objects.
We can review according to usage or make this optional in the future.
- improve the serialization performance for the cache response objects

### Asynchronous

The implementation of the async processor follows exactly the same principle as the sync processor.
The main difference is that:

- The `Request` is the `Message` and contains only data as `[]byte`
- There is no `Response`, so the processor may return an error

    type ProcessorFunc func(context.Context, *Message) error

Everything else is exactly the same.

## Metrics and Tracing

Tracing and metrics are provided by Jaeger's implementation of the OpenTracing project.
Every component has been integrated with the above library and produces traces and metrics.
Metrics are provided with the default HTTP component at the `/metrics` route for Prometheus to scrape.
Tracing will be sent to a jaeger agent which can be setup through environment variables mentioned in the config section. Sane defaults are applied for making the use easy.
We have included some clients inside the trace package which are instrumented and allow propagation of tracing to
downstream systems. The tracing information is added to each implementations header. These clients are:

- HTTP
- AMQP
- Kafka
- SQL

## Correlation ID propagation

Patron receives and propagates a correlation ID. Much like the distributed tracing id, the correlation id is receiver on the entry points of the service e.g. HTTP, Kafka, etc. and is propagated via the provided clients. In case no correlation ID has been received, a new one is created.  
The ID is usually received and sent via a header with key `X-Correlation-Id`.

## Reliability

The reliability package contains the following implementations:

- Circuit Breaker

### Circuit Breaker

The circuit breaker supports a half-open state which allows to probe for successful responses in order to close the circuit again. Every aspect of the circuit breaker is configurable via its settings.

## Clients

The following clients have been implemented:

- http, with distributed tracing and optional circuit breaker
- sql, with distributed tracing
- kafka, with distributed tracing
- amqp, with distributed tracing

## Logging

The log package is designed to be a leveled logger with field support.

The log package defines the logger interface and a factory function type that needs to be implemented in order to set up the logging in this framework.

      // instantiate the implemented factory func type and fields (map[string]interface{})
      err := log.Setup(factory, fields)
      // handle error

`If the setup is omitted the package will not setup any logging!`

From there logging is as simple as

    log.Info("Hello world!")

The implementations should support the following log levels:

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

### Context Logging

Logs can be associated with some contextual data e.g. a request id. Every line logged should contain this id thus grouping the logs together. This is achieved with the usage of the context package as demonstrated below:

    ctx := log.WithContext(r.Context(), log.Sub(map[string]interface{}{"requestID": uuid.New().String()}))

The context travels through the code as an argument and can be acquired as follows:

    logger:=log.FromContext(ctx)
    logger.Infof("request processed")

Benchmarks are provided to show the performance of this.

`Every provided component creates a context logger which is then propagated in the context`

### Logger

The logger interface defines the actual logger.

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

In order to be consistent with the design the implementation of the `Fatal(f)` have to terminate the application with an error and the `Panic(f)` need to panic.

### Factory

The factory function type defines a factory for creating a logger.

    type FactoryFunc func(map[string]interface{}) Logger

## Security

The necessary abstraction is available to implement authentication in the following components:

- HTTP

### HTTP

In order to use authentication, an authenticator has to be implemented following the interface:

    type Authenticator interface {
      Authenticate(req *http.Request) (bool, error)
    }

This authenticator can then be used to set up routes with authentication.

The following authenticator is available:

- API key authenticator, see examples

## HTTP lifecycle endpoints

When creating a new HTTP component, Patron will automatically create a liveness and readiness route, which can be used to know the lifecycle of the application:

    # liveness
    GET /alive
    
    # readiness
    GET /ready

Both can return either a `200 OK` or a `503 Service Unavailable` status code (default: `200 OK`).

It is possible to customize their behaviour by injecting an `http.AliveCheck` and/or an `http.ReadyCheck` `OptionFunc` to the HTTP component constructor.
