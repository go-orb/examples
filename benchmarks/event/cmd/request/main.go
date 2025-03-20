// Package main contains a client which benchmarks requests-per-second (rps) for a go-orb/server.
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
	"time"

	"github.com/go-orb/examples/benchmarks/event/pb/echo"
	"github.com/go-orb/go-orb/cli"
	"github.com/go-orb/go-orb/config"
	"github.com/go-orb/go-orb/event"
	"github.com/go-orb/go-orb/log"
	_ "github.com/go-orb/plugins/codecs/goccyjson"
	_ "github.com/go-orb/plugins/codecs/proto"
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
	eventHandler event.Type,
	logger log.Logger,
	req *echo.Req,
	connectionNum int,
	statsChan chan stats,
) {
	var (
		reqsOk    uint64
		reqsError uint64
	)

	eventHandler = eventHandler.Clone()
	if err := eventHandler.Start(ctx); err != nil {
		logger.Error("Failed to start", "err", err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			logger.Debug("Connection results", "connection", connectionNum, "reqsOk", reqsOk, "reqsError", reqsError)
			wg.Done()

			statsChan <- stats{Ok: reqsOk, Error: reqsError}

			if err := eventHandler.Stop(context.Background()); err != nil {
				logger.Error("Failed to stop", "err", err)
			}

			return
		default:
		}

		// Run the query.
		resp, err := event.Request[echo.Resp](ctx, eventHandler, "echo.echo", req)
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
	svcCtx *cli.ServiceContext,
	logger log.Logger,
	eventHandler event.Type,
) error {
	cfg := &clientConfig{
		Connections: defaultConnections,
		Duration:    defaultDuration,
		Timeout:     defaultTimeout,
		Threads:     defaultThreads,
		PackageSize: defaultPackageSize,
	}

	if err := config.Parse(nil, event.DefaultConfigSection, svcCtx.Config, &cfg); err != nil {
		return err
	}

	logger.Info(
		"Config",
		"connections", cfg.Connections,
		"duration", cfg.Duration,
		"timeout", cfg.Timeout,
		"threads", cfg.Threads,
		"package_size", cfg.PackageSize,
	)

	runtime.GOMAXPROCS(cfg.Threads)

	wCtx, wCancel := context.WithCancel(svcCtx.Context())

	// Create random bytes to ping-pong on each request.
	msg := make([]byte, cfg.PackageSize)
	if _, err := rand.Reader.Read(msg); err != nil {
		logger.Error("Failed to make a request", "error", err)
		wCancel()

		return err
	}

	var wg sync.WaitGroup

	//
	// Warmup
	//

	time.AfterFunc(time.Second*time.Duration(cfg.Duration), func() {
		wCancel()
	})

	logger.Info("Warming up...")

	nullChan := make(chan stats, cfg.Connections)

	// Create a request.
	req := &echo.Req{Payload: msg}

	// Run cfg.Connections go routines which request in a loop.
	for i := 0; i < cfg.Connections; i++ {
		wg.Add(1)

		go connection(wCtx, &wg, eventHandler, logger, req, i, nullChan)
	}

	// Wait for the warmup
	<-wCtx.Done()

	//
	// Bench
	//
	logger.Info("Now running the benchmark")

	ctx, cancel := context.WithCancel(svcCtx.Context())

	// Timer to end requests
	time.AfterFunc(time.Second*time.Duration(cfg.Duration), func() {
		cancel()
	})

	// Statistics channel
	statsChan := make(chan stats, cfg.Connections)

	// Run cfg.Connections go routines which request in a loop.
	for i := 0; i < cfg.Connections; i++ {
		wg.Add(1)

		go connection(ctx, &wg, eventHandler, logger, req, i, statsChan)
	}

	// Blocks until timer/signal happened
	<-ctx.Done()

	// stops requesting
	cancel()

	// Wait for all goroutines to exit properly.
	wg.Wait()

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
		"reqsOk", mStats.Ok,
		"reqsError", mStats.Error,
	)

	return nil
}

func main() {
	app := cli.App{
		Name:     "orb.examples.event.bench_request",
		Version:  "",
		Usage:    "A benchmarking client",
		NoAction: false,
		Flags: []*cli.Flag{
			{
				Name:        "log_level",
				Default:     "INFO",
				EnvVars:     []string{"LOG_LEVEL"},
				ConfigPaths: []cli.FlagConfigPath{{Path: []string{"logger", "level"}}},
				Usage:       "Set the log level, one of TRACE, DEBUG, INFO, WARN, ERROR",
			},
		},
		Commands: []*cli.Command{},
	}
	app.Flags = append(app.Flags, flags()...)

	appContext := cli.NewAppContext(&app)

	_, err := run(appContext, os.Args, bench)
	if err != nil {
		//nolint:forbidigo
		fmt.Printf("run error: %s\n", err)
		os.Exit(1)
	}
}
