<div align="center">

# 🎭 SingerBox

### Parse proxy share links and manage sing-box instances with ease

</div>

[![Go Reference](https://pkg.go.dev/badge/github.com/lunardoesdev/singerbox.svg)](https://pkg.go.dev/github.com/lunardoesdev/singerbox)
[![Go Report Card](https://goreportcard.com/badge/github.com/lunardoesdev/singerbox)](https://goreportcard.com/report/github.com/lunardoesdev/singerbox)
[![License: GPLv3](https://img.shields.io/badge/License-GPLv3-blue.svg)](LICENSE)
[![Coverage](https://img.shields.io/badge/coverage-91.8%25-brightgreen.svg)](.)

```
parser := singerbox.NewParser()
outbound, _ := parser.Parse("vless://uuid@server:443?security=reality&pbk=key...")
pb, _ := singerbox.NewProxyBox(singerbox.ProxyBoxConfig{
    Outbound:   outbound,
    ListenAddr: "127.0.0.1:1080",
})
pb.Start() // 🚀 Proxy is now running!
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
- Mixed proxy mode supporting both SOCKS5 and HTTP on a single port
- Run multiple proxy instances simultaneously
- Full control over configuration

### 💪 Production Ready
- **91.8% test coverage** with 70+ test cases
- Clean, well-documented API
- Fast and efficient
- Battle-tested with sing-box

## 📦 Installation

```bash
go get github.com/lunardoesdev/singerbox
```

## 🚀 Quick Start

### Example 1: Parse a Share Link

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
    link := "vless://550e8400-e29b-41d4-a716-446655440000@example.com:443?type=ws&security=tls"
    outbound, err := parser.Parse(link)
    if err != nil {
        panic(err)
    }

    fmt.Printf("✓ Successfully parsed %s proxy\n", outbound.Type)
    fmt.Printf("  Tag: %s\n", outbound.Tag)
}
```

### Example 2: Start a Local Proxy

```go
package main

import (
    "fmt"
    "os"
    "os/signal"
    "github.com/lunardoesdev/singerbox"
)

func main() {
    // Parse the share link
    parser := singerbox.NewParser()
    outbound, _ := parser.Parse("ss://aes-256-gcm:mypassword@server.com:8388")

    // Create and start the proxy on custom port
    pb, _ := singerbox.NewProxyBox(singerbox.ProxyBoxConfig{
        Outbound:   outbound,
        ListenAddr: "127.0.0.1:1080",  // Listen on port 1080 for SOCKS5 and HTTP
    })

    pb.Start()
    defer pb.Stop()

    fmt.Println("🚀 Proxy is running!")
    fmt.Printf("   Mixed (SOCKS5/HTTP): %s\n", pb.ListenAddr())
    fmt.Println("\nPress Ctrl+C to stop...")

    // Wait for interrupt
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)
    <-c
}
```

## 💡 Examples

### Parse Multiple Proxy Types

```go
parser := singerbox.NewParser()

links := []string{
    "vless://uuid@server:443?security=reality&pbk=key&sid=id",
    "vmess://eyJ2IjoiMiIsInBzIjoidGVzdCIsImFkZCI6InNlcnZlciJ9",
    "ss://aes-256-gcm:password@server:8388",
    "trojan://password@server:443",
    "socks5://user:pass@proxy:1080",
}

for _, link := range links {
    outbound, err := parser.Parse(link)
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
link := "vless://uuid@server:443?" +
        "security=reality&" +
        "pbk=publicKey123&" +
        "sid=shortID&" +
        "sni=www.example.com&" +
        "fp=firefox"  // uTLS fingerprint

parser := singerbox.NewParser()
outbound, _ := parser.Parse(link)

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
parser := singerbox.NewParser()

// Parse two different proxies
shadowsocks, _ := parser.Parse("ss://aes-256-gcm:pass1@server1:8388")
trojan, _ := parser.Parse("trojan://pass2@server2:443")

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
parser := singerbox.NewParser()

// Parse with error handling
outbound, err := parser.Parse("vless://invalid-link")
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
// Create a proxy accessible from network
pb, _ := singerbox.NewProxyBox(singerbox.ProxyBoxConfig{
    Outbound:   outbound,
    ListenAddr: "0.0.0.0:9050",    // Listen on all interfaces
    LogLevel:   "debug",           // Verbose logging
})

pb.Start()
defer pb.Stop()

fmt.Println("✓ Proxy accessible from network!")
fmt.Println("  Other devices can connect to:")
fmt.Printf("  Mixed (SOCKS5/HTTP): YOUR_IP:9050\n")
```

## 📚 API Documentation

### Parser API

#### `NewParser() *Parser`
Creates a new share link parser.

```go
parser := singerbox.NewParser()
```

#### `Parse(link string) (option.Outbound, error)`
Parses any supported share link and returns sing-box outbound config.

```go
outbound, err := parser.Parse("vless://...")
```

#### Protocol-Specific Methods
```go
ParseVLESS(link string) (option.Outbound, error)
ParseVMess(link string) (option.Outbound, error)
ParseShadowsocks(link string) (option.Outbound, error)
ParseTrojan(link string) (option.Outbound, error)
ParseSOCKS(link string) (option.Outbound, error)
ParseHTTP(link string) (option.Outbound, error)
```

### ProxyBox API

#### `NewProxyBox(config ProxyBoxConfig) (*ProxyBox, error)`
Creates a new proxy instance.

```go
type ProxyBoxConfig struct {
    Outbound   option.Outbound  // Required: sing-box outbound config
    ListenAddr string           // Optional: Mixed proxy address (default: "127.0.0.1:1080")
    LogLevel   string           // Optional: "trace", "debug", "info", "warn", "error" (default: "info")
}
```

#### Methods

```go
Start() error              // Start the proxy
Stop() error               // Stop the proxy
IsRunning() bool           // Check if running
ListenAddr() string        // Get mixed proxy address (supports both SOCKS5 and HTTP)
Config() option.Options    // Get sing-box config
Outbound() option.Outbound // Get outbound config
```

## 🔧 Command-Line Tool

A CLI tool is included for quick proxy setup:

```bash
# Build
make build-minimal

# Run with any share link
./proxy-tunnel -link 'vless://uuid@server:443?security=tls&type=ws'

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
  ./proxy-tunnel -link 'vless://uuid@server:443?type=ws&security=tls'
  ./proxy-tunnel -link 'ss://method:password@server:8388'
  ./proxy-tunnel -link 'trojan://password@server:443' -listen 0.0.0.0:1080
```

## 📖 Supported Protocols

| Protocol | Format | Features |
|----------|--------|----------|
| **VLESS** | `vless://uuid@server:port?params` | TLS, Reality, WebSocket, gRPC, HTTP/2, XTLS flow |
| **VMess** | `vmess://base64json` | All transports, TLS, WebSocket, gRPC, HTTP/2 |
| **Shadowsocks** | `ss://method:pass@server:port` | All ciphers, SIP003 plugins |
| **Trojan** | `trojan://pass@server:port?params` | TLS, WebSocket, gRPC |
| **SOCKS5** | `socks5://[user:pass@]server:port` | Authentication, UDP relay |
| **HTTP/HTTPS** | `http[s]://[user:pass@]server:port` | Basic auth, TLS |

## 🧪 Testing

```bash
# Run all tests
make test

# Run with coverage
make test-cover
# Output: coverage: 91.8% of statements

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

- Go 1.21 or higher
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
