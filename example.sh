#!/bin/bash

# Example usage of proxy-tunnel
# Replace the share links below with your actual proxy configurations

echo "Proxy Tunnel Examples"
echo "====================="
echo ""

echo "1. VLESS with TLS and WebSocket:"
echo "   ./proxy-tunnel -link 'vless://your-uuid@example.com:443?type=ws&security=tls&path=/ws&host=example.com'"
echo ""

echo "2. VLESS with Reality:"
echo "   ./proxy-tunnel -link 'vless://your-uuid@example.com:443?security=reality&pbk=your-public-key&sid=your-short-id&sni=www.example.com'"
echo ""

echo "3. VMess (base64 encoded JSON):"
echo "   ./proxy-tunnel -link 'vmess://base64encodedconfig'"
echo ""

echo "4. Shadowsocks:"
echo "   ./proxy-tunnel -link 'ss://aes-256-gcm:your-password@example.com:8388'"
echo ""

echo "5. Trojan with TLS:"
echo "   ./proxy-tunnel -link 'trojan://your-password@example.com:443?sni=example.com'"
echo ""

echo "6. SOCKS5 with authentication:"
echo "   ./proxy-tunnel -link 'socks5://user:pass@proxy.example.com:1080'"
echo ""

echo "7. HTTP/HTTPS proxy:"
echo "   ./proxy-tunnel -link 'http://proxy.example.com:8080'"
echo ""

echo "8. Custom listen address:"
echo "   ./proxy-tunnel -link 'vless://...' -listen '0.0.0.0:1080' -http-port 8080"
echo ""

echo "After starting, test with:"
echo "  curl -x http://127.0.0.1:1081 https://ifconfig.me"
echo "  curl -x socks5://127.0.0.1:1080 https://ifconfig.me"
