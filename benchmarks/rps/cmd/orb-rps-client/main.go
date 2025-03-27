// Package main contains an RPC benchmarking client
package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-orb/go-orb/cli"
	"github.com/go-orb/go-orb/client"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/registry"

	echoproto "github.com/go-orb/examples/benchmarks/rps/proto/echo"

	// Import required plugins.
	_ "github.com/go-orb/plugins/client/orb"
	_ "github.com/go-orb/plugins/codecs/goccyjson"
	_ "github.com/go-orb/plugins/codecs/proto"
	_ "github.com/go-orb/plugins/codecs/yaml"
	_ "github.com/go-orb/plugins/config/source/file"
	_ "github.com/go-orb/plugins/log/slog"
	_ "github.com/go-orb/plugins/registry/consul"
	_ "github.com/go-orb/plugins/registry/mdns"

	// Transport plugins.
	_ "github.com/go-orb/plugins/client/orb_transport/drpc"
	_ "github.com/go-orb/plugins/client/orb_transport/grpc"
	_ "github.com/go-orb/plugins/client/orb_transport/http"
)

const serverName = "benchmarks.rps.server"

type benchStats struct {
	ok    atomic.Uint64
	error atomic.Uint64
}

// runConnection sends benchmark requests until context is canceled.
func runConnection(
	ctx context.Context,
	wg *sync.WaitGroup,
	cli client.Type,
	logger log.Logger,
	msg []byte,
	opts []client.CallOption,
	stats *benchStats,
) {
	defer wg.Done()

	req := &echoproto.Req{Payload: msg}
	client := echoproto.NewEchoClient(cli)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Run the query
			resp, err := client.Echo(
				ctx,
				serverName,
				req,
				opts...,
			)

			if err != nil {
				if !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
					logger.Error("request failed", "error", err)
					stats.error.Add(1)
				}

				continue
			}

			// Verify response matches request
			if !bytes.Equal(req.GetPayload(), resp.GetPayload()) {
				logger.Error("mismatched response data")
				stats.error.Add(1)

				continue
			}

			stats.ok.Add(1)
		}
	}
}

// setupClientOptions prepares client options based on configuration.
func setupClientOptions(ctx context.Context, cfg *clientConfig, cli client.Type, logger log.Logger) ([]client.CallOption, error) {
	// Basic client options
	opts := []client.CallOption{
		client.WithPreferredTransports(cfg.Transport),
		client.WithContentType(cfg.ContentType),
	}

	return opts, nil
}

// runBenchmark executes the benchmark with the given configuration.
func runBenchmark(
	ctx context.Context,
	duration int,
	connections int,
	cli client.Type,
	logger log.Logger,
	msg []byte,
	opts []client.CallOption,
) (uint64, uint64) {
	var wg sync.WaitGroup

	stats := &benchStats{}

	// Create benchmark context with timeout
	benchCtx, cancel := context.WithTimeout(ctx, time.Duration(duration)*time.Second)
	defer cancel()

	// Start worker connections
	for i := 0; i < connections; i++ {
		wg.Add(1)

		go runConnection(benchCtx, &wg, cli, logger, msg, opts, stats)
	}

	// Wait for timeout or cancellation
	<-benchCtx.Done()

	// Wait for all goroutines to finish
	wg.Wait()

	return stats.ok.Load(), stats.error.Load()
}

// bench is the main benchmark function.
func bench(ctx context.Context, cfg *clientConfig, logger log.Logger, cli client.Type) error {
	// Log configuration
	logger.Info("Configuration",
		"connections", cfg.Connections,
		"duration", cfg.Duration,
		"timeout", cfg.Timeout,
		"threads", cfg.Threads,
		"transport", cfg.Transport,
		"package_size", cfg.PackageSize,
		"content_type", cfg.ContentType,
	)

	// Set max threads
	runtime.GOMAXPROCS(cfg.Threads)

	// Setup client options
	opts, err := setupClientOptions(ctx, cfg, cli, logger)
	if err != nil {
		return err
	}

	// Generate random payload data
	msg := make([]byte, cfg.PackageSize)
	if _, err := rand.Reader.Read(msg); err != nil {
		return fmt.Errorf("failed to generate random data: %w", err)
	}

	// Run warmup phase
	logger.Info("Warming up...")
	warmupOk, warmupErr := runBenchmark(ctx, 5, cfg.Connections, cli, logger, msg, opts)
	logger.Debug("Warmup complete", "requests_ok", warmupOk, "requests_error", warmupErr)

	// Run benchmark phase
	logger.Info("Running benchmark...")
	reqsOk, reqsError := runBenchmark(ctx, cfg.Duration, cfg.Connections, cli, logger, msg, opts)

	// Log results
	logger.Info("Summary",
		"requests_ok", reqsOk,
		"requests_error", reqsError,
		"qps", float64(reqsOk)/float64(cfg.Duration),
		"connections", cfg.Connections,
		"duration_seconds", cfg.Duration,
		"package_size", cfg.PackageSize,
		"transport", cfg.Transport,
	)

	return nil
}

func main() {
	app := cli.App{
		Name:    "orb-rps-client",
		Version: "1.0.0",
		Usage:   "A benchmarking client for RPC services",
		Flags: []*cli.Flag{
			{
				Name:        "registry",
				Default:     registry.DefaultRegistry,
				EnvVars:     []string{"REGISTRY"},
				ConfigPaths: [][]string{{"registry", "plugin"}},
				Usage:       "Set the registry plugin, one of mdns, consul, memory",
			},
			{
				Name:        "log_level",
				Default:     "INFO",
				EnvVars:     []string{"LOG_LEVEL"},
				ConfigPaths: [][]string{{"logger", "level"}},
				Usage:       "Set the log level, one of TRACE, DEBUG, INFO, WARN, ERROR",
			},
		},
	}
	app.Flags = append(app.Flags, flags()...)

	appContext := cli.NewAppContext(&app)

	_, err := run(appContext, os.Args, bench)
	if err != nil {
		//nolint:forbidigo
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
