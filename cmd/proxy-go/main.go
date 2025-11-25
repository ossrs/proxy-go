// Copyright (c) 2025 Winlin
//
// SPDX-License-Identifier: MIT
package main

import (
	"context"
	"os"

	"srs-proxy/internal/debug"
	"srs-proxy/internal/env"
	"srs-proxy/internal/errors"
	"srs-proxy/internal/lb"
	"srs-proxy/internal/logger"
	"srs-proxy/internal/protocol"
	"srs-proxy/internal/signal"
	"srs-proxy/internal/utils"
	"srs-proxy/internal/version"
)

func main() {
	ctx := logger.WithContext(context.Background())
	logger.Df(ctx, "%v/%v started", version.Signature(), version.Version())

	// Install signals.
	ctx, cancel := context.WithCancel(ctx)
	signal.InstallSignals(ctx, cancel)

	// Start the main loop, ignore the user cancel error.
	err := doMain(ctx)
	if err != nil && ctx.Err() != context.Canceled {
		logger.Ef(ctx, "main: %+v", err)
		os.Exit(-1)
	}

	logger.Df(ctx, "%v done", version.Signature())
}

func doMain(ctx context.Context) error {
	// Setup the environment variables.
	if err := env.LoadEnvFile(ctx); err != nil {
		return errors.Wrapf(err, "load env")
	}

	env.BuildDefaultEnvironmentVariables(ctx)

	// When cancelled, the program is forced to exit due to a timeout. Normally, this doesn't occur
	// because the main thread exits after the context is cancelled. However, sometimes the main thread
	// may be blocked for some reason, so a forced exit is necessary to ensure the program terminates.
	if err := signal.InstallForceQuit(ctx); err != nil {
		return errors.Wrapf(err, "install force quit")
	}

	// Start the Go pprof if enabled.
	debug.HandleGoPprof(ctx)

	// Initialize the load balancer.
	switch env.EnvLoadBalancerType() {
	case "redis":
		lb.SrsLoadBalancer = lb.NewRedisLoadBalancer(
			env.EnvRedisHost,
			env.EnvRedisPort,
			env.EnvRedisPassword,
			env.EnvRedisDB,
			env.EnvDefaultBackendEnabled,
			env.EnvDefaultBackendIP,
			env.EnvDefaultBackendRTMP,
			env.EnvDefaultBackendHttp,
			env.EnvDefaultBackendAPI,
			env.EnvDefaultBackendRTC,
			env.EnvDefaultBackendSRT,
		)
	default:
		lb.SrsLoadBalancer = lb.NewMemoryLoadBalancer(
			env.EnvDefaultBackendEnabled,
			env.EnvDefaultBackendIP,
			env.EnvDefaultBackendRTMP,
			env.EnvDefaultBackendHttp,
			env.EnvDefaultBackendAPI,
			env.EnvDefaultBackendRTC,
			env.EnvDefaultBackendSRT,
		)
	}

	if err := lb.SrsLoadBalancer.Initialize(ctx); err != nil {
		return errors.Wrapf(err, "initialize srs load balancer")
	}

	// Parse the gracefully quit timeout.
	gracefulQuitTimeout, err := utils.ParseGracefullyQuitTimeout()
	if err != nil {
		return errors.Wrapf(err, "parse gracefully quit timeout")
	}

	// Start the RTMP server.
	srsRTMPServer := protocol.NewSRSRTMPServer()
	defer srsRTMPServer.Close()
	if err := srsRTMPServer.Run(ctx); err != nil {
		return errors.Wrapf(err, "rtmp server")
	}

	// Start the WebRTC server.
	srsWebRTCServer := protocol.NewSRSWebRTCServer()
	defer srsWebRTCServer.Close()
	if err := srsWebRTCServer.Run(ctx); err != nil {
		return errors.Wrapf(err, "rtc server")
	}

	// Start the HTTP API server.
	srsHTTPAPIServer := protocol.NewSRSHTTPAPIServer(gracefulQuitTimeout, srsWebRTCServer)
	defer srsHTTPAPIServer.Close()
	if err := srsHTTPAPIServer.Run(ctx); err != nil {
		return errors.Wrapf(err, "http api server")
	}

	// Start the SRT server.
	srsSRTServer := protocol.NewSRSSRTServer()
	defer srsSRTServer.Close()
	if err := srsSRTServer.Run(ctx); err != nil {
		return errors.Wrapf(err, "srt server")
	}

	// Start the System API server.
	systemAPI := protocol.NewSystemAPI(gracefulQuitTimeout)
	defer systemAPI.Close()
	if err := systemAPI.Run(ctx); err != nil {
		return errors.Wrapf(err, "system api server")
	}

	// Start the HTTP web server.
	srsHTTPStreamServer := protocol.NewSRSHTTPStreamServer(gracefulQuitTimeout)
	defer srsHTTPStreamServer.Close()
	if err := srsHTTPStreamServer.Run(ctx); err != nil {
		return errors.Wrapf(err, "http server")
	}

	// Wait for the main loop to quit.
	<-ctx.Done()
	return nil
}
