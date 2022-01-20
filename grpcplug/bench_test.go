package main_test

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/Devoter/dlplugin_multilib_example/device"
	"github.com/Devoter/dlplugin_multilib_example/grpcplug/shared"
)

type benchDataItem struct {
	api shared.Device
	dev uint64
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
		client := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig: shared.Handshake,
			Plugins:         shared.PluginMap,
			Cmd:             exec.Command("sh", "-c", lib),
			Logger:          hclog.NewNullLogger(),
			AllowedProtocols: []plugin.Protocol{
				plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
		})
		// defer client.Kill()

		rpcClient, err := client.Client()
		if err != nil {
			fmt.Println("Error:", err.Error())
			os.Exit(1)
		}

		raw, err := rpcClient.Dispense("device")
		if err != nil {
			fmt.Println("Error:", err.Error())
			os.Exit(1)
		}

		papi := raw.(shared.Device)
		data = &benchDataItem{api: papi}
		benchData[lib] = data
	}

	if data.dev == 0 {
		dev, err := data.api.CreateDevice()
		if err != nil {
			b.Fatalf("could not instantiate a device by the reason: %v\n", err)
			return nil
		}

		data.dev = dev
	}

	return data
}

func BenchmarkSetGetGoPlug(b *testing.B) {
	benchItem := accessDevice(b, "./device-go-grpc")
	if benchItem == nil {
		b.Fatalf("could not access a device\n")
	}

	api := benchItem.api
	dev := benchItem.dev

	var value int32
	var err error
	i := int32(b.N)

	b.ResetTimer()

	if err = api.DeviceSetValue(dev, i); err != nil {
		b.Fatalf("could not set a device value by the reason: %v\n", err)
	}

	if value, err = api.DeviceValue(dev); err != nil {
		b.Fatalf("could not get a device value by the reason: %v\n", err)
	}

	if i != value {
		b.Errorf("value should be [%d], but got [%d]", i, value)
	}
}

func BenchmarkReadingGoPlug(b *testing.B) {
	benchItem := accessDevice(b, "./device-go-grpc")
	if benchItem == nil {
		b.Fatalf("iteration: %d. could not access a device\n", b.N)
	}

	api := benchItem.api
	dev := benchItem.dev

	var err error
	i := int32(b.N)

	b.ResetTimer()

	api.DeviceSetValue(dev, i)
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
	benchItem := accessDevice(b, "./device-go-grpc")
	if benchItem == nil {
		b.Fatalf("iteration: %d. could not access a device\n", b.N)
	}

	api := benchItem.api
	dev := benchItem.dev

	var err error
	i := int32(b.N)

	b.ResetTimer()

	api.DeviceSetValue(dev, i)
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
