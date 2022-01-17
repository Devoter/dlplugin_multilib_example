package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/Devoter/dlplugin"

	"github.com/Devoter/dlplugin_multilib_example/device"
	"github.com/Devoter/dlplugin_multilib_example/papi"
)

func main() {
	libraryFilename := flag.String("plug1", "", "first plugin library filename")
	library2Filename := flag.String("plug2", "", "second plugin library filename")
	setValue := flag.Int64("val", -120, "value to be set")

	flag.Parse()

	if *libraryFilename == "" {
		fmt.Fprintf(os.Stderr, "empty plugin filename\n")
		os.Exit(2)
	}

	var papi1 papi.DevicePlugin

	plug1, err := dlplugin.Open(*libraryFilename, &papi1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open a library, error=[%v]\n", err)
		os.Exit(1)
	}
	defer plug1.Close()

	var papi2 papi.DevicePlugin

	plug2, err := dlplugin.Open(*library2Filename, &papi2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open a library, error=[%v]\n", err)
		os.Exit(1)
	}
	defer plug2.Close()

	dev1 := papi1.CreateDevice()
	defer papi1.FreeDevice(dev1)
	dev2 := papi1.CreateDevice()
	defer papi1.FreeDevice(dev2)
	dev3 := papi2.CreateDevice()
	defer papi2.FreeDevice(dev3)
	dev4 := papi2.CreateDevice()
	defer papi2.FreeDevice(dev4)

	papi1.Device_SetValue(dev2, 32)
	papi1.Device_Print(dev2)

	value, err := papi1.Device_Value(dev2)

	fmt.Printf("Loaded value: %d, error [%v]\n", value, err)

	papi1.Device_SetValue(dev1, 24)
	papi1.Device_Print(dev1)

	papi1.Device_SetValue(dev1, int32(*setValue))

	value, err = papi1.Device_Value(dev1)

	fmt.Printf("Loaded value: %d, error [%v]\n", value, err)

	encoded, err := papi1.GetDevice(dev1, false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not load an encoded device.Device, error=[%v]\n", err)
	} else {
		decodeBinary(encoded)
	}

	encoded, err = papi1.GetDevice(dev2, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not load an encoded device.Device, error=[%v]\n", err)
	} else {
		decodeJSON(encoded)
	}

	value, err = papi1.Device_Value(uintptr(32))

	fmt.Printf("Loaded value: %d, error [%v]\n", value, err)

	papi2.Device_SetValue(dev3, 46)
	value, err = papi2.Device_Value(dev3)

	fmt.Printf("Loaded value: %d, error [%v]\n", value, err)

	papi2.Device_SetValue(dev4, 42)

	encoded, err = papi2.GetDevice(dev3, false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not load an encoded device.Device, error=[%v]\n", err)
	} else {
		decodeBinary(encoded)
	}

	encoded, err = papi2.GetDevice(dev4, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not load an encoded device.Device, error=[%v]\n", err)
	} else {
		decodeJSON(encoded)
	}
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
