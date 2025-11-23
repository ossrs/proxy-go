# proxy-go Project Guidelines

## Project Summary

proxy-go is a stateless media streaming proxy with built-in load balancing for building scalable origin clusters. It supports RTMP, SRT, WebRTC (WHIP/WHEP), HLS, and HTTP-FLV protocols.

Key characteristics:
- Stateless design enables horizontal scaling
- Built-in load balancer (memory or Redis-based)
- Protocol handlers for RTMP, WebRTC, SRT, HTTP streaming
- Backend origin servers register via System API
- Official solution for SRS Origin Cluster

## Design Overview

See the "Design" section in @README.md for the complete architecture overview, including:
- Stateless proxy architecture with built-in load balancing
- Single Proxy Mode (memory-based)
- Multi-Proxy Mode (Redis sync, AWS NLB)
- Complete Cluster (Edge + Proxy + Origins)

## Configuration

All configuration via environment variables (`.env` file supported):

### Server Listen Ports (client-facing)
- `PROXY_RTMP_SERVER=11935` - RTMP media server
- `PROXY_HTTP_SERVER=18080` - HTTP streaming (HLS, HTTP-FLV)
- `PROXY_WEBRTC_SERVER=18000` - WebRTC server (UDP)
- `PROXY_SRT_SERVER=20080` - SRT server (UDP)
- `PROXY_HTTP_API=11985` - HTTP API (WHIP/WHEP)
- `PROXY_SYSTEM_API=12025` - System API (origin registration)

### Load Balancer Configuration
- `PROXY_LOAD_BALANCER_TYPE=memory` - Use "memory" (single proxy) or "redis" (multi-proxy)
- `PROXY_REDIS_HOST=127.0.0.1`
- `PROXY_REDIS_PORT=6379`
- `PROXY_REDIS_PASSWORD=` (empty for no password)
- `PROXY_REDIS_DB=0`

### Other Settings
- `PROXY_STATIC_FILES=../srs/trunk/research` - Static web files directory
- `PROXY_FORCE_QUIT_TIMEOUT=30s` - Force shutdown timeout
- `PROXY_GRACE_QUIT_TIMEOUT=20s` - Graceful shutdown timeout

## How to Run

When running the project for testing or development, you should:
1. Build and start the proxy server
2. Publish a test stream using FFmpeg
3. Verify the stream is working using ffprobe

### Step 1: Build and Start Proxy Server

```bash
make && env PROXY_RTMP_SERVER=1935 PROXY_HTTP_SERVER=8080 \
    PROXY_HTTP_API=1985 PROXY_WEBRTC_SERVER=8000 PROXY_SRT_SERVER=10080 \
    PROXY_SYSTEM_API=12025 PROXY_LOAD_BALANCER_TYPE=memory ./srs-proxy
```

The proxy server should start and listen on the configured ports.

### Step 2: Publish a Test Stream

In a new terminal, publish a test stream using FFmpeg:

```bash
ffmpeg -stream_loop -1 -re -i ~/git/srs/trunk/doc/source.flv -c copy -f flv rtmp://localhost/live/livestream
```

> Note: `-stream_loop -1` makes FFmpeg loop the input file infinitely, ensuring the stream doesn't quit after the file ends.

### Step 3: Verify Stream with ffprobe

In another terminal, use ffprobe to verify the stream is working:

**Test RTMP stream:**
```bash
ffprobe rtmp://localhost/live/livestream
```

**Test HTTP-FLV stream:**
```bash
ffprobe http://localhost:8080/live/livestream.flv
```

Both commands should successfully detect the stream and display video/audio codec information. If ffprobe shows stream details without errors, the proxy is working correctly.
