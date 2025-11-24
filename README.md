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

Please see [Design](doc/design.md) for details.

William Yang<br/>
June 23, 2025
