#!/bin/bash

# Test script for proxy-tunnel
# This tests with a local SOCKS5 proxy (if you have one)

echo "Testing proxy-tunnel with example configurations..."
echo ""

# Test 1: Simple SOCKS5 proxy (replace with your actual proxy)
echo "Example 1: SOCKS5 proxy"
echo "./proxy-tunnel -link 'socks5://localhost:9050'"
echo ""

# Test 2: HTTP proxy
echo "Example 2: HTTP proxy"
echo "./proxy-tunnel -link 'http://proxy.example.com:8080'"
echo ""

# Test 3: Shadowsocks
echo "Example 3: Shadowsocks"
echo "./proxy-tunnel -link 'ss://aes-256-gcm:your-password@server.example.com:8388'"
echo ""

# Test 4: VLESS with WebSocket + TLS
echo "Example 4: VLESS with WebSocket + TLS"
echo "./proxy-tunnel -link 'vless://your-uuid@server.example.com:443?type=ws&security=tls&path=/ws'"
echo ""

echo "After starting, test with:"
echo "  curl -v -x http://127.0.0.1:1081 https://ifconfig.me"
echo "  curl -v -x socks5://127.0.0.1:1080 https://ifconfig.me"
echo ""

echo "To run a basic test (without actual proxy), run:"
echo "./proxy-tunnel -link 'direct://' 2>&1 | head -20"
