# proxy-go

An common proxy server for any media servers with RTMP/SRT/HLS/HTTP-FLV and 
WebRTC/WHIP/WHEP protocols support.

## Usage

This is a common proxy for all media servers, to enable you to build a Origin Cluster 
for your media server.

However, SRS works with this proxy much better than other media servers, as the proxy
can discover more details from SRS. And this proxy is the official solution to build 
an origin cluster for SRS. Please see [SRS Origin Cluster](https://ossrs.io/lts/en-us/docs/v7/doc/origin-cluster) 
for details.

## Design

**proxy-go** is a stateless media streaming proxy with built-in load balancing that enables building scalable origin clusters. The proxy itself acts as the load balancer, routing streams from clients to backend origin servers.

```
Client → Proxy (with Load Balancer) → Backend Origin Servers
```

Since the proxy is stateless, you can deploy multiple proxies behind a load balancer (like AWS NLB) for horizontal scaling:

```
Client → AWS NLB → Proxy Servers → Backend Origin Servers
```

**Single Proxy Mode**

Use case: Moderate amount of streams requiring multiple origin servers (each stream has few viewers). The total stream count is manageable by a single proxy server. Uses memory-based load balancing (no Redis needed).

```
                                       +--------------------+
                               +-------+ Origin Server A    +
                               +       +--------------------+
                               +
+-----------------------+      +       +--------------------+
+   Proxy Server        +------+-------+ Origin Server B    +
+ (Memory LB)           +      +       +--------------------+
+-----------------------+      +
                               +       +--------------------+
                               +-------+ Origin Server C    +
                                       +--------------------+
```

**Multi-Proxy Mode (Scalable)**

Use case: When a single proxy becomes a bottleneck. Supports a large number of streams across many origin servers, with limited viewers per stream. Redis is required for state synchronization between proxies.

```
                         +-----------------------+
                     +---+ Proxy Server A        +------+
+-----------------+  |   +-----------+-----------+      +
|    AWS NLB      +--+               |                  +
+-----------------+  |          (Redis Sync)            + Origin Servers
                     |   +-----------+-----------+      +
                     +---+ Proxy Server B        +------+
                         +-----------------------+
```

**Complete Cluster (Edge + Proxy + Origins)**

Use case: Very large deployments with both numerous streams AND numerous viewers. Edge servers aggregate upstream connections - fetching one stream from upstream to serve multiple viewers, dramatically reducing load on proxy and origin servers.

```
Edge Servers → Proxy Servers → Origin Servers
(Proxy + Cache)   (Proxy)      (SRS/Media)
```

> **Note**: Future edge servers will be implemented as proxy servers with caching enabled, creating a unified architecture where the same codebase serves both proxy and edge roles. The edge cache aggregates viewer connections, so thousands of viewers can watch the same stream while only requesting it once from upstream.

William Yang<br/>
June 23, 2025
