# SingerBox

A Go library for parsing proxy share links and managing sing-box proxy instances programmatically.

[![Go Reference](https://pkg.go.dev/badge/github.com/lunardoesdev/singerbox.svg)](https://pkg.go.dev/github.com/lunardoesdev/singerbox)

## Features

- ✅ **Parse share links** - VLESS, VMess, Shadowsocks, Trojan, SOCKS5, HTTP/HTTPS
- ✅ **Manage sing-box** - Start/stop proxy instances programmatically
- ✅ **Dual proxy modes** - SOCKS5 and HTTP proxy support
- ✅ **Reality protocol** - Full support with uTLS fingerprinting
- ✅ **Well-tested** - 91.8% test coverage
- ✅ **Clean API** - Simple, documented, easy to use

## Installation

```bash
go get github.com/lunardoesdev/singerbox
```

## Quick Start

### Parse Share Links

```go
package main

import (
    "fmt"
    "github.com/lunardoesdev/singerbox"
)

func main() {
    // Create parser
    parser := singerbox.NewParser()

    // Parse any share link
    outbound, err := parser.Parse("vless://uuid@server:443?security=tls&type=ws")
    if err != nil {
        panic(err)
    }

    fmt.Printf("Parsed %s proxy\n", outbound.Type)
}
```

### Start a Proxy

```go
package main

import (
    "fmt"
    "github.com/lunardoesdev/singerbox"
)

func main() {
    // Parse share link
    parser := singerbox.NewParser()
    outbound, _ := parser.Parse("ss://aes-256-gcm:password@server:8388")

    // Create and start proxy
    pb, _ := singerbox.NewProxyBox(singerbox.ProxyBoxConfig{
        Outbound:   outbound,
        ListenAddr: "127.0.0.1:1080",
        HTTPPort:   1081,
    })

    pb.Start()
    defer pb.Stop()

    fmt.Printf("SOCKS5: %s\n", pb.ListenAddr())
    fmt.Printf("HTTP:   %s\n", pb.HTTPAddr())

    // Proxy is now running...
}
```

## API Overview

### Share Link Parser

**NewParser()** - Create a new share link parser

```go
parser := singerbox.NewParser()
```

**Parse(link)** - Parse any supported share link format

```go
outbound, err := parser.Parse("vless://...")
```

**Protocol-specific parsers:**
- `ParseVLESS(link)` - Parse VLESS links
- `ParseVMess(link)` - Parse VMess links
- `ParseShadowsocks(link)` - Parse Shadowsocks links
- `ParseTrojan(link)` - Parse Trojan links
- `ParseSOCKS(link)` - Parse SOCKS5 links
- `ParseHTTP(link)` - Parse HTTP/HTTPS links

### Proxy Manager

**NewProxyBox(config)** - Create a new proxy instance

```go
pb, err := singerbox.NewProxyBox(singerbox.ProxyBoxConfig{
    Outbound:   outbound,      // sing-box outbound config
    ListenAddr: "127.0.0.1:1080", // SOCKS5 address (optional)
    HTTPPort:   1081,          // HTTP port (optional)
    LogLevel:   "info",        // Log level (optional)
})
```

**Start()** - Start the proxy

```go
err := pb.Start()
```

**Stop()** - Stop the proxy

```go
err := pb.Stop()
```

**Status methods:**
- `IsRunning()` - Check if proxy is running
- `ListenAddr()` - Get SOCKS5 address
- `HTTPAddr()` - Get HTTP address

## Supported Protocols

| Protocol | Format | Features |
|----------|--------|----------|
| **VLESS** | `vless://uuid@server:port?params` | TLS, Reality, WebSocket, gRPC, HTTP/2 |
| **VMess** | `vmess://base64json` | All transports and security options |
| **Shadowsocks** | `ss://method:pass@server:port` | All encryption methods |
| **Trojan** | `trojan://pass@server:port?params` | TLS and transports |
| **SOCKS5** | `socks5://[user:pass@]server:port` | With/without authentication |
| **HTTP/HTTPS** | `http[s]://[user:pass@]server:port` | Basic and authenticated |

## Examples

### Reality Protocol

```go
parser := singerbox.NewParser()
outbound, _ := parser.Parse(
    "vless://uuid@server:443?security=reality&pbk=publicKey&sid=shortID&sni=example.com&fp=chrome"
)

// Reality with custom fingerprint
pb, _ := singerbox.NewProxyBox(singerbox.ProxyBoxConfig{
    Outbound: outbound,
})
pb.Start()
```

### Multiple Proxy Instances

```go
// Create two proxies on different ports
pb1, _ := singerbox.NewProxyBox(singerbox.ProxyBoxConfig{
    Outbound:   outbound1,
    ListenAddr: "127.0.0.1:1080",
    HTTPPort:   1081,
})

pb2, _ := singerbox.NewProxyBox(singerbox.ProxyBoxConfig{
    Outbound:   outbound2,
    ListenAddr: "127.0.0.1:2080",
    HTTPPort:   2081,
})

pb1.Start()
pb2.Start()

// Both proxies running simultaneously
```

### Custom Configuration

```go
pb, _ := singerbox.NewProxyBox(singerbox.ProxyBoxConfig{
    Outbound:   outbound,
    ListenAddr: "0.0.0.0:9050",  // Listen on all interfaces
    HTTPPort:   9051,
    LogLevel:   "debug",          // Verbose logging
})
```

## Command-Line Tool

A proxy-tunnel CLI tool is included in `cmd/proxy-tunnel/`:

```bash
# Build
go build -tags "with_utls" -o proxy-tunnel ./cmd/proxy-tunnel/

# Or use Makefile
make build-minimal

# Run
./proxy-tunnel -link 'vless://uuid@server:443?security=tls&type=ws'
```

## Testing

```bash
# Run all tests
go test

# With coverage
go test -cover
# Output: coverage: 91.8% of statements

# Verbose
go test -v

# Benchmarks
go test -bench=.
```

## Build Tags

When building applications that use this library, include appropriate tags:

```bash
# Minimal (Reality support)
go build -tags "with_utls" your-app.go

# Full features
go build -tags "with_quic,with_utls,with_wireguard,with_dhcp,with_clash_api" your-app.go
```

## Requirements

- Go 1.21 or later
- sing-box library (automatically installed)

## License

This project uses the sing-box library which is licensed under GPLv3 or later.

## Contributing

Contributions are welcome! Please ensure:
- Tests pass: `go test`
- Code is formatted: `go fmt`
- Coverage is maintained: `go test -cover`

## Documentation

See [pkg.go.dev](https://pkg.go.dev/github.com/lunardoesdev/singerbox) for complete API documentation.

## Project Structure

```
singerbox/
├── sharelink.go               # Share link parser
├── proxybox.go                # Proxy manager
├── *_test.go                  # Comprehensive tests
├── cmd/
│   └── proxy-tunnel/          # CLI tool
│       └── main.go
├── go.mod
└── README.md
```

## Credits

Built on top of [sing-box](https://github.com/SagerNet/sing-box) by SagerNet.
