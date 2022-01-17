package main_test

import (
	"sync"
	"testing"

	"github.com/Devoter/dlplugin"
	"github.com/Devoter/dlplugin_multilib_example/device"
	"github.com/Devoter/dlplugin_multilib_example/papi"
)

type benchDataItem struct {
	api  *papi.DevicePlugin
	plug *dlplugin.Plugin
	dev  uintptr
}

var benchData map[string]*benchDataItem
var benchDataMx sync.Mutex

func accessDevice(b *testing.B, lib string) *benchDataItem {
	benchDataMx.Lock()
	defer benchDataMx.Unlock()

	if benchData == nil {
		benchData = map[string]*benchDataItem{}
	}

	data, ok := benchData[lib]
	if !ok {
		var api papi.DevicePlugin

		plug, err := dlplugin.Open(lib, &api)
		if err != nil {
			b.Fatalf("could not load a plugin by the reason: %v\n", err)
			return nil
		}

		data = &benchDataItem{api: &api, plug: plug}
		benchData[lib] = data
	}

	if data.dev == 0 {
		dev := data.api.CreateDevice()
		if dev <= 0 {
			b.Fatalf("could not instantiate a device by the reason: %d\n", dev)
			return nil
		}

		data.dev = dev
	}

	return data
}

func BenchmarkSetGetGoPlug(b *testing.B) {
	benchItem := accessDevice(b, "./goplug/libgoplug.so")
	if benchItem == nil {
		b.Fatalf("could not access a device\n")
	}

	api := benchItem.api
	dev := benchItem.dev

	var value int32
	var err error
	i := int32(b.N)

	b.ResetTimer()

	if err = api.Device_SetValue(dev, i); err != nil {
		b.Fatalf("could not set a device value by the reason: %v\n", err)
	}

	if value, err = api.Device_Value(dev); err != nil {
		b.Fatalf("could not get a device value by the reason: %v\n", err)
	}

	if i != value {
		b.Errorf("value should be [%d], but got [%d]", i, value)
	}
}

func BenchmarkSetGetCPPPlug(b *testing.B) {
	benchItem := accessDevice(b, "./cppplug/libcppplug.so")
	if benchItem == nil {
		b.Fatalf("could not access a device\n")
	}

	api := benchItem.api
	dev := benchItem.dev

	var value int32
	var err error
	i := int32(b.N)

	b.ResetTimer()

	if err = api.Device_SetValue(dev, i); err != nil {
		b.Fatalf("could not set a device value by the reason: %v\n", err)
	}

	if value, err = api.Device_Value(dev); err != nil {
		b.Fatalf("could not get a device value by the reason: %v\n", err)
	}

	if i != value {
		b.Errorf("value should be [%d], but got [%d]", i, value)
	}
}

func BenchmarkSetGetGoPure(b *testing.B) {
	var dev device.Device
	var value int32
	i := int32(b.N)

	b.ResetTimer()

	dev.SetValue(i)
	value = dev.Value()

	if i != value {
		b.Errorf("value should be [%d], but got [%d]", i, value)
	}
}

func BenchmarkReadingGoPlug(b *testing.B) {
	benchItem := accessDevice(b, "./goplug/libgoplug.so")
	if benchItem == nil {
		b.Fatalf("iteration: %d. could not access a device\n", b.N)
	}

	api := benchItem.api
	dev := benchItem.dev

	var err error
	i := int32(b.N)

	b.ResetTimer()

	api.Device_SetValue(dev, i)
	encoded, err := api.GetDevice(dev, false)
	if err != nil {
		b.Fatalf("iteration: %d. could not get an encoded device state by the reason: %v\n", b.N, err)
	}

	var d device.Device

	if err = d.UnmarshalBinary(encoded); err != nil {
		b.Fatalf("iteration: %d. Could not parse an encoded device state by the reason: %v\n", b.N, err)
	}

	if i != d.Value() {
		b.Errorf("iteration: %d. value should be [%d], but got [%d]", b.N, i, d.Value())
	}
}

func BenchmarkReadingCPPPlug(b *testing.B) {
	benchItem := accessDevice(b, "./cppplug/libcppplug.so")
	if benchItem == nil {
		b.Fatalf("iteration: %d. could not access a device\n", b.N)
	}

	api := benchItem.api
	dev := benchItem.dev

	var err error
	i := int32(b.N)

	b.ResetTimer()

	api.Device_SetValue(dev, i)
	encoded, err := api.GetDevice(dev, false)
	if err != nil {
		b.Fatalf("iteration: %d. could not get an encoded device state by the reason: %v\n", b.N, err)
	}

	var d device.Device

	if err = d.UnmarshalBinary(encoded); err != nil {
		b.Fatalf("iteration: %d. Could not parse an encoded device state by the reason: %v\n", b.N, err)
	}

	if i != d.Value() {
		b.Errorf("iteration: %d. value should be [%d], but got [%d]", b.N, i, d.Value())
	}
}

func BenchmarkReadingJSONGoPlug(b *testing.B) {
	benchItem := accessDevice(b, "./goplug/libgoplug.so")
	if benchItem == nil {
		b.Fatalf("iteration: %d. could not access a device\n", b.N)
	}

	api := benchItem.api
	dev := benchItem.dev

	var err error
	i := int32(b.N)

	b.ResetTimer()

	api.Device_SetValue(dev, i)
	encoded, err := api.GetDevice(dev, true)
	if err != nil {
		b.Fatalf("iteration: %d. could not get an encoded device state by the reason: %v\n", b.N, err)
	}

	var d device.Device

	if err = d.UnmarshalJSON(encoded); err != nil {
		b.Fatalf("iteration: %d. Could not parse an encoded device state by the reason: %v\n", b.N, err)
	}

	if i != d.Value() {
		b.Errorf("iteration: %d. value should be [%d], but got [%d]", b.N, i, d.Value())
	}
}

func BenchmarkReadingJSONCPPPlug(b *testing.B) {
	benchItem := accessDevice(b, "./cppplug/libcppplug.so")
	if benchItem == nil {
		b.Fatalf("iteration: %d. could not access a device\n", b.N)
	}

	api := benchItem.api
	dev := benchItem.dev

	var err error
	i := int32(b.N)

	b.ResetTimer()

	api.Device_SetValue(dev, i)
	encoded, err := api.GetDevice(dev, true)
	if err != nil {
		b.Fatalf("iteration: %d. could not get an encoded device state by the reason: %v\n", b.N, err)
	}

	var d device.Device

	if err = d.UnmarshalJSON(encoded); err != nil {
		b.Fatalf("iteration: %d. Could not parse an encoded device state by the reason: %v\n", b.N, err)
	}

	if i != d.Value() {
		b.Errorf("iteration: %d. value should be [%d], but got [%d]", b.N, i, d.Value())
	}
}
