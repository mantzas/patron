# example

## Prerequisites

- Docker
- Docker compose

## setting up environment

To run the full example we need to start [jaeger](https://www.jaegertracing.io/) and [prometheus](https://prometheus.io/). We can startup both of them using docker-compose with the following command.

```shell
docker-compose up -d
```

To tear down the above just:

```shell
docker-compose down
```

## running the example

just run (within the examples folder):

```shell
go run main.go
```

and the use curl to send a request:

```shell
curl -i -H "Content-Type: application/json" http://localhost:50000
```

After that head over to [jaeger](http://localhost:16686/search) and [prometheus](http://localhost:9090/graph).