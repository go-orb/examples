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

### [event/simple](event/simple)

A simple example of RPC requests over the event plugins, currently theres only the natsjs backend.

### [rest/auth](rest/auth)

A simple example of a go-orb service and client with an auth REST API.

This example doesn't utilize the go-orb config system, just to show it's possible to configure a go-orb service without it.

### [rest/middleware](rest/middleware)

A simple example of a go-orb service and client with a REST middleware.

In it's [config](rest/middleware/config) folder you can find a variaty of config files for different transports as well as logging options and registries. All of them run with the same code/binary.

## Authors

### go-orb

- [David Brouwer](https://github.com/Davincible)
- [René Jochum](https://github.com/jochumdev)

## License

go-orb is Apache 2.0 licensed and is based on go-micro.