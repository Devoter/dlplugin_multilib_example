package shared

import (
	plugin "github.com/hashicorp/go-plugin"
	"golang.org/x/net/context"

	"github.com/Devoter/dlplugin_multilib_example/rpcplug/proto"
)

type GRPCClient struct {
	broker *plugin.GRPCBroker
	client proto.DeviceClient
}

func (m *GRPCClient) CreateDevice() (uint64, error) {
	resp, err := m.client.CreateDevice(context.Background(), &proto.CreateDeviceRequest{})
	if err != nil {
		return 0, err
	}

	return resp.Ptr, nil
}

func (m *GRPCClient) FreeDevice(ptr uint64) error {
	_, err := m.client.FreeDevice(context.Background(), &proto.FreeDeviceRequest{Ptr: ptr})
	return err
}

func (m *GRPCClient) GetDevice(ptr uint64, useJson bool) ([]byte, error) {
	resp, err := m.client.GetDevice(context.Background(), &proto.GetDeviceRequest{Ptr: ptr, UseJson: useJson})
	if err != nil {
		return nil, err
	}

	return resp.Encoded, nil
}

func (m *GRPCClient) DevicePrint(ptr uint64) error {
	_, err := m.client.DevicePrint(context.Background(), &proto.DevicePrintRequest{Ptr: ptr})
	return err
}

func (m *GRPCClient) DeviceValue(ptr uint64) (int32, error) {
	resp, err := m.client.DeviceValue(context.Background(), &proto.DeviceValueRequest{Ptr: ptr})
	if err != nil {
		return 0, err
	}

	return resp.Value, nil
}

func (m *GRPCClient) DeviceSetValue(ptr uint64, value int32) error {
	_, err := m.client.DeviceSetValue(context.Background(), &proto.DeviceSetValueRequest{Ptr: ptr, Value: value})
	return err
}

type GRPCServer struct {
	Impl   Device
	broker *plugin.GRPCBroker
	proto.UnimplementedDeviceServer
}

func (m *GRPCServer) CreateDevice(ctx context.Context, req *proto.CreateDeviceRequest) (*proto.CreateDeviceResponse, error) {
	ptr, err := m.Impl.CreateDevice()
	return &proto.CreateDeviceResponse{Ptr: ptr}, err
}

func (m *GRPCServer) FreeDevice(ctx context.Context, req *proto.FreeDeviceRequest) (*proto.FreeDeviceResponse, error) {
	err := m.Impl.FreeDevice(req.Ptr)
	return &proto.FreeDeviceResponse{}, err
}

func (m *GRPCServer) GetDevice(ctx context.Context, req *proto.GetDeviceRequest) (*proto.GetDeviceResponse, error) {
	encoded, err := m.Impl.GetDevice(req.Ptr, req.UseJson)
	return &proto.GetDeviceResponse{Encoded: encoded}, err
}

func (m *GRPCServer) DevicePrint(ctx context.Context, req *proto.DevicePrintRequest) (*proto.DevicePrintResponse, error) {
	err := m.Impl.DevicePrint(req.Ptr)
	return &proto.DevicePrintResponse{}, err
}

func (m *GRPCServer) DeviceValue(ctx context.Context, req *proto.DeviceValueRequest) (*proto.DeviceValueResponse, error) {
	value, err := m.Impl.DeviceValue(req.Ptr)
	return &proto.DeviceValueResponse{Value: value}, err
}

func (m *GRPCServer) DeviceSetValue(ctx context.Context, req *proto.DeviceSetValueRequest) (
	*proto.DeviceSetValueResponse, error) {
	err := m.Impl.DeviceSetValue(req.Ptr, req.Value)
	return &proto.DeviceSetValueResponse{}, err
}
