// Package shared contains shared data between the host and plugins.
package shared

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/hashicorp/go-plugin"

	"github.com/Devoter/dlplugin_multilib_example/grpcplug/proto"
)

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "DEVICE_PLUGIN",
	MagicCookieValue: "device",
}

// PluginMap is the map of plugins we can dispense.
var PluginMap = map[string]plugin.Plugin{
	"device": &DevicePlugin{},
}

// Device is the interface that we're exposing as a plugin.
type Device interface {
	CreateDevice() (uint64, error)
	FreeDevice(ptr uint64) error
	GetDevice(ptr uint64, useJson bool) ([]byte, error)
	DevicePrint(ptr uint64) error
	DeviceValue(ptr uint64) (int32, error)
	DeviceSetValue(ptr uint64, value int32) error
}

// This is the implementation of plugin.Plugin so we can serve/consume this.
// We also implement GRPCPlugin so that this plugin can be served over
// gRPC.
type DevicePlugin struct {
	plugin.NetRPCUnsupportedPlugin
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl Device
}

func (p *DevicePlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterDeviceServer(s, &GRPCServer{
		Impl:   p.Impl,
		broker: broker,
	})
	return nil
}

func (p *DevicePlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{
		client: proto.NewDeviceClient(c),
		broker: broker,
	}, nil
}

var _ plugin.GRPCPlugin = &DevicePlugin{}
