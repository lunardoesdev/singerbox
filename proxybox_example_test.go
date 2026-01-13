package singerbox_test

import (
	"fmt"
	"time"

	"github.com/lunardoesdev/singerbox"

	"github.com/sagernet/sing-box/option"
)

// Example demonstrates basic usage of ProxyBox
func Example() {
	// Create a direct outbound for testing
	outbound := option.Outbound{
		Type:    "direct",
		Tag:     "direct",
		Options: &option.DirectOutboundOptions{},
	}

	// Create proxy box with default settings
	pb, err := singerbox.NewProxyBox(singerbox.ProxyBoxConfig{
		Outbound: outbound,
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Start the proxy
	if err := pb.Start(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Get proxy address (supports both SOCKS5 and HTTP)
	fmt.Printf("Mixed: %s\n", pb.ListenAddr())

	// Stop the proxy
	pb.Stop()

	// Output:
	// Mixed: 127.0.0.1:1080
}

// ExampleNewProxyBox demonstrates creating a new ProxyBox with custom configuration
func ExampleNewProxyBox() {
	outbound := option.Outbound{
		Type:    "direct",
		Tag:     "my-proxy",
		Options: &option.DirectOutboundOptions{},
	}

	pb, err := singerbox.NewProxyBox(singerbox.ProxyBoxConfig{
		Outbound:   outbound,
		ListenAddr: "127.0.0.1:9050",
		LogLevel:   "error",
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Created proxy box with custom port\n")
	fmt.Printf("Mixed proxy will listen on: %s\n", pb.ListenAddr())

	// Output:
	// Created proxy box with custom port
	// Mixed proxy will listen on: 127.0.0.1:9050
}

// ExampleProxyBox_Start demonstrates starting a proxy
func ExampleProxyBox_Start() {
	outbound := option.Outbound{
		Type:    "direct",
		Tag:     "test",
		Options: &option.DirectOutboundOptions{},
	}

	pb, _ := singerbox.NewProxyBox(singerbox.ProxyBoxConfig{
		Outbound:   outbound,
		ListenAddr: "127.0.0.1:19090",
		LogLevel:   "error",
	})

	// Start the proxy
	err := pb.Start()
	if err != nil {
		fmt.Printf("Start failed: %v\n", err)
		return
	}

	fmt.Println("Proxy started successfully")

	// Clean up
	pb.Stop()

	// Output:
	// Proxy started successfully
}

// ExampleProxyBox_Stop demonstrates stopping a proxy
func ExampleProxyBox_Stop() {
	outbound := option.Outbound{
		Type:    "direct",
		Tag:     "test",
		Options: &option.DirectOutboundOptions{},
	}

	pb, _ := singerbox.NewProxyBox(singerbox.ProxyBoxConfig{
		Outbound:   outbound,
		ListenAddr: "127.0.0.1:19092",
		LogLevel:   "error",
	})

	pb.Start()
	time.Sleep(50 * time.Millisecond)

	// Stop the proxy
	err := pb.Stop()
	if err != nil {
		fmt.Printf("Stop failed: %v\n", err)
		return
	}

	fmt.Println("Proxy stopped successfully")

	// Output:
	// Proxy stopped successfully
}

// ExampleProxyBox_IsRunning demonstrates checking proxy status
func ExampleProxyBox_IsRunning() {
	outbound := option.Outbound{
		Type:    "direct",
		Tag:     "test",
		Options: &option.DirectOutboundOptions{},
	}

	pb, _ := singerbox.NewProxyBox(singerbox.ProxyBoxConfig{
		Outbound:   outbound,
		ListenAddr: "127.0.0.1:19094",
		LogLevel:   "error",
	})

	fmt.Printf("Running before Start: %v\n", pb.IsRunning())

	pb.Start()
	fmt.Printf("Running after Start: %v\n", pb.IsRunning())

	pb.Stop()
	fmt.Printf("Running after Stop: %v\n", pb.IsRunning())

	// Output:
	// Running before Start: false
	// Running after Start: true
	// Running after Stop: false
}

// ExampleProxyBox_withSharelink demonstrates integration with sharelink parser
func ExampleProxyBox_withSharelink() {
	// Parse a share link

	outbound, err := singerbox.Parse("socks5://127.0.0.1:1080")
	if err != nil {
		fmt.Printf("Parse error: %v\n", err)
		return
	}

	// Create proxy box with the parsed outbound
	pb, err := singerbox.NewProxyBox(singerbox.ProxyBoxConfig{
		Outbound:   outbound,
		ListenAddr: "127.0.0.1:19096",
		LogLevel:   "error",
	})
	if err != nil {
		fmt.Printf("Create error: %v\n", err)
		return
	}

	fmt.Printf("Created proxy for %s outbound\n", outbound.Type)
	fmt.Printf("Tag: %s\n", pb.Outbound().Tag)

	// Output:
	// Created proxy for socks outbound
	// Tag: proxy
}

// ExampleNewProxyBox_multipleInstances demonstrates running multiple proxies
func ExampleNewProxyBox_multipleInstances() {
	// Create first proxy
	pb1, _ := singerbox.NewProxyBox(singerbox.ProxyBoxConfig{
		Outbound: option.Outbound{
			Type:    "direct",
			Tag:     "proxy1",
			Options: &option.DirectOutboundOptions{},
		},
		ListenAddr: "127.0.0.1:19100",
		LogLevel:   "error",
	})

	// Create second proxy
	pb2, _ := singerbox.NewProxyBox(singerbox.ProxyBoxConfig{
		Outbound: option.Outbound{
			Type:    "direct",
			Tag:     "proxy2",
			Options: &option.DirectOutboundOptions{},
		},
		ListenAddr: "127.0.0.1:19102",
		LogLevel:   "error",
	})

	// Start both
	pb1.Start()
	pb2.Start()

	fmt.Printf("Proxy 1: %s\n", pb1.ListenAddr())
	fmt.Printf("Proxy 2: %s\n", pb2.ListenAddr())

	// Clean up
	pb1.Stop()
	pb2.Stop()

	// Output:
	// Proxy 1: 127.0.0.1:19100
	// Proxy 2: 127.0.0.1:19102
}

// ExampleProxyBox_Config demonstrates accessing the configuration
func ExampleProxyBox_Config() {
	outbound := option.Outbound{
		Type:    "direct",
		Tag:     "test",
		Options: &option.DirectOutboundOptions{},
	}

	pb, _ := singerbox.NewProxyBox(singerbox.ProxyBoxConfig{
		Outbound: outbound,
		LogLevel: "debug",
	})

	cfg := pb.Config()
	fmt.Printf("Log level: %s\n", cfg.Log.Level)
	fmt.Printf("Number of inbounds: %d\n", len(cfg.Inbounds))
	fmt.Printf("Number of outbounds: %d\n", len(cfg.Outbounds))

	// Output:
	// Log level: debug
	// Number of inbounds: 1
	// Number of outbounds: 3
}

// ExampleFromSharedLink demonstrates the simplest way to create a proxy
func ExampleFromSharedLink() {
	// Create and start proxy from a share link in one call
	proxy, err := singerbox.FromSharedLink(
		"ss://aes-256-gcm:password@server.com:8388",
		singerbox.ProxyConfig{
			ListenAddr: "127.0.0.1:19200",
			LogLevel:   "error", // Silent for test
		},
	)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer proxy.Stop()

	fmt.Printf("Proxy running: %v\n", proxy.IsRunning())
	fmt.Printf("Listening on: %s\n", proxy.ListenAddr())

	// Output:
	// Proxy running: true
	// Listening on: 127.0.0.1:19200
}

// ExampleFromSharedLink_minimal demonstrates using default settings
func ExampleFromSharedLink_minimal() {
	// Minimal configuration - all defaults
	proxy, _ := singerbox.FromSharedLink(
		"ss://aes-256-gcm:password@server.com:8388",
		singerbox.ProxyConfig{},
	)
	defer proxy.Stop()

	fmt.Println("Proxy created with defaults")
	fmt.Printf("Default address: %s\n", proxy.ListenAddr())

	// Output:
	// Proxy created with defaults
	// Default address: 127.0.0.1:1080
}
