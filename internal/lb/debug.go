// Copyright (c) 2025 Winlin
//
// SPDX-License-Identifier: MIT
package lb

import (
	"fmt"
	"os"
	"time"

	"srs-proxy/internal/logger"
)

// NewDefaultSRSForDebugging initialize the default SRS media server, for debugging only.
func NewDefaultSRSForDebugging(
	envDefaultBackendEnabled func() string,
	envDefaultBackendIP func() string,
	envDefaultBackendRTMP func() string,
	envDefaultBackendHttp func() string,
	envDefaultBackendAPI func() string,
	envDefaultBackendRTC func() string,
	envDefaultBackendSRT func() string,
) (*SRSServer, error) {
	if envDefaultBackendEnabled() != "on" {
		return nil, nil
	}

	if envDefaultBackendIP() == "" {
		return nil, fmt.Errorf("empty default backend ip")
	}
	if envDefaultBackendRTMP() == "" {
		return nil, fmt.Errorf("empty default backend rtmp")
	}

	server := NewSRSServer(func(srs *SRSServer) {
		srs.IP = envDefaultBackendIP()
		srs.RTMP = []string{envDefaultBackendRTMP()}
		srs.ServerID = fmt.Sprintf("default-%v", logger.GenerateContextID())
		srs.ServiceID = logger.GenerateContextID()
		srs.PID = fmt.Sprintf("%v", os.Getpid())
		srs.UpdatedAt = time.Now()
	})

	if envDefaultBackendHttp() != "" {
		server.HTTP = []string{envDefaultBackendHttp()}
	}
	if envDefaultBackendAPI() != "" {
		server.API = []string{envDefaultBackendAPI()}
	}
	if envDefaultBackendRTC() != "" {
		server.RTC = []string{envDefaultBackendRTC()}
	}
	if envDefaultBackendSRT() != "" {
		server.SRT = []string{envDefaultBackendSRT()}
	}
	return server, nil
}
