# Troubleshooting Guide

This guide covers common issues and their solutions when building and using proxy-tunnel.

## Build Errors

### Error: "uTLS, which is required by reality is not included in this build"

**Problem**: Building without the `with_utls` tag when using Reality protocol.

**Solution**:
```bash
# Rebuild with uTLS support
go build -tags "with_utls" -o proxy-tunnel

# Or use Makefile
make build
```

**Why**: Reality protocol requires uTLS for advanced TLS fingerprinting.

---

### Error: "uTLS is required by reality client"

**Problem**: Reality configuration missing uTLS options (this should be automatic now).

**Solution**: The code automatically enables uTLS when security=reality is detected. This error means:
1. The build doesn't have `with_utls` tag, OR
2. The share link is malformed

Verify your build:
```bash
# Check if built with uTLS
strings proxy-tunnel | grep -i utls

# Rebuild if needed
make build
```

Verify your share link format:
```bash
# Reality link must have:
# - security=reality
# - pbk=<base64-public-key>
# - sid=<short-id>
# - sni=<server-name>

# Example:
vless://uuid@server:443?security=reality&pbk=VALIDBASE64KEY&sid=abc123&sni=example.com
```

---

### Error: "missing endpoint registry in context"

**Problem**: Context not properly initialized with sing-box protocol registries.

**Solution**: This is fixed in the code by using `include.Context(ctx)`. If you're modifying the code, ensure you have:

```go
import "github.com/sagernet/sing-box/include"

ctx := context.Background()
ctx = include.Context(ctx) // This registers all protocols
```

---

### Error: "interface conversion panic"

**Problem**: Options structs not passed as pointers.

**Solution**: Ensure all Options fields use pointers:

```go
// Correct ✓
Options: &option.HTTPMixedInboundOptions{...}

// Wrong ✗
Options: option.HTTPMixedInboundOptions{...}
```

---

### Error: "QUIC is not included in this build"

**Problem**: Using Hysteria or Hysteria2 without QUIC support.

**Solution**:
```bash
go build -tags "with_quic,with_utls" -o proxy-tunnel
```

---

### Error: "cannot find module"

**Problem**: Dependencies not downloaded.

**Solution**:
```bash
go mod download
go mod tidy
```

---

### Error: Build takes too long

**Problem**: First build downloads many dependencies.

**Solution**:
- First build is always slower (downloading deps)
- Subsequent builds are much faster
- Use `go build -v` to see progress

---

## Runtime Errors

### Error: "Error parsing share link"

**Problem**: Invalid share link format.

**Causes**:
1. Missing protocol prefix (vless://, vmess://, etc.)
2. Malformed URL structure
3. Invalid base64 encoding (for vmess://)

**Solutions**:

```bash
# ✓ Correct - has protocol prefix
./proxy-tunnel -link 'vless://uuid@server:443?type=ws&security=tls'

# ✗ Wrong - missing protocol
./proxy-tunnel -link 'uuid@server:443?type=ws&security=tls'

# Use quotes to preserve special characters
./proxy-tunnel -link 'vmess://base64string=='

# Check for proper encoding
echo "your-vmess-config" | base64
```

---

### Error: "Error creating sing-box instance"

**Possible causes**:

1. **Invalid configuration**: Check your share link format
2. **Missing features**: Rebuild with appropriate tags
3. **Port already in use**: Change listen ports

**Solutions**:

```bash
# Try different ports
./proxy-tunnel -link 'your-link' -listen '127.0.0.1:8080' -http-port 8081

# Check if ports are in use
netstat -tlnp | grep -E '1080|1081'
lsof -i :1080
lsof -i :1081
```

---

### Error: "Error starting sing-box"

**Possible causes**:

1. **Ports already in use**
2. **Permission denied** (low port numbers)
3. **Network interface not available**

**Solutions**:

```bash
# Use higher port numbers (above 1024)
./proxy-tunnel -link 'your-link' -listen '127.0.0.1:8080' -http-port 8081

# Check port availability
ss -tlnp | grep -E '1080|1081'

# Run with sudo for low ports (< 1024) - not recommended
sudo ./proxy-tunnel -link 'your-link' -listen '0.0.0.0:80'
```

---

### Error: Connection refused when testing proxy

**Problem**: Proxy not actually started or wrong address.

**Solutions**:

```bash
# 1. Check if proxy is running
ps aux | grep proxy-tunnel

# 2. Verify ports are listening
netstat -tlnp | grep proxy-tunnel

# 3. Test with verbose output
curl -v -x http://127.0.0.1:1081 https://ifconfig.me

# 4. Try SOCKS5 instead
curl -v -x socks5://127.0.0.1:1080 https://ifconfig.me

# 5. Check proxy logs
./proxy-tunnel -link 'your-link' 2>&1 | tee proxy.log
```

---

### Proxy starts but can't connect to upstream

**Problem**: Upstream proxy not reachable or incorrect credentials.

**Solutions**:

1. **Verify upstream proxy is accessible**:
```bash
# For HTTP proxy
curl -v -x http://upstream-proxy:port https://google.com

# For SOCKS5
curl -v -x socks5://upstream-proxy:port https://google.com

# Test connectivity
ping upstream-proxy-server
telnet upstream-proxy-server port
```

2. **Check credentials**:
```bash
# Ensure username/password are correct
./proxy-tunnel -link 'socks5://user:pass@server:1080'

# URL-encode special characters in passwords
# @ becomes %40, : becomes %3A, etc.
```

3. **Check network/firewall**:
```bash
# Test if upstream is blocked
traceroute upstream-server
```

---

## Share Link Format Issues

### VLESS Link Problems

**Common issues**:
- Missing UUID
- Invalid query parameters
- Wrong security type

**Valid format**:
```bash
# Basic VLESS with TLS
vless://uuid@server:443?security=tls&sni=server.com

# VLESS with WebSocket + TLS
vless://uuid@server:443?type=ws&security=tls&path=/ws&host=server.com

# VLESS with Reality
vless://uuid@server:443?security=reality&pbk=publickey&sid=shortid&sni=example.com
```

---

### VMess Link Problems

**Common issues**:
- Invalid base64 encoding
- Corrupted JSON structure

**Debug steps**:
```bash
# Decode vmess link to check JSON
echo "base64string" | base64 -d | jq .

# Valid VMess JSON structure:
{
  "v": "2",
  "ps": "remark",
  "add": "server.com",
  "port": "443",
  "id": "uuid",
  "aid": "0",
  "net": "ws",
  "type": "none",
  "host": "server.com",
  "path": "/path",
  "tls": "tls"
}

# Re-encode if needed
echo '{"v":"2",...}' | base64 -w 0
```

---

### Shadowsocks Link Problems

**Common issues**:
- Invalid method
- Wrong password encoding

**Valid formats**:
```bash
# Method:password@server:port
ss://aes-256-gcm:password@server:8388

# Base64 encoded
ss://base64(method:password)@server:8388

# Supported methods:
# - aes-256-gcm (recommended)
# - aes-128-gcm
# - chacha20-poly1305
# - chacha20-ietf-poly1305
```

---

## Performance Issues

### Slow connection speed

**Possible causes**:
1. Upstream proxy is slow
2. Network congestion
3. CPU limitations
4. Too many connections

**Solutions**:

```bash
# 1. Test upstream speed directly
curl -x socks5://upstream:port -w "@curl-format.txt" -o /dev/null https://speed.cloudflare.com/__down?bytes=100000000

# 2. Monitor CPU usage
top -p $(pgrep proxy-tunnel)

# 3. Check network stats
iftop
nethogs

# 4. Limit to single proxy instance
pkill -f proxy-tunnel
./proxy-tunnel -link 'your-link'
```

---

### High memory usage

**Normal behavior**: sing-box uses memory for connection pooling and caching.

**If excessive**:
```bash
# Monitor memory
watch -n 1 'ps aux | grep proxy-tunnel'

# Restart periodically if needed (in production)
# Add to cron or systemd
```

---

## System-Specific Issues

### Linux: Permission denied on low ports

**Problem**: Ports < 1024 require root privileges.

**Solutions**:

```bash
# Option 1: Use high ports (recommended)
./proxy-tunnel -link 'link' -listen '127.0.0.1:8080' -http-port 8081

# Option 2: Use setcap (allows non-root binding)
sudo setcap 'cap_net_bind_service=+ep' ./proxy-tunnel

# Option 3: Run as root (not recommended)
sudo ./proxy-tunnel -link 'link' -listen '0.0.0.0:80'
```

---

### macOS: Binary won't run ("unidentified developer")

**Problem**: macOS Gatekeeper blocking unsigned binary.

**Solution**:
```bash
# Remove quarantine attribute
xattr -d com.apple.quarantine ./proxy-tunnel

# Or allow in System Preferences:
# System Preferences → Security & Privacy → Allow
```

---

### Windows: Firewall blocking connections

**Problem**: Windows Firewall blocking proxy.

**Solution**:
1. Add firewall exception for proxy-tunnel.exe
2. Or run with administrator privileges
3. Or temporarily disable firewall for testing

---

## Logging and Debugging

### Enable verbose logging

The proxy outputs logs to stderr by default. Capture them:

```bash
# Save to file
./proxy-tunnel -link 'your-link' 2>&1 | tee proxy.log

# View in real-time
./proxy-tunnel -link 'your-link' 2>&1 | tee proxy.log | grep -i error
```

---

### Check sing-box logs

Look for these indicators:

```bash
# ✓ Good - proxy started
"sing-box started"
"tcp server started at 127.0.0.1:1080"
"tcp server started at 127.0.0.1:1081"

# ✗ Bad - errors
"error"
"failed"
"connection refused"
```

---

### Test with curl verbose mode

```bash
# HTTP proxy - very verbose
curl -v -x http://127.0.0.1:1081 https://ifconfig.me

# SOCKS5 proxy - very verbose
curl -v -x socks5://127.0.0.1:1080 https://ifconfig.me

# Show timing
curl -w "@curl-format.txt" -x http://127.0.0.1:1081 -o /dev/null https://google.com
```

Create `curl-format.txt`:
```
time_namelookup:  %{time_namelookup}s\n
time_connect:     %{time_connect}s\n
time_appconnect:  %{time_appconnect}s\n
time_pretransfer: %{time_pretransfer}s\n
time_redirect:    %{time_redirect}s\n
time_starttransfer: %{time_starttransfer}s\n
time_total:       %{time_total}s\n
```

---

## Getting Help

If you're still stuck:

1. **Check logs carefully**: Look for specific error messages
2. **Test upstream proxy directly**: Ensure it works without proxy-tunnel
3. **Try minimal configuration**: Use basic HTTP/SOCKS5 first
4. **Update sing-box**: `go get -u github.com/sagernet/sing-box@dev-next`
5. **File an issue**: Include:
   - Error message
   - Share link format (remove sensitive info)
   - Build command used
   - Operating system
   - sing-box version

## Common Misunderstandings

### "Why do I need to specify a share link?"

proxy-tunnel is not a VPN - it's a proxy tunnel tool. You need an existing upstream proxy (VLESS, VMess, SOCKS5, etc.) to tunnel through.

### "Can I use this without an upstream proxy?"

Not really. The purpose is to create a local proxy that tunnels through an upstream proxy. For direct connections, just don't use a proxy at all.

### "Does this hide my traffic?"

Only if your upstream proxy does. proxy-tunnel just creates a convenient local endpoint.

---

## Quick Diagnostic Commands

```bash
# Check if binary has necessary features
./proxy-tunnel -h

# Test minimal config
./proxy-tunnel -link 'http://httpbin.org:80'

# Check port availability
netstat -tlnp | grep -E '1080|1081'

# Test upstream connectivity
curl -v -x socks5://your-upstream:port https://google.com

# Monitor connections
watch -n 1 'netstat -tn | grep -E "1080|1081"'

# Check DNS resolution
dig server.example.com
nslookup server.example.com
```
