// Copyright (c) 2025 Winlin
//
// SPDX-License-Identifier: MIT
package debug

import (
	"context"
	"net/http"

	"srs-proxy/internal/env"
	"srs-proxy/internal/logger"
)

func HandleGoPprof(ctx context.Context, environment env.Environment) {
	if addr := environment.GoPprof(); addr != "" {
		go func() {
			logger.Df(ctx, "Start Go pprof at %v", addr)
			http.ListenAndServe(addr, nil)
		}()
	}
}
