package shared

import (
	"net/rpc"

	"github.com/Devoter/dlplugin_multilib_example/rpcplug/proto"
)

type RPCClient struct {
	client *rpc.Client
}

func (m *RPCClient) CreateDevice() (uint64, error) {
	var resp uint64
	err := m.client.Call("Plugin.CreateDevice", new(interface{}), &resp)
	if err != nil {
		return 0, err
	}

	return resp, nil
}

func (m *RPCClient) FreeDevice(ptr uint64) error {
	return m.client.Call("Plugin.FreeDevice", ptr, &struct{}{})
}

func (m *RPCClient) GetDevice(ptr uint64, useJson bool) ([]byte, error) {
	var resp []byte

	err := m.client.Call("Plugin.GetDevice", &proto.GetDeviceRequest{Ptr: ptr, UseJson: useJson}, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m *RPCClient) DevicePrint(ptr uint64) error {
	return m.client.Call("Plugin.DevicePrint", ptr, &struct{}{})
}

func (m *RPCClient) DeviceValue(ptr uint64) (int32, error) {
	var resp int32

	err := m.client.Call("Plugin.DeviceValue", ptr, &resp)
	if err != nil {
		return 0, err
	}

	return resp, nil
}

func (m *RPCClient) DeviceSetValue(ptr uint64, value int32) error {
	return m.client.Call("Plugin.DeviceSetValue", &proto.DeviceSetValueRequest{Ptr: ptr, Value: value}, &struct{}{})
}

type RPCServer struct {
	Impl Device
}

func (m *RPCServer) CreateDevice(_ interface{}, resp *uint64) error {
	ptr, err := m.Impl.CreateDevice()
	*resp = ptr

	return err
}

func (m *RPCServer) FreeDevice(ptr uint64, _ *struct{}) error {
	return m.Impl.FreeDevice(ptr)
}

func (m *RPCServer) GetDevice(req *proto.GetDeviceRequest, resp *[]byte) error {
	encoded, err := m.Impl.GetDevice(req.Ptr, req.UseJson)
	*resp = encoded

	return err
}

func (m *RPCServer) DevicePrint(ptr uint64, _ *struct{}) error {
	return m.Impl.DevicePrint(ptr)
}

func (m *RPCServer) DeviceValue(ptr uint64, resp *int32) error {
	value, err := m.Impl.DeviceValue(ptr)
	*resp = value

	return err
}

func (m *RPCServer) DeviceSetValue(req *proto.DeviceSetValueRequest, _ *struct{}) error {
	return m.Impl.DeviceSetValue(req.Ptr, req.Value)
}
