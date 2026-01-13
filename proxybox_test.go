package singerbox_test

import (
	"testing"
	"time"

	"github.com/lunardoesdev/singerbox"

	"github.com/sagernet/sing-box/option"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		config  singerbox.ProxyBoxConfig
		wantErr bool
		check   func(*testing.T, *singerbox.ProxyBox)
	}{
		{
			name: "Valid SOCKS outbound with defaults",
			config: singerbox.ProxyBoxConfig{
				Outbound: option.Outbound{
					Type: "socks",
					Tag:  "test-socks",
					Options: &option.SOCKSOutboundOptions{
						ServerOptions: option.ServerOptions{
							Server:     "127.0.0.1",
							ServerPort: 1080,
						},
					},
				},
			},
			wantErr: false,
			check: func(t *testing.T, pb *singerbox.ProxyBox) {
				if pb == nil {
					t.Fatal("ProxyBox is nil")
				}
				if pb.ListenAddr() != "127.0.0.1:1080" {
					t.Errorf("ListenAddr = %s, want 127.0.0.1:1080", pb.ListenAddr())
				}
			},
		},
		{
			name: "Custom listen address and port",
			config: singerbox.ProxyBoxConfig{
				Outbound: option.Outbound{
					Type: "socks",
					Tag:  "test-socks",
					Options: &option.SOCKSOutboundOptions{
						ServerOptions: option.ServerOptions{
							Server:     "127.0.0.1",
							ServerPort: 1080,
						},
					},
				},
				ListenAddr: "127.0.0.1:9050",
			},
			wantErr: false,
			check: func(t *testing.T, pb *singerbox.ProxyBox) {
				if pb.ListenAddr() != "127.0.0.1:9050" {
					t.Errorf("ListenAddr = %s, want 127.0.0.1:9050", pb.ListenAddr())
				}
			},
		},
		{
			name: "Custom log level",
			config: singerbox.ProxyBoxConfig{
				Outbound: option.Outbound{
					Type: "socks",
					Tag:  "test-socks",
					Options: &option.SOCKSOutboundOptions{
						ServerOptions: option.ServerOptions{
							Server:     "127.0.0.1",
							ServerPort: 1080,
						},
					},
				},
				LogLevel: "debug",
			},
			wantErr: false,
			check: func(t *testing.T, pb *singerbox.ProxyBox) {
				cfg := pb.Config()
				if cfg.Log.Level != "debug" {
					t.Errorf("LogLevel = %s, want debug", cfg.Log.Level)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pb, err := singerbox.NewProxyBox(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, pb)
			}
		})
	}
}

func TestProxyBox_StartStop(t *testing.T) {
	// Create a simple SOCKS outbound for testing
	pb, err := singerbox.NewProxyBox(singerbox.ProxyBoxConfig{
		Outbound: option.Outbound{
			Type: "direct",
			Tag:  "direct",
			Options: &option.DirectOutboundOptions{},
		},
		ListenAddr: "127.0.0.1:19080",
		LogLevel:   "error", // Reduce log noise in tests
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Test initial state
	if pb.IsRunning() {
		t.Error("ProxyBox should not be running initially")
	}

	// Test Start
	err = pb.Start()
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	if !pb.IsRunning() {
		t.Error("ProxyBox should be running after Start()")
	}

	// Give it a moment to actually start
	time.Sleep(100 * time.Millisecond)

	// Test double Start (should fail)
	err = pb.Start()
	if err == nil {
		t.Error("Start() should fail when already running")
	}

	// Test Stop
	err = pb.Stop()
	if err != nil {
		t.Errorf("Stop() error = %v", err)
	}

	if pb.IsRunning() {
		t.Error("ProxyBox should not be running after Stop()")
	}

	// Test double Stop (should fail)
	err = pb.Stop()
	if err == nil {
		t.Error("Stop() should fail when not running")
	}
}

func TestProxyBox_WithSharelink(t *testing.T) {
	// Test integration with sharelink parser
	

	tests := []struct {
		name       string
		link       string
		wantErr    bool
		shouldSkip bool // Some protocols may not work in test environment
		skipReason string
	}{
		{
			name:    "SOCKS5 proxy",
			link:    "socks5://127.0.0.1:1080",
			wantErr: false,
		},
		{
			name:    "HTTP proxy",
			link:    "http://127.0.0.1:8080",
			wantErr: false,
		},
		{
			name:       "VLESS with TLS",
			link:       "vless://550e8400-e29b-41d4-a716-446655440000@example.com:443?security=tls&type=ws",
			wantErr:    false,
			shouldSkip: true,
			skipReason: "Requires network connection",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldSkip {
				t.Skip(tt.skipReason)
			}

			outbound, err := singerbox.Parse(tt.link)
			if err != nil {
				t.Fatalf("Failed to parse link: %v", err)
			}

			pb, err := singerbox.NewProxyBox(singerbox.ProxyBoxConfig{
				Outbound:   outbound,
				ListenAddr: "127.0.0.1:19082",
				LogLevel:   "error",
			})

			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// Try to start (may fail due to network, but shouldn't panic)
			err = pb.Start()
			if err == nil {
				// If it started, stop it
				defer pb.Stop()
				time.Sleep(100 * time.Millisecond)
			}
		})
	}
}

func TestProxyBox_Config(t *testing.T) {
	outbound := option.Outbound{
		Type: "direct",
		Tag:  "test-direct",
		Options: &option.DirectOutboundOptions{},
	}

	pb, err := singerbox.NewProxyBox(singerbox.ProxyBoxConfig{
		Outbound: outbound,
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Test Config() returns valid configuration
	cfg := pb.Config()
	if cfg.Log == nil {
		t.Error("Config.Log is nil")
	}
	if len(cfg.Inbounds) != 1 {
		t.Errorf("Expected 1 inbound (mixed), got %d", len(cfg.Inbounds))
	}
	if len(cfg.Outbounds) != 3 {
		t.Errorf("Expected 3 outbounds (proxy, direct, block), got %d", len(cfg.Outbounds))
	}

	// Test Outbound() returns correct outbound
	if pb.Outbound().Tag != "test-direct" {
		t.Errorf("Outbound().Tag = %s, want test-direct", pb.Outbound().Tag)
	}
}

func TestProxyBox_Addresses(t *testing.T) {
	tests := []struct {
		name           string
		listenAddr     string
		wantListenAddr string
	}{
		{
			name:           "Default addresses",
			listenAddr:     "",
			wantListenAddr: "127.0.0.1:1080",
		},
		{
			name:           "Custom addresses",
			listenAddr:     "127.0.0.1:9050",
			wantListenAddr: "127.0.0.1:9050",
		},
		{
			name:           "High port numbers",
			listenAddr:     "127.0.0.1:19999",
			wantListenAddr: "127.0.0.1:19999",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pb, err := singerbox.NewProxyBox(singerbox.ProxyBoxConfig{
				Outbound: option.Outbound{
					Type:    "direct",
					Tag:     "direct",
					Options: &option.DirectOutboundOptions{},
				},
				ListenAddr: tt.listenAddr,
			})
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}

			if pb.ListenAddr() != tt.wantListenAddr {
				t.Errorf("ListenAddr() = %s, want %s", pb.ListenAddr(), tt.wantListenAddr)
			}
		})
	}
}

func TestProxyBox_MultipleInstances(t *testing.T) {
	// Test that multiple instances can coexist with different ports
	pb1, err := singerbox.NewProxyBox(singerbox.ProxyBoxConfig{
		Outbound: option.Outbound{
			Type:    "direct",
			Tag:     "direct1",
			Options: &option.DirectOutboundOptions{},
		},
		ListenAddr: "127.0.0.1:19100",
		LogLevel:   "error",
	})
	if err != nil {
		t.Fatalf("New(pb1) error = %v", err)
	}

	pb2, err := singerbox.NewProxyBox(singerbox.ProxyBoxConfig{
		Outbound: option.Outbound{
			Type:    "direct",
			Tag:     "direct2",
			Options: &option.DirectOutboundOptions{},
		},
		ListenAddr: "127.0.0.1:19102",
		LogLevel:   "error",
	})
	if err != nil {
		t.Fatalf("New(pb2) error = %v", err)
	}

	// Start both
	if err := pb1.Start(); err != nil {
		t.Fatalf("pb1.Start() error = %v", err)
	}
	defer pb1.Stop()

	time.Sleep(100 * time.Millisecond)

	if err := pb2.Start(); err != nil {
		t.Fatalf("pb2.Start() error = %v", err)
	}
	defer pb2.Stop()

	time.Sleep(100 * time.Millisecond)

	// Both should be running
	if !pb1.IsRunning() {
		t.Error("pb1 should be running")
	}
	if !pb2.IsRunning() {
		t.Error("pb2 should be running")
	}
}
