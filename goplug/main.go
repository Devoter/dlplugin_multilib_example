package main

/*
#include <stddef.h>
#include <stdint.h>

typedef void (*get_device_callback_t)(uintptr_t id, char *, size_t);

static void call_back(get_device_callback_t cb, uintptr_t id, char * data, size_t size)
{
	cb(id, data, size);
}
*/
import "C"

import (
	"fmt"
	"runtime/cgo"
	"sync"
	"unsafe"

	"github.com/Devoter/dlplugin_multilib_example/device"
)

var mx sync.RWMutex

//export create_device
func create_device() C.uintptr_t {
	dev := device.NewDevice()
	h := cgo.NewHandle(dev)

	return C.uintptr_t(h)
}

//export free_device
func free_device(ptr C.uintptr_t) C.int {
	mx.Lock()
	h, _, err := getDeviceHandle(ptr)

	if err != nil {
		mx.Unlock()
		return -1
	}

	h.Delete()
	mx.Unlock()

	return 0
}

//export get_device
func get_device(ptr C.uintptr_t, cbID C.uintptr_t, useJSON C.char, callback C.get_device_callback_t) C.int {
	mx.RLock()
	dev, err := getDevice(ptr)
	if err != nil {
		mx.RUnlock()
		return -1
	}

	var encoded []byte

	if useJSON != 0 {
		encoded, err = dev.MarshalJSON()
	} else {
		encoded, err = dev.MarshalBinary()
	}

	if err != nil {
		mx.RUnlock()
		return -2
	}

	mx.RUnlock()
	C.call_back(callback, cbID, (*C.char)(unsafe.Pointer(&encoded[0])), C.size_t(len(encoded)))

	return 0
}

//export device__print
func device__print(self C.uintptr_t) C.int {
	mx.RLock()
	dev, err := getDevice(self)
	if err != nil {
		mx.RUnlock()
		return -1
	}

	dev.Print()
	mx.RUnlock()

	return 0
}

//export device__value
func device__value(self C.uintptr_t, value *C.int32_t) C.int {
	mx.RLock()
	dev, err := getDevice(self)
	if err != nil {
		*value = 0
		mx.RUnlock()

		return -1
	}

	*value = C.int32_t(dev.Value())
	mx.RUnlock()

	return 0
}

//export device__set_value
func device__set_value(self C.uintptr_t, value C.int32_t) C.int {
	mx.Lock()
	dev, err := getDevice(self)
	if err != nil {
		mx.Unlock()
		return -1
	}

	dev.SetValue(int32(value))
	mx.Unlock()

	return 0
}

func getDeviceHandle(ptr C.uintptr_t) (cgo.Handle, *device.Device, error) {
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

func getDevice(ptr C.uintptr_t) (*device.Device, error) {
	_, dev, err := getDeviceHandle(ptr)

	return dev, err
}

func main() {}
