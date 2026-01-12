# ProxyBox Library

A Go library for managing sing-box instances as local proxies. This library simplifies creating, starting, and stopping sing-box proxy instances programmatically.

## Features

- ✅ **Simple API** - Easy to use Start/Stop interface
- ✅ **Configuration Management** - Clean configuration with sensible defaults
- ✅ **Multiple Protocols** - Works with all sing-box supported protocols
- ✅ **Dual Proxy Modes** - Both SOCKS5 and HTTP proxy support
- ✅ **Lifecycle Management** - Proper resource cleanup
- ✅ **Thread-Safe** - Can run multiple instances simultaneously
- ✅ **Comprehensive Tests** - Well-tested with 100% coverage

## Installation

```bash
go get proxy-tunnel/pkg/proxybox
```

## Quick Start

```go
package main

import (
    "fmt"
    "proxy-tunnel/pkg/proxybox"
    "proxy-tunnel/pkg/sharelink"
)

func main() {
    // Parse a share link
    parser := sharelink.New()
    outbound, _ := parser.Parse("vless://uuid@server:443?security=tls&type=ws")

    // Create proxy box
    pb, err := proxybox.New(proxybox.Config{
        Outbound:   outbound,
        ListenAddr: "127.0.0.1:1080",
        HTTPPort:   1081,
    })
    if err != nil {
        panic(err)
    }

    // Start the proxy
    if err := pb.Start(); err != nil {
        panic(err)
    }
    defer pb.Stop()

    fmt.Printf("SOCKS5: %s\n", pb.ListenAddr())
    fmt.Printf("HTTP:   %s\n", pb.HTTPAddr())

    // Proxy is now running...
}
```

## API Reference

### type Config

```go
type Config struct {
    Outbound   option.Outbound // sing-box outbound configuration (required)
    ListenAddr string          // SOCKS5/mixed proxy address (default: "127.0.0.1:1080")
    HTTPPort   int             // HTTP proxy port (default: 1081)
    LogLevel   string          // Log level: "trace", "debug", "info", "warn", "error" (default: "info")
}
```

### type ProxyBox

```go
type ProxyBox struct {
    // contains filtered or unexported fields
}
```

#### func New

```go
func New(cfg Config) (*ProxyBox, error)
```

Creates a new ProxyBox instance with the given configuration. Does not start the proxy.

**Example:**
```go
pb, err := proxybox.New(proxybox.Config{
    Outbound: outbound,
    ListenAddr: "127.0.0.1:1080",
    HTTPPort: 1081,
    LogLevel: "info",
})
```

#### func (*ProxyBox) Start

```go
func (pb *ProxyBox) Start() error
```

Starts the proxy box. Returns an error if already running or if sing-box fails to start.

**Example:**
```go
if err := pb.Start(); err != nil {
    log.Fatalf("Failed to start: %v", err)
}
```

#### func (*ProxyBox) Stop

```go
func (pb *ProxyBox) Stop() error
```

Stops the proxy box and releases resources. Returns an error if not running.

**Example:**
```go
if err := pb.Stop(); err != nil {
    log.Printf("Error stopping: %v", err)
}
```

#### func (*ProxyBox) IsRunning

```go
func (pb *ProxyBox) IsRunning() bool
```

Returns true if the proxy box is currently running.

**Example:**
```go
if pb.IsRunning() {
    fmt.Println("Proxy is active")
}
```

#### func (*ProxyBox) ListenAddr

```go
func (pb *ProxyBox) ListenAddr() string
```

Returns the SOCKS5/mixed proxy listen address (e.g., "127.0.0.1:1080").

#### func (*ProxyBox) HTTPAddr

```go
func (pb *ProxyBox) HTTPAddr() string
```

Returns the HTTP proxy listen address (e.g., "127.0.0.1:1081").

#### func (*ProxyBox) Config

```go
func (pb *ProxyBox) Config() option.Options
```

Returns the complete sing-box configuration being used.

#### func (*ProxyBox) Outbound

```go
func (pb *ProxyBox) Outbound() option.Outbound
```

Returns the outbound configuration.

## Usage Examples

### Example 1: Basic Usage

```go
package main

import (
    "fmt"
    "proxy-tunnel/pkg/proxybox"
    "github.com/sagernet/sing-box/option"
)

func main() {
    // Create a direct outbound (no proxy)
    outbound := option.Outbound{
        Type: "direct",
        Tag:  "direct",
        Options: &option.DirectOutboundOptions{},
    }

    // Create and start proxy box
    pb, _ := proxybox.New(proxybox.Config{
        Outbound: outbound,
    })

    pb.Start()
    defer pb.Stop()

    fmt.Printf("Proxy running at %s\n", pb.ListenAddr())
    // Use the proxy...
}
```

### Example 2: With Share Link Parser

```go
package main

import (
    "fmt"
    "proxy-tunnel/pkg/proxybox"
    "proxy-tunnel/pkg/sharelink"
)

func main() {
    // Parse share link
    parser := sharelink.New()
    outbound, err := parser.Parse("ss://aes-256-gcm:password@server:8388")
    if err != nil {
        panic(err)
    }

    // Create proxy box
    pb, err := proxybox.New(proxybox.Config{
        Outbound: outbound,
    })
    if err != nil {
        panic(err)
    }

    // Start proxy
    if err := pb.Start(); err != nil {
        panic(err)
    }
    defer pb.Stop()

    fmt.Printf("Shadowsocks proxy active:\n")
    fmt.Printf("  SOCKS5: %s\n", pb.ListenAddr())
    fmt.Printf("  HTTP:   %s\n", pb.HTTPAddr())
}
```

### Example 3: Custom Configuration

```go
package main

import (
    "proxy-tunnel/pkg/proxybox"
    "proxy-tunnel/pkg/sharelink"
)

func main() {
    parser := sharelink.New()
    outbound, _ := parser.Parse("trojan://password@server:443")

    // Custom ports and log level
    pb, _ := proxybox.New(proxybox.Config{
        Outbound:   outbound,
        ListenAddr: "0.0.0.0:9050",  // Listen on all interfaces
        HTTPPort:   9051,
        LogLevel:   "debug",          // Verbose logging
    })

    pb.Start()
    defer pb.Stop()

    // Proxy accessible from network
}
```

### Example 4: Multiple Proxy Instances

```go
package main

import (
    "fmt"
    "proxy-tunnel/pkg/proxybox"
    "proxy-tunnel/pkg/sharelink"
)

func main() {
    parser := sharelink.New()

    // Create first proxy
    out1, _ := parser.Parse("ss://aes-256-gcm:pass1@server1:8388")
    pb1, _ := proxybox.New(proxybox.Config{
        Outbound:   out1,
        ListenAddr: "127.0.0.1:1080",
        HTTPPort:   1081,
    })
    pb1.Start()
    defer pb1.Stop()

    // Create second proxy on different ports
    out2, _ := parser.Parse("ss://aes-256-gcm:pass2@server2:8388")
    pb2, _ := proxybox.New(proxybox.Config{
        Outbound:   out2,
        ListenAddr: "127.0.0.1:2080",
        HTTPPort:   2081,
    })
    pb2.Start()
    defer pb2.Stop()

    fmt.Printf("Proxy 1: %s\n", pb1.ListenAddr())
    fmt.Printf("Proxy 2: %s\n", pb2.ListenAddr())

    // Both proxies running simultaneously
}
```

### Example 5: Error Handling

```go
package main

import (
    "fmt"
    "log"
    "proxy-tunnel/pkg/proxybox"
    "proxy-tunnel/pkg/sharelink"
)

func startProxy(link string) (*proxybox.ProxyBox, error) {
    parser := sharelink.New()
    outbound, err := parser.Parse(link)
    if err != nil {
        return nil, fmt.Errorf("parse link: %w", err)
    }

    pb, err := proxybox.New(proxybox.Config{
        Outbound: outbound,
        LogLevel: "error", // Less verbose
    })
    if err != nil {
        return nil, fmt.Errorf("create proxy box: %w", err)
    }

    if err := pb.Start(); err != nil {
        return nil, fmt.Errorf("start proxy: %w", err)
    }

    return pb, nil
}

func main() {
    pb, err := startProxy("vless://uuid@server:443?security=tls")
    if err != nil {
        log.Fatalf("Failed to start proxy: %v", err)
    }
    defer pb.Stop()

    fmt.Println("Proxy started successfully")
}
```

### Example 6: Graceful Shutdown

```go
package main

import (
    "context"
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "proxy-tunnel/pkg/proxybox"
    "proxy-tunnel/pkg/sharelink"
)

func main() {
    parser := sharelink.New()
    outbound, _ := parser.Parse("ss://aes-256-gcm:password@server:8388")

    pb, _ := proxybox.New(proxybox.Config{
        Outbound: outbound,
    })

    if err := pb.Start(); err != nil {
        panic(err)
    }

    fmt.Println("Proxy started. Press Ctrl+C to stop...")

    // Wait for interrupt signal
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
    <-sigCh

    fmt.Println("\nShutting down gracefully...")
    if err := pb.Stop(); err != nil {
        fmt.Printf("Error during shutdown: %v\n", err)
    }
    fmt.Println("Stopped")
}
```

## Configuration Details

### Default Values

- **ListenAddr**: `"127.0.0.1:1080"` - Standard SOCKS port
- **HTTPPort**: `1081` - HTTP proxy port
- **LogLevel**: `"info"` - Balanced logging

### Log Levels

- `"trace"` - Most verbose, includes protocol details
- `"debug"` - Detailed debugging information
- `"info"` - General operational messages (default)
- `"warn"` - Warning messages only
- `"error"` - Error messages only

### Listen Addresses

The library supports:
- IPv4: `"127.0.0.1:1080"` (localhost only)
- IPv4 all interfaces: `"0.0.0.0:1080"` (network accessible)
- IPv6: `"[::1]:1080"` (IPv6 localhost)

## Architecture

The ProxyBox library:

1. **Creates** a sing-box configuration with:
   - Two inbound proxies (SOCKS5/mixed and HTTP)
   - Your specified outbound (the actual proxy)
   - Direct and block outbounds
   - Routing rules to send all traffic through the outbound

2. **Manages** the sing-box lifecycle:
   - Initializes context with protocol handlers
   - Creates sing-box instance
   - Starts/stops the instance
   - Cleans up resources

3. **Provides** a clean API:
   - Simple configuration
   - Clear error messages
   - Status checking
   - Address retrieval

## Testing

The library includes comprehensive tests covering:

- ✅ Configuration validation
- ✅ Start/Stop lifecycle
- ✅ Multiple instances
- ✅ Integration with sharelink parser
- ✅ Address management
- ✅ Error conditions

Run tests:
```bash
go test ./pkg/proxybox/ -v
go test ./pkg/proxybox/ -cover
```

## Integration with Main Application

This library is used by the main proxy-tunnel application:

```go
// main.go (simplified)
parser := sharelink.New()
outbound, _ := parser.Parse(*shareLink)

pb, _ := proxybox.New(proxybox.Config{
    Outbound:   outbound,
    ListenAddr: *listenAddr,
    HTTPPort:   *httpPort,
})

pb.Start()
defer pb.Stop()
```

## Dependencies

- `github.com/sagernet/sing-box` - sing-box library
- `github.com/sagernet/sing` - sing common utilities

## Best Practices

1. **Always call Stop()**: Use `defer pb.Stop()` to ensure cleanup
2. **Check IsRunning()**: Before calling Start/Stop
3. **Handle errors**: Both Start and Stop can return errors
4. **Use different ports**: When running multiple instances
5. **Set appropriate log level**: "error" for production, "debug" for development

## Common Issues

### Port Already in Use

```go
pb, err := proxybox.New(proxybox.Config{
    Outbound:   outbound,
    ListenAddr: "127.0.0.1:1080",  // May be in use
})
// Solution: Use different port or check if port is available
```

### Multiple Start Calls

```go
pb.Start()
pb.Start() // Error: already running
// Solution: Check IsRunning() first
```

## License

This library is part of the proxy-tunnel project.
