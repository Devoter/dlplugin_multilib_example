package main_test

import (
	"plugin"
	"sync"
	"testing"

	"github.com/Devoter/dlplugin_multilib_example/device"
	"github.com/Devoter/dlplugin_multilib_example/stdplug/papi"
)

type benchDataItem struct {
	api  *papi.DevicePlugin
	plug *plugin.Plugin
	dev  *device.Device
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
		plug, err := plugin.Open(lib)
		if err != nil {
			b.Fatalf("could not load a plugin by the reason: %v\n", err)
			return nil
		}

		var api papi.DevicePlugin

		if err := api.Init(plug.Lookup); err != nil {
			b.Fatalf("could not initialize a plugin by the reason: %v\n", err)
			return nil
		}

		data = &benchDataItem{api: &api, plug: plug}
		benchData[lib] = data
	}

	if data.dev == nil {
		data.dev = data.api.CreateDevice()
	}

	return data
}

func BenchmarkSetGetStdGoPlug(b *testing.B) {
	benchItem := accessDevice(b, "./plug/stdplug.so")
	if benchItem == nil {
		b.Fatalf("could not access a device\n")
	}

	api := benchItem.api
	dev := benchItem.dev

	var value int32
	i := int32(b.N)

	b.ResetTimer()
	api.Device_SetValue(dev, i)
	value = api.Device_Value(dev)

	if i != value {
		b.Errorf("value should be [%d], but got [%d]", i, value)
	}
}

func BenchmarkReadingStdGoPlug(b *testing.B) {
	benchItem := accessDevice(b, "./plug/stdplug.so")
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

func BenchmarkReadingJSONStdGoPlug(b *testing.B) {
	benchItem := accessDevice(b, "./plug/stdplug.so")
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
