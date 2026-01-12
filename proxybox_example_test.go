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

	// Get proxy addresses
	fmt.Printf("SOCKS5: %s\n", pb.ListenAddr())
	fmt.Printf("HTTP:   %s\n", pb.HTTPAddr())

	// Stop the proxy
	pb.Stop()

	// Output:
	// SOCKS5: 127.0.0.1:1080
	// HTTP:   127.0.0.1:1081
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
		HTTPPort:   9051,
		LogLevel:   "error",
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Created proxy box with custom ports\n")
	fmt.Printf("SOCKS5 will listen on: %s\n", pb.ListenAddr())
	fmt.Printf("HTTP will listen on: %s\n", pb.HTTPAddr())

	// Output:
	// Created proxy box with custom ports
	// SOCKS5 will listen on: 127.0.0.1:9050
	// HTTP will listen on: 127.0.0.1:9051
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
		HTTPPort:   19091,
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
		HTTPPort:   19093,
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
		HTTPPort:   19095,
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
	parser := singerbox.NewParser()
	outbound, err := parser.Parse("socks5://127.0.0.1:1080")
	if err != nil {
		fmt.Printf("Parse error: %v\n", err)
		return
	}

	// Create proxy box with the parsed outbound
	pb, err := singerbox.NewProxyBox(singerbox.ProxyBoxConfig{
		Outbound:   outbound,
		ListenAddr: "127.0.0.1:19096",
		HTTPPort:   19097,
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
		HTTPPort:   19101,
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
		HTTPPort:   19103,
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
	// Number of inbounds: 2
	// Number of outbounds: 3
}
