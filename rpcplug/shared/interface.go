package shared

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
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
// We also implement RPCPlugin so that this plugin can be served.
type DevicePlugin struct {
	plugin.NetRPCUnsupportedPlugin
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl Device
}

func (p *DevicePlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &RPCServer{Impl: p.Impl}, nil
}

func (p *DevicePlugin) Client(_ *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &RPCClient{client: c}, nil
}
