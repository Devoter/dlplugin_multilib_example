package papi

import (
	"plugin"

	"github.com/Devoter/dlplugin_multilib_example/device"
)

type DevicePlugin struct {
	createDevice    func() *device.Device
	getDevice       func(ptr *device.Device, useJson bool) (encoded []byte, err error)
	device_Print    func(self *device.Device)
	device_Value    func(self *device.Device) (value int32)
	device_SetValue func(self *device.Device, value int32)
}

func (dev *DevicePlugin) CreateDevice() *device.Device {
	return dev.createDevice()
}

func (dev *DevicePlugin) GetDevice(ptr *device.Device, useJson bool) (encoded []byte, err error) {
	return dev.getDevice(ptr, useJson)
}

func (dev *DevicePlugin) Device_Print(self *device.Device) {
	dev.device_Print(self)
}

func (dev *DevicePlugin) Device_Value(self *device.Device) (value int32) {
	return dev.device_Value(self)
}

func (dev *DevicePlugin) Device_SetValue(self *device.Device, value int32) {
	dev.device_SetValue(self, value)
}

func (dev *DevicePlugin) Init(lookup func(symName string) (plugin.Symbol, error)) error {
	createDeviceSym, err := lookup("CreateDevice")
	if err != nil {
		return err
	}

	getDeviceSym, err := lookup("GetDevice")
	if err != nil {
		return err
	}

	devicePrintSym, err := lookup("Device__Print")
	if err != nil {
		return err
	}

	deviceValueSym, err := lookup("Device__Value")
	if err != nil {
		return err
	}

	deviceSetValueSym, err := lookup("Device__SetValue")
	if err != nil {
		return err
	}

	dev.createDevice = func() *device.Device {
		return createDeviceSym.(func() *device.Device)()
	}

	dev.getDevice = func(ptr *device.Device, useJson bool) ([]byte, error) {
		var encoded []byte
		cb := func(data []byte) { encoded = data }

		err := getDeviceSym.(func(*device.Device, bool, func(data []byte)) error)(ptr, useJson, cb)
		if err != nil {
			return nil, err
		}

		return encoded, nil
	}

	dev.device_Print = func(self *device.Device) {
		devicePrintSym.(func(*device.Device))(self)
	}

	dev.device_Value = func(self *device.Device) int32 {
		return deviceValueSym.(func(*device.Device) int32)(self)
	}

	dev.device_SetValue = func(self *device.Device, value int32) {
		deviceSetValueSym.(func(*device.Device, int32))(self, value)
	}

	return nil
}
