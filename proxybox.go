package singerbox

import (
	"context"
	"fmt"
	"net/netip"
	"strings"

	"github.com/sagernet/sing-box"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/include"
	"github.com/sagernet/sing-box/option"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/common/json/badoption"
)

// ProxyBox manages a sing-box instance for proxying traffic
type ProxyBox struct {
	instance *box.Box
	ctx      context.Context
	cancel   context.CancelFunc
	config   option.Options
	outbound option.Outbound
}

// Config holds configuration for creating a ProxyBox
type ProxyBoxConfig struct {
	// Outbound is the sing-box outbound configuration
	Outbound option.Outbound

	// ListenAddr is the address for SOCKS5/HTTP mixed proxy (default: "127.0.0.1:1080")
	ListenAddr string

	// LogLevel sets the logging level (default: "panic" for silent operation)
	// Available levels: "trace", "debug", "info", "warn", "error", "fatal", "panic"
	LogLevel string
}

// New creates a new ProxyBox with the given configuration
func NewProxyBox(cfg ProxyBoxConfig) (*ProxyBox, error) {
	// Set defaults
	if cfg.ListenAddr == "" {
		cfg.ListenAddr = "127.0.0.1:1080"
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = "panic" // Silent by default - only shows critical errors
	}

	// Create sing-box configuration
	config, err := createConfig(cfg)
	if err != nil {
		return nil, E.Cause(err, "create configuration")
	}

	// Create context with protocol handlers registered
	ctx := context.Background()
	ctx = include.Context(ctx)
	ctx, cancel := context.WithCancel(ctx)

	pb := &ProxyBox{
		ctx:      ctx,
		cancel:   cancel,
		config:   config,
		outbound: cfg.Outbound,
	}

	return pb, nil
}

// Start starts the proxy box
func (pb *ProxyBox) Start() error {
	if pb.instance != nil {
		return E.New("proxy box already started")
	}

	instance, err := box.New(box.Options{
		Context: pb.ctx,
		Options: pb.config,
	})
	if err != nil {
		return E.Cause(err, "create sing-box instance")
	}

	err = instance.Start()
	if err != nil {
		return E.Cause(err, "start sing-box")
	}

	pb.instance = instance
	return nil
}

// Stop stops the proxy box
func (pb *ProxyBox) Stop() error {
	if pb.instance == nil {
		return E.New("proxy box not started")
	}

	err := pb.instance.Close()
	pb.instance = nil
	if pb.cancel != nil {
		pb.cancel()
	}

	return err
}

// IsRunning returns true if the proxy box is currently running
func (pb *ProxyBox) IsRunning() bool {
	return pb.instance != nil
}

// Config returns the current configuration
func (pb *ProxyBox) Config() option.Options {
	return pb.config
}

// Outbound returns the outbound configuration
func (pb *ProxyBox) Outbound() option.Outbound {
	return pb.outbound
}

// ListenAddr returns the mixed proxy listen address (supports both SOCKS5 and HTTP)
func (pb *ProxyBox) ListenAddr() string {
	for _, inbound := range pb.config.Inbounds {
		if inbound.Type == "mixed" {
			if opts, ok := inbound.Options.(*option.HTTPMixedInboundOptions); ok {
				host := "127.0.0.1"
				if opts.Listen != nil {
					addr := netip.Addr(*opts.Listen)
					host = addr.String()
				}
				return fmt.Sprintf("%s:%d", host, opts.ListenPort)
			}
		}
	}
	return ""
}

// createConfig creates a sing-box configuration from the given config
func createConfig(cfg ProxyBoxConfig) (option.Options, error) {
	// Parse listen address
	host := strings.Split(cfg.ListenAddr, ":")[0]
	listenIP, err := netip.ParseAddr(host)
	if err != nil {
		listenIP = netip.MustParseAddr("127.0.0.1")
	}
	listenAddr := (*badoption.Addr)(&listenIP)

	return option.Options{
		Log: &option.LogOptions{
			Level:  cfg.LogLevel,
			Output: "stderr",
		},
		Inbounds: []option.Inbound{
			{
				Type: "mixed",
				Tag:  "mixed-in",
				Options: &option.HTTPMixedInboundOptions{
					ListenOptions: option.ListenOptions{
						Listen:     listenAddr,
						ListenPort: uint16(getPortOrDefault(1080, cfg.ListenAddr)),
					},
				},
			},
		},
		Outbounds: []option.Outbound{
			cfg.Outbound,
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
								Outbound: cfg.Outbound.Tag,
							},
						},
					},
				},
			},
			AutoDetectInterface: true,
		},
	}, nil
}

// getPort extracts port from host:port string
func getPortOrDefault(defaultPort int, hostPort string) int {
	parts := strings.Split(hostPort, ":")
	if len(parts) < 2 {
		return 1080 // Default SOCKS port
	}
	port := 1080
	fmt.Sscanf(parts[len(parts)-1], "%d", &port)
	return port
}
