package singerbox_test

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/lunardoesdev/singerbox"
	"github.com/sagernet/sing-box/option"
)

func TestParseVLESS(t *testing.T) {
	tests := []struct {
		name    string
		link    string
		wantErr bool
		check   func(*testing.T, option.Outbound)
	}{
		{
			name: "VLESS with TLS and WebSocket",
			link: "vless://a1b2c3d4-e5f6-7890-abcd-ef1234567890@example.com:443?type=ws&security=tls&path=/ws&host=example.com&sni=example.com#TestServer",
			wantErr: false,
			check: func(t *testing.T, out option.Outbound) {
				if out.Type != "vless" {
					t.Errorf("Type = %v, want vless", out.Type)
				}
				if out.Tag != "TestServer" {
					t.Errorf("Tag = %v, want TestServer", out.Tag)
				}
				opts, ok := out.Options.(*option.VLESSOutboundOptions)
				if !ok {
					t.Fatal("Options is not *VLESSOutboundOptions")
				}
				if opts.UUID != "a1b2c3d4-e5f6-7890-abcd-ef1234567890" {
					t.Errorf("UUID = %v", opts.UUID)
				}
				if opts.Server != "example.com" {
					t.Errorf("Server = %v, want example.com", opts.Server)
				}
				if opts.ServerPort != 443 {
					t.Errorf("ServerPort = %v, want 443", opts.ServerPort)
				}
				if opts.TLS == nil || !opts.TLS.Enabled {
					t.Error("TLS should be enabled")
				}
				if opts.Transport == nil || opts.Transport.Type != "ws" {
					t.Error("Transport should be WebSocket")
				}
			},
		},
		{
			name: "VLESS with Reality",
			link: "vless://550e8400-e29b-41d4-a716-446655440000@reality.example.com:443?security=reality&pbk=testPublicKey123&sid=abcd1234&sni=www.microsoft.com&fp=chrome",
			wantErr: false,
			check: func(t *testing.T, out option.Outbound) {
				opts, ok := out.Options.(*option.VLESSOutboundOptions)
				if !ok {
					t.Fatal("Options is not *VLESSOutboundOptions")
				}
				if opts.TLS == nil {
					t.Fatal("TLS should not be nil")
				}
				if opts.TLS.Reality == nil || !opts.TLS.Reality.Enabled {
					t.Error("Reality should be enabled")
				}
				if opts.TLS.UTLS == nil || !opts.TLS.UTLS.Enabled {
					t.Error("uTLS should be enabled for Reality")
				}
				if opts.TLS.UTLS.Fingerprint != "chrome" {
					t.Errorf("Fingerprint = %v, want chrome", opts.TLS.UTLS.Fingerprint)
				}
				if opts.TLS.Reality.PublicKey != "testPublicKey123" {
					t.Errorf("PublicKey = %v", opts.TLS.Reality.PublicKey)
				}
				if opts.TLS.Reality.ShortID != "abcd1234" {
					t.Errorf("ShortID = %v", opts.TLS.Reality.ShortID)
				}
			},
		},
		{
			name: "VLESS with gRPC",
			link: "vless://test-uuid@grpc.example.com:443?type=grpc&serviceName=TestService&security=tls&sni=grpc.example.com",
			wantErr: false,
			check: func(t *testing.T, out option.Outbound) {
				opts, ok := out.Options.(*option.VLESSOutboundOptions)
				if !ok {
					t.Fatal("Options is not *VLESSOutboundOptions")
				}
				if opts.Transport == nil || opts.Transport.Type != "grpc" {
					t.Error("Transport should be gRPC")
				}
				if opts.Transport.GRPCOptions.ServiceName != "TestService" {
					t.Errorf("ServiceName = %v, want TestService", opts.Transport.GRPCOptions.ServiceName)
				}
			},
		},
		{
			name:    "VLESS missing UUID",
			link:    "vless://@example.com:443",
			wantErr: true,
		},
		{
			name:    "VLESS missing server",
			link:    "vless://uuid@:443",
			wantErr: true,
		},
	}

	parser := singerbox.NewParser()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := parser.ParseVLESS(tt.link)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseVLESS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, out)
			}
		})
	}
}

func TestParseVMess(t *testing.T) {
	// Create a VMess configuration
	vmessConfig := singerbox.VMessConfig{
		V:    "2",
		Ps:   "TestVMess",
		Add:  "vmess.example.com",
		Port: "443",
		ID:   "b2c3d4e5-f678-90ab-cdef-123456789abc",
		Aid:  "0",
		Net:  "ws",
		Type: "none",
		Host: "vmess.example.com",
		Path: "/vmess",
		TLS:  "tls",
		SNI:  "vmess.example.com",
	}

	configJSON, _ := json.Marshal(vmessConfig)
	encoded := base64.StdEncoding.EncodeToString(configJSON)

	tests := []struct {
		name    string
		link    string
		wantErr bool
		check   func(*testing.T, option.Outbound)
	}{
		{
			name: "VMess with WebSocket and TLS",
			link: "vmess://" + encoded,
			wantErr: false,
			check: func(t *testing.T, out option.Outbound) {
				if out.Type != "vmess" {
					t.Errorf("Type = %v, want vmess", out.Type)
				}
				if out.Tag != "TestVMess" {
					t.Errorf("Tag = %v, want TestVMess", out.Tag)
				}
				opts, ok := out.Options.(*option.VMessOutboundOptions)
				if !ok {
					t.Fatal("Options is not *VMessOutboundOptions")
				}
				if opts.UUID != "b2c3d4e5-f678-90ab-cdef-123456789abc" {
					t.Errorf("UUID = %v", opts.UUID)
				}
				if opts.Server != "vmess.example.com" {
					t.Errorf("Server = %v", opts.Server)
				}
				if opts.TLS == nil || !opts.TLS.Enabled {
					t.Error("TLS should be enabled")
				}
				if opts.Transport == nil || opts.Transport.Type != "ws" {
					t.Error("Transport should be WebSocket")
				}
			},
		},
		{
			name:    "VMess invalid base64",
			link:    "vmess://invalid!!!base64",
			wantErr: true,
		},
		{
			name:    "VMess invalid JSON",
			link:    "vmess://" + base64.StdEncoding.EncodeToString([]byte("{invalid json")),
			wantErr: true,
		},
	}

	parser := singerbox.NewParser()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := parser.ParseVMess(tt.link)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseVMess() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, out)
			}
		})
	}
}

func TestParseShadowsocks(t *testing.T) {
	tests := []struct {
		name    string
		link    string
		wantErr bool
		check   func(*testing.T, option.Outbound)
	}{
		{
			name: "Shadowsocks method:password@server:port",
			link: "ss://aes-256-gcm:mypassword123@ss.example.com:8388#TestSS",
			wantErr: false,
			check: func(t *testing.T, out option.Outbound) {
				if out.Type != "shadowsocks" {
					t.Errorf("Type = %v, want shadowsocks", out.Type)
				}
				if out.Tag != "TestSS" {
					t.Errorf("Tag = %v, want TestSS", out.Tag)
				}
				opts, ok := out.Options.(*option.ShadowsocksOutboundOptions)
				if !ok {
					t.Fatal("Options is not *ShadowsocksOutboundOptions")
				}
				if opts.Method != "aes-256-gcm" {
					t.Errorf("Method = %v, want aes-256-gcm", opts.Method)
				}
				if opts.Password != "mypassword123" {
					t.Errorf("Password = %v", opts.Password)
				}
				if opts.Server != "ss.example.com" {
					t.Errorf("Server = %v", opts.Server)
				}
				if opts.ServerPort != 8388 {
					t.Errorf("ServerPort = %v, want 8388", opts.ServerPort)
				}
			},
		},
		{
			name: "Shadowsocks base64 encoded",
			link: "ss://" + base64.StdEncoding.EncodeToString([]byte("chacha20-poly1305:testpass@192.168.1.1:8388")),
			wantErr: false,
			check: func(t *testing.T, out option.Outbound) {
				opts, ok := out.Options.(*option.ShadowsocksOutboundOptions)
				if !ok {
					t.Fatal("Options is not *ShadowsocksOutboundOptions")
				}
				if opts.Method != "chacha20-poly1305" {
					t.Errorf("Method = %v", opts.Method)
				}
				if opts.Password != "testpass" {
					t.Errorf("Password = %v", opts.Password)
				}
			},
		},
		{
			name:    "Shadowsocks missing method",
			link:    "ss://:password@server:8388",
			wantErr: true,
		},
	}

	parser := singerbox.NewParser()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := parser.ParseShadowsocks(tt.link)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseShadowsocks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, out)
			}
		})
	}
}

func TestParseTrojan(t *testing.T) {
	tests := []struct {
		name    string
		link    string
		wantErr bool
		check   func(*testing.T, option.Outbound)
	}{
		{
			name: "Trojan basic",
			link: "trojan://mypassword@trojan.example.com:443?sni=trojan.example.com#TestTrojan",
			wantErr: false,
			check: func(t *testing.T, out option.Outbound) {
				if out.Type != "trojan" {
					t.Errorf("Type = %v, want trojan", out.Type)
				}
				if out.Tag != "TestTrojan" {
					t.Errorf("Tag = %v, want TestTrojan", out.Tag)
				}
				opts, ok := out.Options.(*option.TrojanOutboundOptions)
				if !ok {
					t.Fatal("Options is not *TrojanOutboundOptions")
				}
				if opts.Password != "mypassword" {
					t.Errorf("Password = %v", opts.Password)
				}
				if opts.Server != "trojan.example.com" {
					t.Errorf("Server = %v", opts.Server)
				}
				if opts.TLS == nil || !opts.TLS.Enabled {
					t.Error("TLS should be enabled")
				}
			},
		},
		{
			name: "Trojan with WebSocket",
			link: "trojan://pass123@ws.trojan.com:443?type=ws&path=/trojan&host=ws.trojan.com",
			wantErr: false,
			check: func(t *testing.T, out option.Outbound) {
				opts, ok := out.Options.(*option.TrojanOutboundOptions)
				if !ok {
					t.Fatal("Options is not *TrojanOutboundOptions")
				}
				if opts.Transport == nil || opts.Transport.Type != "ws" {
					t.Error("Transport should be WebSocket")
				}
				if opts.Transport.WebsocketOptions.Path != "/trojan" {
					t.Errorf("Path = %v", opts.Transport.WebsocketOptions.Path)
				}
			},
		},
		{
			name:    "Trojan missing password",
			link:    "trojan://@server.com:443",
			wantErr: true,
		},
	}

	parser := singerbox.NewParser()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := parser.ParseTrojan(tt.link)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTrojan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, out)
			}
		})
	}
}

func TestParseSOCKS(t *testing.T) {
	tests := []struct {
		name    string
		link    string
		wantErr bool
		check   func(*testing.T, option.Outbound)
	}{
		{
			name: "SOCKS5 with auth",
			link: "socks5://user:pass@socks.example.com:1080",
			wantErr: false,
			check: func(t *testing.T, out option.Outbound) {
				if out.Type != "socks" {
					t.Errorf("Type = %v, want socks", out.Type)
				}
				opts, ok := out.Options.(*option.SOCKSOutboundOptions)
				if !ok {
					t.Fatal("Options is not *SOCKSOutboundOptions")
				}
				if opts.Username != "user" {
					t.Errorf("Username = %v, want user", opts.Username)
				}
				if opts.Password != "pass" {
					t.Errorf("Password = %v, want pass", opts.Password)
				}
				if opts.Server != "socks.example.com" {
					t.Errorf("Server = %v", opts.Server)
				}
				if opts.ServerPort != 1080 {
					t.Errorf("ServerPort = %v, want 1080", opts.ServerPort)
				}
			},
		},
		{
			name: "SOCKS5 without auth",
			link: "socks5://proxy.example.com:1080",
			wantErr: false,
			check: func(t *testing.T, out option.Outbound) {
				opts, ok := out.Options.(*option.SOCKSOutboundOptions)
				if !ok {
					t.Fatal("Options is not *SOCKSOutboundOptions")
				}
				if opts.Username != "" {
					t.Error("Username should be empty")
				}
				if opts.Password != "" {
					t.Error("Password should be empty")
				}
			},
		},
		{
			name:    "SOCKS5 missing server",
			link:    "socks5://:1080",
			wantErr: true,
		},
	}

	parser := singerbox.NewParser()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := parser.ParseSOCKS(tt.link)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSOCKS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, out)
			}
		})
	}
}

func TestParseHTTP(t *testing.T) {
	tests := []struct {
		name    string
		link    string
		wantErr bool
		check   func(*testing.T, option.Outbound)
	}{
		{
			name: "HTTP with auth",
			link: "http://user:pass@proxy.example.com:8080",
			wantErr: false,
			check: func(t *testing.T, out option.Outbound) {
				if out.Type != "http" {
					t.Errorf("Type = %v, want http", out.Type)
				}
				opts, ok := out.Options.(*option.HTTPOutboundOptions)
				if !ok {
					t.Fatal("Options is not *HTTPOutboundOptions")
				}
				if opts.Username != "user" {
					t.Errorf("Username = %v", opts.Username)
				}
				if opts.Password != "pass" {
					t.Errorf("Password = %v", opts.Password)
				}
				if opts.Server != "proxy.example.com" {
					t.Errorf("Server = %v", opts.Server)
				}
				if opts.TLS != nil && opts.TLS.Enabled {
					t.Error("TLS should not be enabled for http://")
				}
			},
		},
		{
			name: "HTTPS with auth",
			link: "https://user:pass@secure.example.com:443",
			wantErr: false,
			check: func(t *testing.T, out option.Outbound) {
				opts, ok := out.Options.(*option.HTTPOutboundOptions)
				if !ok {
					t.Fatal("Options is not *HTTPOutboundOptions")
				}
				if opts.TLS == nil || !opts.TLS.Enabled {
					t.Error("TLS should be enabled for https://")
				}
			},
		},
		{
			name: "HTTP without auth",
			link: "http://proxy.example.com:3128",
			wantErr: false,
			check: func(t *testing.T, out option.Outbound) {
				opts, ok := out.Options.(*option.HTTPOutboundOptions)
				if !ok {
					t.Fatal("Options is not *HTTPOutboundOptions")
				}
				if opts.Username != "" || opts.Password != "" {
					t.Error("Credentials should be empty")
				}
			},
		},
	}

	parser := singerbox.NewParser()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := parser.ParseHTTP(tt.link)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseHTTP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, out)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		link     string
		wantType string
		wantErr  bool
	}{
		{
			name:     "Auto-detect VLESS",
			link:     "vless://uuid@server:443",
			wantType: "vless",
			wantErr:  false,
		},
		{
			name:     "Auto-detect VMess",
			link:     "vmess://" + base64.StdEncoding.EncodeToString([]byte(`{"add":"server","port":"443","id":"uuid","ps":"test"}`)),
			wantType: "vmess",
			wantErr:  false,
		},
		{
			name:     "Auto-detect Shadowsocks",
			link:     "ss://aes-256-gcm:pass@server:8388",
			wantType: "shadowsocks",
			wantErr:  false,
		},
		{
			name:     "Auto-detect Trojan",
			link:     "trojan://pass@server:443",
			wantType: "trojan",
			wantErr:  false,
		},
		{
			name:     "Auto-detect SOCKS5",
			link:     "socks5://server:1080",
			wantType: "socks",
			wantErr:  false,
		},
		{
			name:     "Auto-detect HTTP",
			link:     "http://server:8080",
			wantType: "http",
			wantErr:  false,
		},
		{
			name:     "Unsupported protocol",
			link:     "ftp://server:21",
			wantType: "",
			wantErr:  true,
		},
	}

	parser := singerbox.NewParser()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := parser.Parse(tt.link)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && out.Type != tt.wantType {
				t.Errorf("Parse() type = %v, want %v", out.Type, tt.wantType)
			}
		})
	}
}

// Benchmark tests
func BenchmarkParseVLESS(b *testing.B) {
	parser := singerbox.NewParser()
	link := "vless://uuid@server:443?type=ws&security=tls&path=/ws"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parser.ParseVLESS(link)
	}
}

func BenchmarkParseVMess(b *testing.B) {
	parser := singerbox.NewParser()
	config := `{"v":"2","ps":"test","add":"server","port":"443","id":"uuid","aid":"0","net":"ws","type":"none","host":"server","path":"/","tls":"tls"}`
	link := "vmess://" + base64.StdEncoding.EncodeToString([]byte(config))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parser.ParseVMess(link)
	}
}

func BenchmarkParseShadowsocks(b *testing.B) {
	parser := singerbox.NewParser()
	link := "ss://aes-256-gcm:password@server:8388"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parser.ParseShadowsocks(link)
	}
}
