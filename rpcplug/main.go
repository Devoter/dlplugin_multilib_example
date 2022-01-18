package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/Devoter/dlplugin_multilib_example/device"
	"github.com/Devoter/dlplugin_multilib_example/rpcplug/shared"
)

func main() {
	libraryFilename := flag.String("plug", "", "plugin library filename")
	setValue := flag.Int64("val", -120, "value to be set")

	flag.Parse()

	if *libraryFilename == "" {
		fmt.Fprintf(os.Stderr, "empty plugin filename\n")
		os.Exit(2)
	}

	// We're a host. Start by launching the plugin process.
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         shared.PluginMap,
		Cmd:             exec.Command("sh", "-c", os.Getenv("DEVICE_PLUGIN")),
		Logger:          hclog.NewNullLogger(),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
	})
	defer client.Kill()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("device")
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	papi := raw.(shared.Device)

	dev1, _ := papi.CreateDevice()
	defer papi.FreeDevice(dev1)
	dev2, _ := papi.CreateDevice()
	defer papi.FreeDevice(dev2)

	papi.DeviceSetValue(dev2, 32)
	papi.DevicePrint(dev2)

	value, err := papi.DeviceValue(dev2)

	fmt.Printf("Loaded value: %d, error [%v]\n", value, err)

	papi.DeviceSetValue(dev1, 24)
	papi.DevicePrint(dev1)

	papi.DeviceSetValue(dev1, int32(*setValue))

	value, err = papi.DeviceValue(dev1)

	fmt.Printf("Loaded value: %d, error [%v]\n", value, err)

	encoded, err := papi.GetDevice(dev1, false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not load an encoded device.Device, error=[%v]\n", err)
	} else {
		decodeBinary(encoded)
	}

	encoded, err = papi.GetDevice(dev2, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not load an encoded device.Device, error=[%v]\n", err)
	} else {
		decodeJSON(encoded)
	}

	value, err = papi.DeviceValue(uint64(32))

	fmt.Printf("Loaded value: %d, error [%v]\n", value, err)
}

func decodeBinary(encoded []byte) {
	var d device.Device

	if err := d.UnmarshalBinary(encoded); err != nil {
		fmt.Fprintf(os.Stderr, "could not decode a device.Device, error=[%v]\n", err)
	} else {
		fmt.Println("decoded value", d.Value())
	}
}

func decodeJSON(encoded []byte) {
	var d device.Device

	if err := json.Unmarshal(encoded, &d); err != nil {
		fmt.Fprintf(os.Stderr, "could not decode a device.Device, error=[%v]\n", err)
	} else {
		fmt.Println("decoded value", d.Value())
	}
}
