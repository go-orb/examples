package main

import (
	"errors"
	"runtime"

	"github.com/go-orb/go-orb/config/source/cli"
)

const (
	configSection = "bench_client"

	defaultConnections = 256
	defaultDuration    = 15
	defaultTimeout     = 8
	defaultPackageSize = 1000
	defaultContentType = "application/x-protobuf"
)

//nolint:gochecknoglobals
var (
	defaultThreads = runtime.NumCPU()
)

func init() {
	err := cli.Flags.Add(cli.NewFlag(
		"connections",
		defaultConnections,
		cli.ConfigPathSlice([]string{configSection, "connections"}),
		cli.Usage("Connections to keep open"),
		cli.EnvVars("CONNECTIONS"),
	))
	if err != nil && !errors.Is(err, cli.ErrFlagExists) {
		panic(err)
	}

	err = cli.Flags.Add(cli.NewFlag(
		"duration",
		defaultDuration,
		cli.ConfigPathSlice([]string{configSection, "duration"}),
		cli.Usage("Duration in seconds"),
		cli.EnvVars("DURATION"),
	))
	if err != nil && !errors.Is(err, cli.ErrFlagExists) {
		panic(err)
	}

	err = cli.Flags.Add(cli.NewFlag(
		"timeout",
		defaultTimeout,
		cli.ConfigPathSlice([]string{configSection, "timeout"}),
		cli.Usage("Timeout in seconds"),
		cli.EnvVars("TIMEOUT"),
	))
	if err != nil && !errors.Is(err, cli.ErrFlagExists) {
		panic(err)
	}

	err = cli.Flags.Add(cli.NewFlag(
		"threads",
		defaultThreads,
		cli.ConfigPathSlice([]string{configSection, "threads"}),
		cli.Usage("Number of threads to use = runtime.GOMAXPROCS()"),
		cli.EnvVars("THREADS"),
	))
	if err != nil && !errors.Is(err, cli.ErrFlagExists) {
		panic(err)
	}

	err = cli.Flags.Add(cli.NewFlag(
		"package_size",
		defaultPackageSize,
		cli.ConfigPathSlice([]string{configSection, "packageSize"}),
		cli.Usage("Per request package size"),
		cli.EnvVars("PACKAGE_SIZE"),
	))
	if err != nil && !errors.Is(err, cli.ErrFlagExists) {
		panic(err)
	}
}

type clientConfig struct {
	Connections int `json:"connections" yaml:"connections"`
	Duration    int `json:"duration"    yaml:"duration"`
	Timeout     int `json:"timeout"     yaml:"timeout"`
	Threads     int `json:"threads"     yaml:"threads"`
	PackageSize int `json:"packageSize" yaml:"packageSize"`
}
