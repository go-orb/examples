# go-orb/examples

Contains examples and benchmarks for go-orb.

## Available examples

### [benchmarks](benchmarks/)

#### [benchmarks/event](benchmarks/event)

A benchmark for running RPC requests over the event plugins, currently theres only the natsjs backend.

#### [benchmarks/rps](benchmarks/rps)

A benchmark for running requests-per-second (rps) for a go-orb/server.

The rps benchmark sends X bytes (default `1000`) to server which echoes it to the client.

See the [WIKI](https://github.com/go-orb/go-orb/wiki/RPC-Benchmarks) for results.

### [cmd/foobar](cmd/foobar)

A example on howto run multiple services in a monolith.

### [event/simple](event/simple)

A simple example of RPC requests over the event plugins, currently theres only the natsjs backend.

### [realworld](realworld)

An advanced example which shows howto build a multi-service app that is also able to run as monolith.

This example uses [httpgateway](https://github.com/go-orb/httpgateway) to forward requests.

It is the latest and greatest example and we will work on extending it.

### [rest/auth](rest/auth)

A simple example of a go-orb service and client with an auth REST API.

This example doesn't utilize the go-orb config system, just to show it's possible to configure a go-orb service without it. It has the smallest binary size of all examples.

### [rest/middleware](rest/middleware)

A simple example of a go-orb service and client with a REST middleware.

In it's [config](rest/middleware/config) folder you can find a variaty of config files for different transports as well as logging options and registries. All of them run with the same code/binary.

Run the server with:

```bash
go run ./cmd/server/... --config ./config/grpc.yaml
```

And the client with:

```bash
go run ./cmd/client/... --config ./config/grpc.yaml
```

## Running

To run an example you need to have `go1.23.6` or later installed. You can get it with [gvm](https://github.com/moovweb/gvm).

All the examples have a cmd folder you can run it with `go run "./cmd/handler/..."`.

You can run all examples at once by running `./scripts/test.sh`, this is what we do in CI.

## Development

You need:
- go1.23.6 or later
- docker for dagger
- our custom fork of [wire](https://github.com/go-orb/wire) `go install github.com/go-orb/wire/cmd/wire@latest`
- [protoc-gen-go](https://protobuf.dev/reference/go/go-generated/) `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`
- [protoc-gen-go-orb](https://github.com/go-orb/plugins/server/cmd/protoc-gen-go-orb) `go install github.com/go-orb/plugins/server/cmd/protoc-gen-go-orb@latest`

Now you'r able to run `go generate ./...` in the examples to update `cmd/xyz/wire_gen.go` as well as the `proto` files.

## Authors

### go-orb

- [David Brouwer](https://github.com/Davincible)
- [René Jochum](https://github.com/jochumdev)

## License

go-orb is Apache 2.0 licensed and is based on go-micro.