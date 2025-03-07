package main

import (
	"runtime"

	"github.com/go-orb/go-orb/cli"
)

const (
	configSection = "bench_client"

	defaultBypassRegistry = 1
	defaultConnections    = 256
	defaultDuration       = 15
	defaultTimeout        = 8
	defaultTransport      = "grpc"
	defaultPackageSize    = 1000
	defaultContentType    = "application/x-protobuf"
)

//nolint:gochecknoglobals
var (
	defaultThreads = runtime.NumCPU()
)

func flags() []*cli.Flag {
	var flags []*cli.Flag

	flags = append(flags, cli.NewFlag(
		"bypass_registry",
		defaultBypassRegistry,
		cli.FlagConfigPaths(cli.FlagConfigPath{Path: []string{configSection, "bypassRegistry"}}),
		cli.FlagUsage("Bypasses the registry by caching it, set to 0 to disable"),
		cli.FlagEnvVars("BYPASS_REGISTRY"),
	))

	flags = append(flags, cli.NewFlag(
		"connections",
		defaultConnections,
		cli.FlagConfigPaths(cli.FlagConfigPath{Path: []string{configSection, "connections"}}),
		cli.FlagUsage("Connections to keep open"),
		cli.FlagEnvVars("CONNECTIONS"),
	))

	flags = append(flags, cli.NewFlag(
		"duration",
		defaultDuration,
		cli.FlagConfigPaths(cli.FlagConfigPath{Path: []string{configSection, "duration"}}),
		cli.FlagUsage("Duration in seconds"),
		cli.FlagEnvVars("DURATION"),
	))

	flags = append(flags, cli.NewFlag(
		"timeout",
		defaultTimeout,
		cli.FlagConfigPaths(cli.FlagConfigPath{Path: []string{configSection, "timeout"}}),
		cli.FlagUsage("Timeout in seconds"),
		cli.FlagEnvVars("TIMEOUT"),
	))

	flags = append(flags, cli.NewFlag(
		"threads",
		defaultThreads,
		cli.FlagConfigPaths(cli.FlagConfigPath{Path: []string{configSection, "threads"}}),
		cli.FlagUsage("Number of threads to use = runtime.GOMAXPROCS()"),
		cli.FlagEnvVars("THREADS"),
	))

	flags = append(flags, cli.NewFlag(
		"transport",
		defaultTransport,
		cli.FlagConfigPaths(cli.FlagConfigPath{Path: []string{configSection, "transport"}}),
		cli.FlagUsage("Transport to use (grpc, drpc, http, uvm.)"),
		cli.FlagEnvVars("TRANSPORT"),
	))

	flags = append(flags, cli.NewFlag(
		"package_size",
		defaultPackageSize,
		cli.FlagConfigPaths(cli.FlagConfigPath{Path: []string{configSection, "packageSize"}}),
		cli.FlagUsage("Per request package size"),
		cli.FlagEnvVars("PACKAGE_SIZE"),
	))

	flags = append(flags, cli.NewFlag(
		"content_type",
		defaultContentType,
		cli.FlagConfigPaths(cli.FlagConfigPath{Path: []string{configSection, "contentType"}}),
		cli.FlagUsage("Content-Type (application/x-protobuf, application/x-protobuf+json)"),
		cli.FlagEnvVars("CONTENT_TYPE"),
	))

	return flags
}

type clientConfig struct {
	BypassRegistry int    `json:"bypassRegistry" yaml:"bypassRegistry"`
	Connections    int    `json:"connections"    yaml:"connections"`
	Duration       int    `json:"duration"       yaml:"duration"`
	Timeout        int    `json:"timeout"        yaml:"timeout"`
	Threads        int    `json:"threads"        yaml:"threads"`
	Transport      string `json:"transport"      yaml:"transport"`
	PackageSize    int    `json:"packageSize"    yaml:"packageSize"`
	ContentType    string `json:"contentType"    yaml:"contentType"`
}
