package main

import (
	"context"
	"flag"
	"fmt"
	"net/netip"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"proxy-tunnel/pkg/sharelink"

	"github.com/sagernet/sing-box"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/include"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common/json/badoption"
)

var (
	shareLink  = flag.String("link", "", "Share link (vless://, vmess://, ss://, etc.)")
	listenAddr = flag.String("listen", "127.0.0.1:1080", "Local proxy listen address")
	httpPort   = flag.Int("http-port", 1081, "HTTP proxy port")
)

func main() {
	flag.Parse()

	if *shareLink == "" {
		fmt.Println("Usage: proxy-tunnel -link <share-link> [-listen <addr:port>] [-http-port <port>]")
		fmt.Println("\nSupported protocols: vless, vmess, ss, trojan, http, socks")
		fmt.Println("\nExample:")
		fmt.Println("  proxy-tunnel -link 'vless://uuid@server:443?type=ws&security=tls'")
		os.Exit(1)
	}

	// Parse share link using the library
	parser := sharelink.New()
	outbound, err := parser.Parse(*shareLink)
	if err != nil {
		fmt.Printf("Error parsing share link: %v\n", err)
		os.Exit(1)
	}

	// Create sing-box configuration
	config := createConfig(outbound, *listenAddr, *httpPort)

	// Create and start sing-box instance with proper context
	ctx := context.Background()
	ctx = include.Context(ctx) // Register all protocol handlers
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	instance, err := box.New(box.Options{
		Context: ctx,
		Options: config,
	})
	if err != nil {
		fmt.Printf("Error creating sing-box instance: %v\n", err)
		os.Exit(1)
	}

	// Start the box
	err = instance.Start()
	if err != nil {
		fmt.Printf("Error starting sing-box: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Proxy started successfully!\n")
	fmt.Printf("  SOCKS5: %s\n", *listenAddr)
	fmt.Printf("  HTTP:   127.0.0.1:%d\n", *httpPort)
	fmt.Printf("  Routing through: %s\n", outbound.Tag)
	fmt.Println("\nPress Ctrl+C to stop...")

	// Wait for interrupt signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	fmt.Println("\nStopping proxy...")
	instance.Close()
}

func createConfig(outbound option.Outbound, socksAddr string, httpPort int) option.Options {
	// Split address for HTTP proxy
	host := strings.Split(socksAddr, ":")[0]

	// Parse listen address
	listenIP, err := netip.ParseAddr(host)
	if err != nil {
		listenIP = netip.MustParseAddr("127.0.0.1")
	}
	listenAddr := (*badoption.Addr)(&listenIP)

	return option.Options{
		Log: &option.LogOptions{
			Level:  "info",
			Output: "stderr",
		},
		Inbounds: []option.Inbound{
			{
				Type: "mixed",
				Tag:  "mixed-in",
				Options: &option.HTTPMixedInboundOptions{
					ListenOptions: option.ListenOptions{
						Listen:     listenAddr,
						ListenPort: uint16(getPort(socksAddr)),
					},
				},
			},
			{
				Type: "http",
				Tag:  "http-in",
				Options: &option.HTTPMixedInboundOptions{
					ListenOptions: option.ListenOptions{
						Listen:     listenAddr,
						ListenPort: uint16(httpPort),
					},
				},
			},
		},
		Outbounds: []option.Outbound{
			outbound,
			{
				Type:    "direct",
				Tag:     "direct",
				Options: &option.DirectOutboundOptions{},
			},
			{
				Type:    "block",
				Tag:     "block",
				Options: &option.StubOptions{},
			},
		},
		Route: &option.RouteOptions{
			Rules: []option.Rule{
				{
					Type: C.RuleTypeDefault,
					DefaultOptions: option.DefaultRule{
						RuleAction: option.RuleAction{
							Action: C.RuleActionTypeRoute,
							RouteOptions: option.RouteActionOptions{
								Outbound: outbound.Tag,
							},
						},
					},
				},
			},
			AutoDetectInterface: true,
		},
	}
}

func getPort(hostPort string) int {
	parts := strings.Split(hostPort, ":")
	if len(parts) < 2 {
		return 443
	}
	port := 443
	fmt.Sscanf(parts[len(parts)-1], "%d", &port)
	return port
}
