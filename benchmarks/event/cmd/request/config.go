package main

import (
	"runtime"

	"github.com/go-orb/go-orb/cli"
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

func flags() []*cli.Flag {
	var flags []*cli.Flag

	flags = append(flags, cli.NewFlag(
		"connections",
		defaultConnections,
		cli.FlagConfigPaths([]string{configSection, "connections"}),
		cli.FlagUsage("Connections to keep open"),
		cli.FlagEnvVars("CONNECTIONS"),
	))

	flags = append(flags, cli.NewFlag(
		"duration",
		defaultDuration,
		cli.FlagConfigPaths([]string{configSection, "duration"}),
		cli.FlagUsage("Duration in seconds"),
		cli.FlagEnvVars("DURATION"),
	))

	flags = append(flags, cli.NewFlag(
		"timeout",
		defaultTimeout,
		cli.FlagConfigPaths([]string{configSection, "timeout"}),
		cli.FlagUsage("Timeout in seconds"),
		cli.FlagEnvVars("TIMEOUT"),
	))

	flags = append(flags, cli.NewFlag(
		"threads",
		defaultThreads,
		cli.FlagConfigPaths([]string{configSection, "threads"}),
		cli.FlagUsage("Number of threads to use = runtime.GOMAXPROCS()"),
		cli.FlagEnvVars("THREADS"),
	))

	flags = append(flags, cli.NewFlag(
		"package_size",
		defaultPackageSize,
		cli.FlagConfigPaths([]string{configSection, "packageSize"}),
		cli.FlagUsage("Per request package size"),
		cli.FlagEnvVars("PACKAGE_SIZE"),
	))

	flags = append(flags, cli.NewFlag(
		"content_type",
		defaultContentType,
		cli.FlagConfigPaths([]string{configSection, "contentType"}),
		cli.FlagUsage("Content-Type (application/x-protobuf, application/x-protobuf+json)"),
		cli.FlagEnvVars("CONTENT_TYPE"),
	))

	return flags
}

type clientConfig struct {
	Connections int `json:"connections" yaml:"connections"`
	Duration    int `json:"duration"    yaml:"duration"`
	Timeout     int `json:"timeout"     yaml:"timeout"`
	Threads     int `json:"threads"     yaml:"threads"`
	PackageSize int `json:"packageSize" yaml:"packageSize"`
}
