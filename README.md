<div align="center">

# 🎭 SingerBox

### Parse proxy share links and manage sing-box instances with ease

[![Go Reference](https://pkg.go.dev/badge/github.com/lunardoesdev/singerbox.svg)](https://pkg.go.dev/github.com/lunardoesdev/singerbox)
[![Go Report Card](https://goreportcard.com/badge/github.com/lunardoesdev/singerbox)](https://goreportcard.com/report/github.com/lunardoesdev/singerbox)
[![Coverage](https://img.shields.io/badge/coverage-83.1%25-brightgreen.svg)](.)

</div>


```go
// One-line proxy setup from any share link
proxy, _ := singerbox.FromSharedLink(
    "vless://550e8400-e29b-41d4-a716-446655440000@server:443?security=reality&pbk=key...",
    singerbox.ProxyConfig{
        ListenAddr: "127.0.0.1:1080",  // Optional: default is "127.0.0.1:1080"
        LogLevel:   "info",             // Optional: default is "panic" (silent)
    },
)
defer proxy.Stop() // 🚀 Proxy is now running!
```

<div align="center">

[Features](#-features) • [Installation](#-installation) • [Quick Start](#-quick-start) • [Examples](#-examples) • [Documentation](#-documentation)

</div>

---

## 🎯 What is this?

**SingerBox** is a Go library that makes working with proxy share links dead simple. Parse any proxy format (VLESS, VMess, Shadowsocks, Trojan, etc.) and start a local proxy with just a few lines of code.

Built on top of the powerful [sing-box](https://github.com/SagerNet/sing-box) proxy platform.

## ✨ Features

### 🔗 Parse Any Proxy Link
Parse share links from any proxy protocol:
- **VLESS** - Including Reality protocol with uTLS
- **VMess** - Base64 JSON config support
- **Shadowsocks** - All encryption methods
- **Trojan** - With TLS and transports
- **SOCKS5** - With or without authentication
- **HTTP/HTTPS** - Basic and authenticated proxies

### 🚀 Manage Proxies Programmatically
- Start/stop proxy instances with simple API
- Context-aware methods with timeout/cancellation support
- Mixed proxy mode supporting both SOCKS5 and HTTP on a single port
- Run multiple proxy instances simultaneously
- Full control over configuration

### 💪 Production Ready
- **83% test coverage** with 70+ test cases
- Input validation (UUID format, required fields, size limits)
- Clean, well-documented API
- Fast and efficient
- Battle-tested with sing-box

## 📦 Installation

```bash
go get github.com/lunardoesdev/singerbox
```

## 🚀 Quick Start

### Simplest Way: One-Line Proxy Setup

```go
package main

import (
    "fmt"
    "os"
    "os/signal"
    "github.com/lunardoesdev/singerbox"
)

func main() {
    // Create and start proxy from any share link (replace with your actual server)
    proxy, err := singerbox.FromSharedLink(
        "ss://aes-256-gcm:mypassword@your-server.com:8388",
        singerbox.ProxyConfig{
            ListenAddr: "127.0.0.1:1080",
            // LogLevel: "info",  // Uncomment to see connection logs
        },
    )
    if err != nil {
        panic(err)
    }
    defer proxy.Stop()

    fmt.Println("🚀 Proxy is running!")
    fmt.Printf("   Mixed (SOCKS5/HTTP): %s\n", proxy.ListenAddr())
    fmt.Println("\nPress Ctrl+C to stop...")

    // Wait for interrupt
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)
    <-c
}
```

### Advanced: Parse Links Separately (if you need parsing functionality)

If you need to parse share links without starting a proxy, use the parsing API:

```go
package main

import (
    "fmt"
    "github.com/lunardoesdev/singerbox"
)

func main() {
    // Just parse the link
    link := "vless://550e8400-e29b-41d4-a716-446655440000@example.com:443?type=ws&security=tls"
    outbound, err := singerbox.Parse(link)
    if err != nil {
        panic(err)
    }

    fmt.Printf("✓ Successfully parsed %s proxy\n", outbound.Type)
    fmt.Printf("  Tag: %s\n", outbound.Tag)

    // Now you can manually create and start the proxy if needed
    // pb, _ := singerbox.NewProxyBox(singerbox.ProxyBoxConfig{Outbound: outbound})
    // pb.Start()
}
```

## 💡 Examples

### Quick Proxy Setup

```go
// Start a proxy in one line
proxy, _ := singerbox.FromSharedLink(
    "ss://aes-256-gcm:password@your-server.com:8388",
    singerbox.ProxyConfig{
        ListenAddr: "127.0.0.1:1080",
        LogLevel:   "info",  // Enable logging
    },
)
defer proxy.Stop()

fmt.Printf("✓ Proxy running on %s\n", proxy.ListenAddr())
```

### Multiple Proxies from Different Servers

```go
links := []string{
    "ss://aes-256-gcm:pass1@server1.com:8388",
    "trojan://pass2@server2.com:443",
    "vless://550e8400-e29b-41d4-a716-446655440000@server3.com:443?security=tls",
}

proxies := []*singerbox.ProxyBox{}
for i, link := range links {
    proxy, err := singerbox.FromSharedLink(link, singerbox.ProxyConfig{
        ListenAddr: fmt.Sprintf("127.0.0.1:%d", 1080+i),
    })
    if err != nil {
        fmt.Printf("❌ Failed: %v\n", err)
        continue
    }
    proxies = append(proxies, proxy)
    fmt.Printf("✓ Proxy %d running on %s\n", i+1, proxy.ListenAddr())
}

// Clean up
for _, p := range proxies {
    defer p.Stop()
}
```

### Parse Links Without Starting Proxy

If you just need to parse links (advanced use case):

```go
links := []string{
    "vless://550e8400-e29b-41d4-a716-446655440000@server:443?security=reality&pbk=key&sid=id",
    "vmess://eyJ2IjoiMiIsInBzIjoidGVzdCIsImFkZCI6InNlcnZlciIsImlkIjoiNTUwZTg0MDAtZTI5Yi00MWQ0LWE3MTYtNDQ2NjU1NDQwMDAwIn0=",
    "ss://aes-256-gcm:password@server:8388",
}

for _, link := range links {
    outbound, err := singerbox.Parse(link)
    if err != nil {
        fmt.Printf("❌ Failed to parse: %v\n", err)
        continue
    }
    fmt.Printf("✓ Parsed %s proxy\n", outbound.Type)
}
```

### Reality Protocol with Custom Fingerprint

```go
// Parse VLESS Reality link
// Note: pbk (public key) is required for Reality, sid (short ID) is optional
link := "vless://550e8400-e29b-41d4-a716-446655440000@server:443?" +
        "security=reality&" +
        "pbk=publicKey123&" +   // Required for Reality
        "sid=shortID&" +         // Optional
        "sni=www.example.com&" +
        "fp=firefox"  // uTLS fingerprint

outbound, _ := singerbox.Parse(link)

// Start proxy with Reality support
pb, _ := singerbox.NewProxyBox(singerbox.ProxyBoxConfig{
    Outbound: outbound,
    LogLevel: "info",
})

pb.Start()
defer pb.Stop()

fmt.Println("✓ Reality protocol active with uTLS!")
```

### Multiple Proxy Instances

```go

// Parse two different proxies
shadowsocks, _ := singerbox.Parse("ss://aes-256-gcm:pass1@server1:8388")
trojan, _ := singerbox.Parse("trojan://pass2@server2:443")

// Start first proxy
proxy1, _ := singerbox.NewProxyBox(singerbox.ProxyBoxConfig{
    Outbound:   shadowsocks,
    ListenAddr: "127.0.0.1:1080",
})
proxy1.Start()

// Start second proxy on different port
proxy2, _ := singerbox.NewProxyBox(singerbox.ProxyBoxConfig{
    Outbound:   trojan,
    ListenAddr: "127.0.0.1:2080",
})
proxy2.Start()

fmt.Println("✓ Running 2 proxies simultaneously")
fmt.Printf("  Proxy 1 (Shadowsocks): %s\n", proxy1.ListenAddr())
fmt.Printf("  Proxy 2 (Trojan):      %s\n", proxy2.ListenAddr())

// Don't forget to stop them
defer proxy1.Stop()
defer proxy2.Stop()
```

### Error Handling

```go

// Parse with error handling
outbound, err := singerbox.Parse("vless://invalid-link")
if err != nil {
    fmt.Printf("Parse failed: %v\n", err)
    return
}

// Create proxy with validation
pb, err := singerbox.NewProxyBox(singerbox.ProxyBoxConfig{
    Outbound: outbound,
})
if err != nil {
    fmt.Printf("Configuration error: %v\n", err)
    return
}

// Start with error handling
if err := pb.Start(); err != nil {
    fmt.Printf("Failed to start: %v\n", err)
    return
}
defer pb.Stop()

// Check if running
if pb.IsRunning() {
    fmt.Println("✓ Proxy is active and ready")
}
```

### Custom Configuration

```go
// Create a proxy accessible from network with logging
pb, _ := singerbox.NewProxyBox(singerbox.ProxyBoxConfig{
    Outbound:   outbound,
    ListenAddr: "0.0.0.0:9050",    // Listen on all interfaces
    LogLevel:   "info",             // Enable informational logging (default is "panic" - silent)
})

pb.Start()
defer pb.Stop()

fmt.Println("✓ Proxy accessible from network!")
fmt.Println("  Other devices can connect to:")
fmt.Printf("  Mixed (SOCKS5/HTTP): YOUR_IP:9050\n")
```

### Logging Configuration

By default, the proxy operates silently (LogLevel: `"panic"`). Enable logging for debugging:

```go
// Silent operation (default)
pb, _ := singerbox.NewProxyBox(singerbox.ProxyBoxConfig{
    Outbound: outbound,
    // LogLevel not set - defaults to "panic" (silent)
})

// Enable info logging to see connection activity
pb, _ := singerbox.NewProxyBox(singerbox.ProxyBoxConfig{
    Outbound: outbound,
    LogLevel: "info",  // Shows connection info
})

// Enable debug logging for troubleshooting
pb, _ := singerbox.NewProxyBox(singerbox.ProxyBoxConfig{
    Outbound: outbound,
    LogLevel: "debug",  // Shows detailed debug info
})
```

## 📚 API Documentation

### Quick Setup API (Recommended)

#### `FromSharedLink(link string, cfg ProxyConfig) (*ProxyBox, error)`
**The easiest way to set up a proxy.** Parses the link, creates, and starts the proxy in one call.

```go
proxy, err := singerbox.FromSharedLink(
    "vless://550e8400-e29b-41d4-a716-446655440000@server:443?security=tls",
    singerbox.ProxyConfig{
        ListenAddr: "127.0.0.1:1080",  // Optional: defaults to "127.0.0.1:1080"
        LogLevel:   "info",             // Optional: defaults to "panic" (silent)
    },
)
if err != nil {
    // Handle error
}
defer proxy.Stop()  // Don't forget to stop when done!

// Proxy is already running and ready to use
fmt.Printf("Proxy listening on: %s\n", proxy.ListenAddr())
```

**ProxyConfig fields:**
```go
type ProxyConfig struct {
    ListenAddr string  // Optional: Mixed proxy address (default: "127.0.0.1:1080")
    LogLevel   string  // Optional: Logging level (default: "panic" - silent)
}
```

### Advanced APIs (for special use cases)

#### Parsing API

Use these if you need to parse share links without starting a proxy:

##### `Parse(link string) (option.Outbound, error)`
Parses any supported share link and returns sing-box outbound config.

```go
outbound, err := singerbox.Parse("vless://550e8400-e29b-41d4-a716-446655440000@server:443?security=tls")
if err != nil {
    // Handle error
}
```

##### Protocol-Specific Functions
```go
singerbox.ParseVLESS(link string) (option.Outbound, error)
singerbox.ParseVMess(link string) (option.Outbound, error)
singerbox.ParseShadowsocks(link string) (option.Outbound, error)
singerbox.ParseTrojan(link string) (option.Outbound, error)
singerbox.ParseSOCKS(link string) (option.Outbound, error)
singerbox.ParseHTTP(link string) (option.Outbound, error)
```

#### ProxyBox API (manual control)

#### `NewProxyBox(config ProxyBoxConfig) (*ProxyBox, error)`
Creates a new proxy instance.

```go
type ProxyBoxConfig struct {
    Outbound   option.Outbound  // Required: sing-box outbound config
    ListenAddr string           // Optional: Mixed proxy address (default: "127.0.0.1:1080")
    LogLevel   string           // Optional: Logging level (default: "panic" - silent)
}
```

**Log Levels** (from most to least verbose):
- `"trace"` - Very detailed debugging information
- `"debug"` - Detailed debugging information
- `"info"` - General informational messages
- `"warn"` - Warning messages
- `"error"` - Error messages
- `"fatal"` - Fatal errors (will terminate)
- `"panic"` - Critical errors only (default - silent operation)

#### Methods

```go
// Basic methods
Start() error              // Start the proxy
Stop() error               // Stop the proxy
IsRunning() bool           // Check if running
ListenAddr() string        // Get mixed proxy address (supports both SOCKS5 and HTTP)
Config() option.Options    // Get sing-box config
Outbound() option.Outbound // Get outbound config

// Context-aware methods (for timeout/cancellation control)
StartContext(ctx context.Context) error  // Start with context
StopContext(ctx context.Context) error   // Stop with context
```

#### Using Context for Timeouts

```go
import "context"

// Start with a 5-second timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

err := pb.StartContext(ctx)
if err == context.DeadlineExceeded {
    fmt.Println("Startup timed out")
}

// Stop with cancellation support
stopCtx, stopCancel := context.WithTimeout(context.Background(), 3*time.Second)
defer stopCancel()

err = pb.StopContext(stopCtx)
```

### Input Validation

The library validates input to ensure correct configuration:

- **UUID format**: VLESS and VMess UUIDs must be valid UUID format (e.g., `550e8400-e29b-41d4-a716-446655440000`)
- **Reality public key**: The `pbk` parameter is required when using Reality security
- **Port range**: Ports must be between 1-65535
- **Link size**: Share links are limited to 64KB to prevent abuse

```go
// Invalid UUID will return an error
_, err := singerbox.Parse("vless://invalid-uuid@server:443")
// Error: invalid UUID format in VLESS link

// Missing Reality public key will return an error
_, err := singerbox.Parse("vless://550e8400-e29b-41d4-a716-446655440000@server:443?security=reality")
// Error: missing public key (pbk) for Reality in VLESS link
```

## 🔧 Command-Line Tool

A CLI tool is included for quick proxy setup:

```bash
# Build
make build-minimal

# Run with any share link
./proxy-tunnel -link 'vless://550e8400-e29b-41d4-a716-446655440000@server:443?security=tls&type=ws'

# Custom listen address
./proxy-tunnel -link 'ss://...' -listen 0.0.0.0:9050
```

**Usage:**
```
./proxy-tunnel -link <share-link> [-listen <addr:port>]

Options:
  -link        Proxy share link (required)
  -listen      Mixed proxy listen address (default: 127.0.0.1:1080)

Examples:
  ./proxy-tunnel -link 'vless://550e8400-e29b-41d4-a716-446655440000@server:443?type=ws&security=tls'
  ./proxy-tunnel -link 'ss://method:password@server:8388'
  ./proxy-tunnel -link 'trojan://password@server:443' -listen 0.0.0.0:1080
```

## 📖 Supported Protocols

| Protocol | Format | Features |
|----------|--------|----------|
| **VLESS** | `vless://uuid@server:port?params` | TLS, Reality, WebSocket, gRPC, HTTP/2, XTLS flow |
| **VMess** | `vmess://base64json` | All transports, TLS, WebSocket, gRPC, HTTP/2 |
| **Shadowsocks** | `ss://method:pass@server:port` | All ciphers |
| **Trojan** | `trojan://pass@server:port?params` | TLS, WebSocket, gRPC |
| **SOCKS5** | `socks5://[user:pass@]server:port` | Authentication, UDP relay |
| **HTTP/HTTPS** | `http[s]://[user:pass@]server:port` | Basic auth, TLS |

## 🧪 Testing

```bash
# Run all tests
make test

# Run with coverage
make test-cover
# Output: coverage: 83.1% of statements

# Run specific tests
go test -v -run TestParseVLESS

# Benchmarks
go test -bench=.
```

## 🏗️ Building

```bash
# Build CLI tool with all features
make build

# Minimal build (uTLS only - for Reality support)
make build-minimal

# Basic build (smallest binary)
make build-basic
```

### Build Tags

**Recommended:** Build with all feature tags for full protocol support:

```bash
go build -tags "with_quic,with_utls,with_wireguard,with_dhcp,with_clash_api" your-app.go
```

This enables all sing-box features including Reality, QUIC, WireGuard, DHCP DNS, and Clash API compatibility.

Available tags:
- `with_utls` - **Required for Reality protocol**, uTLS fingerprinting
- `with_quic` - QUIC transport support
- `with_wireguard` - WireGuard protocol
- `with_dhcp` - DHCP DNS server
- `with_clash_api` - Clash API compatibility

**Note:** When using Reality protocol specifically, you must at minimum include `with_utls`:

```bash
go build -tags "with_utls" your-app.go
```

## 🎯 Use Cases

- 🔄 **Proxy Switcher** - Switch between multiple proxies dynamically
- 📋 **Subscription Manager** - Parse and manage subscription URLs
- 🛠️ **Network Tools** - Build custom proxy utilities
- 🧪 **Testing** - Test apps through different proxy configurations
- 🤖 **Automation** - Automate proxy setup and management
- 📱 **VPN Apps** - Create custom VPN client applications

## 📋 Requirements

- Go 1.23 or higher
- sing-box (automatically installed via go modules)

## 🤝 Contributing

Contributions welcome! Please ensure:

1. Tests pass: `make test`
2. Code is formatted: `go fmt`
3. Coverage maintained: `make test-cover`

**Steps:**
1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing`)
5. Open a Pull Request

## 📝 License

Licensed under **GPLv3 or later** (due to sing-box dependency).

## 🙏 Acknowledgments

- Built on [sing-box](https://github.com/SagerNet/sing-box) by SagerNet
- Inspired by the Go proxy community

## 📞 Support

- 🐛 **Bug Reports**: [GitHub Issues](https://github.com/lunardoesdev/singerbox/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/lunardoesdev/singerbox/discussions)
- 📖 **Documentation**: [pkg.go.dev](https://pkg.go.dev/github.com/lunardoesdev/singerbox)

---

<div align="center">

**If you find this useful, please ⭐ star the repo!**

Made with ❤️ by the community

[Report Bug](https://github.com/lunardoesdev/singerbox/issues) • [Request Feature](https://github.com/lunardoesdev/singerbox/issues) • [Documentation](https://pkg.go.dev/github.com/lunardoesdev/singerbox)

</div>
