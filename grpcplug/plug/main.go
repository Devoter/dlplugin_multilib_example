package main

import (
	"fmt"
	"os"
	"runtime/cgo"
	"sync"

	"github.com/hashicorp/go-plugin"

	"github.com/Devoter/dlplugin_multilib_example/device"
	"github.com/Devoter/dlplugin_multilib_example/grpcplug/shared"
)

type Device struct {
	mx sync.RWMutex
}

func (d *Device) CreateDevice() (uint64, error) {
	dev := device.NewDevice()
	h := cgo.NewHandle(dev)

	return uint64(h), nil
}

func (d *Device) FreeDevice(ptr uint64) error {
	d.mx.Lock()
	h, _, err := getDeviceHandle(ptr)

	if err != nil {
		d.mx.Unlock()
		return err
	}

	h.Delete()
	d.mx.Unlock()

	return nil
}

func (d *Device) GetDevice(ptr uint64, useJSON bool) ([]byte, error) {
	d.mx.RLock()
	dev, err := getDevice(ptr)
	if err != nil {
		d.mx.RUnlock()
		return nil, err
	}

	var encoded []byte

	if useJSON {
		encoded, err = dev.MarshalJSON()
	} else {
		encoded, err = dev.MarshalBinary()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "could not encode a device, error=[%v]\n", err)
		d.mx.RUnlock()
		return nil, err
	}

	d.mx.RUnlock()

	return encoded, nil
}

func (d *Device) DevicePrint(self uint64) error {
	d.mx.RLock()
	dev, err := getDevice(self)
	if err != nil {
		d.mx.RUnlock()
		return err
	}

	dev.Print()
	d.mx.RUnlock()

	return nil
}

func (d *Device) DeviceValue(self uint64) (int32, error) {
	d.mx.RLock()
	dev, err := getDevice(self)
	if err != nil {
		d.mx.RUnlock()
		return 0, err
	}

	value := dev.Value()

	d.mx.RUnlock()

	return value, nil
}

func (d *Device) DeviceSetValue(self uint64, value int32) error {
	d.mx.Lock()
	dev, err := getDevice(self)
	if err != nil {
		d.mx.Unlock()
		return err
	}

	dev.SetValue(value)
	d.mx.Unlock()

	return nil
}

func getDeviceHandle(ptr uint64) (cgo.Handle, *device.Device, error) {
	h := cgo.Handle(ptr)
	var dev *device.Device
	var err error

	func() {
		defer func() {
			if msg := recover(); msg != nil {
				err = fmt.Errorf("%v", msg)
			}
		}()

		var ok bool

		dev, ok = h.Value().(*device.Device)
		if !ok {
			err = fmt.Errorf("unexpected value type")
		}
	}()

	if err != nil {
		return h, nil, err
	}

	return h, dev, nil
}

func getDevice(ptr uint64) (*device.Device, error) {
	_, dev, err := getDeviceHandle(ptr)

	return dev, err
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"device": &shared.DevicePlugin{Impl: &Device{}},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
