package main

import (
	"fmt"
	"os"

	"github.com/Devoter/dlplugin_multilib_example/device"
)

func CreateDevice() *device.Device {
	return device.NewDevice()
}

func GetDevice(dev *device.Device, useJSON bool, callback func([]byte)) error {
	var err error
	var encoded []byte

	if useJSON {
		encoded, err = dev.MarshalJSON()
	} else {
		encoded, err = dev.MarshalBinary()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "could not encode a device, error=[%v]\n", err)
		return err
	}

	callback(encoded)

	return nil
}

func Device__Print(dev *device.Device) {
	dev.Print()
}

func Device__Value(dev *device.Device) int32 {
	return dev.Value()
}

func Device__SetValue(dev *device.Device, value int32) {
	dev.SetValue(int32(value))
}
