package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"plugin"

	"github.com/Devoter/dlplugin_multilib_example/device"
	"github.com/Devoter/dlplugin_multilib_example/stdplug/papi"
)

func main() {
	libraryFilename := flag.String("plug", "", "first plugin library filename")
	setValue := flag.Int64("val", -120, "value to be set")

	flag.Parse()

	if *libraryFilename == "" {
		fmt.Fprintf(os.Stderr, "empty plugin filename\n")
		os.Exit(2)
	}

	plug, err := plugin.Open(*libraryFilename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open a library, error=[%v]\n", err)
		os.Exit(1)
	}

	var papi1 papi.DevicePlugin

	if err := papi1.Init(plug.Lookup); err != nil {
		fmt.Fprintf(os.Stderr, "could initialize a plugin, error=[%v]\n", err)
		os.Exit(1)
	}

	dev1 := papi1.CreateDevice()
	dev2 := papi1.CreateDevice()

	fmt.Printf("Setting a value of dev2 to %d\n", 32)
	// papi1.Device_SetValue(dev2, 32)
	dev2.SetValue(32)

	fmt.Printf("Printing a value of dev2\n")
	// papi1.Device_Print(dev2)
	dev2.Print()

	fmt.Printf("Loading a value of dev2\n")
	// value := papi1.Device_Value(dev2)
	value := dev2.Value()
	fmt.Printf("Loaded value: %d\n", value)

	fmt.Printf("Setting a value of dev1 to %d\n", 24)
	// papi1.Device_SetValue(dev1, 24)
	dev1.SetValue(24)

	fmt.Printf("Printing a value of dev1\n")
	// papi1.Device_Print(dev1)
	dev1.Print()

	fmt.Printf("Setting a value of dev1 to %d\n", int32(*setValue))
	// papi1.Device_SetValue(dev1, int32(*setValue))
	dev1.SetValue(int32(*setValue))

	fmt.Printf("Loading a value of dev1\n")
	// value = papi1.Device_Value(dev1)
	value = dev1.Value()
	fmt.Printf("Loaded value: %d\n", value)

	fmt.Printf("Getting a binary state of dev1\n")
	// encoded, err := papi1.GetDevice(dev1, false)
	encoded, err := dev1.MarshalBinary()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not load an encoded device.Device by the reason: %v\n", err)
	} else {
		decodeBinary(encoded)
	}

	fmt.Printf("Getting a JSON state of dev2\n")
	// encoded, err = papi1.GetDevice(dev2, true)
	encoded, err = dev2.MarshalJSON()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not load an encoded device.Device by the reason: %v\n", err)
	} else {
		decodeJSON(encoded)
	}
}

func decodeBinary(encoded []byte) {
	var d device.Device

	if err := d.UnmarshalBinary(encoded); err != nil {
		fmt.Fprintf(os.Stderr, "could not decode a device.Device by the reason: %v\n", err)
	} else {
		fmt.Println("decoded value", d.Value())
	}
}

func decodeJSON(encoded []byte) {
	var d device.Device

	if err := json.Unmarshal(encoded, &d); err != nil {
		fmt.Fprintf(os.Stderr, "could not decode a device.Device by the reason: %v\n", err)
	} else {
		fmt.Println("decoded value", d.Value())
	}
}
