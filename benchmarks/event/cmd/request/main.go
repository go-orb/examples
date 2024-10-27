package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"os"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/go-orb/examples/benchmarks/event/pb/echo"
	"github.com/go-orb/go-orb/config"
	"github.com/go-orb/go-orb/event"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/types"
	_ "github.com/go-orb/plugins/codecs/json"
	_ "github.com/go-orb/plugins/codecs/proto"
	_ "github.com/go-orb/plugins/config/source/cli/urfave"
	_ "github.com/go-orb/plugins/event/natsjs"
	_ "github.com/go-orb/plugins/log/slog"
)

type stats struct {
	Ok    uint64
	Error uint64
}

func connection(
	ctx context.Context,
	wg *sync.WaitGroup,
	eventHandler event.Handler,
	logger log.Logger,
	req *echo.Req,
	connectionNum int,
	statsChan chan stats,
) {
	var (
		reqsOk    uint64
		reqsError uint64
	)

	for {
		select {
		case <-ctx.Done():
			logger.Debug("Connection results", "connection", connectionNum, "reqsOk", reqsOk, "reqsError", reqsError)
			wg.Done()

			statsChan <- stats{Ok: reqsOk, Error: reqsError}

			return
		default:
		}

		// Run the query.
		resp, err := event.Request[echo.Resp](context.Background(), eventHandler, "echo.echo", req)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				continue
			}

			logger.Error("while requesting", "error", err)

			reqsError++

			continue
		}

		// Check if response equals.
		if !bytes.Equal(req.GetPayload(), resp.GetPayload()) {
			logger.Error("request and response are not the same")

			reqsError++

			continue
		}

		reqsOk++
	}
}

// bench.
//
//nolint:funlen
func bench(
	sn types.ServiceName,
	configs types.ConfigData,
	logger log.Logger,
	eventHandler event.Handler,
	done chan os.Signal,
) error {
	cfg := &clientConfig{
		Connections: defaultConnections,
		Duration:    defaultDuration,
		Timeout:     defaultTimeout,
		Threads:     defaultThreads,
		PackageSize: defaultPackageSize,
		ContentType: defaultContentType,
	}

	sections := append(types.SplitServiceName(sn), configSection)
	if err := config.Parse(sections, configs, &cfg); err != nil {
		return err
	}

	logger.Info(
		"Config",
		"connections", cfg.Connections,
		"duration", cfg.Duration,
		"timeout", cfg.Timeout,
		"threads", cfg.Threads,
		"package_size", cfg.PackageSize,
		"content_type", cfg.ContentType,
	)

	runtime.GOMAXPROCS(cfg.Threads)

	wCtx, wCancel := context.WithCancel(context.Background())

	// Create random bytes to ping-pong on each request.
	msg := make([]byte, cfg.PackageSize)
	if _, err := rand.Reader.Read(msg); err != nil {
		logger.Error("Failed to make a request", "error", err)
		wCancel()

		return err
	}

	var wg sync.WaitGroup

	quit := make(chan os.Signal, 1)

	//
	// Warmup
	//

	timer := time.AfterFunc(time.Second*time.Duration(cfg.Duration), func() {
		done <- syscall.SIGINT
	})

	logger.Info("Warming up...")

	nullChan := make(chan stats, cfg.Connections)

	for i := 0; i < cfg.Connections; i++ {
		wg.Add(1)

		// Create a request.
		req := &echo.Req{Payload: msg}

		go connection(wCtx, &wg, eventHandler, logger, req, i, nullChan)
	}

	select {
	case <-done:
		wCancel()
		timer.Stop()
	case <-quit:
		timer.Stop()
		os.Exit(1)
	}

	//
	// Bench
	//
	logger.Info("Now running the benchmark")

	ctx, cancel := context.WithCancel(context.Background())

	// Timer to end requests
	timer = time.AfterFunc(time.Second*time.Duration(cfg.Duration), func() {
		done <- syscall.SIGINT
	})

	// Statistics channel
	statsChan := make(chan stats, cfg.Connections)

	// Run the requests.
	for i := 0; i < cfg.Connections; i++ {
		wg.Add(1)

		// Create a request.
		req := &echo.Req{Payload: msg}

		go connection(ctx, &wg, eventHandler, logger, req, i, statsChan)
	}

	// Blocks until timer/signal happened
	select {
	case <-done:
		timer.Stop()
		// stops requesting
		cancel()

		// Wait for all goroutines to exit properly.
		wg.Wait()
	case <-quit:
		timer.Stop()
		os.Exit(0)
	}

	// Calculate stats
	mStats := stats{}

	for i := 0; i < cfg.Connections; i++ {
		cStat := <-statsChan

		mStats.Ok += cStat.Ok
		mStats.Error += cStat.Error
	}

	logger.Info("Summary",
		"connections", cfg.Connections,
		"duration", cfg.Duration,
		"timeout", cfg.Timeout,
		"threads", cfg.Threads,
		"package_size", cfg.PackageSize,
		"content_type", cfg.ContentType,
		"reqsOk", mStats.Ok,
		"reqsError", mStats.Error,
	)

	return nil
}

func main() {
	var (
		serviceName    = types.ServiceName("orb.examples.event.bench_request")
		serviceVersion = types.ServiceVersion("v0.0.1")
	)

	if _, err := run(serviceName, serviceVersion, bench); err != nil {
		log.Error("while running", "err", err)
	}
}
